import React, { useState, useRef } from "react";
import {
  ImageIcon,
  VideoIcon,
  BarChart2Icon,
  SendIcon,
  X,
  Globe,
  Users,
  Lock,
} from "lucide-react";
import { postAPI } from "../../lib/api";
import UserSearch from "./UserSearch";
import { profileAPI } from "../../lib/api";


const PostCreation = ({ user, onPostCreated }) => {
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [selectedUsers, setSelectedUsers] = useState([]);
  const [showPrivacyDropdown, setShowPrivacyDropdown] = useState(false);
  const fileInputRef = useRef(null);

  const privacyOptions = [
    {
      value: "public",
      label: "Public",
      icon: Globe,
      description: "Anyone can see this post",
    },
    {
      value: "almost_private",
      label: "Followers",
      icon: Users,
      description: "Only your followers can see this",
    },
    {
      value: "private",
      label: "Private",
      icon: Lock,
      description: "Only specific people can see this",
    },
  ];

  const handleImageSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      // Validate file type
      const allowedTypes = ["image/jpeg", "image/png", "image/gif"];
      if (!allowedTypes.includes(file.type)) {
        setError("Only JPEG, PNG, and GIF images are allowed");
        return;
      }

      // Validate file size (20MB limit)
      if (file.size > 20 * 1024 * 1024) {
        setError("Image size must be less than 20MB");
        return;
      }

      setSelectedImage(file);
      setImagePreview(URL.createObjectURL(file));
      setError("");
    }
  };

  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleUserSelect = (user) => {
    setSelectedUsers(prev => [...prev, user]);
  };

  const handleUserRemove = (userId) => {
    setSelectedUsers(prev => prev.filter(user => user.id !== userId));
  };

  const handlePrivacyChange = (newPrivacy) => {
    setPrivacy(newPrivacy);
    setShowPrivacyDropdown(false);
    // Clear selected users if not private
    if (newPrivacy !== "private") {
      setSelectedUsers([]);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!content.trim()) {
      setError("Post content is required");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      const formData = new FormData();
      formData.append("content", content.trim());
      formData.append("privacy", privacy);

      // Add selected users for private posts
      if (privacy === "private" && selectedUsers.length > 0) {
        const viewerIds = selectedUsers.map(user => user.id).join(",");
        formData.append("viewers", viewerIds);
      }

      if (selectedImage) {
        formData.append("image", selectedImage);
      }

      const result = await postAPI.createPost(formData);

      if (result.success) {
        // Reset form
        setContent("");
        setPrivacy("public");
        setSelectedUsers([]);
        setSelectedImage(null);
        setImagePreview(null);
        if (fileInputRef.current) {
          fileInputRef.current.value = "";
        }

        // Notify parent component
        if (onPostCreated) {
          onPostCreated(result.data);
        }
      } else {
        setError(result.error);
      }
    } catch (error) {
      setError("Failed to create post. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const getCurrentPrivacyOption = () => {
    return privacyOptions.find((option) => option.value === privacy);
  };

  return (
    <div
      className="rounded-xl p-4 mb-4"
      style={{ backgroundColor: "var(--primary-background)" }}
    >
      <form onSubmit={handleSubmit}>
        {/* Post Input Area */}
        <div
          className="flex items-start gap-3 rounded-xl p-3 mb-4"
          style={{ backgroundColor: "var(--secondary-background)" }}
        >
          <img
            src={
                profileAPI.fetchProfileImage(user?.avatar || '')
            }
            alt="Profile"
            className="w-10 h-10 rounded-full flex-shrink-0"
          />
          <div className="flex-1 min-w-0">
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Tell your friends about your thoughts..."
              className="w-full focus:outline-none text-sm resize-none min-h-[60px] break-words overflow-wrap-anywhere"
              style={{ backgroundColor: 'transparent', color: 'var(--primary-text)', '--placeholder-color': 'var(--secondary-text)' }}
              rows="3"
              disabled={isSubmitting}
            />

            {/* Image Preview */}
            {imagePreview && (
              <div className="relative mt-3 inline-block">
                <img
                  src={imagePreview}
                  alt="Preview"
                  className="max-w-full max-h-48 rounded-lg"
                />
                <button
                  type="button"
                  onClick={removeImage}
                  className="absolute top-2 right-2 rounded-full p-1"
                  style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)', color: 'var(--primary-text)' }}
                  onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'rgba(0, 0, 0, 0.7)'}
                  onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'rgba(0, 0, 0, 0.5)'}
                  disabled={isSubmitting}
                >
                  <X className="w-4 h-4" />
                </button>
              </div>
            )}
          </div>
        </div>

        {/* User Search for Private Posts */}
        {privacy === "private" && (
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2" style={{ color: 'var(--secondary-text)' }}>
              Select people who can see this post:
            </label>
            <UserSearch
              selectedUsers={selectedUsers}
              onUserSelect={handleUserSelect}
              onUserRemove={handleUserRemove}
            />
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="mb-4 p-3 border rounded-lg text-sm"
            style={{ backgroundColor: 'rgba(var(--danger-color-rgb), 0.2)', borderColor: 'var(--warning-color)', color: 'var(--warning-color)' }}>
            {error}
          </div>
        )}

        {/* Privacy Selector and Action Buttons */}
        <div className="flex justify-between items-center border-t pt-3" style={{ borderColor: 'var(--border-color)' }}>
          <div className="flex items-center gap-2">
            {/* Image Upload Button */}
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg cursor-pointer transition-colors"
              style={{ color: "var(--secondary-text)" }}
              onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
              onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
              disabled={isSubmitting}
            >
              <ImageIcon className="w-4 h-4" />
              <span>Photo</span>
            </button>

            {/* Hidden File Input */}
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/gif"
              onChange={handleImageSelect}
              className="hidden"
              disabled={isSubmitting}
            />

            <button
              className="flex items-center gap-2 text-sm py-1.5 px-3 cursor-pointer rounded-lg"
              style={{ color: "var(--secondary-text)" }}
              onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
              onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            >
              <VideoIcon className="w-4 h-4" />
              <span>Video</span>
            </button>
            
            <button
              className="flex items-center gap-2 text-sm py-1.5 px-3 cursor-pointer rounded-lg"
              style={{ color: "var(--secondary-text)" }}
              onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
              onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
            >
              <BarChart2Icon className="w-4 h-4" />
              <span>Poll</span>
            </button>

            {/* Privacy Selector */}
            <div className="relative">
              <button
                type="button"
                onClick={() => setShowPrivacyDropdown(!showPrivacyDropdown)}
                className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg cursor-pointer transition-colors"
                style={{ color: "var(--secondary-text)" }}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                disabled={isSubmitting}
              >
                {React.createElement(getCurrentPrivacyOption().icon, {
                  className: "w-4 h-4",
                })}
                <span>{getCurrentPrivacyOption().label}</span>
              </button>

              {/* Privacy Dropdown */}
              {showPrivacyDropdown && (
                <div className="absolute bottom-full left-0 mb-2 rounded-lg shadow-lg min-w-[200px] z-10"
                  style={{ backgroundColor: 'var(--secondary-background)', border: '1px solid var(--border-color)' }}>
                  {privacyOptions.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => handlePrivacyChange(option.value)}
                      className={`w-full text-left px-3 py-2 first:rounded-t-lg last:rounded-b-lg transition-colors ${privacy === option.value ? "" : ""}`}
                      style={privacy === option.value ? { backgroundColor: 'var(--hover-background)' } : {}}
                      onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                      onMouseOut={(e) => e.currentTarget.style.backgroundColor = privacy === option.value ? 'var(--hover-background)' : 'transparent'}
                      disabled={isSubmitting}
                    >
                      <div className="flex items-center gap-2 mb-1">
                        {React.createElement(option.icon, {
                          className: "w-4 h-4",
                        })}
                        <span className="text-sm font-medium" style={{ color: 'var(--primary-text)' }}>
                          {option.label}
                        </span>
                      </div>
                      <div className="text-xs ml-6" style={{ color: 'var(--secondary-text)' }}>
                        {option.description}
                      </div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Post Button */}
          <button
            type="submit"
            disabled={!content.trim() || isSubmitting}
            className={`flex items-center gap-2 text-sm py-2 px-4 rounded-lg transition-all`}
            style={{
              cursor: !content.trim() || isSubmitting ? 'not-allowed' : 'pointer',
              backgroundColor: !content.trim() || isSubmitting ? 'var(--secondary-background)' : 'var(--primary-accent)',
              color: !content.trim() || isSubmitting ? 'var(--secondary-text)' : 'var(--quinary-text)'
            }}
            onMouseOver={(e) => e.currentTarget.style.opacity = !content.trim() || isSubmitting ? '1' : '0.8'}
            onMouseOut={(e) => e.currentTarget.style.opacity = '1'}
          >
            <SendIcon className="w-4 h-4" />
            <span>{isSubmitting ? "Posting..." : "Post"}</span>
          </button>
        </div>
      </form>
    </div>
  );
};

export default PostCreation;
