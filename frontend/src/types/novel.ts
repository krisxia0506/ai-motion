import { Timestamps } from './common';

export type NovelStatus = 'uploaded' | 'parsing' | 'parsed' | 'failed';

export interface Novel extends Timestamps {
  id: string;
  title: string;
  author: string;
  content: string;
  status: NovelStatus;
  characterCount: number;
  sceneCount: number;
  metadata?: NovelMetadata;
}

export interface NovelMetadata {
  fileSize?: number;
  wordCount?: number;
  chapterCount?: number;
  language?: string;
}

export interface Chapter {
  id: string;
  novelId: string;
  title: string;
  content: string;
  orderIndex: number;
  wordCount: number;
}

export interface UploadNovelRequest {
  title: string;
  author: string;
  content: string;
}

export interface ParseNovelRequest {
  novelId: string;
}
