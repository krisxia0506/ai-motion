import React, { useRef, useState, useEffect } from 'react';
import {
  MdPlayArrow,
  MdPause,
  MdVolumeUp,
  MdVolumeOff,
  MdFullscreen,
  MdFullscreenExit,
  MdDownload,
} from 'react-icons/md';
import { Button } from '../../common/Button';
import './VideoPlayer.css';

interface VideoPlayerProps {
  src: string;
  poster?: string;
  title?: string;
  autoPlay?: boolean;
  loop?: boolean;
  onDownload?: () => void;
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({
  src,
  poster,
  title,
  autoPlay = false,
  loop = false,
  onDownload,
}) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(1);
  const [isMuted, setIsMuted] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showControls, setShowControls] = useState(true);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const handleTimeUpdate = () => setCurrentTime(video.currentTime);
    const handleDurationChange = () => setDuration(video.duration);
    const handlePlay = () => setIsPlaying(true);
    const handlePause = () => setIsPlaying(false);
    const handleVolumeChange = () => {
      setVolume(video.volume);
      setIsMuted(video.muted);
    };

    video.addEventListener('timeupdate', handleTimeUpdate);
    video.addEventListener('durationchange', handleDurationChange);
    video.addEventListener('play', handlePlay);
    video.addEventListener('pause', handlePause);
    video.addEventListener('volumechange', handleVolumeChange);

    return () => {
      video.removeEventListener('timeupdate', handleTimeUpdate);
      video.removeEventListener('durationchange', handleDurationChange);
      video.removeEventListener('play', handlePlay);
      video.removeEventListener('pause', handlePause);
      video.removeEventListener('volumechange', handleVolumeChange);
    };
  }, []);

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => document.removeEventListener('fullscreenchange', handleFullscreenChange);
  }, []);

  const togglePlay = () => {
    const video = videoRef.current;
    if (!video) return;

    if (isPlaying) {
      video.pause();
    } else {
      video.play();
    }
  };

  const handleSeek = (e: React.ChangeEvent<HTMLInputElement>) => {
    const video = videoRef.current;
    if (!video) return;

    const time = parseFloat(e.target.value);
    video.currentTime = time;
    setCurrentTime(time);
  };

  const handleVolumeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const video = videoRef.current;
    if (!video) return;

    const newVolume = parseFloat(e.target.value);
    video.volume = newVolume;
    setVolume(newVolume);
    setIsMuted(newVolume === 0);
  };

  const toggleMute = () => {
    const video = videoRef.current;
    if (!video) return;

    video.muted = !isMuted;
    setIsMuted(!isMuted);
  };

  const toggleFullscreen = () => {
    const container = containerRef.current;
    if (!container) return;

    if (!isFullscreen) {
      container.requestFullscreen();
    } else {
      document.exitFullscreen();
    }
  };

  const handleDownload = async () => {
    try {
      const response = await fetch(src);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = title || 'video.mp4';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      
      if (onDownload) {
        onDownload();
      }
    } catch (error) {
      console.error('Failed to download video:', error);
    }
  };

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  let controlsTimeout: NodeJS.Timeout;

  const handleMouseMove = () => {
    setShowControls(true);
    clearTimeout(controlsTimeout);
    controlsTimeout = setTimeout(() => {
      if (isPlaying) {
        setShowControls(false);
      }
    }, 3000);
  };

  return (
    <div
      ref={containerRef}
      className={`video-player ${isFullscreen ? 'fullscreen' : ''}`}
      onMouseMove={handleMouseMove}
      onMouseLeave={() => isPlaying && setShowControls(false)}
    >
      <video
        ref={videoRef}
        src={src}
        poster={poster}
        autoPlay={autoPlay}
        loop={loop}
        className="video-element"
        onClick={togglePlay}
      />

      <div className={`video-controls ${showControls ? 'visible' : ''}`}>
        <div className="progress-bar-container">
          <input
            type="range"
            min="0"
            max={duration || 0}
            value={currentTime}
            onChange={handleSeek}
            className="progress-bar"
          />
        </div>

        <div className="controls-bottom">
          <div className="controls-left">
            <Button
              variant="ghost"
              size="small"
              onClick={togglePlay}
              className="control-button"
            >
              {isPlaying ? <MdPause size={24} /> : <MdPlayArrow size={24} />}
            </Button>

            <div className="volume-control">
              <Button
                variant="ghost"
                size="small"
                onClick={toggleMute}
                className="control-button"
              >
                {isMuted || volume === 0 ? <MdVolumeOff size={20} /> : <MdVolumeUp size={20} />}
              </Button>
              <input
                type="range"
                min="0"
                max="1"
                step="0.1"
                value={isMuted ? 0 : volume}
                onChange={handleVolumeChange}
                className="volume-slider"
              />
            </div>

            <div className="time-display">
              <span>{formatTime(currentTime)}</span>
              <span> / </span>
              <span>{formatTime(duration)}</span>
            </div>
          </div>

          <div className="controls-right">
            {title && <span className="video-title">{title}</span>}
            <Button
              variant="ghost"
              size="small"
              onClick={handleDownload}
              className="control-button"
              title="Download"
            >
              <MdDownload size={20} />
            </Button>
            <Button
              variant="ghost"
              size="small"
              onClick={toggleFullscreen}
              className="control-button"
              title={isFullscreen ? 'Exit Fullscreen' : 'Fullscreen'}
            >
              {isFullscreen ? <MdFullscreenExit size={20} /> : <MdFullscreen size={20} />}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
