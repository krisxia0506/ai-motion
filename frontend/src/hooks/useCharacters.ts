import { useState, useEffect } from 'react';
import { Character, ApiError } from '../types';
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

  const fetchCharacters = async () => {
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
  };

  useEffect(() => {
    if (novelId) {
      fetchCharacters();
    }
  }, [novelId]);

  return {
    characters: characters.filter((c) => c.novelId === novelId),
    loading,
    error,
    refetch: fetchCharacters,
  };
};
