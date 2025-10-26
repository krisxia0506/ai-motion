import { apiClient } from './api';

export interface CreateTaskRequest {
  title: string;
  author: string;
  content: string;
}

export type TaskStatusType = 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';

export interface TaskStatus {
  task_id: string;
  status: TaskStatusType;
  progress: {
    current_step: string;
    current_step_index: number;
    total_steps: number;
    percentage: number;
    details?: {
      characters_extracted: number;
      characters_generated: number;
      scenes_divided: number;
      scenes_generated: number;
    };
  };
  result?: {
    novel_id: string;
    title: string;
    character_count: number;
    scene_count: number;
    characters: Array<{
      id: string;
      name: string;
      reference_image_url: string;
    }>;
    scenes: Array<{
      id: string;
      sequence_num: number;
      description: string;
      image_url: string;
    }>;
  };
  error?: {
    code: number;
    message: string;
    retry_able: boolean;
  };
  created_at: string;
  updated_at: string;
  completed_at?: string;
  failed_at?: string;
}

export interface TaskListItem extends TaskStatus {
  title?: string;
  character_count?: number;
  scene_count?: number;
}

export const taskApi = {
  // 创建漫画生成任务
  async createTask(request: CreateTaskRequest) {
    return apiClient.post<{ task_id: string; status: string; created_at: string }>(
      '/manga/generate',
      request
    );
  },

  // 获取任务状态
  async getTaskStatus(taskId: string) {
    return apiClient.get<TaskStatus>(`/manga/task/${taskId}`);
  },

  // 获取任务列表
  async getTaskList(params?: { page?: number; page_size?: number; status?: string }) {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.page_size) queryParams.append('page_size', params.page_size.toString());
    if (params?.status) queryParams.append('status', params.status);

    const query = queryParams.toString();
    return apiClient.get<{
      items: Array<TaskStatus>;
      pagination: {
        page: number;
        page_size: number;
        total: number;
        total_pages: number;
        has_next: boolean;
        has_prev: boolean;
      };
    }>(`/manga/tasks${query ? '?' + query : ''}`);
  },

  // 取消任务
  async cancelTask(taskId: string) {
    return apiClient.post<{ task_id: string; status: string }>(
      `/manga/task/${taskId}/cancel`
    );
  },
};

export default taskApi;
