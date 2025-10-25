import React, { useState } from 'react';
import type { Character, CharacterType } from '../../../types';
import { CharacterCard } from './CharacterCard';
import { Input } from '../../common/Input';
import { LoadingSpinner } from '../../common/LoadingSpinner';
import './CharacterList.css';

interface CharacterListProps {
  characters: Character[];
  loading?: boolean;
  onEdit?: (character: Character) => void;
  onGenerateReference?: (character: Character) => void;
  onDelete?: (characterId: string) => void;
}

const CHARACTER_TYPE_FILTERS: { value: CharacterType | 'all'; label: string }[] = [
  { value: 'all', label: 'All Characters' },
  { value: 'main', label: 'Main' },
  { value: 'supporting', label: 'Supporting' },
  { value: 'minor', label: 'Minor' },
];

export const CharacterList: React.FC<CharacterListProps> = ({
  characters,
  loading = false,
  onEdit,
  onGenerateReference,
  onDelete,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [typeFilter, setTypeFilter] = useState<CharacterType | 'all'>('all');

  const filteredCharacters = characters.filter((character) => {
    const matchesSearch =
      character.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      character.appearance.toLowerCase().includes(searchQuery.toLowerCase()) ||
      character.personality.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = typeFilter === 'all' || character.type === typeFilter;

    return matchesSearch && matchesType;
  });

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="character-list-container">
      <div className="character-list-filters">
        <Input
          type="text"
          placeholder="Search characters..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="character-search"
        />

        <div className="character-type-filters">
          {CHARACTER_TYPE_FILTERS.map((filter) => (
            <button
              key={filter.value}
              className={`type-filter-btn ${typeFilter === filter.value ? 'active' : ''}`}
              onClick={() => setTypeFilter(filter.value)}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      {filteredCharacters.length === 0 ? (
        <div className="character-list-empty">
          <p>No characters found</p>
          {searchQuery && <p>Try adjusting your search or filters</p>}
        </div>
      ) : (
        <div className="character-list-grid">
          {filteredCharacters.map((character) => (
            <CharacterCard
              key={character.id}
              character={character}
              onEdit={onEdit}
              onGenerateReference={onGenerateReference}
              onDelete={onDelete}
            />
          ))}
        </div>
      )}

      <div className="character-list-count">
        Showing {filteredCharacters.length} of {characters.length} characters
      </div>
    </div>
  );
};
