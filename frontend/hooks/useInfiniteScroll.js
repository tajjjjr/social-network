import { useState, useEffect, useCallback, useRef } from 'react';

const useInfiniteScroll = (fetchFunction, options = {}) => {
  const {
    initialPage = 1,
    limit = 15,
    threshold = 200
  } = options;

  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  const [initialLoading, setInitialLoading] = useState(true);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(initialPage);
  
  const observerRef = useRef();
  const loadingRef = useRef(false);

  const loadMore = useCallback(async (pageNum = page, isInitial = false) => {
    if (loadingRef.current || (!hasMore && !isInitial)) return;
    
    loadingRef.current = true;
    if (!isInitial) setLoading(true);
    setError(null);

    try {
      const result = await fetchFunction(pageNum, limit);
      
      if (result.success) {
        const { posts, pagination } = result.data;
        
        setData(prevData => pageNum === 1 ? posts : [...prevData, ...posts]);
        setHasMore(pagination.hasMore);
        setPage(pageNum + 1);
      } else {
        setError(result.error);
      }
    } catch (err) {
      setError('Failed to load posts');
    } finally {
      if (!isInitial) setLoading(false);
      setInitialLoading(false);
      loadingRef.current = false;
    }
  }, [fetchFunction, limit, page, hasMore]);

  const observerCallback = useCallback((entries) => {
    const [entry] = entries;
    if (entry.isIntersecting && hasMore && !loadingRef.current) {
      loadMore();
    }
  }, [loadMore, hasMore]);

  const setObserverRef = useCallback((node) => {
    if (observerRef.current) observerRef.current.disconnect();
    
    if (node) {
      observerRef.current = new IntersectionObserver(observerCallback, {
        rootMargin: `${threshold}px`
      });
      observerRef.current.observe(node);
    }
  }, [observerCallback, threshold]);

  useEffect(() => {
    loadMore(initialPage, true);
  }, [initialPage, loadMore]);

  useEffect(() => {
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, []);

  const refresh = useCallback(() => {
    setData([]);
    setPage(initialPage);
    setHasMore(true);
    setError(null);
    loadMore(initialPage, true);
  }, [loadMore, initialPage]);

  return {
    data,
    loading,
    initialLoading,
    hasMore,
    error,
    loadMore,
    refresh,
    setObserverRef
  };
};

export default useInfiniteScroll;