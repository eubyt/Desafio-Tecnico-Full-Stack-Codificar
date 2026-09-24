package ticket

import (
	"encoding/json"
	"net/http"
	"strconv"

	ticketApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	httpResponse "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	createUC  *ticketApp.CreateUseCase
	updateUC  *ticketApp.UpdateUseCase
	getByIDUC *ticketApp.GetByIDUseCase
	listUC    *ticketApp.ListUseCase
}

func NewHandler(
	createUC *ticketApp.CreateUseCase,
	updateUC *ticketApp.UpdateUseCase,
	getByIDUC *ticketApp.GetByIDUseCase,
	listUC *ticketApp.ListUseCase,
) *Handler {
	return &Handler{
		createUC:  createUC,
		updateUC:  updateUC,
		getByIDUC: getByIDUC,
		listUC:    listUC,
	}
}

// Create godoc
// @Summary      Cria um novo chamado
// @Description  Registra um chamado com atribuição manual ou automática (auto_assign: true)
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        payload  body      ticketApp.CreateTicketInput  true  "Dados do Chamado"
// @Success      201      {object}  ticketApp.TicketOutput
// @Failure      400      {object}  httpResponse.ErrorResponse
// @Failure      404      {object}  httpResponse.ErrorResponse
// @Failure      500      {object}  httpResponse.ErrorResponse
// @Router       /api/tickets [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input ticketApp.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpResponse.JSON(w, http.StatusBadRequest, httpResponse.ErrorResponse{Error: "invalid request payload"})
		return
	}

	output, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusCreated, output)
}

// Update godoc
// @Summary      Atualiza responsável ou status de um chamado
// @Description  Atualiza apenas o status e/ou o responsável (manual ou automático) de um chamado. O conteúdo do chamado (título e descrição) não pode ser editado.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        id       path      string                       true  "ID do Chamado (UUID)"
// @Param        payload  body      ticketApp.UpdateTicketInput  true  "Dados para atualização de status/responsável"
// @Success      200      {object}  ticketApp.TicketOutput
// @Failure      400      {object}  httpResponse.ErrorResponse
// @Failure      404      {object}  httpResponse.ErrorResponse
// @Failure      500      {object}  httpResponse.ErrorResponse
// @Router       /api/tickets/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var input ticketApp.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httpResponse.JSON(w, http.StatusBadRequest, httpResponse.ErrorResponse{Error: "invalid request payload"})
		return
	}

	output, err := h.updateUC.Execute(r.Context(), id, input)
	if err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusOK, output)
}

// GetByID godoc
// @Summary      Busca chamado por ID
// @Description  Retorna os detalhes de um chamado específico pelo seu UUID
// @Tags         Tickets
// @Produce      json
// @Param        id   path      string  true  "ID do Chamado (UUID)"
// @Success      200  {object}  ticketApp.TicketOutput
// @Failure      404  {object}  httpResponse.ErrorResponse
// @Failure      500  {object}  httpResponse.ErrorResponse
// @Router       /api/tickets/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	output, err := h.getByIDUC.Execute(r.Context(), id)
	if err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusOK, output)
}

// List godoc
// @Summary      Lista chamados com paginação, filtros e ordenação
// @Description  Retorna uma lista paginada de chamados de acordo com os filtros especificados
// @Tags         Tickets
// @Produce      json
// @Param        page       query     int     false  "Número da página (padrão: 1)"
// @Param        page_size  query     int     false  "Itens por página (máx: 100, padrão: 10)"
// @Param        sort_by    query     string  false  "Campo de ordenação (created_at, priority, status, title)"
// @Param        order      query     string  false  "Sentido da ordenação (asc, desc)"
// @Param        status     query     string  false  "Filtro por status (open, in_progress, resolved, closed)"
// @Param        priority   query     string  false  "Filtro por prioridade (low, medium, high)"
// @Param        assignee   query     string  false  "Filtro por nome do responsável"
// @Param        search     query     string  false  "Busca textual em título e descrição"
// @Success      200        {object}  ticketApp.PaginatedTicketsOutput
// @Failure      400        {object}  httpResponse.ErrorResponse
// @Failure      500        {object}  httpResponse.ErrorResponse
// @Router       /api/tickets [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	input := ticketApp.ListInput{
		Page:     page,
		PageSize: pageSize,
		SortBy:   q.Get("sort_by"),
		Order:    q.Get("order"),
	}

	if st := q.Get("status"); st != "" {
		input.Status = &st
	}
	if pr := q.Get("priority"); pr != "" {
		input.Priority = &pr
	}
	if as := q.Get("assignee"); as != "" {
		input.Assignee = &as
	}
	if search := q.Get("search"); search != "" {
		input.Search = &search
	}

	output, err := h.listUC.Execute(r.Context(), input)
	if err != nil {
		httpResponse.RespondError(w, err)
		return
	}

	httpResponse.JSON(w, http.StatusOK, output)
}
