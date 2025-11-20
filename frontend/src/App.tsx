import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ToastProvider } from './hooks/useToast';
import ToastContainer from './components/ToastContainer';
import MainLayout from './layouts/MainLayout';
import MessagesPage from './pages/MessagesPage';
import ContactsPage from './pages/ContactsPage';
import TemplatesPage from './pages/TemplatesPage';
import SettingsPage from './pages/SettingsPage';

function App() {
  return (
    <ToastProvider>
      <BrowserRouter>
        <MainLayout>
          <Routes>
            <Route path="/" element={<Navigate to="/messages" replace />} />
            <Route path="/messages" element={<MessagesPage />} />
            <Route path="/contacts" element={<ContactsPage />} />
            <Route path="/templates" element={<TemplatesPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Routes>
        </MainLayout>
        <ToastContainer />
      </BrowserRouter>
    </ToastProvider>
  );
}

export default App;
