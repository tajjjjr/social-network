'use client';

import { useState } from 'react';
import GroupBrowser from '../../components/groups/GroupBrowser';
import GroupCreation from '../../components/groups/GroupCreation';

export default function GroupsPage() {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleGroupCreated = () => {
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto py-8">
        <GroupCreation onGroupCreated={handleGroupCreated} />
        <GroupBrowser key={refreshKey} />
      </div>
    </div>
  );
}