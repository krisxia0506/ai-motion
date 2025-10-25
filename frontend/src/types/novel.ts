import type { Timestamps } from './common';

export type NovelStatus = 'uploaded' | 'parsing' | 'parsed' | 'generating' | 'failed';

export interface Novel extends Timestamps {
  id: string;
  title: string;
  author: string;
  content: string;
  status: NovelStatus;
  characterCount: number;
  sceneCount: number;
  wordCount?: number;
  chapterCount?: number;
  metadata?: NovelMetadata;
}

export interface NovelMetadata {
  fileSize?: number;
  wordCount?: number;
  chapterCount?: number;
  characterCount?: number;
  language?: string;
}

export interface Chapter {
  id: string;
  novelId: string;
  title: string;
  content: string;
  orderIndex: number;
  wordCount: number;
  sceneCount?: number;
}

export interface UploadNovelRequest {
  title: string;
  author: string;
  content: string;
}

export interface ParseNovelRequest {
  novelId: string;
}
