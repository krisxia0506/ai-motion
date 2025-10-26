import { useState, useEffect, useCallback } from 'react';
import type { Character, ApiError } from '../types';
import { characterApi } from '../services';
import { useCharacterStore } from '../store';

interface UseCharactersResult {
  characters: Character[];
  loading: boolean;
  error: ApiError | null;
  refetch: () => Promise<void>;
}

export const useCharacters = (novelId: string): UseCharactersResult => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<ApiError | null>(null);
  const { characters, setCharacters, setError: setStoreError } = useCharacterStore();

  const fetchCharacters = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await characterApi.getCharactersByNovel(novelId);
      setCharacters(response.data);
    } catch (err) {
      const apiError = err as ApiError;
      setError(apiError);
      setStoreError(apiError.message);
    } finally {
      setLoading(false);
    }
  }, [novelId, setCharacters, setStoreError]);

  useEffect(() => {
    if (novelId) {
      fetchCharacters();
    }
  }, [novelId, fetchCharacters]);

  return {
    characters: characters.filter((c: Character) => c.novelId === novelId),
    loading,
    error,
    refetch: fetchCharacters,
  };
};
