'use client';

import { useState, useEffect } from 'react';
import { fetchGroupImage } from '../../lib/api';

const getGroupAvatar = (avatar) => {
  if (!avatar || avatar.trim() === '') {
    return '/default-group-avatar.png';
  }
  return fetchGroupImage(avatar);
};

export default function GroupHeader({ group }) {
  const [stats, setStats] = useState({ members: 0, posts: 0, events: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroupStats = async () => {
      try {
        const [membersRes, postsRes, eventsRes] = await Promise.all([
          fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/members`, { credentials: 'include' }),
          fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/posts?limit=1`, { credentials: 'include' }),
          fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/events`, { credentials: 'include' })
        ]);

        const members = membersRes.ok ? await membersRes.json() : [];
        const posts = postsRes.ok ? await postsRes.json() : [];
        const events = eventsRes.ok ? await eventsRes.json() : [];

        setStats({
          members: (members || []).length,
          posts: (posts || []).length,
          events: (events || []).length
        });
      } catch (error) {
        console.error('Error fetching group stats:', error);
      } finally {
        setLoading(false);
      }
    };

    if (group?.id) {
      fetchGroupStats();
    }
  }, [group?.id]);

  return (
    <div className="relative">
      {/* Cover Photo */}
      <div 
        className="h-64 bg-gradient-to-r from-blue-500 to-purple-600"
        style={{ backgroundColor: 'var(--tertiary-text)' }}
      />
      
      {/* Group Info */}
      <div className="max-w-7xl mx-auto px-4">
        <div className="relative -mt-16 pb-6">
          <div className="flex items-end gap-6">
            {/* Group Avatar */}
            <div className="relative">
              <img
                src={getGroupAvatar(group.avatar)}
                alt={group.title}
                className="w-32 h-32 rounded-full border-4 border-white object-cover"
                style={{ backgroundColor: 'var(--secondary-background)' }}
                onError={(e) => { e.target.src = '/default-group-avatar.png'; }}
              />
            </div>
            
            {/* Group Details */}
            <div className="flex-1 pb-4">
              <h1 className="text-3xl font-bold mb-2" style={{ color: 'var(--primary-text)' }}>
                {group.title}
              </h1>
              <p className="text-lg mb-2" style={{ color: 'var(--secondary-text)' }}>
                {group.description}
              </p>
              <div className="flex items-center gap-4 text-sm" style={{ color: 'var(--secondary-text)' }}>
                <span>Created {new Date(group.created_at).toLocaleDateString()}</span>
                <span>•</span>
                <span className="capitalize">{group.privacy} Group</span>
              </div>
              
              {/* Group Stats */}
              <div className="flex items-center gap-6 mt-3">
                <div className="flex items-center gap-2">
                  <span className="font-semibold" style={{ color: 'var(--primary-text)' }}>
                    {loading ? '...' : stats.members}
                  </span>
                  <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>Members</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="font-semibold" style={{ color: 'var(--primary-text)' }}>
                    {loading ? '...' : stats.posts}
                  </span>
                  <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>Posts</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="font-semibold" style={{ color: 'var(--primary-text)' }}>
                    {loading ? '...' : stats.events}
                  </span>
                  <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>Events</span>
                </div>
              </div>
            </div>
            
            {/* Action Buttons */}
            <div className="flex gap-3 pb-4">
              <button
                className="px-6 py-2 rounded-lg font-medium"
                style={{
                  backgroundColor: 'var(--primary-accent)',
                  color: 'white'
                }}
              >
                Join Group
              </button>
              <button
                className="px-6 py-2 rounded-lg font-medium border"
                style={{
                  backgroundColor: 'var(--secondary-background)',
                  borderColor: 'var(--tertiary-text)',
                  color: 'var(--primary-text)'
                }}
              >
                Message
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}