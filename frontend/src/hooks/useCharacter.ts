import { useEffect, useState } from 'react';
import { characterApi } from '../services/characterApi';
import type { Character } from '../types';

export const useCharacter = (id: string) => {
  const [character, setCharacter] = useState<Character | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCharacter = async () => {
      if (!id) {
        setLoading(false);
        return;
      }

      setLoading(true);
      setError(null);
      
      try {
        const fetchedCharacter = await characterApi.getCharacter(id);
        setCharacter(fetchedCharacter);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch character');
      } finally {
        setLoading(false);
      }
    };

    fetchCharacter();
  }, [id]);

  return { character, loading, error };
};
