import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  fetchTickets,
  fetchAssignees,
  createAssignee,
  type Ticket,
  type Assignee,
  type PaginatedTickets,
} from '../services/api';

export default function TicketListPage() {
  const navigate = useNavigate();

  const [ticketsData, setTicketsData] = useState<PaginatedTickets | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [priorityFilter, setPriorityFilter] = useState<string>('');
  const [search, setSearch] = useState<string>('');
  const [debouncedSearch, setDebouncedSearch] = useState<string>('');

  // Agentes
  const [, setAssigneesList] = useState<Assignee[]>([]);

  // Modal de criação de agente
  const [showAgentModal, setShowAgentModal] = useState(false);
  const [newAgentName, setNewAgentName] = useState('');
  const [savingAgent, setSavingAgent] = useState(false);
  const [agentModalError, setAgentModalError] = useState<string | null>(null);
  const [agentModalSuccess, setAgentModalSuccess] = useState<string | null>(null);

  // Debounce busca
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  // Carrega agentes na inicialização
  useEffect(() => {
    let isMounted = true;
    fetchAssignees()
      .then((data) => {
        if (isMounted) setAssigneesList(data);
      })
      .catch(() => {});

    return () => {
      isMounted = false;
    };
  }, []);

  async function reloadAssignees() {
    try {
      const data = await fetchAssignees();
      setAssigneesList(data);
    } catch {
      // Ignora erro silenciosamente
    }
  }

  // Carrega chamados
  useEffect(() => {
    let isMounted = true;

    async function loadTicketsData() {
      setLoading(true);
      setError(null);
      try {
        const data = await fetchTickets({
          page,
          page_size: pageSize,
          status: statusFilter || undefined,
          priority: priorityFilter || undefined,
          search: debouncedSearch || undefined,
          sort_by: 'created_at',
          order: 'desc',
        });
        if (isMounted) {
          setTicketsData(data);
        }
      } catch (err: unknown) {
        if (isMounted) {
          setError(err instanceof Error ? err.message : 'Erro ao carregar chamados');
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    }

    void loadTicketsData();

    return () => {
      isMounted = false;
    };
  }, [page, pageSize, statusFilter, priorityFilter, debouncedSearch]);

  async function reloadTickets() {
    setLoading(true);
    try {
      const data = await fetchTickets({
        page,
        page_size: pageSize,
        status: statusFilter || undefined,
        priority: priorityFilter || undefined,
        search: debouncedSearch || undefined,
        sort_by: 'created_at',
        order: 'desc',
      });
      setTicketsData(data);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Erro ao carregar chamados');
    } finally {
      setLoading(false);
    }
  }

  const totalPages = ticketsData?.total_pages || 1;
  const totalItems = ticketsData?.total_items || 0;
  const items: Ticket[] = ticketsData?.items || [];

  // Criar novo agente
  async function handleCreateAgent(e: FormEvent) {
    e.preventDefault();
    if (!newAgentName.trim()) {
      setAgentModalError('Informe o nome do agente.');
      return;
    }

    setSavingAgent(true);
    setAgentModalError(null);
    setAgentModalSuccess(null);

    try {
      await createAssignee({ name: newAgentName.trim() });
      setAgentModalSuccess('Agente cadastrado com sucesso!');
      setNewAgentName('');
      await reloadAssignees();
      setTimeout(() => {
        setShowAgentModal(false);
        setAgentModalSuccess(null);
      }, 1000);
    } catch (err: unknown) {
      setAgentModalError(err instanceof Error ? err.message : 'Erro ao cadastrar agente');
    } finally {
      setSavingAgent(false);
    }
  }

  return (
    <div className="min-h-screen bg-white text-black">
      <header className="border-b border-black">
        <div className="max-w-6xl mx-auto px-4 py-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-xl font-bold tracking-tight uppercase">Sistema de Chamados</h1>
            <p className="text-xs text-neutral-500 font-mono mt-0.5">
              Gestão de chamados e atendimento técnico
            </p>
          </div>
          <div className="flex items-center gap-3">
            <button
              type="button"
              onClick={() => {
                setShowAgentModal(true);
                setAgentModalError(null);
                setAgentModalSuccess(null);
              }}
              className="inline-flex items-center justify-center px-4 py-2 border border-black bg-white text-black text-sm font-medium hover:bg-neutral-100 transition-colors cursor-pointer"
            >
              + Novo Agente
            </button>
            <Link
              to="/novo"
              className="inline-flex items-center justify-center px-4 py-2 border border-black bg-white text-black text-sm font-medium hover:bg-neutral-100 transition-colors cursor-pointer"
            >
              + Abrir novo chamado
            </Link>
          </div>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 py-8">
        <div className="border border-black p-4 mb-6 bg-neutral-50 flex flex-col md:flex-row gap-3">
          <div className="flex-1">
            <label htmlFor="search" className="block text-xs font-mono uppercase mb-1">
              Buscar
            </label>
            <input
              id="search"
              type="text"
              placeholder="Buscar por título ou descrição..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full px-3 py-1.5 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
            />
          </div>

          <div className="w-full md:w-44">
            <label htmlFor="status" className="block text-xs font-mono uppercase mb-1">
              Status
            </label>
            <select
              id="status"
              value={statusFilter}
              onChange={(e) => {
                setStatusFilter(e.target.value);
                setPage(1);
              }}
              className="w-full px-3 py-1.5 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
            >
              <option value="">Todos os status</option>
              <option value="open">Aberto</option>
              <option value="in_progress">Em andamento</option>
              <option value="resolved">Resolvido</option>
              <option value="closed">Fechado</option>
            </select>
          </div>

          <div className="w-full md:w-44">
            <label htmlFor="priority" className="block text-xs font-mono uppercase mb-1">
              Prioridade
            </label>
            <select
              id="priority"
              value={priorityFilter}
              onChange={(e) => {
                setPriorityFilter(e.target.value);
                setPage(1);
              }}
              className="w-full px-3 py-1.5 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
            >
              <option value="">Todas as prioridades</option>
              <option value="low">Baixa</option>
              <option value="medium">Média</option>
              <option value="high">Alta</option>
            </select>
          </div>
        </div>

        {error && (
          <div className="border border-black bg-neutral-100 p-4 mb-6">
            <p className="text-sm font-medium">Erro ao carregar chamados: {error}</p>
            <button
              type="button"
              onClick={() => void reloadTickets()}
              className="mt-2 px-3 py-1 border border-black bg-white text-black text-xs uppercase font-mono hover:bg-neutral-100 transition-colors cursor-pointer"
            >
              Tentar novamente
            </button>
          </div>
        )}

        <div className="border border-black overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse text-sm">
              <thead>
                <tr className="border-b border-black bg-neutral-100 font-mono text-xs uppercase">
                  <th className="py-3 px-4">Título do Chamado</th>
                  <th className="py-3 px-4 w-36 text-right">Ação</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-neutral-200">
                {loading ? (
                  <tr>
                    <td
                      colSpan={2}
                      className="py-12 text-center text-neutral-500 font-mono text-sm"
                    >
                      Carregando chamados...
                    </td>
                  </tr>
                ) : items.length === 0 ? (
                  <tr>
                    <td
                      colSpan={2}
                      className="py-12 text-center text-neutral-500 font-mono text-sm"
                    >
                      Nenhum chamado encontrado.
                    </td>
                  </tr>
                ) : (
                  items.map((ticket) => (
                    <tr
                      key={ticket.id}
                      onClick={() => navigate(`/chamados/${ticket.id}`)}
                      className="hover:bg-neutral-50 transition-colors cursor-pointer group"
                    >
                      <td className="py-3.5 px-4 text-xs sm:text-sm font-medium">
                        <span className="group-hover:underline">{ticket.title}</span>
                      </td>
                      <td className="py-3.5 px-4 text-right">
                        <button
                          type="button"
                          onClick={(e) => {
                            e.stopPropagation();
                            navigate(`/chamados/${ticket.id}`);
                          }}
                          className="px-3 py-1 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors cursor-pointer"
                        >
                          Abrir chamado
                        </button>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          <div className="border-t border-black bg-neutral-50 px-4 py-3 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs font-mono">
            <div>
              <span>
                Total de {totalItems} chamado{totalItems === 1 ? '' : 's'} — Página {page} de{' '}
                {totalPages || 1}
              </span>
            </div>

            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1 || loading}
                className="px-3 py-1 border border-black bg-white text-black disabled:opacity-30 disabled:cursor-not-allowed hover:bg-neutral-100 transition-colors cursor-pointer"
              >
                Anterior
              </button>

              <span className="px-2 font-bold">{page}</span>

              <button
                type="button"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages || loading}
                className="px-3 py-1 border border-black bg-white text-black disabled:opacity-30 disabled:cursor-not-allowed hover:bg-neutral-100 transition-colors cursor-pointer"
              >
                Próxima
              </button>
            </div>
          </div>
        </div>
      </main>

      {showAgentModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center p-4 z-50">
          <div className="border border-black bg-white max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between border-b border-black pb-3">
              <h2 className="text-base font-bold uppercase tracking-wide">Cadastrar Novo Agente</h2>
              <button
                type="button"
                onClick={() => setShowAgentModal(false)}
                className="px-2 py-1 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors cursor-pointer"
              >
                ✕
              </button>
            </div>

            {agentModalError && (
              <div className="border border-black bg-neutral-100 p-3 text-xs font-mono">
                <strong>Erro:</strong> {agentModalError}
              </div>
            )}
            {agentModalSuccess && (
              <div className="border border-black bg-neutral-50 p-3 text-xs font-mono font-semibold text-neutral-800">
                ✓ {agentModalSuccess}
              </div>
            )}

            <form onSubmit={handleCreateAgent} className="space-y-4">
              <div>
                <label htmlFor="agent-name" className="block text-xs font-mono uppercase mb-1">
                  Nome do Agente *
                </label>
                <input
                  id="agent-name"
                  type="text"
                  required
                  placeholder="Ex: João Pereira"
                  value={newAgentName}
                  onChange={(e) => setNewAgentName(e.target.value)}
                  className="w-full px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
                />
              </div>

              <div className="pt-2 flex items-center justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowAgentModal(false)}
                  className="px-4 py-2 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors cursor-pointer"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={savingAgent}
                  className="px-5 py-2 border border-black bg-white text-black text-xs font-mono uppercase font-bold hover:bg-neutral-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
                >
                  {savingAgent ? 'Cadastrando...' : 'Cadastrar Agente'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
