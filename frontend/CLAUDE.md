# Frontend CLAUDE.md

This file provides frontend-specific guidance for Claude Code when working with the React + TypeScript frontend.

## Frontend Overview

The AI-Motion frontend is a React 19 + TypeScript application built with Vite 7, designed to provide a user-friendly interface for novel-to-anime generation.

**Framework:** React 19
**Language:** TypeScript
**Build Tool:** Vite 7
**Styling:** CSS (ready for Tailwind/styled-components)
**Current State:** Basic scaffold in place, ready for feature implementation

## Directory Structure

```
frontend/
├── src/
│   ├── components/          # React components
│   │   ├── common/         # Reusable UI components (Button, Input, Modal, etc.)
│   │   └── features/       # Feature-specific components
│   │       ├── novel/      # Novel upload, list, detail components
│   │       ├── character/  # Character display, editing components
│   │       ├── scene/      # Scene viewer, editor components
│   │       └── generation/ # Generation progress, controls
│   ├── pages/              # Page-level components (route targets)
│   │   ├── HomePage.tsx
│   │   ├── NovelListPage.tsx
│   │   ├── NovelDetailPage.tsx
│   │   ├── CharacterPage.tsx
│   │   ├── SceneGenerationPage.tsx
│   │   └── ExportPage.tsx
│   ├── services/           # API client services
│   │   ├── api.ts         # Base API configuration (axios/fetch)
│   │   ├── novelApi.ts    # Novel-related API calls
│   │   ├── characterApi.ts
│   │   ├── sceneApi.ts
│   │   └── generationApi.ts
│   ├── hooks/              # Custom React hooks
│   │   ├── useNovel.ts
│   │   ├── useCharacters.ts
│   │   ├── useGeneration.ts
│   │   └── useWebSocket.ts  # For real-time generation updates
│   ├── store/              # State management (Context API / Zustand / Redux)
│   │   ├── novelStore.ts
│   │   ├── characterStore.ts
│   │   └── uiStore.ts
│   ├── types/              # TypeScript type definitions
│   │   ├── novel.ts
│   │   ├── character.ts
│   │   ├── scene.ts
│   │   ├── media.ts
│   │   └── api.ts
│   ├── utils/              # Utility functions
│   │   ├── formatters.ts  # Date, number formatting
│   │   ├── validators.ts  # Form validation
│   │   └── constants.ts   # App constants
│   ├── styles/             # Global styles
│   │   ├── variables.css  # CSS variables (colors, spacing)
│   │   └── globals.css    # Global styles
│   ├── assets/             # Static assets (images, icons)
│   ├── App.tsx             # Root component with routing
│   └── main.tsx            # Application entry point
├── public/                 # Public static files
├── index.html              # HTML template
├── vite.config.ts          # Vite configuration
├── tsconfig.json           # TypeScript configuration
├── tsconfig.app.json       # App-specific TS config
├── tsconfig.node.json      # Node-specific TS config
├── eslint.config.js        # ESLint configuration
└── package.json            # Dependencies and scripts
```

## Development Commands

### Running the Frontend

```bash
# Development mode with hot reload (from frontend directory)
npm run dev
# Access at http://localhost:5173

# Or from project root
cd frontend && npm run dev

# Build for production
npm run build
# Output to dist/ directory

# Preview production build
npm run preview
```

### Dependencies

```bash
# Install all dependencies
npm install

# Add new dependency
npm install package-name

# Add dev dependency
npm install -D package-name

# Update dependencies
npm update

# Check for outdated packages
npm outdated

# Remove unused dependencies
npm prune
```

### Code Quality

```bash
# Lint code
npm run lint

# Lint and fix
npm run lint -- --fix

# Type check
npx tsc --noEmit

# Format with Prettier (if configured)
npx prettier --write "src/**/*.{ts,tsx}"
```

### Testing

```bash
# Run tests (when configured)
npm test

# Run tests in watch mode
npm test -- --watch

# Generate coverage report
npm test -- --coverage
```

## TypeScript Guidelines

### Type Definitions

Always define explicit types for props, state, and API responses:

```typescript
// src/types/novel.ts
export interface Novel {
  id: string;
  title: string;
  author: string;
  content: string;
  status: NovelStatus;
  characterIds: string[];
  createdAt: Date;
  updatedAt: Date;
}

export type NovelStatus = 'uploaded' | 'parsing' | 'completed' | 'failed';

export interface UploadNovelRequest {
  title: string;
  author: string;
  content: string;
}

export interface NovelResponse {
  data: Novel;
  message?: string;
}

// src/types/character.ts
export interface Character {
  id: string;
  novelId: string;
  name: string;
  appearance: string;
  personality: string;
  referenceImages: string[];
  createdAt: Date;
}

// src/types/api.ts
export interface ApiResponse<T> {
  data: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}
```

### Component Props

```typescript
// Use interface for props
interface NovelCardProps {
  novel: Novel;
  onSelect: (id: string) => void;
  onDelete?: (id: string) => void;
  className?: string;
}

// Functional component with explicit return type
export const NovelCard: React.FC<NovelCardProps> = ({
  novel,
  onSelect,
  onDelete,
  className = '',
}) => {
  return (
    <div className={`novel-card ${className}`}>
      <h3>{novel.title}</h3>
      <p>by {novel.author}</p>
      <button onClick={() => onSelect(novel.id)}>View</button>
      {onDelete && (
        <button onClick={() => onDelete(novel.id)}>Delete</button>
      )}
    </div>
  );
};
```

### Hooks Typing

```typescript
// Custom hook with explicit types
import { useState, useEffect } from 'react';

interface UseNovelResult {
  novel: Novel | null;
  loading: boolean;
  error: Error | null;
  refetch: () => void;
}

export const useNovel = (id: string): UseNovelResult => {
  const [novel, setNovel] = useState<Novel | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  const fetchNovel = async () => {
    try {
      setLoading(true);
      const data = await novelApi.getNovel(id);
      setNovel(data);
      setError(null);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNovel();
  }, [id]);

  return { novel, loading, error, refetch: fetchNovel };
};
```

## Component Patterns

### Page Component Example

```typescript
// src/pages/NovelDetailPage.tsx
import React from 'react';
import { useParams } from 'react-router-dom';
import { useNovel } from '../hooks/useNovel';
import { NovelDetail } from '../components/features/novel/NovelDetail';
import { CharacterList } from '../components/features/character/CharacterList';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';

export const NovelDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { novel, loading, error } = useNovel(id!);

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;
  if (!novel) return <div>Novel not found</div>;

  return (
    <div className="novel-detail-page">
      <NovelDetail novel={novel} />
      <CharacterList novelId={novel.id} />
    </div>
  );
};
```

### Feature Component Example

```typescript
// src/components/features/novel/NovelUpload.tsx
import React, { useState } from 'react';
import { novelApi } from '../../../services/novelApi';
import { UploadNovelRequest } from '../../../types/novel';

interface NovelUploadProps {
  onSuccess: (novelId: string) => void;
}

export const NovelUpload: React.FC<NovelUploadProps> = ({ onSuccess }) => {
  const [formData, setFormData] = useState<UploadNovelRequest>({
    title: '',
    author: '',
    content: '',
  });
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setUploading(true);
      setError(null);

      const response = await novelApi.uploadNovel(formData);
      onSuccess(response.data.id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setUploading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <form onSubmit={handleSubmit} className="novel-upload">
      <input
        type="text"
        name="title"
        value={formData.title}
        onChange={handleChange}
        placeholder="Novel Title"
        required
      />
      <input
        type="text"
        name="author"
        value={formData.author}
        onChange={handleChange}
        placeholder="Author"
      />
      <textarea
        name="content"
        value={formData.content}
        onChange={handleChange}
        placeholder="Paste novel content here..."
        rows={20}
        required
      />
      {error && <div className="error">{error}</div>}
      <button type="submit" disabled={uploading}>
        {uploading ? 'Uploading...' : 'Upload Novel'}
      </button>
    </form>
  );
};
```

### Common Component Example

```typescript
// src/components/common/Button.tsx
import React from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'small' | 'medium' | 'large';
  loading?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'medium',
  loading = false,
  disabled,
  className = '',
  ...props
}) => {
  const classes = `btn btn-${variant} btn-${size} ${className}`;

  return (
    <button
      className={classes}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? 'Loading...' : children}
    </button>
  );
};
```

## API Service Pattern

### Base API Configuration

```typescript
// src/services/api.ts
import axios, { AxiosInstance, AxiosError } from 'axios';

const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const apiClient: AxiosInstance = axios.create({
  baseURL: `${BASE_URL}/api/v1`,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor (for auth tokens, etc.)
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor (for error handling)
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Handle unauthorized (e.g., redirect to login)
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

### Feature API Service

```typescript
// src/services/novelApi.ts
import apiClient from './api';
import { Novel, UploadNovelRequest, ApiResponse } from '../types';

export const novelApi = {
  // Upload novel
  async uploadNovel(data: UploadNovelRequest): Promise<ApiResponse<Novel>> {
    const response = await apiClient.post<ApiResponse<Novel>>('/novel/upload', data);
    return response.data;
  },

  // Get novel by ID
  async getNovel(id: string): Promise<Novel> {
    const response = await apiClient.get<ApiResponse<Novel>>(`/novel/${id}`);
    return response.data.data;
  },

  // List novels
  async listNovels(page = 1, pageSize = 20): Promise<Novel[]> {
    const response = await apiClient.get<ApiResponse<Novel[]>>('/novels', {
      params: { page, pageSize },
    });
    return response.data.data;
  },

  // Parse novel
  async parseNovel(id: string): Promise<void> {
    await apiClient.post(`/novel/${id}/parse`);
  },

  // Delete novel
  async deleteNovel(id: string): Promise<void> {
    await apiClient.delete(`/novel/${id}`);
  },
};
```

```typescript
// src/services/characterApi.ts
import apiClient from './api';
import { Character, ApiResponse } from '../types';

export const characterApi = {
  async getCharactersByNovel(novelId: string): Promise<Character[]> {
    const response = await apiClient.get<ApiResponse<Character[]>>(
      `/characters/${novelId}`
    );
    return response.data.data;
  },

  async generateReferenceImage(characterId: string): Promise<string> {
    const response = await apiClient.post<ApiResponse<{ imageUrl: string }>>(
      `/character/${characterId}/generate-reference`
    );
    return response.data.data.imageUrl;
  },
};
```

## State Management

### Context API Pattern (Simple State)

```typescript
// src/store/NovelContext.tsx
import React, { createContext, useContext, useState, ReactNode } from 'react';
import { Novel } from '../types/novel';

interface NovelContextValue {
  selectedNovel: Novel | null;
  setSelectedNovel: (novel: Novel | null) => void;
}

const NovelContext = createContext<NovelContextValue | undefined>(undefined);

export const NovelProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [selectedNovel, setSelectedNovel] = useState<Novel | null>(null);

  return (
    <NovelContext.Provider value={{ selectedNovel, setSelectedNovel }}>
      {children}
    </NovelContext.Provider>
  );
};

export const useNovelContext = () => {
  const context = useContext(NovelContext);
  if (!context) {
    throw new Error('useNovelContext must be used within NovelProvider');
  }
  return context;
};
```

### Zustand Pattern (Recommended for Complex State)

```typescript
// src/store/novelStore.ts
import { create } from 'zustand';
import { Novel } from '../types/novel';

interface NovelState {
  novels: Novel[];
  selectedNovel: Novel | null;
  loading: boolean;
  error: string | null;

  setNovels: (novels: Novel[]) => void;
  setSelectedNovel: (novel: Novel | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  addNovel: (novel: Novel) => void;
  removeNovel: (id: string) => void;
}

export const useNovelStore = create<NovelState>((set) => ({
  novels: [],
  selectedNovel: null,
  loading: false,
  error: null,

  setNovels: (novels) => set({ novels }),
  setSelectedNovel: (novel) => set({ selectedNovel: novel }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  addNovel: (novel) => set((state) => ({ novels: [...state.novels, novel] })),
  removeNovel: (id) => set((state) => ({
    novels: state.novels.filter((n) => n.id !== id),
  })),
}));
```

## Routing Setup

```typescript
// src/App.tsx
import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { HomePage } from './pages/HomePage';
import { NovelListPage } from './pages/NovelListPage';
import { NovelDetailPage } from './pages/NovelDetailPage';
import { CharacterPage } from './pages/CharacterPage';
import { SceneGenerationPage } from './pages/SceneGenerationPage';
import { ExportPage } from './pages/ExportPage';
import { Layout } from './components/Layout';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<HomePage />} />
          <Route path="novels" element={<NovelListPage />} />
          <Route path="novels/:id" element={<NovelDetailPage />} />
          <Route path="novels/:id/characters" element={<CharacterPage />} />
          <Route path="novels/:id/generate" element={<SceneGenerationPage />} />
          <Route path="novels/:id/export" element={<ExportPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
```

## Environment Variables

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080
VITE_APP_TITLE=AI-Motion
VITE_ENABLE_DEBUG=true

# .env.production
VITE_API_BASE_URL=https://api.ai-motion.com
VITE_APP_TITLE=AI-Motion
VITE_ENABLE_DEBUG=false
```

Access in code:
```typescript
const apiUrl = import.meta.env.VITE_API_BASE_URL;
const isDebug = import.meta.env.VITE_ENABLE_DEBUG === 'true';
```

## Styling Guidelines

### CSS Modules (Recommended)

```typescript
// NovelCard.module.css
.card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 16px;
}

.title {
  font-size: 20px;
  font-weight: bold;
}

// NovelCard.tsx
import styles from './NovelCard.module.css';

export const NovelCard: React.FC<NovelCardProps> = ({ novel }) => {
  return (
    <div className={styles.card}>
      <h3 className={styles.title}>{novel.title}</h3>
    </div>
  );
};
```

### CSS Variables

```css
/* src/styles/variables.css */
:root {
  --color-primary: #3b82f6;
  --color-secondary: #10b981;
  --color-danger: #ef4444;
  --color-text: #1f2937;
  --color-bg: #ffffff;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --border-radius: 8px;
}
```

## Performance Optimization

### Code Splitting

```typescript
// Lazy load pages
import { lazy, Suspense } from 'react';

const NovelDetailPage = lazy(() => import('./pages/NovelDetailPage'));
const CharacterPage = lazy(() => import('./pages/CharacterPage'));

// Use with Suspense
<Suspense fallback={<LoadingSpinner />}>
  <NovelDetailPage />
</Suspense>
```

### Memoization

```typescript
import { memo, useMemo, useCallback } from 'react';

// Memo component
export const NovelCard = memo<NovelCardProps>(({ novel, onSelect }) => {
  return <div onClick={() => onSelect(novel.id)}>{novel.title}</div>;
});

// Memo value
const expensiveValue = useMemo(() => {
  return computeExpensiveValue(data);
}, [data]);

// Memo callback
const handleClick = useCallback(() => {
  doSomething(id);
}, [id]);
```

### Virtual Lists

```typescript
// For long lists, use react-window or react-virtual
import { FixedSizeList } from 'react-window';

const NovelList: React.FC<{ novels: Novel[] }> = ({ novels }) => {
  return (
    <FixedSizeList
      height={600}
      itemCount={novels.length}
      itemSize={80}
      width="100%"
    >
      {({ index, style }) => (
        <div style={style}>
          <NovelCard novel={novels[index]} />
        </div>
      )}
    </FixedSizeList>
  );
};
```

## Testing Patterns

### Component Testing

```typescript
// NovelCard.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { NovelCard } from './NovelCard';

describe('NovelCard', () => {
  const mockNovel = {
    id: '1',
    title: 'Test Novel',
    author: 'Test Author',
    status: 'completed' as const,
  };

  it('renders novel title', () => {
    render(<NovelCard novel={mockNovel} onSelect={() => {}} />);
    expect(screen.getByText('Test Novel')).toBeInTheDocument();
  });

  it('calls onSelect when clicked', () => {
    const handleSelect = jest.fn();
    render(<NovelCard novel={mockNovel} onSelect={handleSelect} />);

    fireEvent.click(screen.getByText('Test Novel'));
    expect(handleSelect).toHaveBeenCalledWith('1');
  });
});
```

### Hook Testing

```typescript
// useNovel.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { useNovel } from './useNovel';

jest.mock('../services/novelApi');

describe('useNovel', () => {
  it('fetches novel data', async () => {
    const { result } = renderHook(() => useNovel('1'));

    expect(result.current.loading).toBe(true);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
      expect(result.current.novel).toBeDefined();
    });
  });
});
```

## Common Patterns

### Error Boundary

```typescript
// src/components/common/ErrorBoundary.tsx
import React, { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>Something went wrong</h2>
          <p>{this.state.error?.message}</p>
        </div>
      );
    }

    return this.props.children;
  }
}
```

### Loading States

```typescript
// src/components/common/LoadingSpinner.tsx
export const LoadingSpinner: React.FC = () => (
  <div className="loading-spinner">
    <div className="spinner"></div>
    <p>Loading...</p>
  </div>
);

// Usage with conditional rendering
{loading && <LoadingSpinner />}
{!loading && data && <DataComponent data={data} />}
```

## Build Configuration

### Vite Config Customization

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@services': path.resolve(__dirname, './src/services'),
      '@types': path.resolve(__dirname, './src/types'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          vendor: ['react', 'react-dom', 'react-router-dom'],
          api: ['axios'],
        },
      },
    },
  },
});
```

## Common Pitfalls to Avoid

1. **Don't** mutate state directly - Use setState or immutable updates
2. **Don't** forget to handle loading and error states
3. **Don't** make API calls without error handling
4. **Don't** use `any` type - Define proper TypeScript types
5. **Don't** forget to cleanup effects (useEffect return function)
6. **Don't** create components in render functions
7. **Don't** use index as key for dynamic lists
8. **Don't** store derived state - Use useMemo instead

## Debugging Tips

```typescript
// React DevTools
// Install browser extension for component inspection

// Console logging with type safety
console.log('Novel data:', novel);
console.table(novels);

// Debug renders
useEffect(() => {
  console.log('Component rendered with:', props);
});

// Network debugging
// Check browser Network tab for API calls
// Use axios interceptors for request/response logging
```

## Useful Packages to Consider

```bash
# UI Component Libraries
npm install @headlessui/react @heroicons/react  # Headless UI
npm install @radix-ui/react-dialog             # Radix UI
npm install react-hot-toast                     # Toast notifications

# Form Handling
npm install react-hook-form zod @hookform/resolvers  # Form + validation

# State Management
npm install zustand                             # Simple state management

# Date/Time
npm install date-fns                           # Date utilities

# Routing
npm install react-router-dom                   # Already included

# HTTP Client
npm install axios                              # Already in use

# Icons
npm install react-icons                        # Icon library
```

## Quick Reference

**API Base URL:** `http://localhost:8080` (dev), configured in `.env`
**Dev Server:** `http://localhost:5173`
**Backend API Docs:** See [/docs/API.md](/docs/API.md)
**Type Definitions:** See `src/types/` directory
**Main App:** See [src/App.tsx](src/App.tsx)
