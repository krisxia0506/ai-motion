import React, { useState } from 'react';
import { Input, Select, EmptyState } from '../../common';
import { CharacterCard } from './CharacterCard';
import type { Character } from '../../../types';
import './CharacterList.css';

interface CharacterListProps {
  characters: Character[];
  onCharacterClick?: (character: Character) => void;
  onEdit?: (character: Character) => void;
  onDelete?: (characterId: string) => void;
  onGenerateImage?: (character: Character) => void;
}

export const CharacterList: React.FC<CharacterListProps> = ({
  characters,
  onCharacterClick,
  onEdit,
  onDelete,
  onGenerateImage,
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [typeFilter, setTypeFilter] = useState('all');

  const filteredCharacters = characters.filter((character) => {
    const matchesSearch = character.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesType = typeFilter === 'all' || character.type === typeFilter;
    return matchesSearch && matchesType;
  });

  return (
    <div className="character-list">
      <div className="character-list-filters">
        <Input
          placeholder="Search characters..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          fullWidth
        />
        <Select
          options={[
            { value: 'all', label: 'All Types' },
            { value: 'main', label: 'Main' },
            { value: 'supporting', label: 'Supporting' },
            { value: 'minor', label: 'Minor' },
          ]}
          value={typeFilter}
          onChange={setTypeFilter}
        />
      </div>

      {filteredCharacters.length > 0 ? (
        <div className="character-list-grid">
          {filteredCharacters.map((character) => (
            <CharacterCard
              key={character.id}
              character={character}
              onClick={() => onCharacterClick?.(character)}
              onEdit={() => onEdit?.(character)}
              onDelete={() => onDelete?.(character.id)}
              onGenerateImage={() => onGenerateImage?.(character)}
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title={searchTerm || typeFilter !== 'all' ? 'No characters found' : 'No characters yet'}
          description={
            searchTerm || typeFilter !== 'all'
              ? 'Try adjusting your filters'
              : 'Characters will appear here once extracted from the novel'
          }
        />
      )}
    </div>
  );
};
