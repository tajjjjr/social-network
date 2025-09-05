'use client';

import { useState, useEffect } from 'react';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

export default function GroupPosts({ groupId }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);

  const fetchPosts = async (pageNum = 1, reset = false) => {
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/groups/${groupId}/posts?limit=10&offset=${(pageNum - 1) * 10}`,
        { credentials: 'include' }
      );
      
      if (response.ok) {
        const newPosts = await response.json();
        const postsArray = newPosts || [];
        if (reset) {
          setPosts(postsArray);
        } else {
          setPosts(prev => [...prev, ...postsArray]);
        }
        setHasMore(postsArray.length === 10);
      }
    } catch (error) {
      console.error('Error fetching group posts:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (groupId) {
      fetchPosts(1, true);
    }
  }, [groupId]);

  const handlePostCreated = () => {
    fetchPosts(1, true);
    setPage(1);
  };

  const loadMore = () => {
    const nextPage = page + 1;
    setPage(nextPage);
    fetchPosts(nextPage, false);
  };

  return (
    <div className="space-y-6">
      {/* Post Creation */}
      <div 
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--secondary-background)' }}
      >
        <PostCreation 
          onPostCreated={handlePostCreated}
          isGroupPost={true}
          groupId={groupId}
        />
      </div>

      {/* Posts Feed */}
      {loading ? (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 mx-auto" style={{ borderColor: 'var(--primary-accent)' }}></div>
        </div>
      ) : (
        <div className="space-y-6">
          <PostList 
            posts={posts} 
            onPostUpdate={handlePostCreated}
            isGroupPost={true}
          />
          
          {hasMore && posts.length > 0 && (
            <div className="text-center">
              <button
                onClick={loadMore}
                className="px-6 py-2 rounded-lg font-medium"
                style={{
                  backgroundColor: 'var(--secondary-background)',
                  color: 'var(--primary-text)',
                  border: '1px solid var(--tertiary-text)'
                }}
              >
                Load More Posts
              </button>
            </div>
          )}
          
          {posts.length === 0 && (
            <div className="text-center py-12">
              <div 
                className="rounded-xl p-8"
                style={{ backgroundColor: 'var(--secondary-background)' }}
              >
                <h3 className="text-lg font-medium mb-2" style={{ color: 'var(--primary-text)' }}>
                  No posts here yet
                </h3>
                <p style={{ color: 'var(--secondary-text)' }}>
                  Be the first to share something with this group!
                </p>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}