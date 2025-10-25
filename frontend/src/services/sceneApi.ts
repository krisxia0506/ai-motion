import { apiClient } from './api';
import {
  Scene,
  GenerateSceneRequest,
  GeneratedMedia,
  ApiResponse,
} from '../types';

export const sceneApi = {
  async getScenesByNovel(novelId: string): Promise<ApiResponse<Scene[]>> {
    return apiClient.get<Scene[]>(`/scenes/${novelId}`);
  },

  async getScenesByChapter(
    novelId: string,
    chapterId: string
  ): Promise<ApiResponse<Scene[]>> {
    return apiClient.get<Scene[]>(`/scenes/${novelId}/chapters/${chapterId}`);
  },

  async getScene(id: string): Promise<ApiResponse<Scene>> {
    return apiClient.get<Scene>(`/scene/${id}`);
  },

  async updateScene(id: string, data: Partial<Scene>): Promise<ApiResponse<Scene>> {
    return apiClient.put<Scene>(`/scene/${id}`, data);
  },

  async deleteScene(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/scene/${id}`);
  },

  async generateScene(request: GenerateSceneRequest): Promise<ApiResponse<GeneratedMedia>> {
    return apiClient.post<GeneratedMedia>(`/scene/${request.sceneId}/generate`, {
      type: request.type,
      config: request.config,
    });
  },

  async getGeneratedMedia(sceneId: string): Promise<ApiResponse<GeneratedMedia[]>> {
    return apiClient.get<GeneratedMedia[]>(`/scene/${sceneId}/media`);
  },
};
