import type { Timestamps, Status } from './common';

export type MediaType = 'image' | 'video' | 'audio';
export type ExportFormat = 'mp4' | 'mov' | 'avi' | 'webm';
export type VideoQuality = '720p' | '1080p' | '4k';

export interface Media extends Timestamps {
  id: string;
  type: MediaType;
  url: string;
  thumbnailUrl?: string;
  status: Status;
  metadata: MediaMetadata;
}

export interface MediaMetadata {
  width: number;
  height: number;
  duration?: number;
  format: string;
  size: number;
  bitrate?: number;
}

export interface ExportConfig {
  format: ExportFormat;
  quality: VideoQuality;
  includeSubtitles: boolean;
  includeAudio: boolean;
  audioConfig?: AudioConfig;
}

export interface AudioConfig {
  backgroundMusic?: string;
  voiceOver?: boolean;
  soundEffects?: boolean;
  volume?: number;
}

export interface ExportRequest {
  novelId: string;
  config: ExportConfig;
}

export interface ExportTask extends Timestamps {
  id: string;
  novelId: string;
  status: Status;
  progress: number;
  config: ExportConfig;
  resultUrl?: string;
  error?: string;
}

export interface MediaFile extends Media {
  filename: string;
  sceneId?: string;
  characterId?: string;
}
