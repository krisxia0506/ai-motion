import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { NovelDetail } from '../components/features/novel';

const NovelDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  if (!id) {
    navigate('/novels');
    return null;
  }

  return (
    <div className="page-container">
      <NovelDetail
        novelId={id}
        onBack={() => navigate('/novels')}
        onEdit={() => console.log('Edit not implemented yet')}
        onDelete={() => {
          if (confirm('Are you sure you want to delete this novel?')) {
            console.log('Delete not implemented yet');
            navigate('/novels');
          }
        }}
      />
    </div>
  );
};

export default NovelDetailPage;
