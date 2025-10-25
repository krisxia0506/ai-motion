import React from 'react';
import { Scene, Status } from '../../../types';
import { Card, CardBody, CardFooter } from '../../common/Card';
import { Button } from '../../common/Button';
import './SceneCard.css';

interface SceneCardProps {
  scene: Scene;
  onGenerate?: (scene: Scene) => void;
  onView?: (scene: Scene) => void;
  onEdit?: (scene: Scene) => void;
}

const getStatusBadgeClass = (status: Status): string => {
  return `scene-status-badge status-${status}`;
};

const getStatusLabel = (status: Status): string => {
  const labels: Record<Status, string> = {
    pending: 'Not Generated',
    processing: 'Generating',
    completed: 'Completed',
    failed: 'Failed',
  };
  return labels[status];
};

export const SceneCard: React.FC<SceneCardProps> = ({
  scene,
  onGenerate,
  onView,
  onEdit,
}) => {
  const latestMedia = scene.generatedMedia?.[scene.generatedMedia.length - 1];

  return (
    <Card className="scene-card" hoverable>
      <CardBody>
        <div className="scene-card-header">
          <div className="scene-number">Scene {scene.orderIndex + 1}</div>
          <span className={getStatusBadgeClass(scene.status)}>
            {getStatusLabel(scene.status)}
          </span>
        </div>

        {latestMedia && latestMedia.url && (
          <div className="scene-thumbnail">
            {latestMedia.type === 'image' ? (
              <img src={latestMedia.url} alt={`Scene ${scene.orderIndex + 1}`} />
            ) : (
              <div className="video-thumbnail">
                <video src={latestMedia.url} />
                <div className="video-overlay">▶</div>
              </div>
            )}
          </div>
        )}

        <div className="scene-description">
          <p>{scene.description}</p>
        </div>

        {scene.characters.length > 0 && (
          <div className="scene-characters">
            <strong>Characters:</strong>
            <div className="character-tags">
              {scene.characters.map((character, index) => (
                <span key={index} className="character-tag">
                  {character}
                </span>
              ))}
            </div>
          </div>
        )}

        {scene.dialogues.length > 0 && (
          <div className="scene-dialogues">
            <strong>Dialogues:</strong>
            <div className="dialogue-count">{scene.dialogues.length} lines</div>
          </div>
        )}
      </CardBody>

      <CardFooter>
        <div className="scene-card-actions">
          {scene.status === 'pending' && onGenerate && (
            <Button variant="primary" size="small" onClick={() => onGenerate(scene)}>
              Generate
            </Button>
          )}
          {scene.status === 'completed' && onView && (
            <Button variant="secondary" size="small" onClick={() => onView(scene)}>
              View
            </Button>
          )}
          {scene.status === 'failed' && onGenerate && (
            <Button variant="primary" size="small" onClick={() => onGenerate(scene)}>
              Retry
            </Button>
          )}
          {onEdit && (
            <Button variant="secondary" size="small" onClick={() => onEdit(scene)}>
              Edit
            </Button>
          )}
        </div>
      </CardFooter>
    </Card>
  );
};
