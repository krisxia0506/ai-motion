import { z } from 'zod';

export const updateCharacterSchema = z.object({
  name: z.string().min(1, 'Name is required').max(100, 'Name is too long').optional(),
  appearance: z.string().min(10, 'Appearance description is too short').max(500, 'Appearance description is too long').optional(),
  personality: z.string().min(10, 'Personality description is too short').max(500, 'Personality description is too long').optional(),
  background: z.string().max(1000, 'Background is too long').optional(),
});

export type UpdateCharacterFormData = z.infer<typeof updateCharacterSchema>;

export const generateReferenceSchema = z.object({
  characterId: z.string().min(1, 'Character ID is required'),
  prompt: z.string().max(500, 'Prompt is too long').optional(),
  style: z.string().max(100, 'Style is too long').optional(),
});

export type GenerateReferenceFormData = z.infer<typeof generateReferenceSchema>;
