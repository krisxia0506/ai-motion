import React, { useState, useMemo } from 'react';
import type { Scene, Status } from '../../../types';
import { SceneCard } from './SceneCard';
import { Input } from '../../common/Input';
import { LoadingSpinner } from '../../common/LoadingSpinner';
import './SceneList.css';

interface SceneListProps {
  scenes: Scene[];
  loading?: boolean;
  groupByChapter?: boolean;
  onGenerate?: (scene: Scene) => void;
  onView?: (scene: Scene) => void;
  onEdit?: (scene: Scene) => void;
}

const STATUS_FILTERS: { value: Status | 'all'; label: string }[] = [
  { value: 'all', label: 'All Scenes' },
  { value: 'pending', label: 'Not Generated' },
  { value: 'processing', label: 'Generating' },
  { value: 'completed', label: 'Completed' },
  { value: 'failed', label: 'Failed' },
];

export const SceneList: React.FC<SceneListProps> = ({
  scenes,
  loading = false,
  groupByChapter = false,
  onGenerate,
  onView,
  onEdit,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<Status | 'all'>('all');

  const filteredScenes = useMemo(() => {
    return scenes.filter((scene) => {
      const matchesSearch = scene.description
        .toLowerCase()
        .includes(searchQuery.toLowerCase());

      const matchesStatus = statusFilter === 'all' || scene.status === statusFilter;

      return matchesSearch && matchesStatus;
    });
  }, [scenes, searchQuery, statusFilter]);

  const groupedScenes = useMemo(() => {
    if (!groupByChapter) return { 'All Scenes': filteredScenes };

    const groups: Record<string, Scene[]> = {};
    filteredScenes.forEach((scene) => {
      const chapterKey = scene.chapterId || 'No Chapter';
      if (!groups[chapterKey]) {
        groups[chapterKey] = [];
      }
      groups[chapterKey].push(scene);
    });

    return groups;
  }, [filteredScenes, groupByChapter]);

  const stats = useMemo(() => {
    return {
      total: scenes.length,
      pending: scenes.filter((s) => s.status === 'pending').length,
      processing: scenes.filter((s) => s.status === 'processing').length,
      completed: scenes.filter((s) => s.status === 'completed').length,
      failed: scenes.filter((s) => s.status === 'failed').length,
    };
  }, [scenes]);

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="scene-list-container">
      <div className="scene-list-stats">
        <div className="stat-item">
          <span className="stat-value">{stats.total}</span>
          <span className="stat-label">Total</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">{stats.completed}</span>
          <span className="stat-label">Completed</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">{stats.processing}</span>
          <span className="stat-label">Generating</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">{stats.pending}</span>
          <span className="stat-label">Pending</span>
        </div>
      </div>

      <div className="scene-list-filters">
        <Input
          type="text"
          placeholder="Search scenes..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="scene-search"
        />

        <div className="scene-status-filters">
          {STATUS_FILTERS.map((filter) => (
            <button
              key={filter.value}
              className={`status-filter-btn ${statusFilter === filter.value ? 'active' : ''}`}
              onClick={() => setStatusFilter(filter.value)}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      {filteredScenes.length === 0 ? (
        <div className="scene-list-empty">
          <p>No scenes found</p>
          {searchQuery && <p>Try adjusting your search or filters</p>}
        </div>
      ) : (
        <div className="scene-list-content">
          {Object.entries(groupedScenes).map(([chapterKey, chapterScenes]) => (
            <div key={chapterKey} className="scene-chapter-group">
              {groupByChapter && chapterScenes.length > 0 && (
                <h3 className="chapter-title">{chapterKey}</h3>
              )}
              <div className="scene-list-grid">
                {chapterScenes.map((scene) => (
                  <SceneCard
                    key={scene.id}
                    scene={scene}
                    onGenerate={onGenerate}
                    onView={onView}
                    onEdit={onEdit}
                  />
                ))}
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="scene-list-count">
        Showing {filteredScenes.length} of {scenes.length} scenes
      </div>
    </div>
  );
};
