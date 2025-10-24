import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import HomePage from './pages/HomePage';
import NovelListPage from './pages/NovelListPage';
import CharacterPage from './pages/CharacterPage';
import GenerationPage from './pages/GenerationPage';
import ExportPage from './pages/ExportPage';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="novels" element={<NovelListPage />} />
          <Route path="characters" element={<CharacterPage />} />
          <Route path="generation" element={<GenerationPage />} />
          <Route path="export" element={<ExportPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
