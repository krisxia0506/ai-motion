export type NovelStatus = 'uploaded' | 'parsing' | 'parsed' | 'generating' | 'completed' | 'failed';

export interface Novel {
  id: string;
  title: string;
  author: string;
  content: string;
  status: NovelStatus;
  characterIds: string[];
  chapterCount: number;
  wordCount: number;
  createdAt: Date | string;
  updatedAt: Date | string;
}

export interface Chapter {
  id: string;
  novelId: string;
  title: string;
  content: string;
  sequence: number;
  sceneCount: number;
  createdAt: Date | string;
}

export interface UploadNovelRequest {
  title: string;
  author: string;
  content: string;
}

export interface UploadNovelResponse {
  novel: Novel;
  message?: string;
}

export interface ParseNovelRequest {
  novelId: string;
}

export interface NovelListQuery {
  page?: number;
  pageSize?: number;
  status?: NovelStatus;
  search?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'title';
  sortOrder?: 'asc' | 'desc';
}
