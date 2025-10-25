import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { MdMovie, MdArrowBack } from 'react-icons/md';
import { SceneList } from '../components/features/scene';
import { GenerationControlPanel, GenerationQueue } from '../components/features/generation';
import { Button } from '../components/common/Button';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { EmptyState } from '../components/common/EmptyState';
import { useScenes } from '../hooks/useScenes';
import { useGenerationStore } from '../store';
import type { Scene, BatchGenerationRequest, GenerationTask } from '../types';
import { generationApi } from '../services';
import './SceneGenerationPage.css';

const SceneGenerationPage: React.FC = () => {
  const { novelId } = useParams<{ novelId: string }>();
  const navigate = useNavigate();
  const { scenes, loading, error, refetch } = useScenes(novelId || '');
  const { tasks, addTask, updateTask, clearCompletedTasks } = useGenerationStore();
  
  const [selectedScenes, setSelectedScenes] = useState<Scene[]>([]);

  const handleSceneSelect = (scene: Scene) => {
    setSelectedScenes((prev) => {
      const isSelected = prev.some((s) => s.id === scene.id);
      if (isSelected) {
        return prev.filter((s) => s.id !== scene.id);
      } else {
        return [...prev, scene];
      }
    });
  };

  const handleSelectAll = () => {
    if (selectedScenes.length === scenes.length) {
      setSelectedScenes([]);
    } else {
      setSelectedScenes(scenes);
    }
  };

  const handleBatchGenerate = async (request: BatchGenerationRequest) => {
    try {
      const response = await generationApi.batchGenerate(request);
      
      response.data.forEach((task: GenerationTask) => {
        addTask(task);
      });

      setSelectedScenes([]);
      
      setTimeout(() => refetch(), 1000);
    } catch (err) {
      console.error('Batch generation failed:', err);
      throw err;
    }
  };

  const handleSingleGenerate = async (scene: Scene) => {
    try {
      const request: BatchGenerationRequest = {
        sceneIds: [scene.id],
        type: 'image',
        config: {
          quality: 'medium',
          aspectRatio: '16:9',
          style: 'anime',
        },
      };
      
      await handleBatchGenerate(request);
    } catch (err) {
      console.error('Scene generation failed:', err);
    }
  };

  const handleCancelTask = async (taskId: string) => {
    try {
      await generationApi.cancelTask(taskId);
      updateTask(taskId, { status: 'failed', error: 'Cancelled by user' });
    } catch (err) {
      console.error('Failed to cancel task:', err);
    }
  };

  const handleRetryTask = async (taskId: string) => {
    const task = tasks.find((t) => t.id === taskId);
    if (!task) return;

    try {
      const request: BatchGenerationRequest = {
        sceneIds: [task.sceneId],
        type: task.type,
      };

      await handleBatchGenerate(request);
    } catch (err) {
      console.error('Failed to retry task:', err);
    }
  };

  const handleClearSelection = () => {
    setSelectedScenes([]);
  };

  if (!novelId) {
    return (
      <div className="scene-generation-page">
        <ErrorMessage 
          message="Novel ID is required" 
          onRetry={() => navigate('/novels')}
        />
      </div>
    );
  }

  if (loading) {
    return (
      <div className="scene-generation-page">
        <LoadingSpinner fullScreen text="Loading scenes..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="scene-generation-page">
        <ErrorMessage 
          message={error.message || 'Failed to load scenes'} 
          onRetry={refetch}
        />
      </div>
    );
  }

  return (
    <div className="scene-generation-page">
      <div className="page-header">
        <div className="header-left">
          <Button
            variant="ghost"
            onClick={() => navigate(`/novels/${novelId}`)}
            className="back-button"
          >
            <MdArrowBack size={20} />
          </Button>
          <div className="header-title">
            <MdMovie size={32} />
            <h1>Scene Generation</h1>
          </div>
        </div>
        <div className="header-actions">
          {scenes.length > 0 && (
            <Button
              variant="secondary"
              onClick={handleSelectAll}
            >
              {selectedScenes.length === scenes.length ? 'Deselect All' : 'Select All'}
            </Button>
          )}
        </div>
      </div>

      {scenes.length === 0 ? (
        <EmptyState
          icon={<MdMovie size={64} />}
          title="No scenes found"
          description="Parse your novel to extract scenes for generation"
          action={{
            label: "Go to Novel",
            onClick: () => navigate(`/novels/${novelId}`)
          }}
        />
      ) : (
        <div className="generation-content">
          <div className="generation-main">
            <SceneList
              scenes={scenes}
              loading={loading}
              onGenerate={handleSingleGenerate}
              onView={(scene) => {
                setSelectedScenes((prev) => {
                  const isSelected = prev.some((s) => s.id === scene.id);
                  return isSelected ? prev.filter((s) => s.id !== scene.id) : [...prev, scene];
                });
              }}
              onEdit={handleSceneSelect}
            />
          </div>

          <div className="generation-sidebar">
            <GenerationControlPanel
              selectedScenes={selectedScenes}
              onBatchGenerate={handleBatchGenerate}
              onClearSelection={handleClearSelection}
            />

            <GenerationQueue
              tasks={tasks.filter((task) => 
                scenes.some((scene) => scene.id === task.sceneId)
              )}
              onCancel={handleCancelTask}
              onRetry={handleRetryTask}
              onClear={clearCompletedTasks}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default SceneGenerationPage;
