import type {
  Character,
  CreateCharacterRequest,
  UpdateCharacterRequest,
  GenerateReferenceImageRequest,
  ReferenceImage,
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

class CharacterApiService {
  async getCharactersByNovel(novelId: string): Promise<Character[]> {
    const response = await fetch(`${API_BASE_URL}/novels/${novelId}/characters`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch characters');
    }
    
    const result = await response.json();
    return result.data;
  }

  async getCharacter(id: string): Promise<Character> {
    const response = await fetch(`${API_BASE_URL}/characters/${id}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch character');
    }
    
    const result = await response.json();
    return result.data;
  }

  async createCharacter(data: CreateCharacterRequest): Promise<Character> {
    const response = await fetch(`${API_BASE_URL}/characters`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error('Failed to create character');
    }
    
    const result = await response.json();
    return result.data;
  }

  async updateCharacter(id: string, data: UpdateCharacterRequest): Promise<Character> {
    const response = await fetch(`${API_BASE_URL}/characters/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error('Failed to update character');
    }
    
    const result = await response.json();
    return result.data;
  }

  async deleteCharacter(id: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/characters/${id}`, {
      method: 'DELETE',
    });
    
    if (!response.ok) {
      throw new Error('Failed to delete character');
    }
  }

  async generateReferenceImage(data: GenerateReferenceImageRequest): Promise<ReferenceImage> {
    const response = await fetch(`${API_BASE_URL}/characters/${data.characterId}/reference-images`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        style: data.style,
        customPrompt: data.customPrompt,
      }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to generate reference image');
    }
    
    const result = await response.json();
    return result.data;
  }

  async getReferenceImageStatus(imageId: string): Promise<ReferenceImage> {
    const response = await fetch(`${API_BASE_URL}/reference-images/${imageId}`);
    
    if (!response.ok) {
      throw new Error('Failed to fetch reference image status');
    }
    
    const result = await response.json();
    return result.data;
  }
}

export const characterApi = new CharacterApiService();
