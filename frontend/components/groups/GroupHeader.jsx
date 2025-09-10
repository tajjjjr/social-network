'use client';

import { useState, useEffect } from 'react';
import { fetchGroupImage } from '../../lib/api';

const getGroupAvatar = (avatar) => {
  if (!avatar || avatar.trim() === '') {
    return '/default-group-avatar.png';
  }
  return fetchGroupImage(avatar);
};

const getMemberAvatar = (avatar) => {
  if (!avatar || avatar.trim() === '') {
    return '/default-avatar.png';
  }
  return `${process.env.NEXT_PUBLIC_API_URL}/avatar?avatar=${encodeURIComponent(avatar)}`;
};

export default function GroupHeader({ group, user }) {
  const [stats, setStats] = useState({ members: 0, posts: 0, events: 0 });
  const [loading, setLoading] = useState(true);
  const [isMember, setIsMember] = useState(false);
  const [membershipLoading, setMembershipLoading] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [confirmationAction, setConfirmationAction] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchGroupStats = async () => {
      try {
        const groupId = group?.id;

        if (!groupId || typeof groupId !== 'string' || groupId.length === 0) {
          setError('Invalid group ID');
          setLoading(false);
          return;
        }

        const [statsRes, membersRes] = await Promise.all([
          fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${groupId}/stats`, { credentials: 'include' }),
          fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${groupId}/members`, { credentials: 'include' })
        ]);

        let stats = { members: 0, posts: 0, events: 0 };
        let members = [];
        
        if (statsRes.ok) {
          stats = await statsRes.json();
        }
        
        if (membersRes.ok) {
          members = await membersRes.json();
        }
        
        // Check if current user is a member (creator is always a member)
        if (user) {
          if (user.id === group.creator_id) {
            setIsMember(true);
          } else if (members && Array.isArray(members)) {
            setIsMember(members.some(member => member.id === user.id));
          }
        }

        setStats(stats);
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
              {user && (
                user.id === group.creator_id ? (
                  <button
                    className="px-6 py-2 rounded-lg font-medium"
                    style={{
                      backgroundColor: 'var(--secondary-background)',
                      color: 'var(--primary-text)',
                      border: '1px solid var(--tertiary-text)'
                    }}
                    disabled
                  >
                    Group Admin
                  </button>
                ) : (
                  <button
                    onClick={() => {
                      setConfirmationAction(isMember ? 'leave' : 'join');
                      setShowConfirmation(true);
                    }}
                    disabled={membershipLoading}
                    className="px-6 py-2 rounded-lg font-medium"
                    style={{
                      backgroundColor: isMember ? 'var(--warning-color)' : 'var(--primary-accent)',
                      color: 'white',
                      opacity: membershipLoading ? 0.6 : 1
                    }}
                  >
                    {membershipLoading ? 'Loading...' : (isMember ? 'Leave Group' : 'Join Group')}
                  </button>
                )
              )}
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
      
      {/* Confirmation Modal */}
      {showConfirmation && (
        <div className="fixed inset-0 flex items-center justify-center z-50" style={{ backgroundColor: 'rgba(0, 0, 0, 0.5)' }}>
          <div className="rounded-lg p-6 max-w-md w-full mx-4" style={{ backgroundColor: 'var(--primary-background)' }}>
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>
              {confirmationAction === 'join' ? 'Join Group' : 'Leave Group'}
            </h3>
            <p className="mb-6" style={{ color: 'var(--secondary-text)' }}>
              {confirmationAction === 'join' 
                ? `Are you sure you want to join "${group.title}"?`
                : `Are you sure you want to leave "${group.title}"?`
              }
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setShowConfirmation(false)}
                className="px-4 py-2 rounded-lg"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
              >
                Cancel
              </button>
              <button
                onClick={handleMembershipAction}
                className="px-4 py-2 rounded-lg"
                style={{
                  backgroundColor: confirmationAction === 'join' ? 'var(--primary-accent)' : 'var(--warning-color)',
                  color: 'white'
                }}
              >
                {confirmationAction === 'join' ? 'Join' : 'Leave'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
  
  async function handleMembershipAction() {
    setMembershipLoading(true);
    setShowConfirmation(false);
    
    try {
      if (confirmationAction === 'join') {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/join-request`, {
          method: 'POST',
          credentials: 'include'
        });
        
        if (response.ok) {
          setIsMember(true);
          // Refresh stats
          const statsRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/stats`, { credentials: 'include' });
          if (statsRes.ok) {
            const newStats = await statsRes.json();
            setStats(newStats);
          }
        } else {
          const errorData = await response.json().catch(() => ({ message: 'Failed to join group' }));
          console.error('Join group error:', errorData);
          alert(errorData.message || 'Failed to join group');
        }
      } else {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/leave`, {
          method: 'DELETE',
          credentials: 'include'
        });
        
        if (response.ok) {
          setIsMember(false);
          // Refresh stats
          const statsRes = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/stats`, { credentials: 'include' });
          if (statsRes.ok) {
            const newStats = await statsRes.json();
            setStats(newStats);
          }
        } else {
          const errorData = await response.json().catch(() => ({ message: 'Failed to leave group' }));
          console.error('Leave group error:', errorData);
          alert(errorData.message || 'Failed to leave group');
        }
      }
    } catch (error) {
      console.error('Error updating membership:', error);
      alert('Network error occurred');
    } finally {
      setMembershipLoading(false);
    }
  }
}