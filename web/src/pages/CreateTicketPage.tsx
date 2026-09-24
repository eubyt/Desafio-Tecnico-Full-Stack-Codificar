import { useEffect, useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { createTicket, fetchAssignees, type Assignee } from '../services/api';

export default function CreateTicketPage() {
  const navigate = useNavigate();

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [priority, setPriority] = useState<'low' | 'medium' | 'high'>('medium');
  const [autoAssign, setAutoAssign] = useState(true);
  const [assignee, setAssignee] = useState('');

  const [assigneesList, setAssigneesList] = useState<Assignee[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAssignees()
      .then((data) => setAssigneesList(data))
      .catch(() => {
        // Ignora erro silenciosamente, permitindo fallback
      });
  }, []);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!title.trim() || !description.trim()) {
      setError('Por favor, preencha o título e a descrição.');
      return;
    }

    setSubmitting(true);
    setError(null);

    try {
      await createTicket({
        title: title.trim(),
        description: description.trim(),
        priority,
        auto_assign: autoAssign,
        assignee: !autoAssign && assignee ? assignee : undefined,
      });

      navigate('/');
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Ocorreu um erro ao cadastrar o chamado');
      }
      setSubmitting(false);
    }
  }

  return (
    <div className="min-h-screen bg-white text-black">
      <header className="border-b border-black">
        <div className="max-w-3xl mx-auto px-4 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold tracking-tight uppercase">Sistema de Chamados</h1>
            <p className="text-xs text-neutral-500 font-mono mt-0.5">Novo Registro de Chamado</p>
          </div>
          <Link
            to="/"
            className="px-3 py-1.5 border border-black bg-white text-black text-xs font-mono uppercase hover:bg-neutral-100 transition-colors"
          >
            Voltar para listagem
          </Link>
        </div>
      </header>

      <main className="max-w-3xl mx-auto px-4 py-8">
        <div className="border border-black p-6 bg-white">
          <div className="border-b border-black pb-4 mb-6">
            <h2 className="text-lg font-bold uppercase tracking-wide">Abrir Novo Chamado</h2>
            <p className="text-xs text-neutral-600 font-mono mt-1">
              Preencha as informações abaixo para encaminhar à equipe técnica.
            </p>
          </div>

          {error && (
            <div className="border border-black bg-neutral-100 p-3 mb-6 text-sm font-mono">
              <strong>Erro:</strong> {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="ticket-title" className="block text-xs font-mono uppercase mb-1">
                Título do Chamado *
              </label>
              <input
                id="ticket-title"
                type="text"
                required
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Ex: Lentidão no fechamento de caixa"
                className="w-full px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
              />
            </div>

            <div>
              <label
                htmlFor="ticket-description"
                className="block text-xs font-mono uppercase mb-1"
              >
                Descrição Detalhada *
              </label>
              <textarea
                id="ticket-description"
                required
                rows={4}
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Descreva o problema com o máximo de detalhes possível..."
                className="w-full px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
              />
            </div>

            <div>
              <label htmlFor="ticket-priority" className="block text-xs font-mono uppercase mb-1">
                Prioridade
              </label>
              <select
                id="ticket-priority"
                value={priority}
                onChange={(e) => setPriority(e.target.value as 'low' | 'medium' | 'high')}
                className="w-full sm:w-60 px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
              >
                <option value="low">Baixa</option>
                <option value="medium">Média</option>
                <option value="high">Alta</option>
              </select>
            </div>

            <div className="border border-neutral-300 p-4 bg-neutral-50 space-y-3">
              <span className="block text-xs font-mono uppercase font-bold text-neutral-800">
                Atribuição do Responsável
              </span>

              <label className="flex items-center gap-2 text-sm cursor-pointer select-none">
                <input
                  type="checkbox"
                  checked={autoAssign}
                  onChange={(e) => setAutoAssign(e.target.checked)}
                  className="w-4 h-4 rounded-none border-black accent-black"
                />
                <span>Atribuir automaticamente ao atendente com menor carga</span>
              </label>

              {!autoAssign && (
                <div className="pt-2">
                  <label
                    htmlFor="ticket-assignee"
                    className="block text-xs font-mono uppercase mb-1"
                  >
                    Selecione o Atendente
                  </label>
                  <select
                    id="ticket-assignee"
                    value={assignee}
                    onChange={(e) => setAssignee(e.target.value)}
                    className="w-full sm:w-72 px-3 py-2 border border-black bg-white text-sm focus:outline-none focus:ring-1 focus:ring-black"
                  >
                    <option value="">-- Não atribuir no momento --</option>
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
                className="px-4 py-2 border border-black bg-white text-black text-sm font-mono uppercase hover:bg-neutral-100 transition-colors"
              >
                Cancelar
              </Link>
              <button
                type="submit"
                disabled={submitting}
                className="px-5 py-2 border border-black bg-white text-black text-sm font-mono uppercase hover:bg-neutral-100 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
              >
                {submitting ? 'Criando...' : 'Criar Chamado'}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}
