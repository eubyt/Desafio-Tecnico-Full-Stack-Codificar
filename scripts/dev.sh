#!/usr/bin/env bash

set -eo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

COLOR_RESET="\033[0m"
COLOR_GREEN="\033[32m"
COLOR_BLUE="\033[34m"
COLOR_YELLOW="\033[33m"
COLOR_RED="\033[31m"
COLOR_CYAN="\033[36m"

# Carrega .env da raiz caso o arquivo exista
load_env() {
  if [ -f "$ROOT_DIR/.env" ]; then
    echo -e "${COLOR_CYAN}[env] Carregando variáveis de ambiente (.env)...${COLOR_RESET}"
    set -a
    # shellcheck disable=SC1090
    source "$ROOT_DIR/.env"
    set +a
  fi
}

# Valida se as variáveis de ambiente obrigatórias estão configuradas
validate_env() {
  local required_vars=(
    "PORT"
    "DB_HOST"
    "DB_PORT"
    "DB_USER"
    "DB_PASSWORD"
    "DB_NAME"
    "VITE_API_URL"
  )

  local missing_vars=()
  for var_name in "${required_vars[@]}"; do
    if [ -z "${!var_name:-}" ]; then
      missing_vars+=("$var_name")
    fi
  done

  if [ ${#missing_vars[@]} -gt 0 ]; then
    echo -e "${COLOR_RED}[env] Erro: As seguintes variáveis de ambiente obrigatórias não estão configuradas:${COLOR_RESET}"
    for var_name in "${missing_vars[@]}"; do
      echo -e "${COLOR_RED}  - ${var_name}${COLOR_RESET}"
    done
    echo -e "${COLOR_YELLOW}[env] Configure-as no ambiente ou através de um arquivo .env (consulte .env.example).${COLOR_RESET}"
    exit 1
  fi
}


check_docker_daemon() {
  if ! docker info >/dev/null 2>&1; then
    echo -e "${COLOR_RED}[docker] Erro: Docker daemon inacessível ou inativo. Verifique se o daemon do Docker está em execução.${COLOR_RESET}"
    exit 1
  fi
}

ensure_docker_services() {
  check_docker_daemon

  local running_containers
  running_containers=$(docker compose -f "$ROOT_DIR/docker-compose.yml" ps --services --filter "status=running" 2>/dev/null || true)

  local needs_up=false
  if ! echo "$running_containers" | grep -q "^postgres$"; then
    needs_up=true
  fi
  if ! echo "$running_containers" | grep -q "^app$"; then
    needs_up=true
  fi

  if [ "$needs_up" = true ]; then
    echo -e "${COLOR_YELLOW}[docker] Provisionando containers de infraestrutura (postgres, app)...${COLOR_RESET}"
    docker compose -f "$ROOT_DIR/docker-compose.yml" up -d
  else
    echo -e "${COLOR_GREEN}[docker] Containers de infraestrutura ativos e em execução (status: running).${COLOR_RESET}"
  fi

  wait_for_postgres
}

wait_for_postgres() {
  local db_user="${DB_USER:-ticket_user}"
  local db_name="${DB_NAME:-ticket_system}"
  local max_retries=30
  local count=0

  while ! docker compose -f "$ROOT_DIR/docker-compose.yml" exec -T postgres pg_isready -U "$db_user" -d "$db_name" >/dev/null 2>&1; do
    count=$((count + 1))
    if [ $count -ge $max_retries ]; then
      echo -e "${COLOR_RED}[docker] Erro: Timeout excedido ao aguardar readiness do PostgreSQL (database: '${db_name}', user: '${db_user}').${COLOR_RESET}"
      exit 1
    fi
    sleep 1
  done
  echo -e "${COLOR_GREEN}[docker] PostgreSQL operacional e pronto para conexões (healthcheck pg_isready: OK).${COLOR_RESET}"
}

find_swag() {
  if command -v swag >/dev/null 2>&1; then
    command -v swag
  elif [ -x "$(go env GOPATH 2>/dev/null)/bin/swag" ]; then
    echo "$(go env GOPATH)/bin/swag"
  else
    echo ""
  fi
}

cmd_build() {
  echo -e "${COLOR_BLUE}[build] Executando análise estática do backend (go vet ./...)...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && go vet ./...)

  echo -e "${COLOR_BLUE}[build] Executando linting do frontend (eslint)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm run lint)

  echo -e "${COLOR_BLUE}[build] Executando suíte de testes do backend (go test -p 1 -v ./...)...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && go test -p 1 -v ./...)

  echo -e "${COLOR_BLUE}[build] Executando suíte de testes do frontend (npm test)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm test)

  echo -e "${COLOR_BLUE}[build] Compilando binário do backend (cmd/api -> bin/api)...${COLOR_RESET}"
  mkdir -p "$ROOT_DIR/backend/bin"
  (cd "$ROOT_DIR/backend" && go build -o bin/api ./cmd/api)

  echo -e "${COLOR_BLUE}[build] Gerando bundle de produção do frontend (vite build)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm run build)

  echo -e "${COLOR_GREEN}[build] Processo de compilação concluído com êxito para todos os alvos.${COLOR_RESET}"
}

cmd_start() {
  ensure_docker_services

  if [ ! -f "$ROOT_DIR/backend/bin/api" ] || [ ! -d "$ROOT_DIR/web/dist" ]; then
    echo -e "${COLOR_RED}[start] Erro: Artefatos de compilação não encontrados (backend/bin/api ou web/dist).${COLOR_RESET}"
    echo -e "${COLOR_YELLOW}[start] Execute './scripts/dev.sh build' (ou 'npm run build') antes de executar o start.${COLOR_RESET}"
    exit 1
  fi

  local port="${PORT:-8080}"
  echo -e "${COLOR_GREEN}[start] Inicializando processo do backend (porta :${port})...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && ./bin/api) &
  local backend_pid=$!

  local retries=30
  while ! curl -s "http://localhost:${port}/health" >/dev/null 2>&1; do
    retries=$((retries - 1))
    if [ $retries -eq 0 ]; then
      echo -e "${COLOR_RED}[start] Erro: Timeout ao aguardar healthcheck do backend (http://localhost:${port}/health).${COLOR_RESET}"
      kill "$backend_pid" 2>/dev/null || true
      exit 1
    fi
    sleep 0.2
  done

  echo -e "${COLOR_GREEN}[start] Backend saudável (HTTP 200). Inicializando servidor web preview (porta :4173)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm run preview -- --host) &
  local web_pid=$!

  cleanup_start() {
    kill "$backend_pid" "$web_pid" 2>/dev/null || true
    wait "$backend_pid" 2>/dev/null || true
    wait "$web_pid" 2>/dev/null || true
  }
  trap cleanup_start SIGINT SIGTERM EXIT
  wait
}

cmd_openapi_gen() {
  local swag_bin
  swag_bin=$(find_swag)

  if [ -z "$swag_bin" ]; then
    echo -e "${COLOR_RED}[openapi] Erro: Binário 'swag' não localizado no PATH nem no GOPATH ($GOPATH/bin/swag). Instale com: go install github.com/swaggo/swag/cmd/swag@latest${COLOR_RESET}"
    exit 1
  fi

  echo -e "${COLOR_BLUE}[openapi] Gerando especificação OpenAPI a partir de anotações Go...${COLOR_RESET}"
  mkdir -p "$ROOT_DIR/api"
  (cd "$ROOT_DIR/backend" && "$swag_bin" init -g cmd/api/main.go -o "$ROOT_DIR/api" -ot yaml --parseDependency --parseInternal --useStructName)
  mv "$ROOT_DIR/api/swagger.yaml" "$ROOT_DIR/api/openapi.yaml"
  echo -e "${COLOR_GREEN}[openapi] Especificação OpenAPI gerada em: api/openapi.yaml.${COLOR_RESET}"

  if [ -d "$ROOT_DIR/web/node_modules" ]; then
    echo -e "${COLOR_BLUE}[openapi] Sincronizando definições de tipos TypeScript no frontend...${COLOR_RESET}"
    (cd "$ROOT_DIR/web" && npm run generate:api)
    echo -e "${COLOR_GREEN}[openapi] Definições de tipos TypeScript sincronizadas com sucesso.${COLOR_RESET}"
  fi
}

cmd_test() {
  echo -e "${COLOR_BLUE}[test] Executando suíte de testes do backend (go test -p 1 -v ./...)...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && go test -p 1 -v ./...)

  echo -e "${COLOR_BLUE}[test] Executando suíte de testes do frontend (npm test)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm test)

  echo -e "${COLOR_GREEN}[test] Baterias de testes finalizadas com êxito.${COLOR_RESET}"
}

cmd_dev() {
  ensure_docker_services

  local port="${PORT:-8080}"
  echo -e "${COLOR_GREEN}[dev] Inicializando processo backend em modo de desenvolvimento (porta :${port})...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && go run ./cmd/api) &
  local backend_pid=$!

  local retries=30
  while ! curl -s "http://localhost:${port}/health" >/dev/null 2>&1; do
    retries=$((retries - 1))
    if [ $retries -eq 0 ]; then
      echo -e "${COLOR_RED}[dev] Erro: Timeout ao aguardar inicialização do backend (healthcheck http://localhost:${port}/health falhou após 30 tentativas).${COLOR_RESET}"
      kill "$backend_pid" 2>/dev/null || true
      exit 1
    fi
    sleep 0.3
  done

  echo -e "${COLOR_GREEN}[dev] Backend operacional em :${port} (healthcheck OK). Inicializando servidor de desenvolvimento frontend (Vite :5173)...${COLOR_RESET}"
  (cd "$ROOT_DIR/web" && npm run dev -- --host) &
  local web_pid=$!

  cleanup_dev() {
    kill "$backend_pid" "$web_pid" 2>/dev/null || true
    wait "$backend_pid" 2>/dev/null || true
    wait "$web_pid" 2>/dev/null || true
  }
  trap cleanup_dev SIGINT SIGTERM EXIT
  wait
}

cmd_db_up() {
  check_docker_daemon
  echo -e "${COLOR_BLUE}[docker] Provisionando containers em modo detached via Docker Compose...${COLOR_RESET}"
  docker compose -f "$ROOT_DIR/docker-compose.yml" up -d
  wait_for_postgres
}

cmd_db_down() {
  check_docker_daemon
  echo -e "${COLOR_YELLOW}[docker] Encerrando serviços gerenciados pelo Docker Compose...${COLOR_RESET}"
  docker compose -f "$ROOT_DIR/docker-compose.yml" down
  echo -e "${COLOR_GREEN}[docker] Containers de infraestrutura encerrados com sucesso.${COLOR_RESET}"
}

cmd_db_status() {
  check_docker_daemon
  docker compose -f "$ROOT_DIR/docker-compose.yml" ps
}

cmd_seed() {
  ensure_docker_services
  echo -e "${COLOR_BLUE}[seed] Executando rotina de seeding do banco de dados (cmd/seed)...${COLOR_RESET}"
  (cd "$ROOT_DIR/backend" && go run ./cmd/seed)
  echo -e "${COLOR_GREEN}[seed] Povoamento da base de dados finalizado com êxito.${COLOR_RESET}"
}

show_help() {
  echo "Uso: ./scripts/dev.sh [comando]"
  echo ""
  echo "Comandos disponíveis:"
  echo "  dev          Provisiona dependências (Docker Compose) e inicializa backend e frontend em modo dev"
  echo "  build        Executa análise estática (lint), suíte de testes e compilação de binários e bundles"
  echo "  start        Provisiona dependências e inicializa aplicação em ambiente de preview/produção"
  echo "  test         Executa suítes de testes unitários e de integração (backend e frontend)"
  echo "  openapi-gen  Gera especificação OpenAPI 3.0 via Swagger e compila definições de tipos TypeScript"
  echo "  seed         Executa rotinas de seed para população de dados no banco relacional"
  echo "  db:up        Provisiona e inicia os containers de infraestrutura em segundo plano"
  echo "  db:down      Interrompe e encerra a execução dos containers de infraestrutura"
  echo "  db:status    Exibe o estado operacional e métricas de status dos containers"
  echo "  help         Exibe informações de uso e lista de comandos disponíveis"
}

main() {
  local cmd="${1:-}"

  case "$cmd" in
    help|--help|-h|"")
      show_help
      ;;
    dev|build|start|openapi-gen|openapi|test|seed|db:up|db-up|db:down|db-down|db:status|db-status)
      load_env
      validate_env
      case "$cmd" in
        dev)
          cmd_dev
          ;;
        build)
          cmd_build
          ;;
        start)
          cmd_start
          ;;
        openapi-gen|openapi)
          cmd_openapi_gen
          ;;
        test)
          cmd_test
          ;;
        seed)
          cmd_seed
          ;;
        db:up|db-up)
          cmd_db_up
          ;;
        db:down|db-down)
          cmd_db_down
          ;;
        db:status|db-status)
          cmd_db_status
          ;;
      esac
      ;;
    *)
      echo -e "${COLOR_RED}[erro] Comando não reconhecido: '$cmd'${COLOR_RESET}"
      show_help
      exit 1
      ;;
  esac
}

main "$@"
