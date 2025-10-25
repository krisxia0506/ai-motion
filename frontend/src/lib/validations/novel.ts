import { z } from 'zod';

export const uploadNovelSchema = z.object({
  title: z.string().min(1, 'Title is required').max(200, 'Title is too long'),
  author: z.string().min(1, 'Author is required').max(100, 'Author name is too long'),
  content: z.string().min(100, 'Content must be at least 100 characters'),
});

export type UploadNovelFormData = z.infer<typeof uploadNovelSchema>;

export const novelFileSchema = z.object({
  file: z
    .instanceof(File)
    .refine((file) => file.size <= 50 * 1024 * 1024, 'File size must be less than 50MB')
    .refine(
      (file) => ['text/plain', 'application/epub+zip'].includes(file.type),
      'Only TXT and EPUB files are supported'
    ),
});

export type NovelFileFormData = z.infer<typeof novelFileSchema>;
