import React, { useState } from 'react';
import { MdImage, MdVideocam, MdDownload, MdClose } from 'react-icons/md';
import { MediaFile } from '../../../types';
import { Card, CardBody, CardHeader } from '../../common/Card';
import { Button } from '../../common/Button';
import { ImageViewer } from './ImageViewer';
import { VideoPlayer } from './VideoPlayer';
import { EmptyState } from '../../common/EmptyState';
import './MediaGallery.css';

interface MediaGalleryProps {
  mediaFiles: MediaFile[];
  title?: string;
  onDownload?: (mediaFile: MediaFile) => void;
  onDelete?: (mediaFileId: string) => void;
}

export const MediaGallery: React.FC<MediaGalleryProps> = ({
  mediaFiles,
  title = 'Media Gallery',
  onDownload,
  onDelete,
}) => {
  const [selectedMedia, setSelectedMedia] = useState<MediaFile | null>(null);
  const [viewerOpen, setViewerOpen] = useState(false);
  const [filter, setFilter] = useState<'all' | 'image' | 'video'>('all');

  const filteredMedia = mediaFiles.filter((media) => {
    if (filter === 'all') return true;
    return media.type === filter;
  });

  const handleMediaClick = (media: MediaFile) => {
    setSelectedMedia(media);
    if (media.type === 'image') {
      setViewerOpen(true);
    }
  };

  const handleCloseViewer = () => {
    setViewerOpen(false);
    setSelectedMedia(null);
  };

  const handleDownload = async (media: MediaFile) => {
    try {
      const response = await fetch(media.url);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = media.filename || `media.${media.type === 'image' ? 'png' : 'mp4'}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      if (onDownload) {
        onDownload(media);
      }
    } catch (error) {
      console.error('Failed to download media:', error);
    }
  };

  const stats = {
    total: mediaFiles.length,
    images: mediaFiles.filter((m) => m.type === 'image').length,
    videos: mediaFiles.filter((m) => m.type === 'video').length,
  };

  return (
    <div className="media-gallery">
      <Card>
        <CardHeader>
          <h3>{title}</h3>
          <div className="gallery-stats">
            <span className="stat-item">
              <strong>{stats.total}</strong> Total
            </span>
            <span className="stat-item">
              <MdImage size={16} />
              <strong>{stats.images}</strong> Images
            </span>
            <span className="stat-item">
              <MdVideocam size={16} />
              <strong>{stats.videos}</strong> Videos
            </span>
          </div>
        </CardHeader>
        <CardBody>
          <div className="gallery-filters">
            <button
              className={`filter-button ${filter === 'all' ? 'active' : ''}`}
              onClick={() => setFilter('all')}
            >
              All Media
            </button>
            <button
              className={`filter-button ${filter === 'image' ? 'active' : ''}`}
              onClick={() => setFilter('image')}
            >
              <MdImage size={16} />
              Images ({stats.images})
            </button>
            <button
              className={`filter-button ${filter === 'video' ? 'active' : ''}`}
              onClick={() => setFilter('video')}
            >
              <MdVideocam size={16} />
              Videos ({stats.videos})
            </button>
          </div>

          {filteredMedia.length === 0 ? (
            <EmptyState
              title="No Media Files"
              description={
                filter === 'all'
                  ? 'No media files have been generated yet.'
                  : `No ${filter}s available.`
              }
            />
          ) : (
            <div className="media-grid">
              {filteredMedia.map((media) => (
                <div
                  key={media.id}
                  className="media-item"
                  onClick={() => handleMediaClick(media)}
                >
                  <div className="media-thumbnail">
                    {media.type === 'image' ? (
                      <img src={media.url} alt={media.filename} />
                    ) : (
                      <div className="video-thumbnail-wrapper">
                        <video src={media.url} />
                        <div className="play-overlay">
                          <MdVideocam size={48} />
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="media-info">
                    <div className="media-type-badge">
                      {media.type === 'image' ? (
                        <MdImage size={14} />
                      ) : (
                        <MdVideocam size={14} />
                      )}
                      {media.type}
                    </div>
                    {media.filename && (
                      <span className="media-filename" title={media.filename}>
                        {media.filename}
                      </span>
                    )}
                  </div>

                  <div className="media-actions" onClick={(e) => e.stopPropagation()}>
                    <Button
                      variant="ghost"
                      size="small"
                      onClick={() => handleDownload(media)}
                      title="Download"
                    >
                      <MdDownload size={16} />
                    </Button>
                    {onDelete && (
                      <Button
                        variant="danger"
                        size="small"
                        onClick={() => onDelete(media.id)}
                        title="Delete"
                      >
                        <MdClose size={16} />
                      </Button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}

          <div className="gallery-footer">
            Showing {filteredMedia.length} of {mediaFiles.length} media files
          </div>
        </CardBody>
      </Card>

      {selectedMedia && selectedMedia.type === 'image' && (
        <ImageViewer
          isOpen={viewerOpen}
          src={selectedMedia.url}
          alt={selectedMedia.filename}
          onClose={handleCloseViewer}
          onDownload={() => handleDownload(selectedMedia)}
        />
      )}

      {selectedMedia && selectedMedia.type === 'video' && (
        <div className="video-modal-overlay" onClick={() => setSelectedMedia(null)}>
          <div className="video-modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="video-modal-header">
              <h4>{selectedMedia.filename}</h4>
              <Button
                variant="ghost"
                size="small"
                onClick={() => setSelectedMedia(null)}
              >
                <MdClose size={20} />
              </Button>
            </div>
            <VideoPlayer
              src={selectedMedia.url}
              title={selectedMedia.filename}
              onDownload={() => handleDownload(selectedMedia)}
            />
          </div>
        </div>
      )}
    </div>
  );
};
