import React, { useState } from 'react';
import { BatchGenerationRequest, GenerationType, Scene } from '../../../types';
import { Card, CardBody, CardHeader } from '../../common/Card';
import { Button } from '../../common/Button';
import './GenerationControlPanel.css';

interface GenerationControlPanelProps {
  selectedScenes: Scene[];
  onBatchGenerate: (request: BatchGenerationRequest) => Promise<void>;
  onClearSelection: () => void;
}

const GENERATION_TYPES: { value: GenerationType; label: string; description: string }[] = [
  { value: 'image', label: 'Images', description: 'Generate static images for scenes' },
  { value: 'video', label: 'Videos', description: 'Generate animated videos for scenes' },
  { value: 'audio', label: 'Audio', description: 'Generate voiceovers for dialogues' },
];

const QUALITY_OPTIONS = [
  { value: 'low', label: 'Low', description: 'Faster, lower quality' },
  { value: 'medium', label: 'Medium', description: 'Balanced speed and quality' },
  { value: 'high', label: 'High', description: 'Slower, highest quality' },
];

const ASPECT_RATIOS = [
  { value: '16:9', label: '16:9 (Widescreen)' },
  { value: '4:3', label: '4:3 (Standard)' },
  { value: '1:1', label: '1:1 (Square)' },
  { value: '9:16', label: '9:16 (Vertical)' },
];

const STYLE_PRESETS = [
  { value: 'anime', label: 'Anime' },
  { value: 'realistic', label: 'Realistic' },
  { value: 'cartoon', label: 'Cartoon' },
  { value: 'cinematic', label: 'Cinematic' },
];

export const GenerationControlPanel: React.FC<GenerationControlPanelProps> = ({
  selectedScenes,
  onBatchGenerate,
  onClearSelection,
}) => {
  const [generationType, setGenerationType] = useState<GenerationType>('image');
  const [quality, setQuality] = useState<'low' | 'medium' | 'high'>('medium');
  const [aspectRatio, setAspectRatio] = useState('16:9');
  const [style, setStyle] = useState('anime');
  const [duration, setDuration] = useState(5);
  const [generating, setGenerating] = useState(false);

  const handleGenerate = async () => {
    const request: BatchGenerationRequest = {
      sceneIds: selectedScenes.map((s) => s.id),
      type: generationType,
      config: {
        quality,
        aspectRatio,
        style,
        duration: generationType === 'video' ? duration : undefined,
      },
    };

    try {
      setGenerating(true);
      await onBatchGenerate(request);
    } catch (error) {
      console.error('Batch generation failed:', error);
    } finally {
      setGenerating(false);
    }
  };

  const isDisabled = selectedScenes.length === 0 || generating;

  return (
    <div className="generation-control-panel">
      <Card>
        <CardHeader>
          <h3>Batch Generation Control</h3>
          {selectedScenes.length > 0 && (
            <span className="selection-count">
              {selectedScenes.length} scene{selectedScenes.length !== 1 ? 's' : ''} selected
            </span>
          )}
        </CardHeader>
        <CardBody>
          <div className="control-section">
            <label className="control-label">Generation Type</label>
            <div className="generation-type-selector">
              {GENERATION_TYPES.map((type) => (
                <button
                  key={type.value}
                  className={`type-option ${generationType === type.value ? 'active' : ''}`}
                  onClick={() => setGenerationType(type.value)}
                  disabled={isDisabled}
                >
                  <div className="type-label">{type.label}</div>
                  <div className="type-description">{type.description}</div>
                </button>
              ))}
            </div>
          </div>

          <div className="control-section">
            <label className="control-label">Quality</label>
            <div className="quality-selector">
              {QUALITY_OPTIONS.map((option) => (
                <button
                  key={option.value}
                  className={`quality-option ${quality === option.value ? 'active' : ''}`}
                  onClick={() => setQuality(option.value as 'low' | 'medium' | 'high')}
                  disabled={isDisabled}
                >
                  <div className="option-label">{option.label}</div>
                  <div className="option-description">{option.description}</div>
                </button>
              ))}
            </div>
          </div>

          <div className="control-section">
            <label className="control-label">Style Preset</label>
            <div className="style-selector">
              {STYLE_PRESETS.map((preset) => (
                <button
                  key={preset.value}
                  className={`style-option ${style === preset.value ? 'active' : ''}`}
                  onClick={() => setStyle(preset.value)}
                  disabled={isDisabled}
                >
                  {preset.label}
                </button>
              ))}
            </div>
          </div>

          <div className="control-section">
            <label className="control-label">Aspect Ratio</label>
            <div className="aspect-ratio-selector">
              {ASPECT_RATIOS.map((ratio) => (
                <button
                  key={ratio.value}
                  className={`aspect-option ${aspectRatio === ratio.value ? 'active' : ''}`}
                  onClick={() => setAspectRatio(ratio.value)}
                  disabled={isDisabled}
                >
                  {ratio.label}
                </button>
              ))}
            </div>
          </div>

          {generationType === 'video' && (
            <div className="control-section">
              <label className="control-label">Video Duration (seconds)</label>
              <div className="duration-control">
                <input
                  type="range"
                  min="3"
                  max="30"
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                  disabled={isDisabled}
                  className="duration-slider"
                />
                <span className="duration-value">{duration}s</span>
              </div>
            </div>
          )}

          <div className="control-actions">
            <Button
              variant="primary"
              onClick={handleGenerate}
              disabled={isDisabled}
              className="generate-button"
            >
              {generating ? 'Generating...' : `Generate ${selectedScenes.length} Scene${selectedScenes.length !== 1 ? 's' : ''}`}
            </Button>
            <Button
              variant="secondary"
              onClick={onClearSelection}
              disabled={selectedScenes.length === 0}
            >
              Clear Selection
            </Button>
          </div>

          {selectedScenes.length === 0 && (
            <div className="no-selection-hint">
              Select scenes from the list to start batch generation
            </div>
          )}
        </CardBody>
      </Card>
    </div>
  );
};
