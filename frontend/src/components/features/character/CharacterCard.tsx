import React from 'react';
import { Card, CardBody } from '../../common';
import type { Character } from '../../../types';
import './CharacterCard.css';

interface CharacterCardProps {
  character: Character;
  onClick?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  onGenerateImage?: () => void;
}

export const CharacterCard: React.FC<CharacterCardProps> = ({
  character,
  onClick,
  onEdit,
  onDelete,
  onGenerateImage,
}) => {
  const getTypeBadge = (type: string) => {
    const config: Record<string, { label: string; className: string }> = {
      main: { label: 'Main', className: 'character-card-badge-main' },
      supporting: { label: 'Supporting', className: 'character-card-badge-supporting' },
      minor: { label: 'Minor', className: 'character-card-badge-minor' },
    };
    const { label, className } = config[type] || config.minor;
    return <span className={`character-card-badge ${className}`}>{label}</span>;
  };

  const handleAction = (e: React.MouseEvent, action?: () => void) => {
    e.stopPropagation();
    action?.();
  };

  const referenceImageCount = character.referenceImages?.length || 0;
  const hasReferenceImage = referenceImageCount > 0;
  const latestImage = character.referenceImages?.[0];

  return (
    <Card className="character-card" onClick={onClick}>
      <CardBody>
        <div className="character-card-header">
          <div className="character-card-avatar">
            {hasReferenceImage && latestImage?.url ? (
              <img src={latestImage.url} alt={character.name} />
            ) : (
              <div className="character-card-avatar-placeholder">
                <svg
                  width="48"
                  height="48"
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
              </div>
            )}
          </div>
          <div className="character-card-info">
            <h3 className="character-card-name">{character.name}</h3>
            {getTypeBadge(character.type)}
          </div>
        </div>

        <div className="character-card-content">
          <div className="character-card-section">
            <h4>Appearance</h4>
            <p>{character.appearance || 'No description'}</p>
          </div>
          <div className="character-card-section">
            <h4>Personality</h4>
            <p>{character.personality || 'No description'}</p>
          </div>
        </div>

        <div className="character-card-footer">
          <div className="character-card-stat">
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
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
              <circle cx="8.5" cy="8.5" r="1.5" />
              <polyline points="21 15 16 10 5 21" />
            </svg>
            <span>{referenceImageCount} reference {referenceImageCount === 1 ? 'image' : 'images'}</span>
          </div>
        </div>

        <div className="character-card-actions">
          <button
            onClick={(e) => handleAction(e, onEdit)}
            className="character-card-action-btn"
            title="Edit"
          >
            Edit
          </button>
          <button
            onClick={(e) => handleAction(e, onGenerateImage)}
            className="character-card-action-btn"
            title="Generate Reference Image"
          >
            Generate Image
          </button>
          <button
            onClick={(e) => handleAction(e, onDelete)}
            className="character-card-action-btn character-card-action-btn-danger"
            title="Delete"
          >
            Delete
          </button>
        </div>
      </CardBody>
    </Card>
  );
};
