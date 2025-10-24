import React from 'react';
import { HashRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import HomePage from './pages/HomePage';
import NovelListPage from './pages/NovelListPage';
import CharacterPage from './pages/CharacterPage';
import GenerationPage from './pages/GenerationPage';
import ExportPage from './pages/ExportPage';

const App: React.FC = () => {
  return (
    <HashRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="novels" element={<NovelListPage />} />
          <Route path="characters" element={<CharacterPage />} />
          <Route path="generation" element={<GenerationPage />} />
          <Route path="export" element={<ExportPage />} />
        </Route>
      </Routes>
    </HashRouter>
  );
};

export default App;
