export interface Novel {
  id: string;
  title: string;
  author: string;
  content: string;
  status: 'uploaded' | 'parsing' | 'completed' | 'failed';
  createdAt: string;
  updatedAt: string;
  metadata?: {
    characterCount?: number;
    chapterCount?: number;
    sceneCount?: number;
    wordCount?: number;
  };
  chapters?: Chapter[];
}

export interface Chapter {
  id: string;
  number: number;
  title: string;
  content: string;
  wordCount: number;
}

export interface UploadNovelRequest {
  title: string;
  author: string;
  content: string;
}

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface Character {
  id: string;
  novelId: string;
  name: string;
  type: 'main' | 'supporting' | 'minor';
  appearance: string;
  personality: string;
  referenceImages?: ReferenceImage[];
  createdAt: string;
  updatedAt: string;
}

export interface ReferenceImage {
  id: string;
  characterId: string;
  url: string;
  status: 'generating' | 'completed' | 'failed';
  prompt?: string;
  style?: string;
  createdAt: string;
}

export interface CreateCharacterRequest {
  novelId: string;
  name: string;
  type: 'main' | 'supporting' | 'minor';
  appearance: string;
  personality: string;
}

export interface UpdateCharacterRequest {
  name?: string;
  type?: 'main' | 'supporting' | 'minor';
  appearance?: string;
  personality?: string;
}

export interface GenerateReferenceImageRequest {
  characterId: string;
  style?: 'anime' | 'realistic' | 'cartoon' | 'semi-realistic';
  customPrompt?: string;
}
