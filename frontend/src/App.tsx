import React, { Suspense, lazy } from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/Layout';
import { LoadingSpinner } from './components/common';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';

const HomePage = lazy(() => import('./pages/HomePage'));
const NovelListPage = lazy(() => import('./pages/NovelListPage'));
const NovelDetailPage = lazy(() => import('./pages/NovelDetailPage'));
const NovelUploadPage = lazy(() => import('./pages/NovelUploadPage'));
const CharacterPage = lazy(() => import('./pages/CharacterPage'));
const GenerationPage = lazy(() => import('./pages/GenerationPage'));
const ExportPage = lazy(() => import('./pages/ExportPage'));
const LoginPage = lazy(() => import('./pages/LoginPage'));
const RegisterPage = lazy(() => import('./pages/RegisterPage'));
const TaskListPage = lazy(() => import('./pages/TaskListPage'));
const TaskDetailPage = lazy(() => import('./pages/TaskDetailPage'));
const CreateMangaPage = lazy(() => import('./pages/CreateMangaPage'));

const LoadingFallback = () => (
  <div style={{ 
    display: 'flex', 
    justifyContent: 'center', 
    alignItems: 'center', 
    minHeight: '400px' 
  }}>
    <LoadingSpinner />
  </div>
);

const App: React.FC = () => {
  return (
    <HashRouter>
      <AuthProvider>
        <Suspense fallback={<LoadingFallback />}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            
            <Route path="/" element={<Layout />}>
              <Route index element={<HomePage />} />
              
              <Route path="novels">
                <Route index element={<ProtectedRoute><NovelListPage /></ProtectedRoute>} />
                <Route path="upload" element={<ProtectedRoute><NovelUploadPage /></ProtectedRoute>} />
                <Route path=":id" element={<ProtectedRoute><NovelDetailPage /></ProtectedRoute>} />
                <Route path=":id/characters" element={<ProtectedRoute><CharacterPage /></ProtectedRoute>} />
                <Route path=":id/generate" element={<ProtectedRoute><GenerationPage /></ProtectedRoute>} />
                <Route path=":id/export" element={<ProtectedRoute><ExportPage /></ProtectedRoute>} />
              </Route>

              <Route path="characters" element={<ProtectedRoute><CharacterPage /></ProtectedRoute>} />
              <Route path="generation" element={<ProtectedRoute><GenerationPage /></ProtectedRoute>} />
              <Route path="export" element={<ProtectedRoute><ExportPage /></ProtectedRoute>} />

              <Route path="tasks" element={<ProtectedRoute><TaskListPage /></ProtectedRoute>} />
              <Route path="task/:taskId" element={<ProtectedRoute><TaskDetailPage /></ProtectedRoute>} />
              <Route path="create" element={<ProtectedRoute><CreateMangaPage /></ProtectedRoute>} />

              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Routes>
        </Suspense>
      </AuthProvider>
    </HashRouter>
  );
};

export default App;
