// @title           Sistema de Chamados
// @version         1.0.0
// @description     API REST para gerenciamento de chamados de suporte
// @host            localhost:8080
// @BasePath        /
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	assigneeApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	ticketApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	infraHttp "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http"
	httpAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/assignee"
	httpTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/ticket"
	infraPg "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres"
)

func main() {
	port := os.Getenv("PORT")
	// NOTE: Talvez pode virar um arquivo apenas de configuração, mas por enquanto deixei assim.
	if port == "" {
		panic("[FATAL] A variável de ambiente PORT não está definida.")
	}

	db, err := infraPg.NewDatabase("")
	if err != nil {
		log.Fatalf("[FATAL] Falha ao conectar ao banco de dados PostgreSQL: %v\n", err)
	}
	defer db.Close()

	ticketRepo := infraPg.NewTicketRepository(db)
	assigneeRepo := infraPg.NewAssigneeRepository(db)

	assigneeCreateUC := assigneeApp.NewCreateUseCase(assigneeRepo)
	assigneeListUC := assigneeApp.NewListUseCase(assigneeRepo, ticketRepo)

	resolveAssigneeUC := ticketApp.NewResolveAssigneeUseCase(assigneeRepo)
	ticketCreateUC := ticketApp.NewCreateUseCase(ticketRepo, resolveAssigneeUC)
	ticketUpdateUC := ticketApp.NewUpdateUseCase(ticketRepo, resolveAssigneeUC)
	ticketGetByIDUC := ticketApp.NewGetByIDUseCase(ticketRepo)
	ticketListUC := ticketApp.NewListUseCase(ticketRepo)

	// Handlers HTTP
	ticketHandler := httpTicket.NewHandler(ticketCreateUC, ticketUpdateUC, ticketGetByIDUC, ticketListUC)
	assigneeHandler := httpAssignee.NewHandler(assigneeCreateUC, assigneeListUC)

	router := infraHttp.NewRouter(ticketHandler, assigneeHandler)

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	server := &http.Server{
		Addr:         host + ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] HTTP server listening on %s:%s\n", host, port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] HTTP server listen error: %v\n", err)
		}
	}()

	<-stopChan
	log.Println("[INFO] Termination signal received, initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[FATAL] HTTP server graceful shutdown failed: %v\n", err)
	}

	log.Println("[INFO] HTTP server stopped.")
}
