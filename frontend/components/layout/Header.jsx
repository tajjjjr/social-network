import React, { useState, useRef, useEffect } from 'react';
import { HomeIcon, BellIcon, UsersIcon, MessageCircleIcon, SearchIcon, ChevronDownIcon, LogOutIcon } from 'lucide-react';
import { handleLogout } from '../../lib/auth';
import { useSimpleNotifications } from '../../hooks/useNotifications';
import { profileAPI } from '../../lib/api';
import ClientDate from '../common/ClientDate';
import { useRouter } from 'next/navigation';

const Header = ({ user = null }) => {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [notificationPanelOpen, setNotificationPanelOpen] = useState(false);
  const profileRef = useRef(null);
  const notificationRef = useRef(null);
  const { unreadCount, notifications, markAllAsRead } = useSimpleNotifications();
  const router = useRouter();

  // Close dropdowns when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (profileRef.current && !profileRef.current.contains(event.target)) {
        setDropdownOpen(false);
      }
      if (notificationRef.current && !notificationRef.current.contains(event.target)) {
        setNotificationPanelOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="flex items-center justify-between w-full px-4 py-2 mt-4 relative">
      
      {/* Logo and Search Bar */}
      <div className="flex items-center gap-4 z-10">
        <div className="bg-white rounded-full p-1 cursor-pointer">
          <div className="w-6 h-6 rounded-full" style={{ backgroundColor: 'var(--primary-background)' }}></div>
        </div>
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
          <input
            type="text"
            placeholder="Search for people, communities ..."
            className="rounded-full py-1.5 pl-10 pr-4 text-sm focus:outline-none"
            style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
          />
        </div>
      </div>

      {/* Navigation Icons - Centered absolutely */}
      <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 flex items-center gap-6 z-0">
        {/* Home Icon */}
        <div
          className="flex flex-col items-center cursor-pointer"
          onClick={() => router.push('/me')}
        >
          <HomeIcon className="w-6 h-6" style={{ color: 'var(--primary-accent)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Home</span>
        </div>

        {/* Notifications Icon */}
        <div className="relative" ref={notificationRef}>
          <div
            className="flex flex-col items-center cursor-pointer"
            onClick={() => {
              setNotificationPanelOpen(!notificationPanelOpen);
              if (!notificationPanelOpen && unreadCount > 0) {
                markAllAsRead();
              }
            }}
          >
            <div className="relative">
              <BellIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
              {unreadCount > 0 && (
                <div className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                  {unreadCount > 9 ? '9+' : unreadCount}
                </div>
              )}
            </div>
            <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Notifications</span>
          </div>

          {/* Notification Panel */}
          {notificationPanelOpen && (
            <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 z-50">
              <div className="p-4 border-b border-gray-200">
                <h3 className="font-semibold text-gray-800">Notifications</h3>
              </div>

              <div className="max-h-96 overflow-y-auto">
                {notifications.length > 0 ? (
                  notifications.map((notification, index) => (
                    <div
                      key={index}
                      className="p-4 border-b border-gray-100 hover:bg-gray-50"
                    >
                      <div className="text-sm text-gray-800">{notification.message}</div>
                      <div className="text-xs text-gray-500 mt-1">
                        <ClientDate dateString={new Date(notification.timestamp * 1000).toISOString()} />
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="p-4 text-center text-gray-500">
                    No notifications
                  </div>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Groups Icon (Disabled) */}
        <div className="flex flex-col items-center cursor-not-allowed opacity-50">
          <UsersIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Groups</span>
        </div>

        {/* Chats Icon */}
        <div
          className="flex flex-col items-center cursor-pointer"
          onClick={() => router.push('/chats')}
        >
          <MessageCircleIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Chats</span>
        </div>
      </div>

      {/* Profile Dropdown */}
      <div className="relative z-10" ref={profileRef}>
        <div className="flex items-center gap-2 cursor-pointer" onClick={() => router.push('/me')}>
          <img src={profileAPI.fetchProfileImage(user.avatar ? user.avatar : '')} alt="Profile" className="w-8 h-8 rounded-full" />
          <span className="text-sm font-medium">{user.nickname}</span>
        </div>
        <button
          className="flex items-center gap-2 focus:outline-none"
          onClick={() => setDropdownOpen((open) => !open)}
        >
          <ChevronDownIcon className="w-4 h-4" />
        </button>
        {dropdownOpen && (
          <div className="absolute right-0 mt-2 w-40 bg-white rounded shadow-lg z-10">
            <button
              className="flex items-center gap-2 px-4 py-2 w-full text-left bg-white rounded transition-colors cursor-pointer hover:bg-blue-100"
              style={{ color: 'var(--primary-background)' }}
              onClick={handleLogout}
            >
              <LogOutIcon className="w-4 h-4" />
              Logout
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Header;
