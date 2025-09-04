'use client';

import { useState } from 'react';
import { api } from '../../lib/api';
import Image from 'next/image';

export default function GroupCreation({ onGroupCreated }) {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    privacy: 'public'
  });
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);
  const [loading, setLoading] = useState(false);
  const [isExpanded, setIsExpanded] = useState(false);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedImage(file);
      const reader = new FileReader();
      reader.onload = (e) => setImagePreview(e.target.result);
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.title.trim()) return;

    setLoading(true);
    try {
      const formDataToSend = new FormData();
      formDataToSend.append('title', formData.title);
      formDataToSend.append('description', formData.description);
      formDataToSend.append('privacy', formData.privacy);
      if (selectedImage) {
        formDataToSend.append('avatar', selectedImage);
      }

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups`, {
        method: 'POST',
        credentials: 'include',
        body: formDataToSend
      });
      if (response.ok) {
        const data = await response.json();
        setFormData({ title: '', description: '', privacy: 'public' });
        setSelectedImage(null);
        setImagePreview(null);
        if (onGroupCreated) onGroupCreated(data);
        // Refresh sidebar groups and group browser
        window.dispatchEvent(new CustomEvent('groupCreated', { detail: data }));
        // Collapse the form after successful creation
        setIsExpanded(false);
      } else {
        console.error('Failed to create group');
      }
    } catch (error) {
      console.error('Error creating group:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mb-6">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="flex items-center justify-between w-full text-xl font-semibold mb-4 p-2 rounded-lg hover:bg-opacity-80 transition-colors"
        style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
      >
        <span>Create New Group</span>
        <span className="transform transition-transform" style={{ transform: isExpanded ? 'rotate(180deg)' : 'rotate(0deg)' }}>▼</span>
      </button>
      {isExpanded && (
        <form onSubmit={handleSubmit} style={{backgroundColor: 'var(--secondary-background)', padding: '1rem', borderRadius: '0.5rem'}}>
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--primary-text)' }}>Group Profile Picture</label>
          <div className="flex items-center gap-4">
            {imagePreview && (
              <Image
                src={imagePreview}
                alt="Group preview"
                width={64}
                height={64}
                className="w-16 h-16 rounded-full object-cover"
              />
            )}
            <input
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              className="block w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold"
              style={{
                color: 'var(--primary-text)'
              }}
            />
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--primary-text)' }}>Group Title</label>
          <input
            type="text"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg"
            style={{ 
              backgroundColor: 'var(--secondary-background)', 
              borderColor: 'var(--tertiary-text)',
              color: 'var(--primary-text)'
            }}
            placeholder="Enter group title"
            required
          />
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--primary-text)' }}>Description</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg h-24"
            style={{ 
              backgroundColor: 'var(--secondary-background)', 
              borderColor: 'var(--tertiary-text)',
              color: 'var(--primary-text)'
            }}
            placeholder="Describe your group"
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium mb-2" style={{ color: 'var(--primary-text)' }}>Privacy</label>
          <select
            value={formData.privacy}
            onChange={(e) => setFormData({ ...formData, privacy: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg"
            style={{ 
              backgroundColor: 'var(--secondary-background)', 
              borderColor: 'var(--tertiary-text)',
              color: 'var(--primary-text)'
            }}
          >
            <option value="public">Public</option>
            <option value="private">Private</option>
          </select>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 rounded-lg disabled:opacity-50"
          style={{
            backgroundColor: 'var(--primary-accent)',
            color: 'white'
          }}
        >
          {loading ? 'Creating...' : 'Create Group'}
        </button>
        </form>
      )}
    </div>
  );
}