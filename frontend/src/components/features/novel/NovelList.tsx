import React, { useState, useEffect } from 'react';
import { Card, LoadingSpinner, ErrorMessage, EmptyState, Select } from '../../common';
import type { Novel, NovelStatus } from '../../../types';
import './NovelList.css';

interface NovelListProps {
  onNovelSelect?: (novel: Novel) => void;
}

export const NovelList: React.FC<NovelListProps> = ({ onNovelSelect }) => {
  const [novels, setNovels] = useState<Novel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<NovelStatus | ''>('');
  const [sortBy, setSortBy] = useState<'createdAt' | 'updatedAt' | 'title'>('createdAt');
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  const pageSize = 12;

  useEffect(() => {
    fetchNovels();
  }, [page, statusFilter, sortBy, searchQuery]);

  const fetchNovels = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams({
        page: page.toString(),
        pageSize: pageSize.toString(),
        sortBy,
        sortOrder: 'desc',
      });

      if (statusFilter) params.append('status', statusFilter);
      if (searchQuery) params.append('search', searchQuery);

      const response = await fetch(`/api/v1/novels?${params}`);
      if (!response.ok) throw new Error('Failed to fetch novels');

      const data = await response.json();
      setNovels(data.novels || []);
      setTotalPages(Math.ceil((data.total || 0) / pageSize));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load novels');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    fetchNovels();
  };

  const getStatusBadgeClass = (status: NovelStatus) => {
    const classMap: Record<NovelStatus, string> = {
      uploaded: 'status-badge-uploaded',
      parsing: 'status-badge-parsing',
      parsed: 'status-badge-parsed',
      generating: 'status-badge-generating',
      failed: 'status-badge-failed',
    };
    return classMap[status] || 'status-badge-default';
  };

  const getStatusLabel = (status: NovelStatus) => {
    const labelMap: Record<NovelStatus, string> = {
      uploaded: 'Uploaded',
      parsing: 'Parsing...',
      parsed: 'Parsed',
      generating: 'Generating...',
      failed: 'Failed',
    };
    return labelMap[status] || status;
  };

  if (loading && novels.length === 0) {
    return (
      <div className="novel-list-loading">
        <LoadingSpinner />
      </div>
    );
  }

  if (error) {
    return (
      <ErrorMessage
        title="Failed to load novels"
        message={error}
        onRetry={fetchNovels}
      />
    );
  }

  return (
    <div className="novel-list">
      <div className="novel-list-header">
        <h1 className="novel-list-title">My Novels</h1>

        <form onSubmit={handleSearch} className="novel-list-search">
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search novels..."
            className="search-input"
          />
          <button type="submit" className="search-button">
            🔍
          </button>
        </form>
      </div>

      <div className="novel-list-filters">
        <Select
          placeholder="All statuses"
          value={statusFilter}
          onChange={(value) => setStatusFilter(value as NovelStatus | '')}
          options={[
            { value: '', label: 'All Statuses' },
            { value: 'uploaded', label: 'Uploaded' },
            { value: 'parsing', label: 'Parsing' },
            { value: 'parsed', label: 'Parsed' },
            { value: 'generating', label: 'Generating' },
            { value: 'completed', label: 'Completed' },
            { value: 'failed', label: 'Failed' },
          ]}
        />

        <Select
          value={sortBy}
          onChange={(value) => setSortBy(value as typeof sortBy)}
          options={[
            { value: 'createdAt', label: 'Sort by: Created Date' },
            { value: 'updatedAt', label: 'Sort by: Updated Date' },
            { value: 'title', label: 'Sort by: Title' },
          ]}
        />
      </div>

      {novels.length === 0 ? (
        <EmptyState
          title="No novels found"
          description="Upload your first novel to get started"
          action={{
            label: 'Upload Novel',
            onClick: () => window.location.href = '#/novels/upload',
          }}
        />
      ) : (
        <>
          <div className="novel-grid">
            {novels.map((novel) => (
              <Card
                key={novel.id}
                className="novel-card"
                onClick={() => onNovelSelect?.(novel)}
              >
                <div className="novel-card-header">
                  <h3 className="novel-card-title">{novel.title}</h3>
                  <span className={`status-badge ${getStatusBadgeClass(novel.status)}`}>
                    {getStatusLabel(novel.status)}
                  </span>
                </div>

                <div className="novel-card-body">
                  {novel.author && (
                    <p className="novel-card-author">by {novel.author}</p>
                  )}

                  <div className="novel-card-stats">
                    <div className="stat-item">
                      <span className="stat-label">Chapters:</span>
                      <span className="stat-value">{novel.chapterCount || 0}</span>
                    </div>
                    <div className="stat-item">
                      <span className="stat-label">Words:</span>
                      <span className="stat-value">
                        {(novel.wordCount || 0).toLocaleString()}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="novel-card-footer">
                  <span className="novel-card-date">
                    {new Date(novel.createdAt).toLocaleDateString()}
                  </span>
                </div>
              </Card>
            ))}
          </div>

          {totalPages > 1 && (
            <div className="novel-list-pagination">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="pagination-btn"
              >
                Previous
              </button>
              <span className="pagination-info">
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="pagination-btn"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};
