import type { ApiResponse, ApiError } from '../types';
import { supabase } from '../lib/supabase';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
const API_VERSION = '/api/v1';

export class ApiClient {
  private baseURL: string;
  private version: string;

  constructor(baseURL: string, version: string) {
    this.baseURL = baseURL;
    this.version = version;
  }

  private getFullURL(endpoint: string): string {
    return `${this.baseURL}${this.version}${endpoint}`;
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      // 从 Supabase 获取 Token
      const { data: { session } } = await supabase.auth.getSession();
      const token = session?.access_token;

      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      };

      if (options.headers) {
        Object.entries(options.headers).forEach(([key, value]) => {
          headers[key] = String(value);
        });
      }

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(this.getFullURL(endpoint), {
        ...options,
        headers,
      });

      const data = await response.json();

      if (!response.ok) {
        const error: ApiError = {
          error: data.error || 'Request failed',
          message: data.message || `HTTP error! status: ${response.status}`,
          statusCode: response.status,
          timestamp: new Date().toISOString(),
        };
        throw error;
      }

      return {
        data: data.data || data,
        message: data.message,
        success: true,
      };
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(endpoint: string, body?: unknown): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  async uploadFile<T>(
    endpoint: string,
    file: File,
    additionalData?: Record<string, string>
  ): Promise<ApiResponse<T>> {
    const formData = new FormData();
    formData.append('file', file);

    if (additionalData) {
      Object.entries(additionalData).forEach(([key, value]) => {
        formData.append(key, value);
      });
    }

    // 从 Supabase 获取 Token
    const { data: { session } } = await supabase.auth.getSession();
    const token = session?.access_token;

    const headers: Record<string, string> = {};
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(this.getFullURL(endpoint), {
      method: 'POST',
      headers,
      body: formData,
    });

    const data = await response.json();

    if (!response.ok) {
      const error: ApiError = {
        error: data.error || 'Upload failed',
        message: data.message || `HTTP error! status: ${response.status}`,
        statusCode: response.status,
        timestamp: new Date().toISOString(),
      };
      throw error;
    }

    return {
      data: data.data || data,
      message: data.message,
      success: true,
    };
  }
}

export const apiClient = new ApiClient(API_BASE_URL, API_VERSION);
export default apiClient;
