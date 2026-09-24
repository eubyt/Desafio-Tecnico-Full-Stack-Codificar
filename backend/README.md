# Sistema de Chamados - Back-end

API REST para gerenciamento de chamados de suporte e distribuição de atendimentos desenvolvida em Go.

## Como Rodar

### 1. Pré-requisitos

- Go: 1.24 ou superior
- Docker e Docker Compose (para subir o PostgreSQL)
- sqlc (opcional, necessário apenas se alterar arquivos .sql): `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

### 2. Configurar Variáveis de Ambiente

Copie o arquivo de exemplo `.env.example` para criar seu .env.

Variáveis necessárias:
| Variável | Obrigatória | Descrição | Exemplo Padrão |
| :--- | :--- | :--- | :--- |
| `DB_HOST` | **Sim** | Host do PostgreSQL | `localhost` |
| `DB_PORT` | **Sim** | Porta do PostgreSQL | `5432` |
| `DB_USER` | **Sim** | Usuário do banco de dados | `ticket_user` |
| `DB_PASSWORD` | **Sim** | Senha do usuário | `ticket_password` |
| `DB_NAME` | **Sim** | Nome do banco de dados | `ticket_system` |
| `DB_SSLMODE` | Não | Modo de conexão SSL | `disable` |
| `PORT` | Não | Porta do servidor HTTP | `8080` |

### 3. Subir o Banco de Dados

A partir da raiz do projeto:

```bash
make db-up
```

### 4. Executar a Aplicação

No diretório `backend`:

```bash
# Exportar variáveis (ou carregar seu .env)
export DB_HOST=localhost DB_PORT=5432 DB_USER=ticket_user DB_PASSWORD=ticket_password DB_NAME=ticket_system

# Iniciar servidor
make run
# ou: go run ./cmd/api
```

O servidor iniciará em `http://localhost:8080`.

## Comandos e Testes

```bash
# Popular o banco de dados com atendentes e chamados de demonstração
make seed

# Executar todos os testes da aplicação sequencialmente
make test

# Executar testes apenas da camada de domínio
make test-domain

# Executar testes com o detector de race condition do Go
make test-race

# Executar testes de integração ponta a ponta com Testcontainers (PostgreSQL real)
make test-e2e

# Compilar o binário na pasta bin/
make build

# Regenerar código Go a partir dos arquivos SQL (sqlc)
make sqlc

# Gerar especificação OpenAPI automaticamente na pasta api/
make openapi

# Limpar artefatos gerados
make clean
```

## Endpoints da API

A especificação completa OpenAPI está disponível na pasta [`api/openapi.yaml`](file:///Volumes/Arquivos/Projetos/demo-sistema-chamado/api/openapi.yaml).

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| `GET` | `/health` | Checagem de integridade da API |
| `POST` | `/api/tickets` | Cria um novo chamado (com atribuição manual ou automática) |
| `GET` | `/api/tickets` | Lista chamados com paginação, filtros e ordenação |
| `GET` | `/api/tickets/{id}` | Busca os detalhes de um chamado específico por UUID |
| `PUT` | `/api/tickets/{id}` | Atualiza título, descrição, prioridade, status ou responsável |
| `GET` | `/api/assignees` | Lista todos os responsáveis com o total de chamados em aberto |
| `POST` | `/api/assignees` | Cadastra um novo membro responsável |

## Estrutura de Pastas

```
backend/
├── cmd/               # Ponto de entrada da aplicação (main.go, bootstrap e shutdown)
├── internal/
│   ├── domain/        # Regras de negócio puras, entidades, value objects e interfaces
│   ├── application/   # Casos de uso da aplicação (Use Cases) e DTOs
│   ├── infra/         # Adaptadores externos: HTTP (rotas e handlers) e PostgreSQL (queries e pool)
│   └── util/          # Utilitários compartilhados (paginação genérica, helpers de teste e mocks)
└── bin/               # Binários compilados da aplicação
```

## Tecnologias e Versões

- **Linguagem**: Go (`1.24+` / `1.25.0`)
- **Banco de Dados**: PostgreSQL `16`
- **Driver de Banco / Connection Pool**: `pgx/v5` (`v5.11.0`) com `pgxpool`
- **Compilador SQL (Type-Safe)**: `sqlc` (`v1.31+`)
- **Roteador HTTP / Middlewares**: `go-chi/chi/v5` (`v5.3.2`) e `go-chi/cors` (`v1.2.2`)
- **Identificadores Únicos**: `google/uuid` (`v1.6.0`) gerando UUIDv7 ordenável
- **Testes e Asserções**: `stretchr/testify` (`v1.12.1`) com mocks em memória
- **Testes E2E / Integração**: `testcontainers/testcontainers-go` (`v0.44.0`) com módulo PostgreSQL
