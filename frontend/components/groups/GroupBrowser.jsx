'use client';

import { useState, useEffect } from 'react';
import { api } from '../../lib/api';

export default function GroupBrowser() {
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

  const fetchAllGroups = async () => {
    setLoading(true);
    try {
      const response = await api.get('/groups');
      setGroups(response.data || []);
    } catch (error) {
      console.error('Error fetching groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchMyGroups = async () => {
    setLoading(true);
    try {
      const response = await api.get('/groups/my-groups');
      setMyGroups(response.data || []);
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
      const response = await api.get(`/groups/search?query=${encodeURIComponent(searchQuery)}`);
      setSearchResults(response.data || []);
    } catch (error) {
      console.error('Error searching groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const joinGroup = async (groupId) => {
    try {
      await api.post(`/groups/${groupId}/join-request`);
      alert('Join request sent successfully!');
    } catch (error) {
      console.error('Error joining group:', error);
      alert('Failed to send join request');
    }
  };

  const GroupCard = ({ group, showJoinButton = true }) => (
    <div className="bg-white rounded-lg shadow p-4 mb-4">
      <h3 className="text-lg font-semibold mb-2">{group.title}</h3>
      <p className="text-gray-600 mb-3">{group.description}</p>
      <div className="flex justify-between items-center">
        <span className="text-sm text-gray-500">
          Created: {new Date(group.created_at).toLocaleDateString()}
        </span>
        {showJoinButton && (
          <button
            onClick={() => joinGroup(group.id)}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            Join Group
          </button>
        )}
      </div>
    </div>
  );

  return (
    <div className="max-w-4xl mx-auto p-4">
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
          />
          <button
            onClick={searchGroups}
            className="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600"
          >
            Search
          </button>
        </div>
      </div>

      {/* Tabs */}
      <div className="mb-6">
        <div className="flex border-b">
          <button
            onClick={() => setActiveTab('browse')}
            className={`px-4 py-2 ${activeTab === 'browse' ? 'border-b-2 border-blue-500 text-blue-500' : 'text-gray-500'}`}
          >
            Browse All Groups
          </button>
          <button
            onClick={() => setActiveTab('my-groups')}
            className={`px-4 py-2 ${activeTab === 'my-groups' ? 'border-b-2 border-blue-500 text-blue-500' : 'text-gray-500'}`}
          >
            My Groups
          </button>
        </div>
      </div>

      {/* Loading */}
      {loading && (
        <div className="text-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
        </div>
      )}

      {/* Search Results */}
      {searchQuery && searchResults.length > 0 && (
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-4">Search Results</h2>
          {searchResults.map(group => (
            <GroupCard key={group.id} group={group} />
          ))}
        </div>
      )}

      {/* Browse All Groups */}
      {activeTab === 'browse' && !loading && (
        <div>
          <h2 className="text-xl font-semibold mb-4">All Public Groups</h2>
          {groups.length === 0 ? (
            <p className="text-gray-500">No public groups found.</p>
          ) : (
            groups.map(group => (
              <GroupCard key={group.id} group={group} />
            ))
          )}
        </div>
      )}

      {/* My Groups */}
      {activeTab === 'my-groups' && !loading && (
        <div>
          <h2 className="text-xl font-semibold mb-4">My Groups</h2>
          {myGroups.length === 0 ? (
            <p className="text-gray-500">You haven't joined any groups yet.</p>
          ) : (
            myGroups.map(group => (
              <GroupCard key={group.id} group={group} showJoinButton={false} />
            ))
          )}
        </div>
      )}
    </div>
  );
}