'use client';

import { useState, useEffect } from 'react';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';
import { groupAPI } from '../../lib/api';

export default function GroupPosts({ groupId, user, group }) {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState(null);
  const [isMember, setIsMember] = useState(false);

  const fetchPosts = async (pageNum = 1, reset = false) => {
    try {
      const result = await groupAPI.getGroupPosts(groupId, 10, (pageNum - 1) * 10);
      
      if (result.success) {
        const postsArray = result.data || [];
        if (reset) {
          setPosts(postsArray);
        } else {
          setPosts(prev => [...prev, ...postsArray]);
        }
        setHasMore(postsArray.length === 10);
        setIsMember(true);
        setError(null);
      } else {
        if (result.error.includes('403') || result.error.includes('Forbidden')) {
          setIsMember(false);
          setError('You must be a group member to view posts');
          setPosts([]);
        } else {
          setError('Failed to load posts');
        }
      }
    } catch (error) {
      console.error('Error fetching group posts:', error);
      setError('Failed to load posts');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (groupId) {
      // Check if user is creator (always a member)
      if (user && group && user.id === group.creator_id) {
        setIsMember(true);
      }
      fetchPosts(1, true);
    }
  }, [groupId, user, group]);

  const handlePostCreated = (postData, isGroupContext) => {
    // Only refresh if this is actually a group post
    if (isGroupContext) {
      fetchPosts(1, true);
      setPage(1);
    }
  };

  const loadMore = () => {
    const nextPage = page + 1;
    setPage(nextPage);
    fetchPosts(nextPage, false);
  };

  if (!groupId || typeof groupId !== 'string' || groupId.length === 0) {
    setError('Invalid group ID');
    setLoading(false);
    return;
  }

  if (error && !isMember) {
    return (
      <div className="space-y-6">
        <div 
          className="rounded-xl p-8 text-center"
          style={{ backgroundColor: 'var(--secondary-background)' }}
        >
          <h3 className="text-lg font-medium mb-2" style={{ color: 'var(--primary-text)' }}>
            Join the group to see posts
          </h3>
          <p style={{ color: 'var(--secondary-text)' }}>
            You need to be a member of this group to view and create posts.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Post Creation - Only show for members */}
      {isMember && (
        <div 
          className="rounded-xl p-6"
          style={{ backgroundColor: 'var(--secondary-background)' }}
        >
          <PostCreation 
            user={user}
            onPostCreated={handlePostCreated}
            isGroupPost={true}
            groupId={groupId}
          />
        </div>
      )}

      {/* Posts Feed */}
      {loading ? (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 mx-auto" style={{ borderColor: 'var(--primary-accent)' }}></div>
        </div>
      ) : (
        <div className="space-y-6">
          <PostList 
            initialPosts={posts}
            user={user}
            onPostUpdate={handlePostCreated}
            isGroupPost={true}
            groupId={groupId}
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