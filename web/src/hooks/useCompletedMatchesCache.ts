import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  cacheCompletedMatchData,
  clearExpiredCache,
  clearMatchCache,
  clearAllCache,
} from '@/store/reducers/matchSlice';

export const useCompletedMatchesCache = () => {
  const dispatch = useAppDispatch();
  const { completedMatchesCache } = useAppSelector(state => state.match);

  const getCachedData = useCallback(
    (matchId: string) => {
      const cachedData = completedMatchesCache[matchId];
      const now = Date.now();

      if (cachedData && cachedData.expiresAt > now) {
        return cachedData.data;
      }

      return null;
    },
    [completedMatchesCache]
  );

  const setCachedData = useCallback(
    (matchId: string, data: unknown, cacheDuration = 10 * 60 * 1000) => {
      dispatch(cacheCompletedMatchData({
        matchId,
        data,
        cacheDuration,
      }));
    },
    [dispatch]
  );

  const clearExpired = useCallback(() => {
    dispatch(clearExpiredCache());
  }, [dispatch]);

  const clearMatch = useCallback(
    (matchId: string) => {
      dispatch(clearMatchCache(matchId));
    },
    [dispatch]
  );

  const clearAll = useCallback(() => {
    dispatch(clearAllCache());
  }, [dispatch]);

  const isCached = useCallback(
    (matchId: string) => {
      const cachedData = completedMatchesCache[matchId];
      const now = Date.now();
      return cachedData && cachedData.expiresAt > now;
    },
    [completedMatchesCache]
  );

  const getCacheStats = useCallback(() => {
    const now = Date.now();
    const totalEntries = Object.keys(completedMatchesCache).length;
    const expiredEntries = Object.values(completedMatchesCache).filter(
      entry => entry.expiresAt < now
    ).length;
    const validEntries = totalEntries - expiredEntries;

    return {
      totalEntries,
      validEntries,
      expiredEntries,
      cacheSize: JSON.stringify(completedMatchesCache).length,
    };
  }, [completedMatchesCache]);

  return {
    getCachedData,
    setCachedData,
    clearExpired,
    clearMatch,
    clearAll,
    isCached,
    getCacheStats,
    cache: completedMatchesCache,
  };
};
