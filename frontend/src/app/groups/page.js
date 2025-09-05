'use client';

import { useState } from 'react';
import withAuth from '../../../lib/withAuth';
import Header from '../../../components/layout/Header';
import ProfileSidebar from '../../../components/layout/ProfileSidebar';
import ActivitySidebar from '../../../components/layout/ActivitySidebar';
import GroupBrowser from '../../../components/groups/GroupBrowser';
import GroupCreation from '../../../components/groups/GroupCreation';

const GroupsPage = ({ user, profile }) => {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleGroupCreated = () => {
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <div className="w-2/3 flex flex-col text-white">
          <Header user={user} />
          <div className="flex flex-1 w-full max-w-7xl mx-auto gap-4 p-4">
            <ProfileSidebar profile={profile.profile_details} connectionStatus="disconnected" />
            <div className="flex-1 flex flex-col">
              <div className="bg-secondary rounded-lg p-6" style={{ color: 'var(--primary-text)' }}>
                <GroupCreation onGroupCreated={handleGroupCreated} />
                <GroupBrowser key={refreshKey} user={user} />
              </div>
            </div>
            <ActivitySidebar user={user} />
          </div>
        </div>
      </main>
    </div>
  );
};

export default withAuth(GroupsPage);