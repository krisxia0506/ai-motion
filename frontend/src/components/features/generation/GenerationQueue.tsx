import React from 'react';
import { GenerationTask } from '../../../types';
import { Card, CardBody, CardHeader } from '../../common/Card';
import { Button } from '../../common/Button';
import './GenerationQueue.css';

interface GenerationQueueProps {
  tasks: GenerationTask[];
  onCancel?: (taskId: string) => void;
  onRetry?: (taskId: string) => void;
  onClear?: () => void;
}

const getStatusColor = (status: GenerationTask['status']): string => {
  const colors = {
    pending: '#6c757d',
    processing: '#ffa500',
    completed: '#28a745',
    failed: '#dc3545',
  };
  return colors[status];
};

export const GenerationQueue: React.FC<GenerationQueueProps> = ({
  tasks,
  onCancel,
  onRetry,
  onClear,
}) => {
  const activeTasks = tasks.filter((t) => t.status === 'processing' || t.status === 'pending');
  const completedTasks = tasks.filter((t) => t.status === 'completed');
  const failedTasks = tasks.filter((t) => t.status === 'failed');

  return (
    <div className="generation-queue">
      <Card>
        <CardHeader>
          <h3>Generation Queue</h3>
          <div className="queue-stats">
            <span className="queue-stat">
              <strong>{activeTasks.length}</strong> active
            </span>
            <span className="queue-stat">
              <strong>{completedTasks.length}</strong> completed
            </span>
            {failedTasks.length > 0 && (
              <span className="queue-stat error">
                <strong>{failedTasks.length}</strong> failed
              </span>
            )}
          </div>
        </CardHeader>
        <CardBody>
          {tasks.length === 0 ? (
            <div className="queue-empty">
              <p>No generation tasks</p>
              <p className="queue-empty-hint">Tasks will appear here when you start generating</p>
            </div>
          ) : (
            <div className="task-list">
              {activeTasks.map((task) => (
                <div key={task.id} className="task-item">
                  <div className="task-info">
                    <div className="task-type-badge">{task.type}</div>
                    <div className="task-status">
                      <span style={{ color: getStatusColor(task.status) }}>
                        {task.status === 'processing' ? 'Generating...' : 'In Queue'}
                      </span>
                    </div>
                  </div>
                  {task.status === 'processing' && (
                    <div className="task-progress">
                      <div className="progress-bar">
                        <div
                          className="progress-fill"
                          style={{ width: `${task.progress}%` }}
                        ></div>
                      </div>
                      <span className="progress-text">{task.progress}%</span>
                    </div>
                  )}
                  {onCancel && (
                    <Button
                      variant="danger"
                      size="small"
                      onClick={() => onCancel(task.id)}
                    >
                      Cancel
                    </Button>
                  )}
                </div>
              ))}

              {completedTasks.map((task) => (
                <div key={task.id} className="task-item completed">
                  <div className="task-info">
                    <div className="task-type-badge">{task.type}</div>
                    <div className="task-status">
                      <span style={{ color: getStatusColor(task.status) }}>✓ Completed</span>
                    </div>
                  </div>
                  {task.resultUrl && (
                    <a
                      href={task.resultUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="view-result-link"
                    >
                      View Result
                    </a>
                  )}
                </div>
              ))}

              {failedTasks.map((task) => (
                <div key={task.id} className="task-item failed">
                  <div className="task-info">
                    <div className="task-type-badge">{task.type}</div>
                    <div className="task-status">
                      <span style={{ color: getStatusColor(task.status) }}>✗ Failed</span>
                    </div>
                  </div>
                  {task.error && <p className="task-error">{task.error}</p>}
                  {onRetry && (
                    <Button
                      variant="primary"
                      size="small"
                      onClick={() => onRetry(task.id)}
                    >
                      Retry
                    </Button>
                  )}
                </div>
              ))}
            </div>
          )}

          {tasks.length > 0 && onClear && (
            <div className="queue-actions">
              <Button variant="secondary" size="small" onClick={onClear}>
                Clear Completed
              </Button>
            </div>
          )}
        </CardBody>
      </Card>
    </div>
  );
};
