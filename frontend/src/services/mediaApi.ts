import { apiClient } from './api';
import {
  Media,
  ExportRequest,
  ExportTask,
  ApiResponse,
} from '../types';

export const mediaApi = {
  async getMedia(id: string): Promise<ApiResponse<Media>> {
    return apiClient.get<Media>(`/media/${id}`);
  },

  async deleteMedia(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/media/${id}`);
  },

  async getMediaByScene(sceneId: string): Promise<ApiResponse<Media[]>> {
    return apiClient.get<Media[]>(`/media/scene/${sceneId}`);
  },

  async exportNovel(request: ExportRequest): Promise<ApiResponse<ExportTask>> {
    return apiClient.post<ExportTask>('/export', request);
  },

  async getExportStatus(taskId: string): Promise<ApiResponse<ExportTask>> {
    return apiClient.get<ExportTask>(`/export/task/${taskId}`);
  },

  async cancelExport(taskId: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/export/task/${taskId}`);
  },

  async downloadExport(taskId: string): Promise<Blob> {
    const response = await fetch(
      `${import.meta.env.VITE_API_BASE_URL}/api/v1/export/task/${taskId}/download`
    );
    if (!response.ok) {
      throw new Error('Download failed');
    }
    return response.blob();
  },
};
