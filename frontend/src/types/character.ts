export interface Character {
  id: string;
  novelId: string;
  name: string;
  appearance: string;
  personality: string;
  role: CharacterRole;
  referenceImages: string[];
  hasReferenceImage: boolean;
  firstAppearanceChapter?: number;
  createdAt: Date | string;
  updatedAt: Date | string;
}

export type CharacterRole = 'protagonist' | 'antagonist' | 'supporting' | 'minor';

export interface CharacterUpdateRequest {
  name?: string;
  appearance?: string;
  personality?: string;
  role?: CharacterRole;
}

export interface GenerateReferenceImageRequest {
  characterId: string;
  style?: 'anime' | 'realistic' | 'cartoon';
  seed?: number;
}

export interface GenerateReferenceImageResponse {
  imageUrl: string;
  characterId: string;
  message?: string;
}

export interface CharacterListQuery {
  novelId: string;
  role?: CharacterRole;
  hasReferenceImage?: boolean;
  search?: string;
}
