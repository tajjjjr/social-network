import React, { useState, useRef } from "react";
import { ImageIcon, SendIcon, X } from "lucide-react";
import { postAPI } from "../../lib/api";
import { profileAPI } from "../../lib/api";

const CommentForm = ({ postId, user, onCommentCreated }) => {
  const [content, setContent] = useState("");
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const fileInputRef = useRef(null);

  // Validate image file
  const validateImage = (file) => {
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif'];
    const maxSize = 10 * 1024 * 1024; // 10MB

    if (!allowedTypes.includes(file.type)) {
      return "Please select a valid image file (JPEG, PNG, or GIF)";
    }

    if (file.size > maxSize) {
      return "Image size must be less than 10MB";
    }

    return null;
  };

  // Handle image selection
  const handleImageSelect = (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const validationError = validateImage(file);
    if (validationError) {
      setError(validationError);
      return;
    }

    setError("");
    setSelectedImage(file);

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      setImagePreview(e.target.result);
    };
    reader.readAsDataURL(file);
  };

  // Remove selected image
  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!content.trim() && !selectedImage) {
      setError("Please enter a comment or select an image");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      const formData = new FormData();
      if (content.trim()) {
        formData.append("content", content.trim());
      }

      if (selectedImage) {
        formData.append("image", selectedImage);
      }

      const result = await postAPI.createComment(postId, formData);

      if (result.success) {
        // Reset form
        setContent("");
        setSelectedImage(null);
        setImagePreview(null);
        if (fileInputRef.current) {
          fileInputRef.current.value = "";
        }

        // Enhance the comment data with user info and timestamp before notifying parent
        const enhancedComment = {
          ...result.data,
          author: user, // Add the current user as the author
          created_at: result.data.created_at || new Date().toISOString() // Ensure we have a timestamp
        };

        // Notify parent component
        if (onCommentCreated) {
          onCommentCreated(enhancedComment);
        }
      } else {
        setError(result.error);
      }
    } catch (error) {
      console.error('Error creating comment:', error);
      setError('Network error occurred');
    } finally {
      setIsSubmitting(false);
    }
  };

  // Handle key press for quick submit
  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  return (
    <div className="rounded-xl p-4 w-full" style={{ backgroundColor: 'var(--primary-background)' }}>
      <form onSubmit={handleSubmit}>
        <div className="flex items-start gap-3 rounded-xl p-3" style={{ backgroundColor: 'var(--secondary-background)' }}>
          {/* User Avatar */}
          <img 
            src={profileAPI.fetchProfileImage(user?.avatar || '')}
            alt={user?.nickname || `${user?.first_name || ''} ${user?.last_name || ''}`.trim() || 'User'}
            className="w-10 h-10 rounded-full object-cover flex-shrink-0" 
          />

          {/* Input Container */}
          <div className="flex-1 min-w-0">
            {/* Text Input */}
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Write your comment..."
              className="bg-transparent w-full focus:outline-none text-sm resize-none break-words overflow-wrap-anywhere"
              style={{ color: 'var(--primary-text)' }}
              rows={1}
              disabled={isSubmitting}
              onInput={(e) => {
                e.target.style.height = 'auto';
                e.target.style.height = e.target.scrollHeight + 'px';
              }}
            />

            {/* Image Preview */}
            {imagePreview && (
              <div className="mt-3 relative inline-block">
                <img 
                  src={imagePreview} 
                  alt="Preview" 
                  className="max-w-full max-h-32 rounded-lg object-cover"
                />
                <button
                  type="button"
                  onClick={removeImage}
                  className="absolute top-2 right-2 rounded-full p-1 transition-colors"
                  style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)', color: 'var(--primary-text)' }}
                  onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'rgba(0, 0, 0, 0.7)'}
                  onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'rgba(0, 0, 0, 0.5)'}
                  disabled={isSubmitting}
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            )}

            {/* Error Message */}
            {error && (
              <div className="mt-2 text-xs" style={{ color: 'var(--warning-color)' }}>
                {error}
              </div>
            )}
          </div>

          {/* Action Buttons */}
          <div className="flex items-center gap-2 flex-shrink-0">
            {/* Image Upload Button */}
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="cursor-pointer hover:opacity-70 transition-opacity"
              style={{ color: 'var(--secondary-text)' }}
              disabled={isSubmitting}
            >
              <ImageIcon className="w-5 h-5" />
            </button>

            {/* Submit Button */}
            <button
              type="submit"
              className="cursor-pointer hover:opacity-70 transition-opacity disabled:opacity-50"
              style={{ color: 'var(--secondary-text)' }}
              disabled={isSubmitting || (!content.trim() && !selectedImage)}
            >
              <SendIcon className="w-5 h-5" />
            </button>
          </div>

          {/* Hidden File Input */}
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png,image/gif"
            onChange={handleImageSelect}
            className="hidden"
            disabled={isSubmitting}
          />
        </div>
      </form>
    </div>
  );
};

export default CommentForm;
