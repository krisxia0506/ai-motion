import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { NovelList } from '../components/features/novel';
import type { Novel } from '../types';

const NovelListPage: React.FC = () => {
  const navigate = useNavigate();

  const handleNovelSelect = (novel: Novel) => {
    navigate(`/novels/${novel.id}`);
  };

  return (
    <div className="page-container">
      <NovelList onNovelSelect={handleNovelSelect} />
    </div>
  );
};

export default NovelListPage;
