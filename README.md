# Sistema de Chamados

Solução fullstack desenvolvida com back-end em Go e front-end em React.

Para detalhes específicos de arquitetura, decisões de negócio e regras de cada camada:

- [Documentação do Back-end (`backend/README.md`)](backend/README.md)
- [Documentação do Front-end (`web/README.md`)](web/README.md)
- [Decisões Técnicas e Arquitetura (`doc/decisoes_tecnicas.md`)](doc/decisoes_tecnicas.md)

## Pré-requisitos

Certifique-se de ter instalado em seu ambiente:

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/)
- [Node.js](https://nodejs.org/) (v18 ou superior) e `npm`
- [Go](https://go.dev/) (v1.24 ou superior)

## Como Rodar

### 1. Configurar variáveis de ambiente

Defina as variáveis de ambiente necessárias para a execução da aplicação:

**Obrigatórias:**

- `PORT`: Porta da API HTTP do back-end (ex: `8080`)
- `DB_HOST`: Host de conexão com o banco PostgreSQL (ex: `localhost`)
- `DB_PORT`: Porta do banco de dados PostgreSQL (ex: `5432`)
- `DB_USER`: Usuário do banco de dados PostgreSQL (ex: `ticket_user`)
- `DB_PASSWORD`: Senha do banco de dados PostgreSQL (ex: `ticket_password`)
- `DB_NAME`: Nome do banco de dados (ex: `ticket_system`)
- `VITE_API_URL`: URL base da API consumida pelo front-end (ex: `http://localhost:8080/api`)

**Opcionais:**

- `HOST`: Host de escuta da API (padrão: `0.0.0.0`)
- `DB_SSLMODE`: Modo de conexão SSL com o PostgreSQL (padrão: `disable`)
- `VITE_PORT`: Porta do servidor de desenvolvimento do front-end (padrão: `5173`)

### 2. Instalar as dependências

Instale os pacotes da raiz e do front-end:

```bash
npm install && npm --prefix web install
```

### 3. Iniciar a aplicação

```bash
./scripts/dev.sh dev
# ou: npm run dev
```

Após a inicialização, os serviços estarão acessíveis em:

- **Front-end (Web):** [http://localhost:5173](http://localhost:5173)
- **API Back-end:** [http://localhost:8080](http://localhost:8080/api/tickets)

## Comandos Disponíveis

| Comando       | Descrição                                                                                        | Atalho npm            |
| :------------ | :----------------------------------------------------------------------------------------------- | :-------------------- |
| `dev`         | Inicia os containers Docker, back-end e front-end                                                | `npm run dev`         |
| `build`       | Executa linter, suítes de testes e gera a build de produção                                      | `npm run build`       |
| `start`       | Garante containers ativos e executa a aplicação previamente compilada em modo preview | `npm start`           |
| `test`        | Executa as suítes de testes unitários do back-end (`go test`) e validação de tipos do front-end  | `npm test`            |
| `openapi-gen` | Regenera o contrato OpenAPI e sincroniza as tipagens TypeScript do front-end                     | `npm run openapi-gen` |
| `seed`        | Popula o banco de dados com dados iniciais de demonstração                                       | —                     |
| `db:up`       | Provisiona e inicia os containers de banco de dados em segundo plano via Docker Compose          | —                     |
| `db:down`     | Interrompe e encerra os containers da aplicação                                                  | —                     |
| `db:status`   | Exibe o status e métricas operacionais dos containers Docker                                     | —                     |
| `help`        | Exibe a mensagem de ajuda com as instruções de uso dos comandos                                  | —                     |
