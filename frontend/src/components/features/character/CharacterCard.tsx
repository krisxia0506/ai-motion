import React from 'react';
import { Character, CharacterType } from '../../../types';
import { Card, CardBody, CardFooter } from '../../common/Card';
import { Button } from '../../common/Button';
import './CharacterCard.css';

interface CharacterCardProps {
  character: Character;
  onEdit?: (character: Character) => void;
  onGenerateReference?: (character: Character) => void;
  onDelete?: (characterId: string) => void;
}

const getCharacterTypeLabel = (type: CharacterType): string => {
  const labels: Record<CharacterType, string> = {
    main: 'Main Character',
    supporting: 'Supporting',
    minor: 'Minor',
  };
  return labels[type];
};

const getCharacterTypeBadgeClass = (type: CharacterType): string => {
  return `character-type-badge character-type-${type}`;
};

export const CharacterCard: React.FC<CharacterCardProps> = ({
  character,
  onEdit,
  onGenerateReference,
  onDelete,
}) => {
  const latestReferenceImage = character.referenceImages[character.referenceImages.length - 1];

  return (
    <Card className="character-card" hoverable>
      <CardBody>
        <div className="character-card-header">
          {latestReferenceImage ? (
            <img
              src={latestReferenceImage.url}
              alt={character.name}
              className="character-avatar"
            />
          ) : (
            <div className="character-avatar-placeholder">
              <span>{character.name.charAt(0).toUpperCase()}</span>
            </div>
          )}
          <div className="character-card-info">
            <h3 className="character-name">{character.name}</h3>
            <span className={getCharacterTypeBadgeClass(character.type)}>
              {getCharacterTypeLabel(character.type)}
            </span>
          </div>
        </div>

        <div className="character-description">
          <div className="character-field">
            <strong>Appearance:</strong>
            <p>{character.appearance.substring(0, 100)}...</p>
          </div>
          <div className="character-field">
            <strong>Personality:</strong>
            <p>{character.personality.substring(0, 100)}...</p>
          </div>
        </div>

        <div className="character-stats">
          <span className="character-stat">
            {character.referenceImages.length} reference{' '}
            {character.referenceImages.length === 1 ? 'image' : 'images'}
          </span>
        </div>
      </CardBody>

      <CardFooter>
        <div className="character-card-actions">
          {onEdit && (
            <Button variant="secondary" size="small" onClick={() => onEdit(character)}>
              Edit
            </Button>
          )}
          {onGenerateReference && (
            <Button variant="primary" size="small" onClick={() => onGenerateReference(character)}>
              Generate Image
            </Button>
          )}
          {onDelete && (
            <Button variant="danger" size="small" onClick={() => onDelete(character.id)}>
              Delete
            </Button>
          )}
        </div>
      </CardFooter>
    </Card>
  );
};
