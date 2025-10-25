import { apiClient } from './api';
import type {
  Character,
  CreateCharacterRequest,
  GenerateReferenceRequest,
  UpdateCharacterRequest,
  ReferenceImage,
  ApiResponse,
} from '../types';

export const characterApi = {
  async getCharactersByNovel(novelId: string): Promise<ApiResponse<Character[]>> {
    return apiClient.get<Character[]>(`/characters/${novelId}`);
  },

  async getCharacter(id: string): Promise<ApiResponse<Character>> {
    return apiClient.get<Character>(`/character/${id}`);
  },

  async createCharacter(data: CreateCharacterRequest): Promise<ApiResponse<Character>> {
    return apiClient.post<Character>(`/character`, data);
  },

  async updateCharacter(
    id: string,
    data: UpdateCharacterRequest
  ): Promise<ApiResponse<Character>> {
    return apiClient.put<Character>(`/character/${id}`, data);
  },

  async deleteCharacter(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/character/${id}`);
  },

  async generateReferenceImage(
    request: GenerateReferenceRequest
  ): Promise<ApiResponse<ReferenceImage>> {
    return apiClient.post<ReferenceImage>(
      `/character/${request.characterId}/generate-reference`,
      { prompt: request.prompt, style: request.style }
    );
  },

  async getReferenceImages(characterId: string): Promise<ApiResponse<ReferenceImage[]>> {
    return apiClient.get<ReferenceImage[]>(`/character/${characterId}/references`);
  },

  async deleteReferenceImage(
    characterId: string,
    imageId: string
  ): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/character/${characterId}/references/${imageId}`);
  },
};
