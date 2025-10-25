import React, { useState } from 'react';
import { Character, UpdateCharacterRequest } from '../../../types';
import { Card, CardBody, CardHeader } from '../../common/Card';
import { Button } from '../../common/Button';
import { Input } from '../../common/Input';
import './CharacterDetail.css';

interface CharacterDetailProps {
  character: Character;
  onUpdate?: (id: string, data: UpdateCharacterRequest) => Promise<void>;
  onGenerateReference?: () => void;
  onClose?: () => void;
}

export const CharacterDetail: React.FC<CharacterDetailProps> = ({
  character,
  onUpdate,
  onGenerateReference,
  onClose,
}) => {
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<UpdateCharacterRequest>({
    name: character.name,
    appearance: character.appearance,
    personality: character.personality,
    background: character.background,
  });
  const [saving, setSaving] = useState(false);

  const handleInputChange = (field: keyof UpdateCharacterRequest, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleSave = async () => {
    if (!onUpdate) return;

    try {
      setSaving(true);
      await onUpdate(character.id, formData);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to update character:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = () => {
    setFormData({
      name: character.name,
      appearance: character.appearance,
      personality: character.personality,
      background: character.background,
    });
    setIsEditing(false);
  };

  return (
    <div className="character-detail">
      <div className="character-detail-header">
        <div className="character-detail-title">
          <h2>{character.name}</h2>
          <span className={`character-type-badge character-type-${character.type}`}>
            {character.type}
          </span>
        </div>
        <div className="character-detail-actions">
          {!isEditing && onUpdate && (
            <Button variant="secondary" onClick={() => setIsEditing(true)}>
              Edit
            </Button>
          )}
          {onGenerateReference && (
            <Button variant="primary" onClick={onGenerateReference}>
              Generate Reference Image
            </Button>
          )}
          {onClose && (
            <Button variant="secondary" onClick={onClose}>
              Close
            </Button>
          )}
        </div>
      </div>

      <div className="character-detail-content">
        <Card>
          <CardHeader>
            <h3>Character Information</h3>
          </CardHeader>
          <CardBody>
            <div className="character-detail-form">
              <div className="form-field">
                <label htmlFor="name">Name</label>
                {isEditing ? (
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => handleInputChange('name', e.target.value)}
                  />
                ) : (
                  <p>{character.name}</p>
                )}
              </div>

              <div className="form-field">
                <label htmlFor="appearance">Appearance</label>
                {isEditing ? (
                  <textarea
                    id="appearance"
                    value={formData.appearance}
                    onChange={(e) => handleInputChange('appearance', e.target.value)}
                    rows={4}
                    className="character-textarea"
                  />
                ) : (
                  <p>{character.appearance}</p>
                )}
              </div>

              <div className="form-field">
                <label htmlFor="personality">Personality</label>
                {isEditing ? (
                  <textarea
                    id="personality"
                    value={formData.personality}
                    onChange={(e) => handleInputChange('personality', e.target.value)}
                    rows={4}
                    className="character-textarea"
                  />
                ) : (
                  <p>{character.personality}</p>
                )}
              </div>

              <div className="form-field">
                <label htmlFor="background">Background (Optional)</label>
                {isEditing ? (
                  <textarea
                    id="background"
                    value={formData.background || ''}
                    onChange={(e) => handleInputChange('background', e.target.value)}
                    rows={4}
                    className="character-textarea"
                  />
                ) : (
                  <p>{character.background || 'No background information'}</p>
                )}
              </div>

              {isEditing && (
                <div className="form-actions">
                  <Button variant="primary" onClick={handleSave} disabled={saving}>
                    {saving ? 'Saving...' : 'Save Changes'}
                  </Button>
                  <Button variant="secondary" onClick={handleCancel} disabled={saving}>
                    Cancel
                  </Button>
                </div>
              )}
            </div>
          </CardBody>
        </Card>

        <Card>
          <CardHeader>
            <h3>Reference Images</h3>
          </CardHeader>
          <CardBody>
            {character.referenceImages.length > 0 ? (
              <div className="reference-images-grid">
                {character.referenceImages.map((image) => (
                  <div key={image.id} className="reference-image-item">
                    <img src={image.url} alt={`Reference for ${character.name}`} />
                    <div className="reference-image-info">
                      <span className={`status-badge status-${image.status}`}>
                        {image.status}
                      </span>
                      <small>{new Date(image.generatedAt).toLocaleDateString()}</small>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="no-references">
                <p>No reference images yet</p>
                <Button variant="primary" onClick={onGenerateReference}>
                  Generate First Reference Image
                </Button>
              </div>
            )}
          </CardBody>
        </Card>
      </div>
    </div>
  );
};
