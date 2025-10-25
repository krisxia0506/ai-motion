import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Button, LoadingSpinner, ErrorMessage } from '../components/common';
import { CharacterList } from '../components/features/character/CharacterList';
import { CharacterEditor } from '../components/features/character/CharacterEditor';
import { ReferenceImageGenerator } from '../components/features/character/ReferenceImageGenerator';
import { useCharacters } from '../hooks/useCharacters';
import { useCharacterStore } from '../store/characterStore';
import { characterApi } from '../services/characterApi';
import type { Character, CreateCharacterRequest, UpdateCharacterRequest, ReferenceImage, GenerateReferenceRequest } from '../types';
import './CharacterPage.css';

const CharacterPage: React.FC = () => {
  const { id: novelId } = useParams<{ id: string }>();
  const { characters, loading, error } = useCharacters(novelId || '');
  const { addCharacter, updateCharacter, removeCharacter } = useCharacterStore();
  
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [isGeneratorOpen, setIsGeneratorOpen] = useState(false);
  const [editingCharacter, setEditingCharacter] = useState<Character | null>(null);
  const [generatingCharacter, setGeneratingCharacter] = useState<Character | null>(null);

  const handleCreate = () => {
    setEditingCharacter(null);
    setIsEditorOpen(true);
  };

  const handleEdit = (character: Character) => {
    setEditingCharacter(character);
    setIsEditorOpen(true);
  };

  const handleSave = async (data: CreateCharacterRequest | UpdateCharacterRequest) => {
    if (editingCharacter) {
      const response = await characterApi.updateCharacter(editingCharacter.id, data as UpdateCharacterRequest);
      updateCharacter(editingCharacter.id, response.data);
    } else {
      const response = await characterApi.createCharacter(data as CreateCharacterRequest);
      addCharacter(response.data);
    }
  };

  const handleDelete = async (characterId: string) => {
    if (window.confirm('Are you sure you want to delete this character?')) {
      await characterApi.deleteCharacter(characterId);
      removeCharacter(characterId);
    }
  };

  const handleGenerateImage = (character: Character) => {
    setGeneratingCharacter(character);
    setIsGeneratorOpen(true);
  };

  const handleImageGenerated = async (request: GenerateReferenceRequest): Promise<ReferenceImage> => {
    const response = await characterApi.generateReferenceImage(request);
    const image = response.data;

    if (generatingCharacter) {
      const updatedImages = [...(generatingCharacter.referenceImages || []), image];
      updateCharacter(generatingCharacter.id, { referenceImages: updatedImages });
    }
    setIsGeneratorOpen(false);
    setGeneratingCharacter(null);

    return image;
  };

  if (loading) {
    return (
      <div className="character-page">
        <div className="character-page-loading">
          <LoadingSpinner />
          <p>Loading characters...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="character-page">
        <ErrorMessage
          title="Failed to load characters"
          message={error.message || error.error}
          onRetry={() => window.location.reload()}
        />
      </div>
    );
  }

  return (
    <div className="character-page">
      <div className="character-page-header">
        <h1 className="character-page-title">Characters</h1>
        <Button variant="primary" onClick={handleCreate}>
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
          Create Character
        </Button>
      </div>

      <CharacterList
        characters={characters}
        onEdit={handleEdit}
        onDelete={handleDelete}
        onGenerateReference={handleGenerateImage}
      />

      <CharacterEditor
        character={editingCharacter || undefined}
        novelId={novelId}
        isOpen={isEditorOpen}
        onClose={() => setIsEditorOpen(false)}
        onSave={handleSave}
      />

      {generatingCharacter && (
        <ReferenceImageGenerator
          characterId={generatingCharacter.id}
          characterName={generatingCharacter.name}
          isOpen={isGeneratorOpen}
          onClose={() => {
            setIsGeneratorOpen(false);
            setGeneratingCharacter(null);
          }}
          onGenerate={handleImageGenerated}
        />
      )}
    </div>
  );
};

export default CharacterPage;
