import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { Layout } from './components/Layout';
import { Login } from './pages/Login';
import { CharacterList } from './pages/CharacterList';
import { CharacterSheet } from './pages/CharacterSheet';
import { Combat } from './pages/Combat';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--ash)' }}>Loading…</div>;
  if (!user) return <Navigate to="/login" replace />;
  return <Layout>{children}</Layout>;
}

function DMRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  if (loading) return null;
  if (!user) return <Navigate to="/login" replace />;
  if (user.role !== 'dm') return <Navigate to="/characters" replace />;
  return <Layout>{children}</Layout>;
}

function AppRoutes() {
  const { user, loading } = useAuth();
  if (loading) return null;
  return (
    <Routes>
      <Route path="/login" element={user ? <Navigate to="/characters" replace /> : <Login />} />
      <Route path="/characters" element={<ProtectedRoute><CharacterList /></ProtectedRoute>} />
      <Route path="/characters/:id" element={<ProtectedRoute><CharacterSheet /></ProtectedRoute>} />
      <Route path="/combat" element={<DMRoute><Combat /></DMRoute>} />
      <Route path="*" element={<Navigate to={user ? '/characters' : '/login'} replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  );
}
