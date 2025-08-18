import React, { useState, useEffect, useRef } from 'react';
// import { fetchPostsPaginated, deletePost, updatePost } from '../../lib/auth';
import { MoreHorizontalIcon, MessageCircleIcon, Globe, Users, Lock, Edit, Trash2, UserPlus, User } from 'lucide-react';
import VerifiedBadge from '../homepage/VerifiedBadge';
import CommentForm from './CommentForm';
import CommentList from './CommentList';
import ReactionButtons from './ReactionButtons';
import ClientDate from '../common/ClientDate';
import { useRouter } from 'next/navigation';
import { generateInitialSkeletons, generatePaginationSkeletons } from '../../lib/skeletonUtils';
import { profileAPI } from '../../lib/api';
import { postAPI } from '../../lib/api';

const PostList = ({ refreshTrigger, user, posts: initialPosts, profileView = false }) => {
  const router = useRouter();
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState('');
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [expandedComments, setExpandedComments] = useState(new Set());
  const [newComments, setNewComments] = useState({});
  const [openDropdown, setOpenDropdown] = useState(null);
  const [deleteConfirmation, setDeleteConfirmation] = useState(null);
  const [editModal, setEditModal] = useState(null);
  const [editContent, setEditContent] = useState('');
  const [editImage, setEditImage] = useState(null);
  const [editLoading, setEditLoading] = useState(false);
  const loadingRef = useRef(false);

  const loadPosts = async (pageNum = 1, append = false) => {
    if (append) {
      setLoadingMore(true);
    } else {
      setLoading(true);
      setPosts([]);
      setPage(1);
      setHasMore(true);
    }
    setError('');

    try {
      const result = await postAPI.fetchPostsPaginated(pageNum, 10);

      if (result.success) {
        const { posts: newPosts, pagination } = result.data;
        if (append) {
          setPosts(prev => [...prev, ...newPosts]);
        } else {
          setPosts(newPosts);
        }
        setHasMore(pagination.hasMore);
        setPage(pageNum + 1);
      } else {
        setError(result.error);
      }
    } catch (error) {
      setError('Failed to load posts');
    } finally {
      setLoading(false);
      setLoadingMore(false);
      loadingRef.current = false;
    }
  };

  useEffect(() => {
    if (initialPosts) {
      setPosts(initialPosts);
      setLoading(false);
    } else {
      loadPosts();
    }
  }, [refreshTrigger, initialPosts]);

  const loadMorePosts = () => {
    if (loadingRef.current || loadingMore || !hasMore) return;
    loadingRef.current = true;
    loadPosts(page, true);
  };

  useEffect(() => {
    // Don't add scroll listener in profile view
    if (profileView) return;
    
    const handleScroll = () => {
      if (loadingMore || !hasMore) return;
      if (window.innerHeight + document.documentElement.scrollTop >= document.documentElement.offsetHeight - 1000) {
        loadMorePosts();
      }
    };

    const debouncedScroll = debounce(handleScroll, 200);
    window.addEventListener('scroll', debouncedScroll);
    return () => window.removeEventListener('scroll', debouncedScroll);
  }, [page, loadingMore, hasMore, profileView]);

  const debounce = (func, wait) => {
    let timeout;
    return function executedFunction(...args) {
      const later = () => {
        clearTimeout(timeout);
        func(...args);
      };
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
    };
  };

  const getPrivacyIcon = (privacy) => {
    switch (privacy) {
      case 'public':
        return Globe;
      case 'almost_private':
        return Users;
      case 'private':
        return Lock;
      default:
        return Globe;
    }
  };

  const getPrivacyLabel = (privacy) => {
    switch (privacy) {
      case 'public':
        return 'Public';
      case 'almost_private':
        return 'Followers';
      case 'private':
        return 'Private';
      default:
        return 'Public';
    }
  };



  // Toggle comments visibility for a post
  const toggleComments = (postId) => {
    setExpandedComments(prev => {
      const newSet = new Set(prev);
      if (newSet.has(postId)) {
        newSet.delete(postId);
      } else {
        newSet.add(postId);
      }
      return newSet;
    });
  };

  // Handle new comment creation
  const handleCommentCreated = (postId, comment) => {
    setNewComments(prev => ({
      ...prev,
      [postId]: comment
    }));

    // Clear the new comment after a short delay to allow CommentList to process it
    setTimeout(() => {
      setNewComments(prev => {
        const updated = { ...prev };
        delete updated[postId];
        return updated;
      });
    }, 100);
  };

  // Handle dropdown toggle
  const toggleDropdown = (postId) => {
    setOpenDropdown(openDropdown === postId ? null : postId);
  };

  // Handle edit post
  const handleEditPost = (post) => {
    setEditModal(post.id);
    setEditContent(post.content);
    setEditImage(null);
    setOpenDropdown(null);
  };

  // Handle edit post submission
  const handleEditSubmit = async (postId) => {
    if (!editContent.trim()) return;

    setEditLoading(true);
    try {
      const result = await postAPI.updatePost(postId, editContent, editImage);
      if (result.success) {
        // Update the post in the local state
        setPosts(prev => prev.map(post =>
          post.id === postId ? result.data : post
        ));
        setEditModal(null);
        setEditContent('');
        setEditImage(null);
      } else {
        console.error('Failed to update post:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error updating post:', error);
    } finally {
      setEditLoading(false);
    }
  };

  // Handle delete post
  const handleDeletePost = async (postId) => {
    try {
      const result = await postAPI.deletePost(postId);
      if (result.success) {
        // Remove the post from the local state
        setPosts(prev => prev.filter(post => post.id !== postId));
        setDeleteConfirmation(null);
        setOpenDropdown(null);
      } else {
        console.error('Failed to delete post:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error deleting post:', error);
    }
  };

  // Handle view profile
  const handleViewProfile = (userId) => {
    router.push(`/profile/${userId}`);
    setOpenDropdown(null);
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      // Check if the clicked element is part of any dropdown
      const isDropdownClick = event.target.closest('.dropdown-container');
      if (!isDropdownClick && openDropdown !== null) {
        setOpenDropdown(null);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [openDropdown]);

  if (loading && (!posts || posts.length === 0)) {
    try {
      return (
        <div className="space-y-4">
          {generateInitialSkeletons()}
        </div>
      );
    } catch (error) {
      console.warn('Skeleton loading failed, using fallback:', error);
      return (
        <div className="flex justify-center items-center py-8">
          <div style={{ color: 'var(--secondary-text)' }}>Loading posts...</div>
        </div>
      );
    }
  }

  if (error) {
    return (
      <div className="rounded-xl p-4 mb-4" style={{ backgroundColor: 'rgba(var(--danger-color-rgb), 0.2)', border: '1px solid var(--warning-color)' }}>
        <div style={{ color: 'var(--warning-color)' }} className="text-center">{error}</div>
        <button
          onClick={loadPosts}
          className="mt-2 w-full py-2 px-4 rounded-lg transition-colors"
          style={{ backgroundColor: 'var(--warning-color)', color: 'var(--primary-text)' }}
          onMouseOver={(e) => e.currentTarget.style.opacity = '0.8'}
          onMouseOut={(e) => e.currentTarget.style.opacity = '1'}
        >
          Try Again
        </button>
      </div>
    );
  }

  if (!posts || posts.length === 0) {
    return (
      <div className="rounded-xl p-8 text-center" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div style={{ color: 'var(--secondary-text)' }} className="mb-2">No posts yet</div>
        <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Be the first to share something!</div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <div key={post.id} className="rounded-xl p-4 post-content" style={{ backgroundColor: 'var(--primary-background)' }}>
          {/* Post Header */}
          <div className="flex justify-between mb-3">
            <div className="flex items-center gap-2">
              <div className="relative">
                <img
                  src={profileAPI.fetchProfileImage(post.author?.avatar || '')}
                  alt={post.author?.nickname || `${post.author?.first_name || ''} ${post.author?.last_name || ''}`.trim() || 'User'}
                  className="w-10 h-10 rounded-full"
                />
                <div className="absolute -bottom-1 -right-1">
                  <VerifiedBadge />
                </div>
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-medium break-words" style={{ color: 'var(--primary-text)' }}>
                    {post.author?.nickname || `${post.author?.first_name || ''} ${post.author?.last_name || ''}`.trim() || 'User'}
                  </span>
                  <div className="flex items-center gap-1 text-xs" style={{ color: 'var(--secondary-text)' }}>
                    {React.createElement(getPrivacyIcon(post.privacy), { className: "w-3 h-3" })}
                    <span>{getPrivacyLabel(post.privacy)}</span>
                  </div>
                </div>
                <div className="text-xs" style={{ color: 'var(--secondary-text)' }}>
                  <ClientDate dateString={post.created_at} />
                  {post.is_edited && (
                    <span className="ml-2 text-xs" style={{ color: 'var(--secondary-text)' }}>
                      • edited
                    </span>
                  )}
                </div>
              </div>
            </div>
            {/* Post Actions Dropdown */}
            <div className="relative dropdown-container">
              <button
                onClick={() => toggleDropdown(post.id)}
                className="transition-colors"
                style={{ color: 'var(--secondary-text)' }}
                onMouseOver={(e) => e.currentTarget.style.color = 'var(--primary-text)'}
                onMouseOut={(e) => e.currentTarget.style.color = 'var(--secondary-text)'}
              >
                <MoreHorizontalIcon className="w-5 h-5" />
              </button>

              {/* Dropdown Menu */}
              {openDropdown === post.id && (
                <div
                  className="absolute right-0 top-8 w-48 rounded-lg shadow-lg z-10"
                  style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
                >
                  <div className="py-1">
                    {user && user.id === post.user_id ? (
                      <>
                        {/* Edit Option */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEditPost(post);
                          }}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <Edit className="w-4 h-4" />
                          Edit
                        </button>

                        {/* Delete Option */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setDeleteConfirmation(post.id);
                            setOpenDropdown(null); // Close dropdown when opening modal
                          }}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--warning-color)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <Trash2 className="w-4 h-4" />
                          Delete
                        </button>
                      </>
                    ) : (
                      <>
                        {/* View Profile Option */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleViewProfile(post.user_id);
                          }}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <User className="w-4 h-4" />
                          View Profile
                        </button>

                        {/* Follow Option */}
                        <button
                          onClick={() => profileAPI.follow(post.user_id)}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <UserPlus className="w-4 h-4" />
                          Follow
                        </button>
                      </>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Post Content */}
          <div className="mb-4">
            <p className="whitespace-pre-wrap break-words overflow-wrap-anywhere" style={{ color: 'var(--primary-text)' }}>{post.content}</p>

            {/* Post Image */}
            {post.image && (
              <div className="mt-3">
                <img
                  src={`http://localhost:9000/avatar?avatar=${post.image}`}
                  alt="Post image"
                  className="max-w-full rounded-lg"
                  onError={(e) => {
                    e.target.style.display = 'none';
                  }}
                />
              </div>
            )}
          </div>

          {/* Post Actions */}
          <div className="flex items-center justify-between pt-3 border-t" style={{ borderColor: 'var(--border-color)' }}>
            <div className="flex items-center gap-4">
              <ReactionButtons post={post} user={user} />
              <button
                onClick={() => toggleComments(post.id)}
                className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg transition-colors"
                style={{ color: 'var(--secondary-text)' }}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}>
                <MessageCircleIcon className="w-4 h-4" />
                <span>{expandedComments.has(post.id) ? 'Hide Comments' : 'Comment'}</span>
              </button>
            </div>
          </div>

          {/* Comments Section */}
          {expandedComments.has(post.id) && (
            <div className="mt-4 pt-3 border-t space-y-4" style={{ borderColor: 'var(--border-color)' }}>
              {/* Comment Form */}
              <CommentForm
                postId={post.id}
                user={user}
                onCommentCreated={(comment) => handleCommentCreated(post.id, comment)}
              />

              {/* Comments List */}
              <CommentList
                postId={post.id}
                newComment={newComments[post.id]}
                user={user}
              />
            </div>
          )}
        </div>
      ))}

      {/* Loading more indicator - Only show in non-profile view */}
      {!profileView && loadingMore && (
        (() => {
          try {
            return (
              <div className="space-y-4">
                {generatePaginationSkeletons()}
              </div>
            );
          } catch (error) {
            console.warn('Pagination skeleton loading failed, using fallback:', error);
            return (
              <div className="flex justify-center py-4">
                <div style={{ color: 'var(--secondary-text)' }}>Loading more posts...</div>
              </div>
            );
          }
        })()
      )}

      {/* End of posts - Only show in non-profile view */}
      {!profileView && !hasMore && posts.length > 0 && (
        <div className="text-center py-4" style={{ color: 'var(--secondary-text)' }}>
          No more posts to show
        </div>
      )}

      {/* Delete Confirmation Modal */}
      {deleteConfirmation && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Delete Post</h3>
            <p className="mb-6" style={{ color: 'var(--primary-text)' }}>Are you sure you want to delete this post? This action cannot be undone.</p>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setDeleteConfirmation(null)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'var(--secondary-background)'}
              >
                Cancel
              </button>
              <button
                onClick={() => handleDeletePost(deleteConfirmation)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--warning-color)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.opacity = '0.8'}
                onMouseLeave={(e) => e.currentTarget.style.opacity = '1'}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Post Modal */}
      {editModal && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Edit Post</h3>

            <div className="mb-4">
              <textarea
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                placeholder="What's on your mind?"
                className="w-full p-3 rounded-lg resize-none break-words overflow-wrap-anywhere"
                style={{
                  backgroundColor: 'var(--secondary-background)',
                  border: '1px solid var(--border-color)',
                  minHeight: '100px',
                  color: 'var(--primary-text)'
                }}
                rows={4}
              />
            </div>

            <div className="mb-4">
              <input
                type="file"
                accept="image/*"
                onChange={(e) => setEditImage(e.target.files[0])}
                className="w-full"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
              />
            </div>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => {
                  setEditModal(null);
                  setEditContent('');
                  setEditImage(null);
                }}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'var(--secondary-background)'}
                disabled={editLoading}
              >
                Cancel
              </button>
              <button
                onClick={() => handleEditSubmit(editModal)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--primary-accent)', color: 'var(--quinary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.opacity = '0.8'}
                onMouseLeave={(e) => e.currentTarget.style.opacity = '1'}
                disabled={editLoading || !editContent.trim()}
              >
                {editLoading ? 'Saving...' : 'Save Changes'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PostList;
