import React, { useState } from 'react';
import { GenerateReferenceRequest, ReferenceImage } from '../../../types';
import { Button } from '../../common/Button';
import { Modal } from '../../common/Modal';
import './ReferenceImageGenerator.css';

interface ReferenceImageGeneratorProps {
  characterId: string;
  characterName: string;
  isOpen: boolean;
  onClose: () => void;
  onGenerate: (request: GenerateReferenceRequest) => Promise<ReferenceImage>;
}

const STYLE_OPTIONS = [
  { value: 'anime', label: 'Anime' },
  { value: 'realistic', label: 'Realistic' },
  { value: 'cartoon', label: 'Cartoon' },
  { value: 'semi-realistic', label: 'Semi-Realistic' },
];

export const ReferenceImageGenerator: React.FC<ReferenceImageGeneratorProps> = ({
  characterId,
  characterName,
  isOpen,
  onClose,
  onGenerate,
}) => {
  const [customPrompt, setCustomPrompt] = useState('');
  const [selectedStyle, setSelectedStyle] = useState('anime');
  const [generating, setGenerating] = useState(false);
  const [generatedImage, setGeneratedImage] = useState<ReferenceImage | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleGenerate = async () => {
    try {
      setGenerating(true);
      setError(null);
      setGeneratedImage(null);

      const request: GenerateReferenceRequest = {
        characterId,
        prompt: customPrompt || undefined,
        style: selectedStyle,
      };

      const result = await onGenerate(request);
      setGeneratedImage(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to generate image');
    } finally {
      setGenerating(false);
    }
  };

  const handleAccept = () => {
    setGeneratedImage(null);
    setCustomPrompt('');
    onClose();
  };

  const handleRegenerate = () => {
    setGeneratedImage(null);
    handleGenerate();
  };

  const handleReset = () => {
    setGeneratedImage(null);
    setCustomPrompt('');
    setError(null);
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={`Generate Reference Image for ${characterName}`}>
      <div className="reference-generator">
        {!generatedImage && !generating && (
          <div className="generator-form">
            <div className="form-group">
              <label htmlFor="style">Art Style</label>
              <div className="style-options">
                {STYLE_OPTIONS.map((style) => (
                  <button
                    key={style.value}
                    className={`style-option ${selectedStyle === style.value ? 'active' : ''}`}
                    onClick={() => setSelectedStyle(style.value)}
                  >
                    {style.label}
                  </button>
                ))}
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="customPrompt">
                Custom Prompt (Optional)
                <span className="label-hint">Add specific details or modifications</span>
              </label>
              <textarea
                id="customPrompt"
                value={customPrompt}
                onChange={(e) => setCustomPrompt(e.target.value)}
                placeholder="e.g., 'wearing a red jacket', 'smiling', 'with long hair'"
                rows={4}
                className="custom-prompt-textarea"
              />
            </div>

            <Button variant="primary" onClick={handleGenerate} className="generate-btn">
              Generate Reference Image
            </Button>
          </div>
        )}

        {generating && (
          <div className="generating-state">
            <div className="spinner"></div>
            <h3>Generating reference image...</h3>
            <p>This may take 30-60 seconds</p>
            <div className="progress-bar">
              <div className="progress-bar-fill"></div>
            </div>
          </div>
        )}

        {generatedImage && (
          <div className="generated-result">
            <div className="generated-image-container">
              <img src={generatedImage.url} alt={`Generated reference for ${characterName}`} />
              <div className="image-status">
                <span className={`status-badge status-${generatedImage.status}`}>
                  {generatedImage.status}
                </span>
              </div>
            </div>

            <div className="result-info">
              <p className="result-meta">
                <strong>Model Used:</strong> {generatedImage.modelUsed}
              </p>
              {generatedImage.prompt && (
                <p className="result-meta">
                  <strong>Prompt:</strong> {generatedImage.prompt}
                </p>
              )}
            </div>

            <div className="result-actions">
              <Button variant="primary" onClick={handleAccept}>
                Accept & Save
              </Button>
              <Button variant="secondary" onClick={handleRegenerate}>
                Regenerate
              </Button>
              <Button variant="secondary" onClick={handleReset}>
                Start Over
              </Button>
            </div>
          </div>
        )}

        {error && (
          <div className="error-state">
            <div className="error-icon">⚠️</div>
            <h3>Generation Failed</h3>
            <p>{error}</p>
            <div className="error-actions">
              <Button variant="primary" onClick={handleGenerate}>
                Try Again
              </Button>
              <Button variant="secondary" onClick={handleReset}>
                Start Over
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
};
