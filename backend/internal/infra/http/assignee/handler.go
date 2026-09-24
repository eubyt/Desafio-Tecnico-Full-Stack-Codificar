package assignee

import (
	"encoding/json"
	"net/http"

	assigneeApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	httpResponse "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/response"
)

type Handler struct {
	createUC *assigneeApp.CreateUseCase
	listUC   *assigneeApp.ListUseCase
}

func NewHandler(createUC *assigneeApp.CreateUseCase, listUC *assigneeApp.ListUseCase) *Handler {
	return &Handler{
		createUC: createUC,
		listUC:   listUC,
	}
}

// List godoc
// @Summary      Lista todos os responsáveis
// @Description  Retorna todos os atendentes e a contagem de chamados em aberto (status open ou in_progress)
// @Tags         Assignees
// @Produce      json
// @Success      200  {array}   assigneeApp.AssigneeOutput
// @Failure      500  {object}  httpResponse.ErrorResponse
// @Router       /api/assignees [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	output, err := h.listUC.Execute(r.Context())
	if err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusOK, output)
}

// Create godoc
// @Summary      Cadastra um novo responsável
// @Description  Registra um novo membro da equipe de suporte
// @Tags         Assignees
// @Accept       json
// @Produce      json
// @Param        payload  body      assigneeApp.CreateAssigneeInput  true  "Dados do Responsável"
// @Success      201      {object}  map[string]string        "Mensagem de sucesso"
// @Failure      400      {object}  httpResponse.ErrorResponse
// @Failure      409      {object}  httpResponse.ErrorResponse
// @Failure      500      {object}  httpResponse.ErrorResponse
// @Router       /api/assignees [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input assigneeApp.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpResponse.JSON(w, http.StatusBadRequest, httpResponse.ErrorResponse{Error: "invalid request payload"})
		return
	}

	if err := h.createUC.Execute(r.Context(), input); err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusCreated, map[string]string{"message": "assignee created successfully"})
}
