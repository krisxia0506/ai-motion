import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { MdFileDownload, MdCheck } from 'react-icons/md';
import { useNovel } from '../hooks/useNovel';
import { mediaApi } from '../services';
import type { ExportConfig, ExportFormat, VideoQuality, ExportTask } from '../types';
import { Card, CardBody, CardHeader } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { EmptyState } from '../components/common/EmptyState';
import { ProgressBar } from '../components/common/ProgressBar';

const EXPORT_FORMATS: { value: ExportFormat; label: string; description: string }[] = [
  { value: 'mp4', label: 'MP4', description: 'Most compatible format' },
  { value: 'mov', label: 'MOV', description: 'High quality, larger file size' },
  { value: 'webm', label: 'WebM', description: 'Web-optimized format' },
  { value: 'avi', label: 'AVI', description: 'Legacy format' },
];

const VIDEO_QUALITIES: { value: VideoQuality; label: string; description: string }[] = [
  { value: '720p', label: '720p (HD)', description: '1280x720, smaller file size' },
  { value: '1080p', label: '1080p (Full HD)', description: '1920x1080, balanced' },
  { value: '4k', label: '4K (Ultra HD)', description: '3840x2160, largest file size' },
];

function ExportPage() {
  const { novelId } = useParams<{ novelId: string }>();
  const navigate = useNavigate();
  const { novel, loading: novelLoading } = useNovel(novelId || '');
  
  const [exportConfig, setExportConfig] = useState<ExportConfig>({
    format: 'mp4',
    quality: '1080p',
    includeSubtitles: true,
    includeAudio: true,
    audioConfig: {
      voiceOver: true,
      soundEffects: true,
      volume: 0.8,
    },
  });
  
  const [exporting, setExporting] = useState(false);
  const [exportTask, setExportTask] = useState<ExportTask | null>(null);

  if (!novelId) {
    return (
      <div className="container" style={{ padding: '48px 0', textAlign: 'center' }}>
        <EmptyState
          title="No Novel Selected"
          description="Please select a novel to export."
          action={{
            label: "Go to Novels",
            onClick: () => navigate('/novels')
          }}
        />
      </div>
    );
  }

  const handleExport = async () => {
    try {
      setExporting(true);
      const response = await mediaApi.exportNovel({
        novelId,
        config: exportConfig,
      });
      setExportTask(response.data);
    } catch (error) {
      console.error('Export failed:', error);
    } finally {
      setExporting(false);
    }
  };

  const handleDownload = async () => {
    if (!exportTask?.resultUrl) return;
    
    try {
      const response = await fetch(exportTask.resultUrl);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${novel?.title || 'export'}.${exportConfig.format}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Download failed:', error);
    }
  };

  const handleConfigChange = (field: keyof ExportConfig, value: unknown) => {
    setExportConfig((prev) => ({ ...prev, [field]: value }));
  };

  const handleAudioConfigChange = (field: string, value: unknown) => {
    setExportConfig((prev) => ({
      ...prev,
      audioConfig: {
        ...prev.audioConfig,
        [field]: value,
      },
    }));
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
          <MdFileDownload size={32} style={{ color: 'var(--color-primary)' }} />
          <h1 style={{ margin: 0 }}>Export Video</h1>
        </div>
      </div>

      {novelLoading ? (
        <div style={{ textAlign: 'center', padding: '3rem' }}>
          <p>Loading novel...</p>
        </div>
      ) : !novel ? (
        <EmptyState
          title="Novel Not Found"
          description="The requested novel could not be found."
          action={{
            label: "Go to Novels",
            onClick: () => navigate('/novels')
          }}
        />
      ) : (
        <div style={{ maxWidth: '800px', margin: '0 auto', display: 'flex', flexDirection: 'column', gap: '24px' }}>
          <Card>
            <CardHeader>
              <h3>Novel Information</h3>
            </CardHeader>
            <CardBody>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                <div>
                  <strong>Title:</strong> {novel.title}
                </div>
                <div>
                  <strong>Author:</strong> {novel.author || 'Unknown'}
                </div>
                <div>
                  <strong>Status:</strong>{' '}
                  <span style={{
                    color: novel.status === 'parsed' ? 'var(--color-success)' : 'var(--color-warning)'
                  }}>
                    {novel.status}
                  </span>
                </div>
              </div>
            </CardBody>
          </Card>

          <Card>
            <CardHeader>
              <h3>Export Configuration</h3>
            </CardHeader>
            <CardBody>
              <div style={{ display: 'flex', flexDirection: 'column', gap: '24px' }}>
                <div>
                  <label style={{ display: 'block', marginBottom: '12px', fontWeight: 600 }}>
                    Video Format
                  </label>
                  <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '12px' }}>
                    {EXPORT_FORMATS.map((format) => (
                      <button
                        key={format.value}
                        onClick={() => handleConfigChange('format', format.value)}
                        style={{
                          padding: '16px',
                          border: `2px solid ${exportConfig.format === format.value ? 'var(--color-primary)' : 'var(--color-border)'}`,
                          borderRadius: '8px',
                          background: exportConfig.format === format.value ? 'var(--color-primary-light)' : 'white',
                          cursor: 'pointer',
                          textAlign: 'left',
                          transition: 'all 0.2s ease',
                        }}
                      >
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                          {exportConfig.format === format.value && <MdCheck size={20} style={{ color: 'var(--color-primary)' }} />}
                          <strong>{format.label}</strong>
                        </div>
                        <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                          {format.description}
                        </div>
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label style={{ display: 'block', marginBottom: '12px', fontWeight: 600 }}>
                    Video Quality
                  </label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    {VIDEO_QUALITIES.map((quality) => (
                      <button
                        key={quality.value}
                        onClick={() => handleConfigChange('quality', quality.value)}
                        style={{
                          padding: '16px',
                          border: `2px solid ${exportConfig.quality === quality.value ? 'var(--color-primary)' : 'var(--color-border)'}`,
                          borderRadius: '8px',
                          background: exportConfig.quality === quality.value ? 'var(--color-primary-light)' : 'white',
                          cursor: 'pointer',
                          textAlign: 'left',
                          transition: 'all 0.2s ease',
                        }}
                      >
                        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
                          {exportConfig.quality === quality.value && <MdCheck size={20} style={{ color: 'var(--color-primary)' }} />}
                          <strong>{quality.label}</strong>
                        </div>
                        <div style={{ fontSize: '0.875rem', color: 'var(--color-text-secondary)' }}>
                          {quality.description}
                        </div>
                      </button>
                    ))}
                  </div>
                </div>

                <div>
                  <label style={{ display: 'block', marginBottom: '12px', fontWeight: 600 }}>
                    Additional Options
                  </label>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                      <input
                        type="checkbox"
                        checked={exportConfig.includeSubtitles}
                        onChange={(e) => handleConfigChange('includeSubtitles', e.target.checked)}
                        style={{ width: '18px', height: '18px' }}
                      />
                      <span>Include Subtitles</span>
                    </label>
                    <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                      <input
                        type="checkbox"
                        checked={exportConfig.includeAudio}
                        onChange={(e) => handleConfigChange('includeAudio', e.target.checked)}
                        style={{ width: '18px', height: '18px' }}
                      />
                      <span>Include Audio</span>
                    </label>
                  </div>
                </div>

                {exportConfig.includeAudio && exportConfig.audioConfig && (
                  <div>
                    <label style={{ display: 'block', marginBottom: '12px', fontWeight: 600 }}>
                      Audio Configuration
                    </label>
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                        <input
                          type="checkbox"
                          checked={exportConfig.audioConfig.voiceOver ?? false}
                          onChange={(e) => handleAudioConfigChange('voiceOver', e.target.checked)}
                          style={{ width: '18px', height: '18px' }}
                        />
                        <span>Voice Over</span>
                      </label>
                      <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                        <input
                          type="checkbox"
                          checked={exportConfig.audioConfig.soundEffects ?? false}
                          onChange={(e) => handleAudioConfigChange('soundEffects', e.target.checked)}
                          style={{ width: '18px', height: '18px' }}
                        />
                        <span>Sound Effects</span>
                      </label>
                      <div>
                        <label style={{ display: 'block', marginBottom: '8px' }}>
                          Volume: {Math.round((exportConfig.audioConfig.volume || 0.8) * 100)}%
                        </label>
                        <input
                          type="range"
                          min="0"
                          max="1"
                          step="0.1"
                          value={exportConfig.audioConfig.volume || 0.8}
                          onChange={(e) => handleAudioConfigChange('volume', parseFloat(e.target.value))}
                          style={{ width: '100%' }}
                        />
                      </div>
                    </div>
                  </div>
                )}

                <div style={{ display: 'flex', gap: '12px', marginTop: '16px' }}>
                  <Button
                    variant="primary"
                    onClick={handleExport}
                    disabled={exporting || (exportTask?.status === 'processing')}
                    style={{ flex: 1 }}
                  >
                    {exporting ? 'Starting Export...' : exportTask ? 'Export Again' : 'Start Export'}
                  </Button>
                  <Button
                    variant="secondary"
                    onClick={() => navigate(`/novels/${novelId}`)}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            </CardBody>
          </Card>

          {exportTask && (
            <Card>
              <CardHeader>
                <h3>Export Progress</h3>
              </CardHeader>
              <CardBody>
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                      <span>Status: <strong>{exportTask.status}</strong></span>
                      <span><strong>{exportTask.progress}%</strong></span>
                    </div>
                    <ProgressBar value={exportTask.progress} />
                  </div>

                  {exportTask.status === 'completed' && exportTask.resultUrl && (
                    <div style={{ 
                      padding: '16px', 
                      background: 'var(--color-success-light)', 
                      borderRadius: '8px',
                      border: '1px solid var(--color-success)'
                    }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '12px' }}>
                        <MdCheck size={24} style={{ color: 'var(--color-success)' }} />
                        <strong>Export Completed!</strong>
                      </div>
                      <Button variant="primary" onClick={handleDownload}>
                        <MdFileDownload size={20} />
                        Download Video
                      </Button>
                    </div>
                  )}

                  {exportTask.status === 'failed' && (
                    <div style={{ 
                      padding: '16px', 
                      background: '#fef2f2', 
                      borderRadius: '8px',
                      border: '1px solid var(--color-danger)',
                      color: 'var(--color-danger)'
                    }}>
                      <strong>Export Failed</strong>
                      {exportTask.error && <p style={{ margin: '8px 0 0' }}>{exportTask.error}</p>}
                    </div>
                  )}
                </div>
              </CardBody>
            </Card>
          )}
        </div>
      )}
    </div>
  );
}

export default ExportPage;
