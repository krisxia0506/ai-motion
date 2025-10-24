import { Timestamps, Status } from './common';

export interface Scene extends Timestamps {
  id: string;
  novelId: string;
  chapterId: string;
  orderIndex: number;
  description: string;
  characters: string[];
  dialogues: Dialogue[];
  visualPrompt: string;
  status: Status;
  generatedMedia?: GeneratedMedia[];
}

export interface Dialogue {
  character: string;
  text: string;
  emotion?: string;
}

export interface GeneratedMedia {
  id: string;
  sceneId: string;
  type: 'image' | 'video';
  url: string;
  status: Status;
  metadata?: MediaMetadata;
  generatedAt: string;
}

export interface MediaMetadata {
  width?: number;
  height?: number;
  duration?: number;
  format?: string;
  size?: number;
}

export interface GenerateSceneRequest {
  sceneId: string;
  type: 'image' | 'video';
  config?: GenerationConfig;
}

export interface GenerationConfig {
  style?: string;
  aspectRatio?: string;
  quality?: 'low' | 'medium' | 'high';
  duration?: number;
}
