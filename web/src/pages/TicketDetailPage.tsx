import { useEffect, useState, type FormEvent } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import {
  fetchTicketById,
  fetchAssignees,
  updateTicket,
  type Ticket,
  type Assignee,
} from '../services/api';

const STATUS_LABELS: Record<string, string> = {
  open: 'Aberto',
  in_progress: 'Em andamento',
  resolved: 'Resolvido',
  closed: 'Fechado',
};

const PRIORITY_LABELS: Record<string, string> = {
  low: 'Baixa',
  medium: 'Média',
  high: 'Alta',
};

function getPriorityLabel(priority?: string) {
  if (!priority) return '-';
  return PRIORITY_LABELS[priority] || priority;
}

function getStatusLabel(status?: string) {
  if (!status) return '-';
  return STATUS_LABELS[status] || status;
}

export default function TicketDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [ticket, setTicket] = useState<Ticket | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [editStatus, setEditStatus] = useState<string>('open');
  const [editAssignee, setEditAssignee] = useState<string>('');
  const [editAutoAssign, setEditAutoAssign] = useState<boolean>(false);
  const [savingTicket, setSavingTicket] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState<string | null>(null);
  const [saveError, setSaveError] = useState<string | null>(null);

  const [assigneesList, setAssigneesList] = useState<Assignee[]>([]);

  useEffect(() => {
    let isMounted = true;

    async function loadData() {
      if (!id) {
        setError('ID do chamado não fornecido');
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const [ticketData, assigneesData] = await Promise.all([
          fetchTicketById(id),
          fetchAssignees().catch(() => [] as Assignee[]),
        ]);

        if (isMounted) {
          setTicket(ticketData);
          setEditStatus(ticketData.status || 'open');
          setEditAssignee(ticketData.assignee || '');
          setEditAutoAssign(false);
          setAssigneesList(assigneesData);
        }
      } catch (err: unknown) {
        if (isMounted) {
          setError(err instanceof Error ? err.message : 'Erro ao carregar detalhes do chamado');
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    }

    void loadData();

    return () => {
      isMounted = false;
    };
  }, [id]);

  async function handleSaveTicket(e: FormEvent) {
    e.preventDefault();
    if (!id || !ticket) return;

    setSavingTicket(true);
    setSaveError(null);
    setSaveSuccess(null);

    try {
      const updated = await updateTicket(id, {
        status: editStatus,
        assignee: !editAutoAssign && editAssignee ? editAssignee : undefined,
        auto_assign: editAutoAssign,
      });

      setTicket(updated);
      setEditStatus(updated.status || 'open');
      setEditAssignee(updated.assignee || '');
      setEditAutoAssign(false);
      setSaveSuccess('Chamado atualizado com sucesso!');

      // Atualiza lista de atendentes para sincronizar contagem
      try {
        const newAssignees = await fetchAssignees();
        setAssigneesList(newAssignees);
      } catch {
        // Ignora erro silenciosamente
      }
    } catch (err: unknown) {
      setSaveError(err instanceof Error ? err.message : 'Erro ao atualizar chamado');
    } finally {
      setSavingTicket(false);
    }
  }

  function formatDate(isoString?: string) {
    if (!isoString) return '-';
    try {
      const date = new Date(isoString);
      return new Intl.DateTimeFormat('pt-BR', {
        dateStyle: 'short',
        timeStyle: 'short',
      }).format(date);
    } catch {
      return isoString;
    }
  }

  return (
    <div className="min-h-screen bg-white text-black">
      <header className="border-b border-black">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold tracking-tight uppercase">Sistema de Chamados</h1>
            <p className="text-xs text-neutral-500 font-mono mt-0.5">Detalhes do Chamado</p>
          </div>
          <Link
            to="/"
            className="px-3 py-1.5 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors"
          >
            Voltar para listagem
          </Link>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 py-8">
        {loading ? (
          <div className="border border-black p-12 text-center font-mono text-sm text-neutral-600 bg-neutral-50">
            Carregando detalhes do chamado...
          </div>
        ) : error || !ticket ? (
          <div className="border border-black p-6 bg-neutral-100 space-y-4">
            <h2 className="text-base font-bold uppercase">Não foi possível carregar o chamado</h2>
            <p className="text-sm font-mono text-neutral-700">
              {error || 'Chamado não encontrado.'}
            </p>
            <div>
              <button
                type="button"
                onClick={() => navigate('/')}
                className="px-4 py-2 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors cursor-pointer"
              >
                Voltar para listagem
              </button>
            </div>
          </div>
        ) : (
          <div className="space-y-6">
            {saveError && (
              <div className="border border-black bg-neutral-100 p-4 text-xs font-mono">
                <strong>Erro:</strong> {saveError}
              </div>
            )}
            {saveSuccess && (
              <div className="border border-black bg-neutral-50 p-4 text-xs font-mono font-semibold text-neutral-800">
                ✓ {saveSuccess}
              </div>
            )}

            <div className="border border-black p-6 bg-white space-y-6">
              <div className="border-b border-black pb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div>
                  <span className="text-[10px] font-mono uppercase tracking-widest text-neutral-500">
                    Chamado #{ticket.id}
                  </span>
                  <h2 className="text-xl font-bold text-black mt-1">{ticket.title}</h2>
                </div>

                <div className="flex items-center gap-2 flex-wrap">
                  <span className="text-xs font-mono px-2.5 py-1 border border-black uppercase bg-neutral-50">
                    Prioridade: {getPriorityLabel(ticket.priority)}
                  </span>
                  <span className="text-xs font-mono px-2.5 py-1 border border-black uppercase bg-neutral-50">
                    Status: {getStatusLabel(ticket.status)}
                  </span>
                </div>
              </div>

              <div>
                <span className="block text-xs font-mono uppercase text-neutral-500 mb-1">
                  Descrição
                </span>
                <div className="border border-neutral-200 bg-neutral-50 p-4">
                  <p className="text-sm text-neutral-800 whitespace-pre-wrap leading-relaxed">
                    {ticket.description}
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-2 border-t border-neutral-200 text-xs font-mono">
                <div>
                  <span className="text-neutral-500 block">Responsável Atual:</span>
                  <span className="font-semibold text-black mt-0.5 block">
                    {ticket.assignee || 'Não atribuído'}
                  </span>
                </div>
                <div>
                  <span className="text-neutral-500 block">Data de Abertura:</span>
                  <span className="font-semibold text-black mt-0.5 block">
                    {formatDate(ticket.created_at)}
                  </span>
                </div>
                <div>
                  <span className="text-neutral-500 block">Última Atualização:</span>
                  <span className="font-semibold text-black mt-0.5 block">
                    {formatDate(ticket.updated_at)}
                  </span>
                </div>
              </div>
            </div>

            <div className="border border-black p-6 bg-white">
              <div className="border-b border-black pb-3 mb-6">
                <h3 className="text-sm font-bold uppercase tracking-wide">Atualizar Chamado</h3>
                <p className="text-xs text-neutral-500 font-mono mt-0.5">
                  Altere o status ou atribua a um atendente responsável
                </p>
              </div>

              <form onSubmit={handleSaveTicket} className="space-y-6">
                <div>
                  <label htmlFor="edit-status" className="block text-xs font-mono uppercase mb-1">
                    Status do Chamado
                  </label>
                  <select
                    id="edit-status"
                    value={editStatus}
                    onChange={(e) => setEditStatus(e.target.value)}
                    className="w-full sm:w-72 px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
                  >
                    <option value="open">Aberto</option>
                    <option value="in_progress">Em andamento</option>
                    <option value="resolved">Resolvido</option>
                    <option value="closed">Fechado</option>
                  </select>
                </div>

                <div className="border border-neutral-300 p-4 bg-neutral-50 space-y-3">
                  <span className="block text-xs font-mono uppercase font-bold text-neutral-800">
                    Atribuição de Responsável
                  </span>

                  <label className="flex items-center gap-2 text-xs font-mono cursor-pointer select-none">
                    <input
                      type="checkbox"
                      checked={editAutoAssign}
                      onChange={(e) => setEditAutoAssign(e.target.checked)}
                      className="w-4 h-4 rounded-none border-black accent-black"
                    />
                    <span>Atribuir automaticamente ao atendente com menor carga</span>
                  </label>

                  {!editAutoAssign && (
                    <div className="pt-2">
                      <label
                        htmlFor="edit-assignee"
                        className="block text-xs font-mono uppercase mb-1"
                      >
                        Responsável
                      </label>
                      <select
                        id="edit-assignee"
                        value={editAssignee}
                        onChange={(e) => setEditAssignee(e.target.value)}
                        className="w-full sm:w-80 px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
                      >
                        <option value="">-- Não atribuído --</option>
                        {assigneesList.map((item) => (
                          <option key={item.id} value={item.name}>
                            {item.name} ({item.open_tickets_count ?? 0} chamados abertos)
                          </option>
                        ))}
                      </select>
                    </div>
                  )}
                </div>

                <div className="pt-4 border-t border-black flex items-center justify-end gap-3">
                  <Link
                    to="/"
                    className="px-4 py-2 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors"
                  >
                    Voltar para listagem
                  </Link>
                  <button
                    type="submit"
                    disabled={savingTicket}
                    className="px-5 py-2 border border-black bg-white text-black text-xs font-mono uppercase font-bold hover:bg-neutral-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
                  >
                    {savingTicket ? 'Salvando...' : 'Salvar Alterações'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
