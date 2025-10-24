import type { Novel, UploadNovelRequest } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

class NovelApiService {
  async uploadNovel(data: UploadNovelRequest): Promise<Novel> {
    const response = await fetch(`${API_BASE_URL}/novels`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error('Failed to upload novel');
    }
    
    const result = await response.json();
    return result.data;
  }

  async getNovel(id: string): Promise<Novel> {
    const response = await fetch(`${API_BASE_URL}/novels/${id}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch novel');
    }
    
    const result = await response.json();
    return result.data;
  }

  async listNovels(page = 1, size = 20): Promise<Novel[]> {
    const response = await fetch(`${API_BASE_URL}/novels?page=${page}&size=${size}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch novels');
    }
    
    const result = await response.json();
    return result.data;
  }

  async deleteNovel(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/novels/${id}`, {
      method: 'DELETE',
    });
    
    if (!response.ok) {
      throw new Error('Failed to delete novel');
    }
  }
}

export const novelApi = new NovelApiService();
