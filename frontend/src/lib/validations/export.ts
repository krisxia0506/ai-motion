import { z } from 'zod';

export const audioConfigSchema = z.object({
  backgroundMusic: z.string().optional(),
  voiceOver: z.boolean().default(false),
  soundEffects: z.boolean().default(false),
  volume: z.number().min(0).max(100).default(80),
});

export const exportConfigSchema = z.object({
  format: z.enum(['mp4', 'mov', 'avi', 'webm']).default('mp4'),
  quality: z.enum(['720p', '1080p', '4k']).default('1080p'),
  includeSubtitles: z.boolean().default(true),
  includeAudio: z.boolean().default(true),
  audioConfig: audioConfigSchema.optional(),
});

export const exportNovelSchema = z.object({
  novelId: z.string().min(1, 'Novel ID is required'),
  config: exportConfigSchema,
});

export type ExportNovelFormData = z.infer<typeof exportNovelSchema>;
