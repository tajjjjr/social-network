import React, { useState } from 'react';
import Header from '../layout/Header';
import ProfileSidebar from '../layout/ProfileSidebar';
import Feed from '../layout/Feed';
import ActivitySidebar from '../layout/ActivitySidebar';
import UserListModal from '../chat/UserListModal';

const MainHomepage = ({ user, profile, connectionStatus, connectedUsers = [] }) => {
  const [showUserListModal, setShowUserListModal] = useState(false);
  return (
    <div className="w-2/3 flex flex-col text-white">
      <Header user={user} />
      <div className="flex flex-1 w-full max-w-7xl mx-auto gap-4 p-4">
        <ProfileSidebar profile={profile} connectionStatus={connectionStatus} />
        <div className="flex-1 flex flex-col">
          <Feed user={user} connectedUsers={connectedUsers} />
        </div>
        <ActivitySidebar user={user} />
      </div>
    </div>
  );
};

export default MainHomepage; 