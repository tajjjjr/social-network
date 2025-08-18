import React, { useState, useEffect } from "react";
// import { fetchComments, updateComment, deleteComment, followUser } from "../../lib/auth";
import { MoreHorizontalIcon, Edit, Trash2, UserPlus } from 'lucide-react';
import CommentReactionButtons from './CommentReactionButtons';
import ClientDate from '../common/ClientDate';
import { profileAPI } from "../../lib/api";
import { postAPI } from "../../lib/api";

const CommentList = ({ postId, newComment, user }) => {
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [openDropdown, setOpenDropdown] = useState(null);
  const [deleteConfirmation, setDeleteConfirmation] = useState(null);
  const [editModal, setEditModal] = useState(null);
  const [editContent, setEditContent] = useState('');
  const [editImage, setEditImage] = useState(null);
  const [editLoading, setEditLoading] = useState(false);



  // Get display name for user
  const getDisplayName = (author) => {
    if (author?.nickname) {
      return author.nickname;
    }
    const firstName = author?.first_name || '';
    const lastName = author?.last_name || '';
    const fullName = `${firstName} ${lastName}`.trim();
    return fullName || 'User';
  };

  // Load comments
  const loadComments = async () => {
    try {
      setLoading(true);
      const result = await postAPI.fetchComments(postId);

      if (result.success) {
        setComments(result.data || []);
        setError("");
      } else {
        setError(result.error);
      }
    } catch (error) {
      console.error('Error loading comments:', error);
      setError('Failed to load comments');
    } finally {
      setLoading(false);
    }
  };

  // Load comments on mount and when postId changes
  useEffect(() => {
    if (postId) {
      loadComments();
    }
  }, [postId]);

  // Add new comment to the list when one is created
  useEffect(() => {
    if (newComment) {
      setComments(prev => [newComment, ...prev]); // Add to beginning for newest first
    }
  }, [newComment]);

  // Handle dropdown toggle
  const toggleDropdown = (commentId) => {
    setOpenDropdown(openDropdown === commentId ? null : commentId);
  };

  // Handle edit comment
  const handleEditComment = (comment) => {
    setEditModal(comment.id);
    setEditContent(comment.content);
    setEditImage(null);
    setOpenDropdown(null);
  };

  // Handle edit comment submission
  const handleEditSubmit = async (commentId) => {
    if (!editContent.trim()) return;

    setEditLoading(true);
    try {
      const result = await postAPI.updateComment(postId, commentId, editContent, editImage);
      if (result.success) {
        // Update the comment in the local state
        setComments(prev => prev.map(comment =>
          comment.id === commentId ? result.data : comment
        ));
        setEditModal(null);
        setEditContent('');
        setEditImage(null);
      } else {
        console.error('Failed to update comment:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error updating comment:', error);
    } finally {
      setEditLoading(false);
    }
  };

  // Handle delete comment
  const handleDeleteComment = async (commentId) => {
    try {
      const result = await postAPI.deleteComment(postId, commentId);
      if (result.success) {
        // Remove the comment from the local state
        setComments(prev => prev.filter(comment => comment.id !== commentId));
        setDeleteConfirmation(null);
        setOpenDropdown(null);
      } else {
        console.error('Failed to delete comment:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error deleting comment:', error);
    }
  };

  // Handle follow user
  const handleFollowUser = async (userId) => {
    try {
      const result = await followUser(userId);
      if (result.success) {
        // You could add a toast notification here
        console.log('Followed user successfully');
      } else {
        console.error('Failed to follow user:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error following user:', error);
    }
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

  if (loading) {
    return (
      <div className="text-center py-4" style={{ color: 'var(--secondary-text)' }}>
        <div className="text-sm">Loading comments...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-4">
        <div className="text-sm" style={{ color: 'var(--warning-color)' }}>{error}</div>
        <button
          onClick={loadComments}
          className="text-sm mt-2 hover:underline"
          style={{ color: 'var(--primary-accent)' }}
        >
          Try again
        </button>
      </div>
    );
  }

  if (comments.length === 0) {
    return (
      <div className="text-center py-4" style={{ color: 'var(--secondary-text)' }}>
        <div className="text-sm">No comments yet. Be the first to comment!</div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {comments.map((comment) => (
        <div key={comment.id} className="flex gap-3">
          {/* Comment Author Avatar */}
          <img
            src={
              profileAPI.fetchProfileImage(comment.author?.avatar || '')
            }
            alt={getDisplayName(comment.author)}
            className="w-8 h-8 rounded-full object-cover flex-shrink-0"
          />

          {/* Comment Content */}
          <div className="flex-1 min-w-0">
            <div
              className="rounded-lg p-3 relative"
              style={{ backgroundColor: 'var(--secondary-background)' }}
            >
              {/* Author Name */}
              <div className="font-medium text-sm mb-1 break-words" style={{ color: 'var(--primary-text)' }}>
                {getDisplayName(comment.author)}
              </div>

              {/* Comment Text */}
              {comment.content && (
                <div className="text-sm mb-2" style={{ color: 'var(--primary-text)' }}>
                  <p className="whitespace-pre-wrap break-words overflow-wrap-anywhere">{comment.content}</p>
                </div>
              )}

              {/* Comment Image */}
              {comment.image && (
                <div className="mt-2">
                  <img
                    src={`http://localhost:9000/avatar?avatar=${comment.image}`}
                    alt="Comment image"
                    className="max-w-full max-h-48 rounded-lg object-cover"
                    onError={(e) => {
                      e.target.style.display = 'none';
                    }}
                  />
                </div>
              )}
              {/* Comment Actions */}
              <div className="flex items-center justify-between pt-3 border-t" style={{ borderColor: 'var(--border-color)' }}>
                <div className="flex items-center gap-4">
                  <CommentReactionButtons comment={comment} user={user} />
                </div>
              </div>

              {/* Comment Actions Dropdown */}
              <div className="absolute top-2 right-2 dropdown-container">
                <button
                  onClick={() => toggleDropdown(comment.id)}
                  className="transition-colors"
                  style={{ color: 'var(--secondary-text)' }}
                  onMouseOver={(e) => e.currentTarget.style.color = 'var(--primary-text)'}
                  onMouseOut={(e) => e.currentTarget.style.color = 'var(--secondary-text)'}
                >
                  <MoreHorizontalIcon className="w-5 h-5" />
                </button>

                {/* Dropdown Menu */}
                {openDropdown === comment.id && (
                  <div
                    className="absolute right-0 top-8 w-48 rounded-lg shadow-lg z-10"
                    style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
                  >
                    <div className="py-1">
                      {user && user.id === comment.user_id ? (
                        <>
                          {/* Edit Option */}
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleEditComment(comment);
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
                              setDeleteConfirmation(comment.id);
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
                        <button
                          onClick={() => profileAPI.follow(comment.user_id)}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <UserPlus className="w-4 h-4" />
                          Follow
                        </button>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Comment Timestamp */}
            <div
              className="text-xs mt-1 ml-3"
              style={{ color: 'var(--secondary-text)' }}
            >
              <ClientDate dateString={comment.created_at} />
              {comment.is_edited && (
                <span className="ml-2 text-xs" style={{ color: 'var(--secondary-text)' }}>
                  • edited
                </span>
              )}
            </div>
          </div>
        </div>
      ))}

      {/* Delete Confirmation Modal */}
      {deleteConfirmation && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Delete Comment</h3>
            <p className="mb-6" style={{ color: 'var(--primary-text)' }}>Are you sure you want to delete this comment? This action cannot be undone.</p>

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
                onClick={() => handleDeleteComment(deleteConfirmation)}
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

      {/* Edit Comment Modal */}
      {editModal && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Edit Comment</h3>

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

export default CommentList;