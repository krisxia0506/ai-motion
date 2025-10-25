import type { Timestamps } from './common';

export type CharacterType = 'main' | 'supporting' | 'minor';

export interface Character extends Timestamps {
  id: string;
  novelId: string;
  name: string;
  type: CharacterType;
  appearance: string;
  personality: string;
  background?: string;
  relationships?: string[];
  referenceImages: ReferenceImage[];
}

export interface ReferenceImage {
  id: string;
  characterId: string;
  url: string;
  prompt: string;
  modelUsed: string;
  status: 'generating' | 'completed' | 'failed';
  generatedAt: string;
}

export interface GenerateReferenceRequest {
  characterId: string;
  prompt?: string;
  style?: string;
}

export interface CreateCharacterRequest {
  novelId: string;
  name: string;
  type: CharacterType;
  appearance: string;
  personality: string;
  background?: string;
}

export interface UpdateCharacterRequest {
  name?: string;
  appearance?: string;
  personality?: string;
  background?: string;
}
