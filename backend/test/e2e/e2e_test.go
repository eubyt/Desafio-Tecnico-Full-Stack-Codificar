//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	assigneeApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	ticketApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	infraHttp "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http"
	httpAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/assignee"
	httpTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/ticket"
	infraPg "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startE2EEnvironment(t *testing.T) (*httptest.Server, func()) {
	ctx := context.Background()

	// 1. Inicia o container PostgreSQL temporário via Testcontainers
	pgContainer, err := tcPostgres.Run(ctx,
		"postgres:16-alpine",
		tcPostgres.WithDatabase("ticket_system"),
		tcPostgres.WithUsername("ticket_user"),
		tcPostgres.WithPassword("ticket_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// 2. Conecta o pgxpool e roda o schema.sql DDL real automaticamente
	pool, err := infraPg.NewDatabase(connStr)
	require.NoError(t, err)

	// 3. Monta todos os componentes reais do backend
	ticketRepo := infraPg.NewTicketRepository(pool)
	assigneeRepo := infraPg.NewAssigneeRepository(pool)

	assigneeCreateUC := assigneeApp.NewCreateUseCase(assigneeRepo)
	assigneeListUC := assigneeApp.NewListUseCase(assigneeRepo, ticketRepo)

	resolveAssigneeUC := ticketApp.NewResolveAssigneeUseCase(assigneeRepo)
	ticketCreateUC := ticketApp.NewCreateUseCase(ticketRepo, resolveAssigneeUC)
	ticketUpdateUC := ticketApp.NewUpdateUseCase(ticketRepo, resolveAssigneeUC)
	ticketGetByIDUC := ticketApp.NewGetByIDUseCase(ticketRepo)
	ticketListUC := ticketApp.NewListUseCase(ticketRepo)

	ticketHandler := httpTicket.NewHandler(ticketCreateUC, ticketUpdateUC, ticketGetByIDUC, ticketListUC)
	assigneeHandler := httpAssignee.NewHandler(assigneeCreateUC, assigneeListUC)
	router := infraHttp.NewRouter(ticketHandler, assigneeHandler)

	// 4. Inicia servidor HTTP real
	ts := httptest.NewServer(router)

	cleanup := func() {
		ts.Close()
		pool.Close()
		_ = pgContainer.Terminate(ctx)
	}

	return ts, cleanup
}

func TestE2E_TicketAndAssigneeFullFlow(t *testing.T) {
	ts, cleanup := startE2EEnvironment(t)
	defer cleanup()

	client := ts.Client()

	// 1. Healthcheck
	res, err := client.Get(ts.URL + "/health")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// 2. Cadastrar atendentes "Ana Souza" e "Carlos Silva" via POST /api/assignees
	for _, name := range []string{"Ana Souza", "Carlos Silva"} {
		body, _ := json.Marshal(map[string]string{"name": name})
		res, err := client.Post(ts.URL+"/api/assignees", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	}

	// 3. Listar atendentes via GET /api/assignees (ambos com 0 chamados abertos)
	res, err = client.Get(ts.URL + "/api/assignees")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var assignees []map[string]any
	err = json.NewDecoder(res.Body).Decode(&assignees)
	require.NoError(t, err)
	assert.Len(t, assignees, 2)
	assert.Equal(t, "Ana Souza", assignees[0]["name"])
	assert.Equal(t, float64(0), assignees[0]["open_tickets_count"])

	// 4. Criar Chamado 1 com auto_assign: true -> vai para "Ana Souza" (desempate alfabético)
	t1Body, _ := json.Marshal(map[string]any{
		"title":       "Falha na emissão de NF",
		"description": "Ao emitir nota fiscal retorna erro 502",
		"priority":    "high",
		"auto_assign": true,
	})
	res, err = client.Post(ts.URL+"/api/tickets", "application/json", bytes.NewBuffer(t1Body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var t1Created map[string]any
	err = json.NewDecoder(res.Body).Decode(&t1Created)
	require.NoError(t, err)
	assert.NotEmpty(t, t1Created["id"])
	assert.Equal(t, "Ana Souza", t1Created["assignee"])
	assert.Equal(t, "open", t1Created["status"])
	t1ID := t1Created["id"].(string)

	// 5. Criar Chamado 2 com auto_assign: true -> deve ir para "Carlos Silva" (pois Ana tem 1 chamado em aberto)
	t2Body, _ := json.Marshal(map[string]any{
		"title":       "Lentidão no painel de relatórios",
		"description": "Carregamento demora mais de 20 segundos",
		"priority":    "medium",
		"auto_assign": true,
	})
	res, err = client.Post(ts.URL+"/api/tickets", "application/json", bytes.NewBuffer(t2Body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var t2Created map[string]any
	err = json.NewDecoder(res.Body).Decode(&t2Created)
	require.NoError(t, err)
	assert.Equal(t, "Carlos Silva", t2Created["assignee"])

	// 6. Consultar Chamado 1 por ID via GET /api/tickets/{id}
	res, err = client.Get(ts.URL + "/api/tickets/" + t1ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var t1Get map[string]any
	err = json.NewDecoder(res.Body).Decode(&t1Get)
	require.NoError(t, err)
	assert.Equal(t, "Falha na emissão de NF", t1Get["title"])
	assert.Equal(t, "Ana Souza", t1Get["assignee"])

	// 7. Atualizar Chamado 1 via PUT /api/tickets/{id} (mover para in_progress e reatribuir para Carlos Silva)
	updatePayload, _ := json.Marshal(map[string]any{
		"title":       "Falha na emissão de NF (Em Análise)",
		"description": "Timeout na SEFAZ identificado pelo time de infra",
		"priority":    "high",
		"status":      "in_progress",
		"assignee":    "Carlos Silva",
	})
	req, _ := http.NewRequest(http.MethodPut, ts.URL+"/api/tickets/"+t1ID, bytes.NewBuffer(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	res, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var t1Updated map[string]any
	err = json.NewDecoder(res.Body).Decode(&t1Updated)
	require.NoError(t, err)
	assert.Equal(t, "in_progress", t1Updated["status"])
	assert.Equal(t, "Carlos Silva", t1Updated["assignee"])

	// 8. Listar chamados com filtros e ordenação via GET /api/tickets?priority=high&order=desc
	res, err = client.Get(ts.URL + "/api/tickets?priority=high&sort_by=created_at&order=desc")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var listResult map[string]any
	err = json.NewDecoder(res.Body).Decode(&listResult)
	require.NoError(t, err)
	assert.Equal(t, float64(1), listResult["total_items"])
	items := listResult["items"].([]any)
	assert.Len(t, items, 1)

	// 9. Verificar contagem de abertos atualizada no GET /api/assignees
	// Carlos agora tem 2 abertos (o que criou e o que foi transferido) e Ana tem 0 abertos
	res, err = client.Get(ts.URL + "/api/assignees")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var assigneesUpdated []map[string]any
	err = json.NewDecoder(res.Body).Decode(&assigneesUpdated)
	require.NoError(t, err)
	assert.Equal(t, float64(0), assigneesUpdated[0]["open_tickets_count"]) // Ana Souza
	assert.Equal(t, float64(2), assigneesUpdated[1]["open_tickets_count"]) // Carlos Silva
}
