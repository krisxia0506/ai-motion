import { create } from 'zustand';
import type { Character } from '../types';

interface CharacterState {
  characters: Character[];
  selectedCharacter: Character | null;
  loading: boolean;
  error: string | null;
  
  setCharacters: (characters: Character[]) => void;
  addCharacter: (character: Character) => void;
  updateCharacter: (id: string, updates: Partial<Character>) => void;
  removeCharacter: (id: string) => void;
  setSelectedCharacter: (character: Character | null) => void;
  getCharactersByNovel: (novelId: string) => Character[];
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useCharacterStore = create<CharacterState>((set, get) => ({
  characters: [],
  selectedCharacter: null,
  loading: false,
  error: null,

  setCharacters: (characters) => set({ characters }),
  
  addCharacter: (character) =>
    set((state) => ({ characters: [...state.characters, character] })),
  
  updateCharacter: (id, updates) =>
    set((state) => ({
      characters: state.characters.map((char) =>
        char.id === id ? { ...char, ...updates } : char
      ),
      selectedCharacter:
        state.selectedCharacter?.id === id
          ? { ...state.selectedCharacter, ...updates }
          : state.selectedCharacter,
    })),
  
  removeCharacter: (id) =>
    set((state) => ({
      characters: state.characters.filter((char) => char.id !== id),
      selectedCharacter:
        state.selectedCharacter?.id === id ? null : state.selectedCharacter,
    })),
  
  setSelectedCharacter: (character) => set({ selectedCharacter: character }),
  
  getCharactersByNovel: (novelId) => {
    const { characters } = get();
    return characters.filter((char) => char.novelId === novelId);
  },
  
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  clearError: () => set({ error: null }),
}));
