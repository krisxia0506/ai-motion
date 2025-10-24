import { useEffect } from 'react';
import { useCharacterStore } from '../store/characterStore';
import { characterApi } from '../services/characterApi';

export const useCharacters = (novelId: string) => {
  const { characters, loading, error, setCharacters, setLoading, setError } = useCharacterStore();

  useEffect(() => {
    const fetchCharacters = async () => {
      if (!novelId) {
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);
      
      try {
        const fetchedCharacters = await characterApi.getCharactersByNovel(novelId);
        setCharacters(fetchedCharacters);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch characters');
      } finally {
        setLoading(false);
      }
    };

    fetchCharacters();
  }, [novelId]);

  return { characters, loading, error };
};
