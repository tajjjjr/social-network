'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import GroupHeader from '../../../../components/groups/GroupHeader';
import GroupSidebar from '../../../../components/groups/GroupSidebar';
import GroupPosts from '../../../../components/groups/GroupPosts';
import withAuth from '../../../../lib/withAuth';

function GroupPage({ user }) {
  const params = useParams();
  const groupId = params.groupid;
  const [group, setGroup] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchGroup = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${groupId}`, {
          credentials: 'include'
        });
        
        if (response.ok) {
          const groupData = await response.json();
          setGroup(groupData);
        } else {
          setError('Failed to load group');
        }
      } catch (error) {
        console.error('Error fetching group:', error);
        setError('Failed to load group');
      } finally {
        setLoading(false);
      }
    };

    if (groupId) {
      fetchGroup();
    }
  }, [groupId]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: 'var(--primary-accent)' }}></div>
      </div>
    );
  }

  if (error || !group) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold mb-4" style={{ color: 'var(--primary-text)' }}>
            {error || 'Group not found'}
          </h1>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
      <GroupHeader group={group} user={user} />
      
      <div className="max-w-7xl mx-auto px-4 py-6">
        <div className="flex gap-6">
          <div className="w-80 flex-shrink-0">
            <GroupSidebar group={group} />
          </div>
          
          <div className="flex-1">
            <GroupPosts groupId={groupId} user={user} group={group} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default withAuth(GroupPage);