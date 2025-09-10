'use client';

import { useState, useEffect } from 'react';
import { api, fetchGroupImage } from '../../lib/api';
import { useRouter } from 'next/navigation';

export default function GroupBrowser({ user }) {
  const router = useRouter();
  const [groups, setGroups] = useState([]);
  const [myGroups, setMyGroups] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [activeTab, setActiveTab] = useState('browse');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (activeTab === 'browse') {
      fetchAllGroups();
    } else if (activeTab === 'my-groups') {
      fetchMyGroups();
    }
  }, [activeTab]);

  useEffect(() => {
    // Listen for group creation events
    const handleGroupCreated = () => {
      if (activeTab === 'my-groups') {
        fetchMyGroups();
      }
      fetchAllGroups();
    };
    window.addEventListener('groupCreated', handleGroupCreated);

    return () => {
      window.removeEventListener('groupCreated', handleGroupCreated);
    };
  }, [activeTab]);

  const fetchAllGroups = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups`, {
        credentials: 'include'
      });
      if (response.ok) {
        const data = await response.json();
        setGroups(data || []);
      }
    } catch (error) {
      console.error('Error fetching groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMyGroups = async () => {
    setLoading(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/my-groups`, {
        credentials: 'include'
      });
      if (response.ok) {
        const data = await response.json();
        setMyGroups(data || []);
      }
    } catch (error) {
      console.error('Error fetching my groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const searchGroups = async () => {
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }
    
    setLoading(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/search?query=${encodeURIComponent(searchQuery)}`, {
        credentials: 'include'
      });
      if (response.ok) {
        const data = await response.json();
        setSearchResults(data || []);
      }
    } catch (error) {
      console.error('Error searching groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const joinGroup = async (groupId) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${groupId}/join-request`, {
        method: 'POST',
        credentials: 'include'
      });
      if (response.ok) {
        // Refresh groups list
        fetchAllGroups();
        fetchMyGroups();
        // Dispatch event to refresh sidebar
        window.dispatchEvent(new CustomEvent('groupJoined'));
      }
    } catch (error) {
      console.error('Error joining group:', error);
    }
  };

  const GroupCard = ({ group, showJoinButton = true }) => {
    const isCreator = user && group.creator_id === user.id;
    return (
      <div className="rounded-lg shadow p-4 mb-4 cursor-pointer hover:opacity-90 transition-opacity" 
           style={{ backgroundColor: 'var(--secondary-background)' }}
           onClick={() => router.push(`/groups/${group.public_id}`)}>
        <div className="flex items-center gap-4 mb-2">
          <img
            src={group.avatar && group.avatar.trim() !== '' ? fetchGroupImage(group.avatar) : '/default-group-avatar.png'}
            alt={group.title + ' avatar'}
            className="w-12 h-12 rounded-full object-cover border"
            onError={e => { e.target.src = '/default-group-avatar.png'; }}
          />
          <h3 className="text-lg font-semibold" style={{ color: 'var(--primary-text)' }}>{group.title}</h3>
        </div>
        <p className="mb-3" style={{ color: 'var(--secondary-text)' }}>{group.description}</p>
        <div className="flex justify-between items-center">
          <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>
            Created: {new Date(group.created_at).toLocaleDateString()}
            {isCreator && <span className="ml-2 text-xs" style={{ color: 'var(--primary-accent)' }}>(Creator)</span>}
          </span>
          <div className="flex gap-2">
            {showJoinButton && !isCreator && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  joinGroup(group.public_id);
                }}
                className="px-4 py-2 rounded"
                style={{
                  backgroundColor: 'var(--primary-accent)',
                  color: 'white'
                }}
              >
                Join Group
              </button>
            )}
            {!showJoinButton && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  router.push(`/groups/${group.public_id}`);
                }}
                className="px-4 py-2 rounded border"
                style={{
                  backgroundColor: 'var(--secondary-background)',
                  borderColor: 'var(--tertiary-text)',
                  color: 'var(--primary-text)'
                }}
              >
                View Group
              </button>
            )}
          </div>
        </div>
      </div>
    );
  };

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Groups</h1>
      
      {/* Search Bar */}
      <div className="mb-6">
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="Search groups..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="flex-1 px-4 py-2 border rounded-lg"
            style={{ 
              backgroundColor: 'var(--secondary-background)', 
              borderColor: 'var(--tertiary-text)',
              color: 'var(--primary-text)'
            }}
          />
          <button
            onClick={searchGroups}
            className="px-6 py-2 rounded-lg"
            style={{
              backgroundColor: 'var(--primary-accent)',
              color: 'white'
            }}
          >
            Search
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="mb-6">
        <div className="flex border-b" style={{ borderColor: 'var(--tertiary-text)' }}>
          <button
            onClick={() => setActiveTab('browse')}
            className="px-4 py-2"
            style={{
              borderBottom: activeTab === 'browse' ? '2px solid var(--primary-accent)' : 'none',
              color: activeTab === 'browse' ? 'var(--primary-accent)' : 'var(--secondary-text)'
            }}
          >
            Browse All Groups
          </button>
          <button
            onClick={() => setActiveTab('my-groups')}
            className="px-4 py-2"
            style={{
              borderBottom: activeTab === 'my-groups' ? '2px solid var(--primary-accent)' : 'none',
              color: activeTab === 'my-groups' ? 'var(--primary-accent)' : 'var(--secondary-text)'
            }}
          >
            My Groups
          </button>
        </div>
      </div>

      {/* Loading */}
      {loading && (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 mx-auto" style={{ borderColor: 'var(--primary-accent)' }}></div>
        </div>
      )}

      {/* Search Results */}
      {searchQuery && searchResults.length > 0 && (
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-4">Search Results</h2>
          {searchResults.map(group => (
            <GroupCard key={group.public_id} group={group} showJoinButton={true} user={user} joinGroup={joinGroup} />
          ))}
        </div>
      )}

      {/* Browse All Groups */}
      {activeTab === 'browse' && !loading && (
        <div>
          <h2 className="text-xl font-semibold mb-4">All Public Groups</h2>
          {groups.length === 0 ? (
            <p style={{ color: 'var(--secondary-text)' }}>No public groups found.</p>
          ) : (
            groups.map(group => (
              <GroupCard key={group.public_id} group={group} showJoinButton={true} user={user} joinGroup={joinGroup} />
            ))
          )}
        </div>
      )}

      {/* My Groups */}
      {activeTab === 'my-groups' && !loading && (
        <div>
          <h2 className="text-xl font-semibold mb-4">My Groups</h2>
          {myGroups.length === 0 ? (
            <p style={{ color: 'var(--secondary-text)' }}>You haven&apos;t joined any groups yet.</p>
          ) : (
            myGroups.map(group => (
              <GroupCard key={group.public_id} group={group} showJoinButton={false} user={user} joinGroup={joinGroup} />
            ))
          )}
        </div>
      )}
    </div>
  );
}