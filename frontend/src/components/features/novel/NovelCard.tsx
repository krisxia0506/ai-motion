import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardBody, Button } from '../../common';
import type { Novel } from '../../../types';
import './NovelCard.css';

interface NovelCardProps {
  novel: Novel;
  onDelete?: (id: string) => void;
}

export const NovelCard: React.FC<NovelCardProps> = ({ novel, onDelete }) => {
  const navigate = useNavigate();

  const getStatusBadge = (status: string) => {
    const statusConfig: Record<string, { label: string; className: string }> = {
      uploaded: { label: 'Uploaded', className: 'novel-card-badge-uploaded' },
      parsing: { label: 'Parsing', className: 'novel-card-badge-parsing' },
      completed: { label: 'Completed', className: 'novel-card-badge-completed' },
      failed: { label: 'Failed', className: 'novel-card-badge-failed' },
    };

    const config = statusConfig[status] || statusConfig.uploaded;
    return (
      <span className={`novel-card-badge ${config.className}`}>
        {config.label}
      </span>
    );
  };

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const handleView = () => {
    navigate(`/novels/${novel.id}`);
  };

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (window.confirm(`Are you sure you want to delete "${novel.title}"?`)) {
      onDelete?.(novel.id);
    }
  };

  return (
    <Card className="novel-card" onClick={handleView}>
      <CardBody>
        <div className="novel-card-header">
          <div>
            <h3 className="novel-card-title">{novel.title}</h3>
            <p className="novel-card-author">by {novel.author}</p>
          </div>
          {getStatusBadge(novel.status)}
        </div>

        <div className="novel-card-meta">
          <div className="novel-card-meta-item">
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
              <line x1="16" y1="2" x2="16" y2="6" />
              <line x1="8" y1="2" x2="8" y2="6" />
              <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
            <span>{formatDate(novel.createdAt)}</span>
          </div>

          {novel.metadata && (
            <>
              {novel.metadata.characterCount !== undefined && (
                <div className="novel-card-meta-item">
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                    <circle cx="12" cy="7" r="4" />
                  </svg>
                  <span>{novel.metadata.characterCount} characters</span>
                </div>
              )}

              {novel.metadata.chapterCount !== undefined && (
                <div className="novel-card-meta-item">
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  >
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
                  </svg>
                  <span>{novel.metadata.chapterCount} chapters</span>
                </div>
              )}
            </>
          )}
        </div>

        <div className="novel-card-actions">
          <Button size="small" variant="primary" onClick={handleView}>
            View Details
          </Button>
          <Button
            size="small"
            variant="danger"
            onClick={handleDelete}
          >
            Delete
          </Button>
        </div>
      </CardBody>
    </Card>
  );
};
