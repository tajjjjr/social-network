'use client';

import { useState } from 'react';
import { api } from '../../lib/api';

export default function GroupCreation({ onGroupCreated }) {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    privacy: 'public'
  });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!formData.title.trim()) return;

    setLoading(true);
    try {
      const response = await api.post('/groups', formData);
      setFormData({ title: '', description: '', privacy: 'public' });
      if (onGroupCreated) onGroupCreated(response.data);
      alert('Group created successfully!');
    } catch (error) {
      console.error('Error creating group:', error);
      alert('Failed to create group');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow p-6 mb-6">
      <h2 className="text-xl font-semibold mb-4">Create New Group</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2">Group Title</label>
          <input
            type="text"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg"
            placeholder="Enter group title"
            required
          />
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2">Description</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg h-24"
            placeholder="Describe your group"
          />
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium mb-2">Privacy</label>
          <select
            value={formData.privacy}
            onChange={(e) => setFormData({ ...formData, privacy: e.target.value })}
            className="w-full px-3 py-2 border rounded-lg"
          >
            <option value="public">Public</option>
            <option value="private">Private</option>
          </select>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600 disabled:opacity-50"
        >
          {loading ? 'Creating...' : 'Create Group'}
        </button>
      </form>
    </div>
  );
}