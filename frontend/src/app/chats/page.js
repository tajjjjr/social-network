'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '../../../lib/withAuth';
import Header from '../../../components/layout/Header';
import { chatAPI, profileAPI } from '../../../lib/api';
import { MessageCircleIcon, ArrowLeftIcon, UsersIcon, MessageSquarePlus } from 'lucide-react';
import Image from 'next/image';
import ChatSwitcher from '../../../components/chat/ChatSwitcher';
import UserListModal from '../../../components/chat/UserListModal';

const ChatsPage = ({ user }) => {
  const [messageableUsers, setMessageableUsers] = useState([]);
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [unreadCount, setUnreadCount] = useState(0);
  const [currentView, setCurrentView] = useState('All Chats'); // All Chats, Unread, Groups
  const [showUserListModal, setShowUserListModal] = useState(false);
  const router = useRouter();

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      let unread = { count: 0 };
      try {
        unread = await chatAPI.getUnreadChatCount();
      } catch (e) {
        console.error("Failed to fetch unread chat count:", e);
      }
      setUnreadCount(unread.count);

      if (currentView === 'All Chats') {
        const users = await chatAPI.getMessageableUsers();
        setMessageableUsers(users || []);
        const groupData = await chatAPI.getGroups();
        setGroups(groupData || []);
      } else if (currentView === 'Unread') {
        let unreadChats = [];
        try {
          unreadChats = await chatAPI.getUnreadChats();
        } catch (e) {
          console.error("Failed to fetch unread chats:", e);
        }
        setMessageableUsers(unreadChats || []);
        setGroups([]);
      } else if (currentView === 'Groups') {
        const groupData = await chatAPI.getGroups();
        setGroups(groupData || []);
        setMessageableUsers([]);
      }
    } catch (error) {
      console.error(`Failed to load data for ${currentView}:`, error);
    } finally {
      setLoading(false);
    }
  }, [currentView]);

  useEffect(() => {
    loadData();
  }, [currentView]);

  const handleUserSelect = (selectedUser) => {
    router.push(`/chats/${selectedUser.id}?nickname=${encodeURIComponent(selectedUser.nickname)}`);
  };

  const handleGroupSelect = (selectedGroup) => {
    router.push(`/chats/${selectedGroup.id}?name=${encodeURIComponent(selectedGroup.name)}&type=group`);
  };

  const handleBackToHome = () => {
    router.push('/me');
  };

  const renderContent = () => {
    if (loading) {
      return (
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p style={{ color: 'var(--primary-text)' }}>Loading...</p>
          </div>
        </div>
      );
    }

    const allChats = [...messageableUsers, ...groups];

    if (currentView === 'All Chats' && allChats.length > 0) {
      return (
        <div className="grid gap-4">
          {messageableUsers.map((messageableUser) => (
            <ChatItem key={`user-${messageableUser.id}`} item={messageableUser} type="user" onSelect={handleUserSelect} />
          ))}
          {groups.map((group) => (
            <ChatItem key={`group-${group.id}`} item={group} type="group" onSelect={handleGroupSelect} />
          ))}
        </div>
      );
    } else if (currentView === 'Groups' && groups.length > 0) {
      return (
        <div className="grid gap-4">
          {groups.map((group) => (
            <ChatItem key={`group-${group.id}`} item={group} type="group" onSelect={handleGroupSelect} />
          ))}
        </div>
      );
    } else {
      return (
        <div className="text-center py-12">
          <MessageCircleIcon className="w-16 h-16 mx-auto mb-4 opacity-50" style={{ color: 'var(--secondary-text)' }} />
          <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--primary-text)' }}>No chats available</h3>
          <p style={{ color: 'var(--secondary-text)' }}>
            {currentView === 'Unread' ? 'You have no unread messages.' : 'Start a new conversation or create a group.'}
          </p>
        </div>
      );
    }
  };

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
      <Header user={user} />
      <div className="max-w-4xl mx-auto p-6">
        <div className="flex items-center mb-6">
          <button
            onClick={handleBackToHome}
            className="flex items-center gap-2 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
            style={{ color: 'var(--primary-text)' }}
          >
            <ArrowLeftIcon className="w-5 h-5" />
            Back to Home
          </button>
        </div>
        <div className="flex items-center gap-3 mb-8">
          <MessageCircleIcon className="w-8 h-8" style={{ color: 'var(--primary-accent)' }} />
          <h1 className="text-3xl font-bold" style={{ color: 'var(--primary-text)' }}>
            Chats
          </h1>
        </div>
        <ChatSwitcher currentView={currentView} setCurrentView={setCurrentView} unreadCount={unreadCount} />
        {renderContent()}
      </div>

      {/* Start Conversation Button */}
      <button
        className="fixed bottom-8 right-8 bg-blue-500 text-white p-4 rounded-full shadow-lg hover:bg-blue-600 transition-all duration-300 group"
        onClick={() => setShowUserListModal(true)}
      >
        <div className="flex items-center">
          <MessageSquarePlus className="w-6 h-6" />
          <span className="hidden group-hover:inline ml-2">Start Conversation</span>
        </div>
      </button>

      {showUserListModal && (
        <UserListModal user={user} onClose={() => setShowUserListModal(false)} />
      )}
    </div>
  );
};

const ChatItem = ({ item, type, onSelect }) => {
  const isGroup = type === 'group';
  const avatar = isGroup ? chatAPI.fetchGroupImage(item.avatar || '') : profileAPI.fetchProfileImage(item.avatar || '');
  const name = isGroup ? item.name : item.nickname;

  return (
    <div
      onClick={() => onSelect(item)}
      className="flex items-center gap-4 p-4 rounded-lg border cursor-pointer hover:shadow-md transition-all duration-200 hover:scale-[1.02]"
      style={{
        backgroundColor: 'var(--secondary-background)',
        borderColor: 'var(--border-color)',
      }}
    >
      <Image
        src={avatar}
        alt={`${name}'s avatar`}
        width={48}
        height={48}
        className="w-12 h-12 rounded-full object-cover"
        onError={(e) => {
          e.target.src = "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNGM0Y0RjYiLz4KPGNpcmNsZSBjeD0iMjAiIGN5PSIxNiIgcj0iNiIgZmlsbD0iIzlDQTNBRiIvPgo8cGF0aCBkPSJNMzIgMzJDMzIgMjYuNDc3MiAyNy41MjI4IDIyIDIyIDIySDE4QzEyLjQ3NzIgMjIgOCAyNi40NzcyIDggMzJWMzJIMzJWMzJaIiBmaWxsPSIjOUNBM0FGIi8+Cjwvc3ZnPgo=";
        }}
      />
      <div className="flex-1">
        <h3 className="font-semibold text-lg" style={{ color: 'var(--primary-text)' }}>
          {name}
        </h3>
        <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
          {isGroup ? 'Group Chat' : 'Click to start chatting'}
        </p>
      </div>
      {isGroup ? (
        <UsersIcon className="w-6 h-6" style={{ color: 'var(--primary-accent)' }} />
      ) : (
        <MessageCircleIcon className="w-6 h-6" style={{ color: 'var(--primary-accent)' }} />
      )}
    </div>
  );
};

export default withAuth(ChatsPage);
