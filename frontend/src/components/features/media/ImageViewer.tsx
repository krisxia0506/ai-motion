import React, { useState } from 'react';
import { MdZoomIn, MdZoomOut, MdClose, MdDownload, MdFullscreen } from 'react-icons/md';
import { Button } from '../../common/Button';
import './ImageViewer.css';

interface ImageViewerProps {
  src: string;
  alt?: string;
  isOpen: boolean;
  onClose: () => void;
  onDownload?: () => void;
}

export const ImageViewer: React.FC<ImageViewerProps> = ({
  src,
  alt = 'Image',
  isOpen,
  onClose,
  onDownload,
}) => {
  const [zoom, setZoom] = useState(1);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  const handleZoomIn = () => {
    setZoom((prev) => Math.min(prev + 0.25, 3));
  };

  const handleZoomOut = () => {
    setZoom((prev) => Math.max(prev - 0.25, 0.5));
  };

  const handleReset = () => {
    setZoom(1);
    setPosition({ x: 0, y: 0 });
  };

  const handleMouseDown = (e: React.MouseEvent) => {
    if (zoom > 1) {
      setIsDragging(true);
      setDragStart({
        x: e.clientX - position.x,
        y: e.clientY - position.y,
      });
    }
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (isDragging) {
      setPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleFullscreen = () => {
    const element = document.querySelector('.image-viewer-content') as HTMLElement;
    if (element && document.fullscreenEnabled) {
      element.requestFullscreen();
    }
  };

  const handleDownload = async () => {
    try {
      const response = await fetch(src);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = alt || 'image.png';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      if (onDownload) {
        onDownload();
      }
    } catch (error) {
      console.error('Failed to download image:', error);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="image-viewer-overlay" onClick={onClose}>
      <div className="image-viewer-content" onClick={(e) => e.stopPropagation()}>
        <div className="image-viewer-toolbar">
          <div className="toolbar-left">
            <span className="image-title">{alt}</span>
          </div>
          <div className="toolbar-right">
            <Button
              variant="ghost"
              size="small"
              onClick={handleZoomOut}
              disabled={zoom <= 0.5}
              title="Zoom Out"
            >
              <MdZoomOut size={20} />
            </Button>
            <span className="zoom-level">{Math.round(zoom * 100)}%</span>
            <Button
              variant="ghost"
              size="small"
              onClick={handleZoomIn}
              disabled={zoom >= 3}
              title="Zoom In"
            >
              <MdZoomIn size={20} />
            </Button>
            <Button
              variant="ghost"
              size="small"
              onClick={handleReset}
              title="Reset"
            >
              Reset
            </Button>
            <Button
              variant="ghost"
              size="small"
              onClick={handleFullscreen}
              title="Fullscreen"
            >
              <MdFullscreen size={20} />
            </Button>
            <Button
              variant="ghost"
              size="small"
              onClick={handleDownload}
              title="Download"
            >
              <MdDownload size={20} />
            </Button>
            <Button
              variant="ghost"
              size="small"
              onClick={onClose}
              title="Close"
            >
              <MdClose size={20} />
            </Button>
          </div>
        </div>

        <div
          className="image-viewer-container"
          onMouseDown={handleMouseDown}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUp}
          onMouseLeave={handleMouseUp}
          style={{ cursor: zoom > 1 ? (isDragging ? 'grabbing' : 'grab') : 'default' }}
        >
          <img
            src={src}
            alt={alt}
            className="viewer-image"
            style={{
              transform: `scale(${zoom}) translate(${position.x / zoom}px, ${position.y / zoom}px)`,
              transition: isDragging ? 'none' : 'transform 0.2s ease',
            }}
            draggable={false}
          />
        </div>
      </div>
    </div>
  );
};
