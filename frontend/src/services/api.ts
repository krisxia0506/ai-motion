const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

interface ApiResponse<T = unknown> {
  data?: T;
  error?: string;
  message?: string;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${this.baseURL}${endpoint}`, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return { data };
    } catch (error) {
      console.error('API request failed:', error);
      return { error: error instanceof Error ? error.message : 'Unknown error' };
    }
  }

  // 小说相关 API
  async uploadNovel(file: File): Promise<ApiResponse> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${this.baseURL}/api/v1/novel/upload`, {
      method: 'POST',
      body: formData,
    });

    return response.json();
  }

  async parseNovel(novelId: string): Promise<ApiResponse> {
    return this.request(`/api/v1/novel/${novelId}/parse`, {
      method: 'POST',
    });
  }

  // 角色相关 API
  async getCharacters(novelId: string): Promise<ApiResponse> {
    return this.request(`/api/v1/characters/${novelId}`);
  }

  // 生成相关 API
  async generateScene(sceneData: Record<string, unknown>): Promise<ApiResponse> {
    return this.request('/api/v1/generate/scene', {
      method: 'POST',
      body: JSON.stringify(sceneData),
    });
  }

  async generateVoice(voiceData: Record<string, unknown>): Promise<ApiResponse> {
    return this.request('/api/v1/generate/voice', {
      method: 'POST',
      body: JSON.stringify(voiceData),
    });
  }

  // 导出相关 API
  async exportAnime(animeId: string): Promise<ApiResponse> {
    return this.request(`/api/v1/anime/${animeId}/export`, {
      method: 'POST',
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
export default apiClient;
