import { apiClient } from './api';
import type {
  GenerationTask,
  BatchGenerationRequest,
  GenerationProgress,
  ApiResponse,
} from '../types';

export const generationApi = {
  async batchGenerate(request: BatchGenerationRequest): Promise<ApiResponse<GenerationTask[]>> {
    return apiClient.post<GenerationTask[]>('/generation/batch', request);
  },

  async getTaskStatus(taskId: string): Promise<ApiResponse<GenerationProgress>> {
    return apiClient.get<GenerationProgress>(`/generation/task/${taskId}`);
  },

  async cancelTask(taskId: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/generation/task/${taskId}`);
  },

  async getActiveTasks(): Promise<ApiResponse<GenerationTask[]>> {
    return apiClient.get<GenerationTask[]>('/generation/tasks/active');
  },

  async getCompletedTasks(): Promise<ApiResponse<GenerationTask[]>> {
    return apiClient.get<GenerationTask[]>('/generation/tasks/completed');
  },

  async retryTask(taskId: string): Promise<ApiResponse<GenerationTask>> {
    return apiClient.post<GenerationTask>(`/generation/task/${taskId}/retry`);
  },
};
