import React, { useState, useEffect, useCallback } from 'react';
import { Button, LoadingSpinner, ErrorMessage, Card } from '../../common';
import type { Novel, Chapter } from '../../../types';
import './NovelDetail.css';

interface NovelDetailProps {
  novelId: string;
  onBack?: () => void;
  onEdit?: (novel: Novel) => void;
  onDelete?: (novel: Novel) => void;
}

export const NovelDetail: React.FC<NovelDetailProps> = ({
  novelId,
  onBack,
  onEdit,
  onDelete,
}) => {
  const [novel, setNovel] = useState<Novel | null>(null);
  const [chapters, setChapters] = useState<Chapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTab, setSelectedTab] = useState<'info' | 'chapters'>('info');

  const fetchNovelDetail = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(`/api/v1/novels/${novelId}`);
      if (!response.ok) throw new Error('Failed to fetch novel');

      const data = await response.json();
      setNovel(data.novel);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load novel');
    } finally {
      setLoading(false);
    }
  }, [novelId]);

  const fetchChapters = useCallback(async () => {
    try {
      const response = await fetch(`/api/v1/novels/${novelId}/chapters`);
      if (!response.ok) throw new Error('Failed to fetch chapters');

      const data = await response.json();
      setChapters(data.chapters || []);
    } catch (err) {
      console.error('Failed to load chapters:', err);
    }
  }, [novelId]);

  useEffect(() => {
    fetchNovelDetail();
    fetchChapters();
  }, [fetchNovelDetail, fetchChapters]);

  const handleParse = async () => {
    try {
      const response = await fetch(`/api/v1/novels/${novelId}/parse`, {
        method: 'POST',
      });
      if (!response.ok) throw new Error('Failed to start parsing');
      
      fetchNovelDetail();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to start parsing');
    }
  };

  if (loading) {
    return (
      <div className="novel-detail-loading">
        <LoadingSpinner />
      </div>
    );
  }

  if (error || !novel) {
    return (
      <ErrorMessage
        title="Failed to load novel"
        message={error || 'Novel not found'}
        onRetry={fetchNovelDetail}
      />
    );
  }

  return (
    <div className="novel-detail">
      <div className="novel-detail-header">
        {onBack && (
          <button onClick={onBack} className="back-btn">
            ← Back
          </button>
        )}
        <div className="novel-detail-actions">
          {novel.status === 'uploaded' && (
            <Button onClick={handleParse} variant="primary">
              Start Parsing
            </Button>
          )}
          {onEdit && (
            <Button onClick={() => onEdit(novel)} variant="secondary">
              Edit
            </Button>
          )}
          {onDelete && (
            <Button onClick={() => onDelete(novel)} variant="danger">
              Delete
            </Button>
          )}
        </div>
      </div>

      <div className="novel-detail-content">
        <Card className="novel-detail-main">
          <div className="novel-title-section">
            <h1 className="novel-title">{novel.title}</h1>
            {novel.author && (
              <p className="novel-author">by {novel.author}</p>
            )}
          </div>

          <div className="novel-metadata">
            <div className="metadata-item">
              <span className="metadata-label">Status:</span>
              <span className={`status-badge status-${novel.status}`}>
                {novel.status.charAt(0).toUpperCase() + novel.status.slice(1)}
              </span>
            </div>
            <div className="metadata-item">
              <span className="metadata-label">Chapters:</span>
              <span className="metadata-value">{novel.chapterCount || 0}</span>
            </div>
            <div className="metadata-item">
              <span className="metadata-label">Word Count:</span>
              <span className="metadata-value">
                {(novel.wordCount || 0).toLocaleString()}
              </span>
            </div>
            <div className="metadata-item">
              <span className="metadata-label">Created:</span>
              <span className="metadata-value">
                {new Date(novel.createdAt).toLocaleString()}
              </span>
            </div>
            <div className="metadata-item">
              <span className="metadata-label">Updated:</span>
              <span className="metadata-value">
                {new Date(novel.updatedAt).toLocaleString()}
              </span>
            </div>
          </div>

          <div className="novel-tabs">
            <button
              className={`tab-btn ${selectedTab === 'info' ? 'active' : ''}`}
              onClick={() => setSelectedTab('info')}
            >
              Information
            </button>
            <button
              className={`tab-btn ${selectedTab === 'chapters' ? 'active' : ''}`}
              onClick={() => setSelectedTab('chapters')}
            >
              Chapters ({chapters.length})
            </button>
          </div>

          <div className="novel-tab-content">
            {selectedTab === 'info' && (
              <div className="novel-info">
                <h3>Content Preview</h3>
                <div className="content-preview">
                  {novel.content?.substring(0, 500)}
                  {novel.content && novel.content.length > 500 && '...'}
                </div>
              </div>
            )}

            {selectedTab === 'chapters' && (
              <div className="novel-chapters">
                {chapters.length === 0 ? (
                  <p className="no-chapters">
                    No chapters available. Parse the novel to generate chapters.
                  </p>
                ) : (
                  <div className="chapters-list">
                    {chapters.map((chapter, index) => (
                      <div key={chapter.id} className="chapter-item">
                        <div className="chapter-number">{index + 1}</div>
                        <div className="chapter-details">
                          <h4 className="chapter-title">{chapter.title}</h4>
                          <p className="chapter-info">
                            {chapter.sceneCount || 0} scenes
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            )}
          </div>
        </Card>
      </div>
    </div>
  );
};
