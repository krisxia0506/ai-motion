export interface Scene {
  id: string;
  novelId: string;
  chapterId: string;
  description: string;
  dialogue: string[];
  characters: string[];
  sequence: number;
  imageUrl?: string;
  videoUrl?: string;
  status: SceneStatus;
  createdAt: Date | string;
  updatedAt: Date | string;
}

export type SceneStatus = 'pending' | 'generating_image' | 'generating_video' | 'completed' | 'failed';

export interface SceneGenerationConfig {
  style?: 'anime' | 'realistic' | 'cartoon';
  aspectRatio?: '16:9' | '9:16' | '1:1' | '4:3';
  quality?: 'low' | 'medium' | 'high';
  useCharacterConsistency?: boolean;
}

export interface GenerateSceneRequest {
  sceneId: string;
  config: SceneGenerationConfig;
}

export interface BatchGenerateRequest {
  sceneIds: string[];
  config: SceneGenerationConfig;
}

export interface GenerationProgress {
  sceneId: string;
  status: SceneStatus;
  progress: number;
  estimatedTimeRemaining?: number;
  error?: string;
}

export interface SceneListQuery {
  novelId?: string;
  chapterId?: string;
  status?: SceneStatus;
  page?: number;
  pageSize?: number;
}
