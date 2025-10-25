import React from 'react';
import { useNavigate } from 'react-router-dom';
import { NovelUpload } from '../components/features/novel';
import type { Novel } from '../types';

const NovelUploadPage: React.FC = () => {
  const navigate = useNavigate();

  const handleUploadSuccess = (novel: Novel) => {
    navigate(`/novels/${novel.id}`);
  };

  return (
    <div className="page-container">
      <NovelUpload onUploadSuccess={handleUploadSuccess} />
    </div>
  );
};

export default NovelUploadPage;
