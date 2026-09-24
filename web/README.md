# Front-end

Interface para o gerenciamento de chamados de suporte

## Tecnologias Utilizadas

- **Framework:** [React 19](https://react.dev/) + [TypeScript](https://www.typescriptlang.org/)
- **Build Tool:** [Vite](https://vite.dev/)
- **Roteamento:** [React Router v7](https://reactrouter.com/)
- **Estilizacao:** [Tailwind CSS v4](https://tailwindcss.com/)
- **Contratos & Tipagem:** [openapi-typescript](https://github.com/openapi-ts/openapi-typescript) gerando tipos a partir do `api/openapi.yaml`
- **Qualidade de Codigo:** [ESLint](https://eslint.org/) (Flat Config) + [Prettier](https://prettier.io/)

## Como Rodar

### 1. Pre-requisitos

- **Node.js:** Versao 18+ (recomendado 20 ou superior)
- **npm:** Versao 9+
- Backend rodando em `http://localhost:8080` (veja instrucoes no README da raiz)

### 2. Configurar Variaveis de Ambiente

| Variavel       | Obrigatoria | Descricao                     | Exemplo Padrao              |
| :------------- | :---------- | :---------------------------- | :-------------------------- |
| `VITE_API_URL` | **Sim**     | URL base dos endpoints da API | `http://localhost:8080/api` |
| `VITE_PORT`    | Não         | Porta do servidor Vite        | `5173`                      |

### 3. Instalar Dependencias

```bash
npm install
```

### 4. Executar em Modo de Desenvolvimento

```bash
npm run dev
```

A aplicacao estara disponivel em `http://localhost:5173`.

## Scripts Disponiveis

| Comando                | Descricao                                                                                 |
| :--------------------- | :---------------------------------------------------------------------------------------- |
| `npm run dev`          | Inicia o servidor de desenvolvimento com Hot Module Replacement (HMR)                     |
| `npm run build`        | Valida tipagens com `tsc` e gera os arquivos otimizados para producao em `dist/`          |
| `npm run preview`      | Executa localmente a build de producao para testes pre-deploy (porta 4173)                |
| `npm run lint`         | Executa o linter ESLint em todo o codigo TypeScript/React                                 |
| `npm run lint:fix`     | Corrige problemas automaticos apontados pelo ESLint                                       |
| `npm run format`       | Formata todo o codigo utilizando o Prettier                                               |
| `npm run format:check` | Verifica se os arquivos estao formatados de acordo com o padrao do Prettier               |
| `npm run generate:api` | Regenera as definicoes de tipo TypeScript a partir da especificacao `../api/openapi.yaml` |

## Type-Safety com OpenAPI

Todas as requisicoes e respostas da API utilizam contratos gerados automaticamente:

1. A especificacao em `../api/openapi.yaml` e convertida e tipada em `src/types/api.ts`.
2. O cliente HTTP em `src/services/api.ts` exporta as entidades e funcoes tipadas:
   - `Ticket` (`TicketOutput`)
   - `PaginatedTickets` (`PaginatedTicketsOutput`)
   - `CreateTicketInput`
   - `Assignee` (`AssigneeOutput`)
   - `fetchTickets()`, `createTicket()`, `fetchAssignees()`

Sempre que a API for alterada no backend, execute:

```bash
npm run generate:api
```
