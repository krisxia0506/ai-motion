import React, { useState } from 'react';
import { Button, Input, Select, Modal } from '../../common';
import type { Character, CreateCharacterRequest, UpdateCharacterRequest } from '../../../types';
import './CharacterEditor.css';

interface CharacterEditorProps {
  character?: Character;
  novelId?: string;
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: CreateCharacterRequest | UpdateCharacterRequest) => Promise<void>;
}

export const CharacterEditor: React.FC<CharacterEditorProps> = ({
  character,
  novelId,
  isOpen,
  onClose,
  onSave,
}) => {
  const [formData, setFormData] = useState({
    name: character?.name || '',
    type: character?.type || 'supporting' as const,
    appearance: character?.appearance || '',
    personality: character?.personality || '',
  });
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim()) {
      setError('Name is required');
      return;
    }

    setIsSaving(true);
    setError(null);

    try {
      if (character) {
        await onSave(formData);
      } else if (novelId) {
        await onSave({ ...formData, novelId });
      }
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save character');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={character ? 'Edit Character' : 'Create Character'}>
      <form onSubmit={handleSubmit} className="character-editor">
        <Input
          label="Name"
          value={formData.name}
          onChange={(e) => handleChange('name', e.target.value)}
          placeholder="Enter character name"
          required
          fullWidth
        />

        <Select
          label="Type"
          options={[
            { value: 'main', label: 'Main Character' },
            { value: 'supporting', label: 'Supporting Character' },
            { value: 'minor', label: 'Minor Character' },
          ]}
          value={formData.type}
          onChange={(value) => handleChange('type', value)}
          fullWidth
        />

        <div className="character-editor-textarea-wrapper">
          <label htmlFor="appearance" className="character-editor-label">
            Appearance Description
          </label>
          <textarea
            id="appearance"
            value={formData.appearance}
            onChange={(e) => handleChange('appearance', e.target.value)}
            placeholder="Describe the character's physical appearance..."
            className="character-editor-textarea"
            rows={4}
          />
        </div>

        <div className="character-editor-textarea-wrapper">
          <label htmlFor="personality" className="character-editor-label">
            Personality Description
          </label>
          <textarea
            id="personality"
            value={formData.personality}
            onChange={(e) => handleChange('personality', e.target.value)}
            placeholder="Describe the character's personality traits..."
            className="character-editor-textarea"
            rows={4}
          />
        </div>

        {error && <div className="character-editor-error">{error}</div>}

        <div className="character-editor-actions">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" loading={isSaving}>
            {isSaving ? 'Saving...' : 'Save'}
          </Button>
        </div>
      </form>
    </Modal>
  );
};
