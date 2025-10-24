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
