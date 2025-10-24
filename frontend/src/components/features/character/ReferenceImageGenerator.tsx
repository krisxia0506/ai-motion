import React, { useState } from 'react';
import { Button, Select, Modal, ProgressBar } from '../../common';
import { characterApi } from '../../../services/characterApi';
import type { Character, ReferenceImage } from '../../../types';
import './ReferenceImageGenerator.css';

interface ReferenceImageGeneratorProps {
  character: Character;
  isOpen: boolean;
  onClose: () => void;
  onGenerated: (image: ReferenceImage) => void;
}

export const ReferenceImageGenerator: React.FC<ReferenceImageGeneratorProps> = ({
  character,
  isOpen,
  onClose,
  onGenerated,
}) => {
  const [style, setStyle] = useState<'anime' | 'realistic' | 'cartoon' | 'semi-realistic'>('anime');
  const [customPrompt, setCustomPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [generatedImage, setGeneratedImage] = useState<ReferenceImage | null>(null);

  const handleGenerate = async () => {
    setIsGenerating(true);
    setError(null);
    setProgress(0);

    try {
      const progressInterval = setInterval(() => {
        setProgress((prev) => Math.min(prev + 5, 95));
      }, 500);

      const image = await characterApi.generateReferenceImage({
        characterId: character.id,
        style,
        customPrompt: customPrompt || undefined,
      });

      clearInterval(progressInterval);
      setProgress(100);
      setGeneratedImage(image);
      
      setTimeout(() => {
        onGenerated(image);
      }, 500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate image');
      setProgress(0);
    } finally {
      setIsGenerating(false);
    }
  };

  const handleClose = () => {
    setStyle('anime');
    setCustomPrompt('');
    setProgress(0);
    setError(null);
    setGeneratedImage(null);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title="Generate Reference Image">
      <div className="reference-image-generator">
        <div className="reference-image-generator-character">
          <h3>{character.name}</h3>
          <p className="reference-image-generator-type">{character.type} character</p>
        </div>

        {!generatedImage && (
          <>
            <Select
              label="Style"
              options={[
                { value: 'anime', label: 'Anime' },
                { value: 'realistic', label: 'Realistic' },
                { value: 'cartoon', label: 'Cartoon' },
                { value: 'semi-realistic', label: 'Semi-Realistic' },
              ]}
              value={style}
              onChange={(value) => setStyle(value as typeof style)}
              fullWidth
            />

            <div className="reference-image-generator-textarea-wrapper">
              <label htmlFor="customPrompt" className="reference-image-generator-label">
                Custom Prompt (Optional)
              </label>
              <textarea
                id="customPrompt"
                value={customPrompt}
                onChange={(e) => setCustomPrompt(e.target.value)}
                placeholder="Add additional details for the image generation..."
                className="reference-image-generator-textarea"
                rows={3}
              />
              <span className="reference-image-generator-hint">
                The character's appearance and personality will be automatically included
              </span>
            </div>

            {error && <div className="reference-image-generator-error">{error}</div>}

            {isGenerating && (
              <ProgressBar
                value={progress}
                label="Generating reference image..."
                showPercentage
              />
            )}

            <div className="reference-image-generator-actions">
              <Button type="button" variant="outline" onClick={handleClose} disabled={isGenerating}>
                Cancel
              </Button>
              <Button onClick={handleGenerate} loading={isGenerating}>
                {isGenerating ? 'Generating...' : 'Generate'}
              </Button>
            </div>
          </>
        )}

        {generatedImage && (
          <div className="reference-image-generator-result">
            <div className="reference-image-generator-preview">
              {generatedImage.url ? (
                <img src={generatedImage.url} alt={character.name} />
              ) : (
                <div className="reference-image-generator-placeholder">
                  Image generated successfully
                </div>
              )}
            </div>
            <div className="reference-image-generator-result-actions">
              <Button variant="outline" onClick={handleGenerate}>
                Regenerate
              </Button>
              <Button onClick={handleClose}>
                Done
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
};
