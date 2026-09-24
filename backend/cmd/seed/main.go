package main

import (
	"context"
	"log"
	"os"

	assigneeApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	ticketApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	infraPg "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres"
)

func main() {
	log.Println("[INFO] Iniciando processo de seed do banco de dados...")

	db, err := infraPg.NewDatabase("")
	if err != nil {
		log.Fatalf("[FATAL] Falha ao conectar ao banco de dados PostgreSQL: %v\n", err)
	}
	defer db.Close()

	ctx := context.Background()

	assigneeRepo := infraPg.NewAssigneeRepository(db)
	ticketRepo := infraPg.NewTicketRepository(db)

	assigneeCreateUC := assigneeApp.NewCreateUseCase(assigneeRepo)
	assigneeListUC := assigneeApp.NewListUseCase(assigneeRepo, ticketRepo)

	resolveAssigneeUC := ticketApp.NewResolveAssigneeUseCase(assigneeRepo)
	ticketCreateUC := ticketApp.NewCreateUseCase(ticketRepo, resolveAssigneeUC)
	ticketListUC := ticketApp.NewListUseCase(ticketRepo)

	// 1. Seed de atendentes padrão
	assignees, err := assigneeListUC.Execute(ctx)
	if err != nil {
		log.Fatalf("[FATAL] Falha ao consultar atendentes: %v\n", err)
	}

	if len(assignees) == 0 {
		initialAssignees := []string{"Carlos Silva", "Ana Souza", "Bruno Santos"}
		for _, name := range initialAssignees {
			if err := assigneeCreateUC.Execute(ctx, assigneeApp.CreateInput{Name: name}); err != nil {
				log.Printf("[ERRO] Falha ao criar atendente %s: %v\n", name, err)
				os.Exit(1)
			}
			log.Printf("[INFO] Atendente criado: %s\n", name)
		}
	} else {
		log.Printf("[INFO] Banco já possui %d atendente(s) cadastrado(s). Pulando seed de atendentes.\n", len(assignees))
	}

	// 2. Seed de chamados de demonstração
	tickets, err := ticketListUC.Execute(ctx, ticketApp.ListInput{Page: 1, PageSize: 1})
	if err != nil {
		log.Fatalf("[FATAL] Falha ao consultar chamados: %v\n", err)
	}

	if tickets.TotalItems == 0 {
		sampleTickets := []ticketApp.CreateInput{
			{
				Title:       "Lentidão ao processar pagamentos Pix",
				Description: "Clientes relatam demora excessiva de mais de 1 minuto para confirmação",
				Priority:    "high",
				Assignee:    "Carlos Silva",
			},
			{
				Title:       "Correção de texto na tela de perfil",
				Description: "Corrigir erro ortográfico na instrução do formulário",
				Priority:    "low",
				Assignee:    "Ana Souza",
			},
			{
				Title:       "Falha na sincronização do estoque",
				Description: "Integração do ERP retornou código 500 no lote das 14h",
				Priority:    "medium",
				AutoAssign:  true,
			},
		}

		for _, st := range sampleTickets {
			ticket, err := ticketCreateUC.Execute(ctx, st)
			if err != nil {
				log.Printf("[ERRO] Falha ao criar chamado '%s': %v\n", st.Title, err)
				os.Exit(1)
			}
			assigneeName := "Nenhum"
			if ticket.Assignee != "" {
				assigneeName = ticket.Assignee
			}
			log.Printf("[INFO] Chamado criado: '%s' [Prioridade: %s, Responsável: %s]\n", ticket.Title, ticket.Priority, assigneeName)
		}
	} else {
		log.Printf("[INFO] Banco já possui %d chamado(s) cadastrado(s). Pulando seed de chamados.\n", tickets.TotalItems)
	}

	log.Println("[INFO] Seed concluído com sucesso!")
}
