import { apiClient } from './api';
import type {
  Novel,
  UploadNovelRequest,
  ParseNovelRequest,
  ApiResponse,
  PaginatedResponse,
  ListQueryParams,
  Chapter,
} from '../types';

export const novelApi = {
  async uploadNovel(data: UploadNovelRequest): Promise<ApiResponse<Novel>> {
    return apiClient.post<Novel>('/novel/upload', data);
  },

  async uploadNovelFile(file: File): Promise<ApiResponse<Novel>> {
    return apiClient.uploadFile<Novel>('/novel/upload', file);
  },

  async getNovel(id: string): Promise<ApiResponse<Novel>> {
    return apiClient.get<Novel>(`/novel/${id}`);
  },

  async listNovels(params?: ListQueryParams): Promise<ApiResponse<PaginatedResponse<Novel>>> {
    const query = new URLSearchParams();
    if (params?.page) query.append('page', params.page.toString());
    if (params?.pageSize) query.append('pageSize', params.pageSize.toString());
    if (params?.sortBy) query.append('sortBy', params.sortBy);
    if (params?.sortOrder) query.append('sortOrder', params.sortOrder);
    if (params?.search) query.append('search', params.search);

    const endpoint = `/novels${query.toString() ? `?${query}` : ''}`;
    return apiClient.get<PaginatedResponse<Novel>>(endpoint);
  },

  async parseNovel(request: ParseNovelRequest): Promise<ApiResponse<void>> {
    return apiClient.post<void>(`/novel/${request.novelId}/parse`);
  },

  async updateNovel(id: string, data: Partial<Novel>): Promise<ApiResponse<Novel>> {
    return apiClient.put<Novel>(`/novel/${id}`, data);
  },

  async deleteNovel(id: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/novel/${id}`);
  },

  async getChapters(novelId: string): Promise<ApiResponse<Chapter[]>> {
    return apiClient.get<Chapter[]>(`/novel/${novelId}/chapters`);
  },

  async getChapter(novelId: string, chapterId: string): Promise<ApiResponse<Chapter>> {
    return apiClient.get<Chapter>(`/novel/${novelId}/chapters/${chapterId}`);
  },
};
