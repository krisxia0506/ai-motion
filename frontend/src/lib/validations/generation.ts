import { z } from 'zod';

export const generationConfigSchema = z.object({
  style: z.string().max(100, 'Style is too long').optional(),
  quality: z.enum(['low', 'medium', 'high']).default('medium'),
  aspectRatio: z.string().max(20).optional(),
  duration: z.number().min(1).max(60).optional(),
});

export const generateSceneSchema = z.object({
  sceneId: z.string().min(1, 'Scene ID is required'),
  type: z.enum(['image', 'video']),
  config: generationConfigSchema.optional(),
});

export type GenerateSceneFormData = z.infer<typeof generateSceneSchema>;

export const batchGenerationSchema = z.object({
  sceneIds: z.array(z.string()).min(1, 'At least one scene must be selected'),
  type: z.enum(['image', 'video', 'audio']),
  config: generationConfigSchema.optional(),
});

export type BatchGenerationFormData = z.infer<typeof batchGenerationSchema>;
