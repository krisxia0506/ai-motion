import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { MdMovie, MdCheckBox, MdCheckBoxOutlineBlank } from 'react-icons/md';
import { useScenes } from '../hooks/useScenes';
import { useGenerationStore } from '../store';
import { generationApi } from '../services';
import { Scene, BatchGenerationRequest } from '../types';
import { SceneList } from '../components/features/scene/SceneList';
import { GenerationControlPanel } from '../components/features/generation/GenerationControlPanel';
import { GenerationQueue } from '../components/features/generation/GenerationQueue';
import { EmptyState } from '../components/common/EmptyState';
import { Button } from '../components/common/Button';

function GenerationPage() {
  const { novelId } = useParams<{ novelId: string }>();
  const navigate = useNavigate();
  const { scenes, loading, refetch } = useScenes(novelId || '');
  const { tasks, addTask, updateTask, removeTask } = useGenerationStore();
  
  const [selectedScenes, setSelectedScenes] = useState<Scene[]>([]);
  const [selectionMode, setSelectionMode] = useState(false);

  if (!novelId) {
    return (
      <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
        <EmptyState
          title="No Novel Selected"
          description="Please select a novel to generate scenes."
          actionLabel="Go to Novels"
          onAction={() => navigate('/novels')}
        />
      </div>
    );
  }

  const handleSceneSelect = (scene: Scene) => {
    setSelectedScenes((prev) => {
      const isSelected = prev.find((s) => s.id === scene.id);
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
      setSelectedScenes([...scenes]);
    }
  };

  const handleClearSelection = () => {
    setSelectedScenes([]);
    setSelectionMode(false);
  };

  const handleBatchGenerate = async (request: BatchGenerationRequest) => {
    try {
      const response = await generationApi.batchGenerate(request);
      
      response.data.forEach((task) => {
        addTask(task);
      });
      
      setSelectedScenes([]);
      setSelectionMode(false);
      await refetch();
    } catch (error) {
      console.error('Batch generation failed:', error);
      throw error;
    }
  };

  const handleCancelTask = async (taskId: string) => {
    try {
      await generationApi.cancelTask(taskId);
      updateTask(taskId, { status: 'failed', error: 'Cancelled by user' });
    } catch (error) {
      console.error('Failed to cancel task:', error);
    }
  };

  const handleRetryTask = async (taskId: string) => {
    try {
      const response = await generationApi.retryTask(taskId);
      updateTask(taskId, response.data);
    } catch (error) {
      console.error('Failed to retry task:', error);
    }
  };

  const handleClearQueue = () => {
    const completedTasks = tasks.filter((t) => t.status === 'completed');
    completedTasks.forEach((t) => removeTask(t.id));
  };

  const handleGenerateSingle = (scene: Scene) => {
    setSelectedScenes([scene]);
    setSelectionMode(true);
  };

  const novelTasks = tasks.filter((t) => 
    scenes.some((s) => s.id === t.sceneId)
  );

  return (
    <div className="container" style={{ padding: '48px 0' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: '32px' 
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <MdMovie size={32} style={{ color: 'var(--color-primary)' }} />
          <h1 style={{ margin: 0 }}>Scene Generation</h1>
        </div>
        <div style={{ display: 'flex', gap: '12px' }}>
          {!selectionMode ? (
            <Button 
              variant="primary" 
              onClick={() => setSelectionMode(true)}
              disabled={scenes.length === 0}
            >
              Select Scenes for Batch
            </Button>
          ) : (
            <>
              <Button 
                variant="secondary" 
                onClick={handleSelectAll}
                style={{ display: 'flex', alignItems: 'center', gap: '8px' }}
              >
                {selectedScenes.length === scenes.length ? (
                  <MdCheckBox size={20} />
                ) : (
                  <MdCheckBoxOutlineBlank size={20} />
                )}
                {selectedScenes.length === scenes.length ? 'Deselect All' : 'Select All'}
              </Button>
              <Button 
                variant="secondary" 
                onClick={() => setSelectionMode(false)}
              >
                Cancel Selection
              </Button>
            </>
          )}
        </div>
      </div>

      {scenes.length === 0 && !loading ? (
        <EmptyState
          title="No Scenes Available"
          description="Parse your novel first to extract scenes for generation."
          actionLabel="Go to Novel Details"
          onAction={() => navigate(`/novels/${novelId}`)}
        />
      ) : (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 400px', gap: '24px' }}>
          <div>
            <div style={{ 
              marginBottom: '16px', 
              padding: '12px 16px', 
              background: 'var(--color-bg)', 
              borderRadius: '8px',
              border: '1px solid var(--color-border)' 
            }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                  {selectionMode ? (
                    <>
                      <strong>{selectedScenes.length}</strong> of <strong>{scenes.length}</strong> scenes selected
                    </>
                  ) : (
                    <>
                      <strong>{scenes.length}</strong> scenes available for generation
                    </>
                  )}
                </span>
                {selectionMode && selectedScenes.length > 0 && (
                  <Button 
                    variant="secondary" 
                    size="small"
                    onClick={handleClearSelection}
                  >
                    Clear Selection
                  </Button>
                )}
              </div>
            </div>

            <div style={{ position: 'relative' }}>
              {selectionMode && (
                <div 
                  style={{ 
                    position: 'absolute', 
                    top: 0, 
                    left: 0, 
                    right: 0, 
                    bottom: 0, 
                    zIndex: 1,
                    pointerEvents: 'none'
                  }}
                />
              )}
              <div onClick={(e) => {
                if (selectionMode) {
                  const target = e.target as HTMLElement;
                  const sceneCard = target.closest('[data-scene-id]');
                  if (sceneCard) {
                    const sceneId = sceneCard.getAttribute('data-scene-id');
                    const scene = scenes.find((s) => s.id === sceneId);
                    if (scene) {
                      handleSceneSelect(scene);
                    }
                  }
                }
              }}>
                <SceneList
                  scenes={scenes.map((scene) => ({
                    ...scene,
                    _selected: selectedScenes.find((s) => s.id === scene.id) !== undefined,
                  }))}
                  loading={loading}
                  groupByChapter={true}
                  onGenerate={!selectionMode ? handleGenerateSingle : undefined}
                />
              </div>
            </div>
          </div>

          <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
            <GenerationControlPanel
              selectedScenes={selectedScenes}
              onBatchGenerate={handleBatchGenerate}
              onClearSelection={handleClearSelection}
            />

            <GenerationQueue
              tasks={novelTasks}
              onCancel={handleCancelTask}
              onRetry={handleRetryTask}
              onClear={handleClearQueue}
            />
          </div>
        </div>
      )}
    </div>
  );
}

export default GenerationPage;
