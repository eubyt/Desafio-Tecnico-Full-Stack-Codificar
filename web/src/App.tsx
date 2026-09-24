import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import TicketListPage from './pages/TicketListPage';
import CreateTicketPage from './pages/CreateTicketPage';
import TicketDetailPage from './pages/TicketDetailPage';

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<TicketListPage />} />
        <Route path="/novo" element={<CreateTicketPage />} />
        <Route path="/chamados/:id" element={<TicketDetailPage />} />
        <Route path="/tickets/:id" element={<TicketDetailPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
