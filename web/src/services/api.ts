import type { components, paths } from '../types/api';

export type Ticket = components['schemas']['TicketOutput'];
export type PaginatedTickets = components['schemas']['PaginatedTicketsOutput'];
export type CreateTicketInput = components['schemas']['CreateTicketInput'];
export type UpdateTicketInput = components['schemas']['UpdateTicketInput'];
export type Assignee = components['schemas']['AssigneeOutput'];
export type CreateAssigneeInput = components['schemas']['CreateAssigneeInput'];
export type GetTicketsQuery = NonNullable<paths['/api/tickets']['get']['parameters']['query']>;

const API_BASE = import.meta.env.VITE_API_URL.replace(/\/+$/, '');

export async function fetchTickets(params: GetTicketsQuery = {}): Promise<PaginatedTickets> {
  const query = new URLSearchParams();
  if (params.page) query.set('page', String(params.page));
  if (params.page_size) query.set('page_size', String(params.page_size));
  if (params.status) query.set('status', params.status);
  if (params.priority) query.set('priority', params.priority);
  if (params.search) query.set('search', params.search);
  if (params.sort_by) query.set('sort_by', params.sort_by);
  if (params.order) query.set('order', params.order);
  if (params.assignee) query.set('assignee', params.assignee);

  const qs = query.toString();
  const url = `${API_BASE}/tickets${qs ? `?${qs}` : ''}`;
  const response = await fetch(url);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao buscar chamados' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }

  return response.json();
}

export async function fetchTicketById(id: string): Promise<Ticket> {
  const response = await fetch(`${API_BASE}/tickets/${encodeURIComponent(id)}`);
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao buscar chamado' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }
  return response.json();
}

export async function fetchAssignees(): Promise<Assignee[]> {
  const response = await fetch(`${API_BASE}/assignees`);
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao buscar responsáveis' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }
  return response.json();
}

export async function createTicket(payload: CreateTicketInput): Promise<Ticket> {
  const response = await fetch(`${API_BASE}/tickets`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao criar chamado' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }

  return response.json();
}

export async function updateTicket(id: string, payload: UpdateTicketInput): Promise<Ticket> {
  const response = await fetch(`${API_BASE}/tickets/${encodeURIComponent(id)}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao atualizar chamado' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }

  return response.json();
}

export async function createAssignee(payload: CreateAssigneeInput): Promise<{ message: string }> {
  const response = await fetch(`${API_BASE}/assignees`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Erro ao cadastrar agente' }));
    throw new Error(errorData.error || `Erro HTTP ${response.status}`);
  }

  return response.json();
}
