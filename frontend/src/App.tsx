import React, { Suspense, lazy, useEffect } from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Layout } from './components/Layout';
import { LoadingSpinner } from './components/common';
import { useAuthStore } from './store/authStore';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';

const HomePage = lazy(() => import('./pages/HomePage'));
const NovelListPage = lazy(() => import('./pages/NovelListPage'));
const NovelDetailPage = lazy(() => import('./pages/NovelDetailPage'));
const NovelUploadPage = lazy(() => import('./pages/NovelUploadPage'));
const CharacterPage = lazy(() => import('./pages/CharacterPage'));
const GenerationPage = lazy(() => import('./pages/GenerationPage'));
const ExportPage = lazy(() => import('./pages/ExportPage'));

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
  const { initialize } = useAuthStore();

  useEffect(() => {
    initialize();
  }, [initialize]);

  return (
    <HashRouter>
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          
          <Route path="/" element={<Layout />}>
            <Route index element={<HomePage />} />
            
            <Route path="novels">
              <Route index element={<NovelListPage />} />
              <Route path="upload" element={<NovelUploadPage />} />
              <Route path=":id" element={<NovelDetailPage />} />
              <Route path=":id/characters" element={<CharacterPage />} />
              <Route path=":id/generate" element={<GenerationPage />} />
              <Route path=":id/export" element={<ExportPage />} />
            </Route>

            <Route path="characters" element={<CharacterPage />} />
            <Route path="generation" element={<GenerationPage />} />
            <Route path="export" element={<ExportPage />} />
            
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </Suspense>
    </HashRouter>
  );
};

export default App;
