export type Status = 'pending' | 'processing' | 'completed' | 'failed';

export interface Timestamps {
  createdAt: string;
  updatedAt: string;
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
}

export interface SelectOption {
  label: string;
  value: string;
}
