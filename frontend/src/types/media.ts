export interface MediaFile {
  id: string;
  type: MediaType;
  url: string;
  thumbnail?: string;
  filename: string;
  size: number;
  mimeType: string;
  width?: number;
  height?: number;
  duration?: number;
  relatedId?: string;
  relatedType?: 'scene' | 'character' | 'novel';
  createdAt: Date | string;
}

export type MediaType = 'image' | 'video' | 'audio';

export interface MediaUploadRequest {
  file: File;
  type: MediaType;
  relatedId?: string;
  relatedType?: 'scene' | 'character' | 'novel';
}

export interface MediaUploadResponse {
  media: MediaFile;
  message?: string;
}

export interface MediaListQuery {
  type?: MediaType;
  relatedId?: string;
  relatedType?: 'scene' | 'character' | 'novel';
  page?: number;
  pageSize?: number;
}

export interface ExportConfig {
  format: ExportFormat;
  quality: ExportQuality;
  includeAudio?: boolean;
  watermark?: WatermarkConfig;
}

export type ExportFormat = 'mp4' | 'gif' | 'images' | 'pdf';
export type ExportQuality = 'low' | 'medium' | 'high' | 'ultra';

export interface WatermarkConfig {
  enabled: boolean;
  text?: string;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right' | 'center';
  opacity?: number;
}

export interface ExportRequest {
  novelId: string;
  sceneIds?: string[];
  config: ExportConfig;
}

export interface ExportResponse {
  exportId: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  downloadUrl?: string;
  message?: string;
}
