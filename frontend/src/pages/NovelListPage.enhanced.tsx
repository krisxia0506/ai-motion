import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Input, Select, LoadingSpinner, EmptyState } from '../components/common';
import { NovelCard } from '../components/features/novel/NovelCard';
import { useNovelStore } from '../store/novelStore';
import { useNovels } from '../hooks/useNovels';
import { novelApi } from '../services/novelApi';
import './NovelListPage.css';

const NovelListPage: React.FC = () => {
  const navigate = useNavigate();
  const { novels, loading, error } = useNovels();
  const removeNovel = useNovelStore((state) => state.removeNovel);
  
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [sortBy, setSortBy] = useState('createdAt');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 12;

  const filteredNovels = novels.filter((novel) => {
    const matchesSearch =
      novel.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      novel.author.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesStatus = statusFilter === 'all' || novel.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  const sortedNovels = [...filteredNovels].sort((a, b) => {
    if (sortBy === 'createdAt') {
      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    }
    if (sortBy === 'title') {
      return a.title.localeCompare(b.title);
    }
    if (sortBy === 'author') {
      return a.author.localeCompare(b.author);
    }
    return 0;
  });

  const totalPages = Math.ceil(sortedNovels.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const paginatedNovels = sortedNovels.slice(startIndex, startIndex + itemsPerPage);

  useEffect(() => {
    setCurrentPage(1);
  }, [searchTerm, statusFilter, sortBy]);

  const handleDelete = async (id: string) => {
    try {
      await novelApi.deleteNovel(id);
      removeNovel(id);
    } catch (err) {
      console.error('Failed to delete novel:', err);
    }
  };

  if (loading) {
    return (
      <div className="novel-list-page">
        <div className="novel-list-loading">
          <LoadingSpinner />
          <p>Loading novels...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="novel-list-page">
      <div className="novel-list-container">
        <div className="novel-list-header">
          <h1 className="novel-list-title">My Novels</h1>
          <Button variant="primary" onClick={() => navigate('/novels/upload')}>
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              style={{ marginRight: '0.5rem' }}
            >
              <line x1="12" y1="5" x2="12" y2="19" />
              <line x1="5" y1="12" x2="19" y2="12" />
            </svg>
            Upload Novel
          </Button>
        </div>

        <div className="novel-list-filters">
          <div className="novel-list-search">
            <Input
              placeholder="Search by title or author..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              fullWidth
            />
          </div>
          
          <div className="novel-list-filter-controls">
            <Select
              options={[
                { value: 'all', label: 'All Status' },
                { value: 'uploaded', label: 'Uploaded' },
                { value: 'parsing', label: 'Parsing' },
                { value: 'completed', label: 'Completed' },
                { value: 'failed', label: 'Failed' },
              ]}
              value={statusFilter}
              onChange={setStatusFilter}
              placeholder="Filter by status"
            />
            
            <Select
              options={[
                { value: 'createdAt', label: 'Latest' },
                { value: 'title', label: 'Title' },
                { value: 'author', label: 'Author' },
              ]}
              value={sortBy}
              onChange={setSortBy}
              placeholder="Sort by"
            />
          </div>
        </div>

        {error && (
          <div className="novel-list-error">
            <p>Error loading novels: {error}</p>
          </div>
        )}

        {paginatedNovels.length > 0 ? (
          <>
            <div className="novel-list-grid">
              {paginatedNovels.map((novel) => (
                <NovelCard key={novel.id} novel={novel} onDelete={handleDelete} />
              ))}
            </div>

            {totalPages > 1 && (
              <div className="novel-list-pagination">
                <Button
                  variant="outline"
                  size="small"
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  Previous
                </Button>
                
                <div className="novel-list-pagination-info">
                  Page {currentPage} of {totalPages}
                </div>
                
                <Button
                  variant="outline"
                  size="small"
                  onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                >
                  Next
                </Button>
              </div>
            )}
          </>
        ) : (
          <EmptyState
            title={searchTerm || statusFilter !== 'all' ? 'No novels found' : 'No novels yet'}
            description={
              searchTerm || statusFilter !== 'all'
                ? 'Try adjusting your filters'
                : 'Upload your first novel to get started!'
            }
            action={{
              label: 'Upload Novel',
              onClick: () => navigate('/novels/upload'),
            }}
          />
        )}
      </div>
    </div>
  );
};

export default NovelListPage;
