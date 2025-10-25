import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { MdAdd, MdPerson } from 'react-icons/md';
import { useCharacters } from '../hooks/useCharacters';
import { useCharacterStore } from '../store';
import { characterApi } from '../services';
import { Character, CreateCharacterRequest, UpdateCharacterRequest } from '../types';
import { CharacterList } from '../components/features/character';
import { Modal } from '../components/common/Modal';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { EmptyState } from '../components/common/EmptyState';
import { CharacterDetailModal } from '../components/features/character/CharacterDetailModal';
import { ReferenceImageGenerator } from '../components/features/character/ReferenceImageGenerator';

function CharacterPage() {
  const { novelId } = useParams<{ novelId: string }>();
  const navigate = useNavigate();
  const { characters, loading, refetch } = useCharacters(novelId || '');
  const { addCharacter, updateCharacter, removeCharacter } = useCharacterStore();
  
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedCharacter, setSelectedCharacter] = useState<Character | null>(null);
  const [showDetailModal, setShowDetailModal] = useState(false);
  const [showReferenceGenerator, setShowReferenceGenerator] = useState(false);
  const [characterForReference, setCharacterForReference] = useState<Character | null>(null);
  
  const [createForm, setCreateForm] = useState<CreateCharacterRequest>({
    novelId: novelId || '',
    name: '',
    type: 'main',
    appearance: '',
    personality: '',
    background: '',
  });
  const [creating, setCreating] = useState(false);

  if (!novelId) {
    return (
      <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
        <EmptyState
          title="No Novel Selected"
          description="Please select a novel to view its characters."
          actionLabel="Go to Novels"
          onAction={() => navigate('/novels')}
        />
      </div>
    );
  }

  const handleCreateCharacter = async () => {
    try {
      setCreating(true);
      const response = await characterApi.createCharacter({
        ...createForm,
        novelId: novelId,
      });
      addCharacter(response.data);
      setShowCreateModal(false);
      setCreateForm({
        novelId: novelId,
        name: '',
        type: 'main',
        appearance: '',
        personality: '',
        background: '',
      });
      await refetch();
    } catch (error) {
      console.error('Failed to create character:', error);
    } finally {
      setCreating(false);
    }
  };

  const handleUpdateCharacter = async (id: string, data: UpdateCharacterRequest) => {
    try {
      const response = await characterApi.updateCharacter(id, data);
      updateCharacter(id, response.data);
      await refetch();
    } catch (error) {
      console.error('Failed to update character:', error);
    }
  };

  const handleDeleteCharacter = async (characterId: string) => {
    if (!window.confirm('Are you sure you want to delete this character?')) {
      return;
    }

    try {
      await characterApi.deleteCharacter(characterId);
      removeCharacter(characterId);
      await refetch();
    } catch (error) {
      console.error('Failed to delete character:', error);
    }
  };

  const handleEdit = (character: Character) => {
    setSelectedCharacter(character);
    setShowDetailModal(true);
  };

  const handleGenerateReference = (character: Character) => {
    setCharacterForReference(character);
    setShowReferenceGenerator(true);
  };

  const handleReferenceGenerated = async () => {
    setShowReferenceGenerator(false);
    setCharacterForReference(null);
    await refetch();
  };

  return (
    <div className="container" style={{ padding: '48px 0' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: '32px' 
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <MdPerson size={32} style={{ color: 'var(--color-primary)' }} />
          <h1 style={{ margin: 0 }}>Characters</h1>
        </div>
        <Button 
          variant="primary" 
          onClick={() => setShowCreateModal(true)}
          style={{ display: 'flex', alignItems: 'center', gap: '8px' }}
        >
          <MdAdd size={20} />
          Add Character
        </Button>
      </div>

      {characters.length === 0 && !loading ? (
        <EmptyState
          title="No Characters Yet"
          description="Start by adding characters for your novel. You can extract them from the novel text or create them manually."
          actionLabel="Add Character"
          onAction={() => setShowCreateModal(true)}
        />
      ) : (
        <CharacterList
          characters={characters}
          loading={loading}
          onEdit={handleEdit}
          onGenerateReference={handleGenerateReference}
          onDelete={handleDeleteCharacter}
        />
      )}

      <Modal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        title="Create New Character"
      >
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              Name *
            </label>
            <Input
              value={createForm.name}
              onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
              placeholder="Enter character name"
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              Type *
            </label>
            <select
              value={createForm.type}
              onChange={(e) => setCreateForm({ ...createForm, type: e.target.value as 'main' | 'supporting' | 'minor' })}
              style={{
                width: '100%',
                padding: '12px',
                border: '1px solid var(--color-border)',
                borderRadius: '8px',
                fontSize: '1rem',
              }}
            >
              <option value="main">Main Character</option>
              <option value="supporting">Supporting Character</option>
              <option value="minor">Minor Character</option>
            </select>
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              Appearance *
            </label>
            <textarea
              value={createForm.appearance}
              onChange={(e) => setCreateForm({ ...createForm, appearance: e.target.value })}
              placeholder="Describe the character's physical appearance..."
              rows={4}
              style={{
                width: '100%',
                padding: '12px',
                border: '1px solid var(--color-border)',
                borderRadius: '8px',
                fontSize: '1rem',
                fontFamily: 'inherit',
                resize: 'vertical',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              Personality *
            </label>
            <textarea
              value={createForm.personality}
              onChange={(e) => setCreateForm({ ...createForm, personality: e.target.value })}
              placeholder="Describe the character's personality traits..."
              rows={4}
              style={{
                width: '100%',
                padding: '12px',
                border: '1px solid var(--color-border)',
                borderRadius: '8px',
                fontSize: '1rem',
                fontFamily: 'inherit',
                resize: 'vertical',
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontWeight: 500 }}>
              Background (Optional)
            </label>
            <textarea
              value={createForm.background}
              onChange={(e) => setCreateForm({ ...createForm, background: e.target.value })}
              placeholder="Character's backstory and history..."
              rows={4}
              style={{
                width: '100%',
                padding: '12px',
                border: '1px solid var(--color-border)',
                borderRadius: '8px',
                fontSize: '1rem',
                fontFamily: 'inherit',
                resize: 'vertical',
              }}
            />
          </div>

          <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end', marginTop: '8px' }}>
            <Button 
              variant="secondary" 
              onClick={() => setShowCreateModal(false)}
              disabled={creating}
            >
              Cancel
            </Button>
            <Button 
              variant="primary" 
              onClick={handleCreateCharacter}
              disabled={creating || !createForm.name || !createForm.appearance || !createForm.personality}
            >
              {creating ? 'Creating...' : 'Create Character'}
            </Button>
          </div>
        </div>
      </Modal>

      {selectedCharacter && (
        <CharacterDetailModal
          isOpen={showDetailModal}
          character={selectedCharacter}
          onClose={() => {
            setShowDetailModal(false);
            setSelectedCharacter(null);
          }}
          onUpdate={handleUpdateCharacter}
          onGenerateReference={() => {
            setCharacterForReference(selectedCharacter);
            setShowReferenceGenerator(true);
            setShowDetailModal(false);
          }}
        />
      )}

      {characterForReference && (
        <ReferenceImageGenerator
          isOpen={showReferenceGenerator}
          character={characterForReference}
          onClose={() => {
            setShowReferenceGenerator(false);
            setCharacterForReference(null);
          }}
          onGenerated={handleReferenceGenerated}
        />
      )}
    </div>
  );
};

export default CharacterPage;
