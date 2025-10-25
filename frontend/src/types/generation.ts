import { Status } from './common';

export type GenerationType = 'image' | 'video' | 'audio';

export interface GenerationTask {
  id: string;
  type: GenerationType;
  sceneId: string;
  status: Status;
  progress: number;
  startedAt: string;
  completedAt?: string;
  error?: string;
  resultUrl?: string;
}

export interface BatchGenerationRequest {
  sceneIds: string[];
  type: GenerationType;
  config?: GenerationConfig;
}

export interface GenerationConfig {
  style?: string;
  quality?: 'low' | 'medium' | 'high';
  aspectRatio?: string;
  duration?: number;
}

export interface GenerationProgress {
  taskId: string;
  status: Status;
  progress: number;
  currentStep?: string;
  message?: string;
}
