import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button, Card, CardBody, LoadingSpinner, ErrorMessage } from '../components/common';
import { useNovel } from '../hooks/useNovel';
import './NovelDetailPage.css';

const NovelDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { novel, loading, error } = useNovel(id || '');
  const [expandedChapters, setExpandedChapters] = useState<Set<number>>(new Set());

  if (loading) {
    return (
      <div className="novel-detail-page">
        <div className="novel-detail-loading">
          <LoadingSpinner />
          <p>Loading novel details...</p>
        </div>
      </div>
    );
  }

  if (error || !novel) {
    return (
      <div className="novel-detail-page">
        <ErrorMessage
          title="Failed to load novel"
          message={error || 'Novel not found'}
          onRetry={() => window.location.reload()}
        />
      </div>
    );
  }

  const toggleChapter = (index: number) => {
    setExpandedChapters((prev) => {
      const next = new Set(prev);
      if (next.has(index)) {
        next.delete(index);
      } else {
        next.add(index);
      }
      return next;
    });
  };

  const getStatusBadge = (status: string) => {
    const config: Record<string, { label: string; color: string }> = {
      uploaded: { label: 'Uploaded', color: '#3b82f6' },
      parsing: { label: 'Parsing', color: '#f59e0b' },
      completed: { label: 'Completed', color: '#10b981' },
      failed: { label: 'Failed', color: '#ef4444' },
    };
    const { label, color } = config[status] || config.uploaded;
    return (
      <span className="novel-detail-badge" style={{ backgroundColor: `${color}20`, color }}>
        {label}
      </span>
    );
  };

  return (
    <div className="novel-detail-page">
      <div className="novel-detail-container">
        <div className="novel-detail-header">
          <Button variant="outline" size="small" onClick={() => navigate('/novels')}>
            ← Back to Novels
          </Button>
        </div>

        <Card className="novel-detail-card">
          <CardBody>
            <div className="novel-detail-title-row">
              <div>
                <h1 className="novel-detail-title">{novel.title}</h1>
                <p className="novel-detail-author">by {novel.author}</p>
              </div>
              {getStatusBadge(novel.status)}
            </div>

            <div className="novel-detail-metadata">
              <div className="novel-detail-meta-item">
                <span className="novel-detail-meta-label">Created:</span>
                <span className="novel-detail-meta-value">
                  {new Date(novel.createdAt).toLocaleDateString('en-US', {
                    year: 'numeric',
                    month: 'long',
                    day: 'numeric',
                  })}
                </span>
              </div>

              {novel.metadata?.characterCount !== undefined && (
                <div className="novel-detail-meta-item">
                  <span className="novel-detail-meta-label">Characters:</span>
                  <span className="novel-detail-meta-value">{novel.metadata.characterCount}</span>
                </div>
              )}

              {novel.metadata?.chapterCount !== undefined && (
                <div className="novel-detail-meta-item">
                  <span className="novel-detail-meta-label">Chapters:</span>
                  <span className="novel-detail-meta-value">{novel.metadata.chapterCount}</span>
                </div>
              )}

              {novel.metadata?.sceneCount !== undefined && (
                <div className="novel-detail-meta-item">
                  <span className="novel-detail-meta-label">Scenes:</span>
                  <span className="novel-detail-meta-value">{novel.metadata.sceneCount}</span>
                </div>
              )}
            </div>

            {novel.content && (
              <div className="novel-detail-preview">
                <h3>Content Preview</h3>
                <div className="novel-detail-preview-text">
                  {novel.content.substring(0, 500)}
                  {novel.content.length > 500 && '...'}
                </div>
              </div>
            )}

            <div className="novel-detail-actions">
              <Button
                variant="primary"
                onClick={() => navigate(`/novels/${novel.id}/characters`)}
                disabled={novel.status !== 'completed'}
              >
                View Characters
              </Button>
              <Button
                variant="primary"
                onClick={() => navigate(`/novels/${novel.id}/generate`)}
                disabled={novel.status !== 'completed'}
              >
                Generate Scenes
              </Button>
              <Button
                variant="primary"
                onClick={() => navigate(`/novels/${novel.id}/export`)}
                disabled={novel.status !== 'completed'}
              >
                Export
              </Button>
              <Button variant="danger" onClick={() => {
                if (window.confirm('Are you sure you want to delete this novel?')) {
                  navigate('/novels');
                }
              }}>
                Delete
              </Button>
            </div>
          </CardBody>
        </Card>

        {novel.chapters && novel.chapters.length > 0 && (
          <Card className="novel-detail-chapters">
            <CardBody>
              <h2 className="novel-detail-chapters-title">Chapters</h2>
              <div className="novel-detail-chapters-list">
                {novel.chapters.map((chapter, index) => (
                  <div key={chapter.id} className="novel-detail-chapter">
                    <div
                      className="novel-detail-chapter-header"
                      onClick={() => toggleChapter(index)}
                    >
                      <div className="novel-detail-chapter-info">
                        <h3 className="novel-detail-chapter-title">
                          Chapter {chapter.number}: {chapter.title}
                        </h3>
                        <p className="novel-detail-chapter-meta">
                          {chapter.wordCount} words
                        </p>
                      </div>
                      <svg
                        className={`novel-detail-chapter-icon ${expandedChapters.has(index) ? 'expanded' : ''}`}
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="2"
                        strokeLinecap="round"
                        strokeLinejoin="round"
                      >
                        <polyline points="6 9 12 15 18 9" />
                      </svg>
                    </div>
                    {expandedChapters.has(index) && (
                      <div className="novel-detail-chapter-content">
                        <p>{chapter.content}</p>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </CardBody>
          </Card>
        )}
      </div>
    </div>
  );
};

export default NovelDetailPage;
