export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

export interface AsyncState<T = any> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  status: LoadingState;
}

export interface FormFieldError {
  field: string;
  message: string;
}

export interface ValidationError {
  errors: FormFieldError[];
  message: string;
}

export type Nullable<T> = T | null;
export type Optional<T> = T | undefined;
export type Maybe<T> = T | null | undefined;

export interface SelectOption<T = string> {
  value: T;
  label: string;
  disabled?: boolean;
}

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title?: string;
  message: string;
  duration?: number;
  dismissible?: boolean;
}

export interface ModalProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  size?: 'small' | 'medium' | 'large' | 'full';
}

export interface Breadcrumb {
  label: string;
  path: string;
  active?: boolean;
}

export type Theme = 'light' | 'dark' | 'system';

export interface UserPreferences {
  theme: Theme;
  language: string;
  defaultQuality: string;
  autoSave: boolean;
}
