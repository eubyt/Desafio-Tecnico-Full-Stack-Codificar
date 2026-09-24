# Documentação Técnica e Decisões de Arquitetura

Este documento consolida as decisões técnicas, regras de negócio, justificativas arquiteturais e trade-offs adotados no **Sistema de Chamados** (back-end e front-end).

## 1. Decisão Técnica: O que significa um Chamado "Em Aberto"?

Conforme solicitado no desafio técnico:

> _"Considere como 'em aberto' os chamados que ainda não foram concluídos. Cabe a você definir o que isso significa dentro do seu modelo de status e justificar a decisão."_

### Modelo de Status Adotado

O agregado `Ticket` adota quatro status canônicos:

- `open`
- `in_progress`
- `resolved`
- `closed`

### Definição

Um chamado é considerado **"em aberto"** se, e somente se, o seu status for **`open`** ou **`in_progress`**.

```go
func (s Status) IsOpen() bool {
    return s == StatusOpen || s == StatusInProgress
}
```

### Justificativas:

1. **Significado dos status:**
    - **`open`**: Representa chamados aguardando triagem ou início de atendimento pelo atendente atribuído.
    - **`in_progress`**: Representa chamados ativamente sob investigação e atuação do técnico, consumindo seu esforço e capacidade de trabalho.
    - **`resolved`**: Significa que a equipe técnica entregou a solução e o chamado aguarda apenas homologação do usuário ou encerramento automático. O atendente **não** está mais gastando esforço ativo de trabalho nesse chamado.
    - **`closed`**: Encerra formalmente o ciclo do chamado (arquivado).

## 2. Decisão Técnica: Escolha do `sqlc` sobre ORMs Tradicionais

### Motivação:

- ORMs convencionais (como GORM) introduzem overhead em runtime, geram queries genéricas muitas vezes imprevisíveis e dificultam o controle em queries complexas.
- O **`sqlc`** foi adotado por ser **SQL-First**: escrevemos SQL puro em arquivos `.sql` e o compilador do `sqlc` gera código Go com tipagem estática (type-safe).

## 3. Decisão Técnica: Algoritmo de Distribuição Automática

Para distribuir chamados de maneira determinística quando (`auto_assign: true`):

1. **Agregação e Seleção Diretamente no Banco de Dados (PostgreSQL):**
    - Em vez de buscar todos os atendentes em memória e contar os chamados em laços Go, a seleção é resolvida em uma única consulta SQL no PostgreSQL:
    ```sql
    SELECT a.id, a.name, a.created_at, a.updated_at
    FROM assignees a
    LEFT JOIN tickets t ON t.assignee_id = a.id AND t.status IN ('open', 'in_progress')
    GROUP BY a.id, a.name, a.created_at, a.updated_at
    ORDER BY COUNT(t.id) ASC, a.name ASC
    LIMIT 1;
    ```
2. **Vantagens Técnicas:**
    - **Roundtrip Único:** Uma única query indexada resolve agregação e seleção.
    - **Desempate Determinístico:** Caso dois ou mais atendentes tenham a mesma quantidade de chamados em aberto, o critério de desempate alfabético pelo nome (`a.name ASC`) garante comportamento previsível e consistente.
    - **Atendentes Recém-adicionados:** Atendentes sem chamados recebem contagem `0` via `LEFT JOIN` e são priorizados automaticamente.

3. **Cenário de Empate com Nomes Idênticos e Unicidade:**
    - **E se dois atendentes tiverem o mesmo nome e a mesma quantidade de chamados?**
        - No modelo atual, este cenário é **impossível de ocorrer**, pois a coluna `name` possui restrição estrita de unicidade no banco de dados (`name VARCHAR(100) NOT NULL UNIQUE` na tabela `assignees`), reforçada pela validação de domínio que rejeita duplicidades (`ErrAssigneeAlreadyExists`).

## 4. Decisão Técnica: Paginação e Ordenação Genéricas no `util/pagination`

- Paginação e ordenação foram abstraídas como structs genéricas Go (`PaginatedResult[T]`, `PaginationRequest`, `SortOption`, `SortOrder`) no pacote `internal/util/pagination`.
- Evita duplicação de lógica (DRY) de cálculo de `offset`, `limit`, páginas totais e limites máximos (`page_size` limitado entre 1 e 100), pronta para ser reaproveitada em qualquer nova entidade do sistema.

## 5. Trade-offs de Escopo e Segurança da API (Pontos Não Elaborados)

Para priorizar os objetivos essenciais do desafio técnico certos pontos foram deixados de fora.

### 5.1. Segurança nos Endpoints

- **Estado Atual:** Os endpoints REST (`/api/tickets`, `/api/assignees`) são públicos e acessíveis sem credenciais.

### 5.2. Chaves de Acesso para Serviços (API Keys)

- **Estado Atual:** Não há suporte a autenticação via cabeçalhos como `X-API-Key`.

### 5.3. Rate Limiting e Throttling (Proteção contra Abuso e DoS)

- **Estado Atual:** O servidor atende requisições sem limitação volumétrica por IP ou cliente.

### 5.4. Política de CORS (Cross-Origin Resource Sharing) Permissiva

- **Estado Atual:** O roteador utiliza `AllowedOrigins: []string{"*"}` (`backend/internal/infra/http/router.go`).

## 6. Decisão Técnica de Front-end: SPA (Client-Side Rendering) vs. Server-Side Rendering (SSR)

O front-end do projeto foi desenvolvido como uma **SPA pura (Vite + React)**, sem uso de frameworks SSR (como Next.js ou Remix).

### Por que NÃO estamos utilizando SSR?

O motivo principal é **evitar complexidade desnecessária (over-engineering)** para o propósito do projeto:

1. **Escopo de Desafio Técnico Demonstrativo:**
    - Adicionar SSR aumentaria drasticamente a complexidade do projeto (runtime Node.js em servidor, hidratação, deploys híbridos, sincronização de estado servidor/cliente) sem agregar valor real.
