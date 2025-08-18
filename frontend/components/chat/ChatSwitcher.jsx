
"use client";

import React from 'react';

const ChatSwitcher = ({ currentView, setCurrentView, unreadCount }) => {
  const views = ["All Chats", "Unread", "Groups"];

  return (
    <div className="flex justify-center mb-8">
      <div className="flex rounded-lg p-1" style={{ backgroundColor: 'var(--secondary-background)' }}>
        {views.map((view) => (
          <button
            key={view}
            onClick={() => setCurrentView(view)}
            className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
              currentView === view
                ? 'text-white'
                : 'text-gray-400 hover:text-white'
            }`}
            style={{
              backgroundColor: currentView === view ? 'var(--primary-accent)' : 'transparent',
            }}
          >
            {view}
            {view === 'Unread' && unreadCount > 0 && (
              <span className="ml-2 px-2 py-0.5 bg-red-500 text-white text-xs font-bold rounded-full">
                {unreadCount}
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  );
};

export default ChatSwitcher;
