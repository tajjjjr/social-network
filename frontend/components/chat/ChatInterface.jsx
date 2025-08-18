'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { wsService } from '../../lib/websocket';
import { chatAPI } from '../../lib/api';
import { notificationService } from '../../lib/notificationService';
import ClientDate from '../common/ClientDate';
import Picker from 'emoji-picker-react';

const ChatInterface = ({ user, connectionStatus = 'disconnected', recipient = null, showSidebar = true }) => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  // { type: 'private', id: userId } or { type: 'group', id: groupId }
  const [activeChat, setActiveChat] = useState(null);
  const [activeChatName, setActiveChatName] = useState('');
  const [onlineUsers, setOnlineUsers] = useState(new Set());
  const [messageableUsers, setMessageableUsers] = useState([]);
  const [groups, setGroups] = useState([]);
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const messagesEndRef = useRef(null);
  const activeChatRef = useRef(activeChat);
  const userRef = useRef(user);

  useEffect(() => {
    activeChatRef.current = activeChat;
    userRef.current = user;
  }, [activeChat, user]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handlePrivateMessage = useCallback((message) => {
    console.log('Received private message:', message);

    // Use refs to get current values and avoid stale closures
    const currentActiveChat = activeChatRef.current;
    const currentUser = userRef.current;

    console.log('Current activeChat:', currentActiveChat);
    console.log('Message from:', message.from, 'to:', message.to);
    console.log('User ID:', currentUser.id);

    // A private message is relevant if the active chat is private
    // and the message is either from me to the active user, or from the active user to me.
    if (currentActiveChat?.type === 'private' &&
        ((message.from === currentUser.id && message.to === currentActiveChat.id) ||
         (message.from === currentActiveChat.id && message.to === currentUser.id))) {

      console.log('Message is relevant, adding to chat');

      // remove any optimistic message with the same content and timestamp (within 2 seconds)
      setMessages(prev => {
        const filteredMessages = prev.filter(msg => {
          if (msg.isOptimistic &&
              msg.content === message.content &&
              msg.from === message.from &&
              Math.abs(msg.timestamp - message.timestamp) < 2) {
            console.log('Removing optimistic message');
            return false; // Remove optimistic message
          }
          return true;
        });

        return [...filteredMessages, message];
      });
    } else {
      console.log('Message not relevant for current chat');
    }
  }, []);

  const handleGroupMessage = useCallback((message) => {
    const currentActiveChat = activeChatRef.current;
    if (currentActiveChat?.type === 'group' && message.group_id === currentActiveChat.id) {
      setMessages(prev => [...prev, message]);
    }
  }, []);

  const handleBroadcastMessage = useCallback((message) => {
    // Show broadcast messages in all chats
    setMessages(prev => [...prev, { ...message, isBroadcast: true }]);
  }, []);

  const handleNotification = useCallback((notification) => {
    // Pass notification to notification service for processing
    notificationService.handleNotification(notification);
  }, []);

  const handleUserConnected = useCallback((notification) => {
    setOnlineUsers(prev => new Set([...prev, notification.user_id]));
  }, []);

  const handleUserDisconnected = useCallback((notification) => {
    setOnlineUsers(prev => {
      const newSet = new Set(prev);
      newSet.delete(notification.user_id);
      return newSet;
    });
  }, []);

  useEffect(() => {
    // Set up message handlers only (connection is managed by parent)
    wsService.onMessage('private', handlePrivateMessage);
    wsService.onMessage('group', handleGroupMessage);
    wsService.onMessage('broadcast', handleBroadcastMessage);
    wsService.onMessage('notification', handleNotification);

    // Set up notification handlers for user connection tracking
    notificationService.onNotification('user_connected', handleUserConnected);
    notificationService.onNotification('user_disconnected', handleUserDisconnected);

    // Load initial online users
    loadOnlineUsers();
    loadMessageableUsers();
    loadGroups();

    // If recipient is provided, automatically load that chat
    if (recipient) {
      loadChatHistory(recipient.type, recipient.id);
    }

    return () => {
      // Clean up WebSocket message handlers
      wsService.removeHandler('private', handlePrivateMessage);
      wsService.removeHandler('group', handleGroupMessage);
      wsService.removeHandler('broadcast', handleBroadcastMessage);
      wsService.removeHandler('notification', handleNotification);

      // Clean up notification handlers
      notificationService.removeHandler('user_connected', handleUserConnected);
      notificationService.removeHandler('user_disconnected', handleUserDisconnected);
    };
  }, [recipient, handlePrivateMessage, handleGroupMessage, handleBroadcastMessage, handleNotification, handleUserConnected, handleUserDisconnected]);

  const loadOnlineUsers = async () => {
    try {
      const users = await notificationService.loadOnlineUsers();
      const userIds = new Set(users.map(user => user.user_id));
      setOnlineUsers(userIds);
    } catch (error) {
      console.error('Failed to load online users:', error);
    }
  };

  const loadMessageableUsers = async () => {
    try {
      const users = await chatAPI.getMessageableUsers();
      setMessageableUsers(users || []);
    } catch (error) {
      console.error('Failed to load messageable users:', error);
    }
  };

  const loadGroups = async () => {
    try {
      const groupData = await chatAPI.getGroups();
      setGroups(groupData || []);
    } catch (error) {
      console.error('Failed to load groups:', error);
    }
  };

  const sendMessage = () => {
    if (!newMessage.trim() || !activeChat) return;

    // Check if WebSocket is connected
    if (!wsService.isConnected()) {
      alert('Not connected to chat server. Please wait for connection to be established.');
      return;
    }

    // Create optimistic message for immediate UI update
    const optimisticMessage = {
      from: user.id,
      to: activeChat.type === 'private' ? activeChat.id : null,
      group_id: activeChat.type === 'group' ? activeChat.id : null,
      content: newMessage,
      timestamp: Date.now() / 1000, // Convert to seconds to match backend
      type: activeChat.type,
      isOptimistic: true // Flag to identify optimistic messages
    };

    // Immediately add message to UI (optimistic update)
    setMessages(prev => [...prev, optimisticMessage]);

    // Send message via WebSocket
    if (activeChat.type === 'private') {
      wsService.sendMessage('private', newMessage, activeChat.id);
    } else if (activeChat.type === 'group') {
      wsService.sendMessage('group', newMessage, null, activeChat.id);
    }

    setNewMessage('');
  };

  const onEmojiClick = (emojiObject) => {
    setNewMessage(prevInput => prevInput + emojiObject.emoji);
    setShowEmojiPicker(false);
  };

  const loadChatHistory = async (chatType, chatId) => {
    try {
      let history;
      let chatName = '';
      if (chatType === 'private') {
        history = await chatAPI.getPrivateMessages(chatId);
        const user = messageableUsers.find(u => u.id === chatId);
        if (user) {
          chatName = user.nickname;
        }
      } else if (chatType === 'group') {
        history = await chatAPI.getGroupMessages(chatId);
        const group = groups.find(g => g.id === chatId);
        if (group) {
          chatName = group.name;
        }
      }
      let finalHistory = history || [];

      // If it's a group chat and no history, add a welcome message
      if (chatType === 'group' && finalHistory.length === 0) {
        finalHistory = [{
          id: 'welcome',
          content: `Welcome to ${chatName}! This is the beginning of your group chat.`, 
          from: 'system',
          timestamp: Date.now() / 1000,
          isWelcome: true,
        }];
      }

      setMessages(finalHistory);
      setActiveChat({ type: chatType, id: chatId });
      setActiveChatName(chatName);
    } catch (error) {
      console.error('Failed to load chat history:', error);
      // If we get an error, still allow the chat to open with empty history
      // This enables users to start new conversations
      setMessages([]);
      setActiveChat({ type: chatType, id: chatId });
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <div className={`flex h-96 bg-white rounded-lg shadow-lg ${!showSidebar ? 'h-[500px]' : ''}`}>
      {/* Chat Sidebar - Only show if showSidebar is true */}
      {showSidebar && (
        <div className="w-1/3 border-r border-gray-200 p-4">
          <div className="flex items-center justify-between mb-4">
            <h3 className="font-semibold">Chats</h3>
            <div className={`w-3 h-3 rounded-full ${
              connectionStatus === 'connected' ? 'bg-green-500' :
              connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
            }`} title={`WebSocket ${connectionStatus}`}></div>
          </div>

          {/* Private Chats */}
          <div className="mb-4">
            <h4 className="text-sm text-gray-600 mb-2">Private Messages</h4>
            {messageableUsers.map((chatUser) => (
              <button
                key={chatUser.id}
                onClick={() => loadChatHistory('private', chatUser.id)}
                className="w-full text-left p-2 hover:bg-gray-100 rounded flex items-center justify-between"
              >
                <span>{chatUser.nickname}</span>
                <div className={`w-2 h-2 rounded-full ${
                  onlineUsers.has(chatUser.id) ? 'bg-green-500' : 'bg-gray-300'
                }`} title={onlineUsers.has(chatUser.id) ? 'Online' : 'Offline'}></div>
              </button>
            ))}
          </div>

          {/* Group Chats */}
          <div>
            <h4 className="text-sm text-gray-600 mb-2">Groups</h4>
            {groups.map((group) => (
              <button
                key={group.id}
                onClick={() => loadChatHistory('group', group.id)} // Example group ID
                className="w-full text-left p-2 hover:bg-gray-100 rounded"
              >
                {group.name}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Chat Messages */}
      <div className={`flex-1 flex flex-col ${!showSidebar ? 'w-full' : ''}`}>
        {activeChat ? (
          <>
            {/* Chat Header - Only show when no sidebar */}
            {!showSidebar && (
              <div className="border-b border-gray-200 p-3 bg-gray-50">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-semibold">
                    {recipient?.name?.charAt(0)?.toUpperCase() || 'U'}
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-800">
                      {recipient?.name || 'User'}
                    </h3>
                    <p className="text-xs text-gray-500">
                      {connectionStatus === 'connected' ? 'Online' : 'Offline'}
                    </p>
                  </div>
                </div>
              </div>
            )}

            {/* Messages Area */}
            <div className="flex-1 p-4 overflow-y-auto">
              {messages.map((message, index) => {
                const isSender = message.from === user.id;
                const isBroadcast = message.isBroadcast;
                const isOptimistic = message.isOptimistic;
                const isWelcome = message.isWelcome;

                if (isWelcome) {
                  return (
                    <div key={message.id} className="text-center text-gray-500 italic mb-4">
                      {message.content}
                    </div>
                  );
                }

                return (
                  <div
                    key={`${message.timestamp}-${message.from}-${index}`}
                    className={`flex mb-4 ${isSender ? 'justify-end' : 'justify-start'} ${isBroadcast ? 'justify-center' : ''} ${isOptimistic ? 'opacity-70' : ''}`}
                  >
                    {!isSender && !isBroadcast && <div className="w-8 h-8 rounded-full bg-gray-300 mr-3"></div> /* Avatar placeholder */}
                    <div className={isBroadcast ? 'text-center text-blue-600' : ''}>
                      <div className={`text-xs mb-1 ${isSender ? 'text-right' : 'text-left'} text-gray-500`}>
                        <ClientDate dateString={new Date(message.timestamp * 1000).toISOString()} format="time" />
                        {isOptimistic && <span className="ml-1 text-yellow-500">⏳</span>}
                      </div>
                      <div className={`p-3 rounded-lg ${isSender ? 'bg-blue-500 text-white' : 'bg-gray-200'} ${isOptimistic ? 'border-2 border-dashed border-yellow-400' : ''}`}>
                        {message.content}
                      </div>
                    </div>
                  </div>
                );
              })}
              <div ref={messagesEndRef} />
            </div>
            
            {/* Message Input */}
            <div className="p-4 border-t border-gray-200">
              {connectionStatus !== 'connected' && (
                <div className="mb-2 text-sm text-center">
                  <span className={`${
                    connectionStatus === 'connecting' ? 'text-yellow-600' : 'text-red-600'
                  }`}>
                    {connectionStatus === 'connecting' ? 'Connecting to chat...' : 'Disconnected from chat server'}
                  </span>
                </div>
              )}
              <div className="relative flex">
                <div className="relative flex-grow">
                  <input
                    type="text"
                    value={newMessage}
                    onChange={(e) => setNewMessage(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                    placeholder={connectionStatus === 'connected' ? "Type a message..." : "Waiting for connection..."}
                    disabled={connectionStatus !== 'connected'}
                    className={`w-full p-2 border border-gray-300 rounded-l ${
                      connectionStatus !== 'connected' ? 'bg-gray-100 cursor-not-allowed' : ''
                    }`}
                  />
                  <button
                    onClick={() => setShowEmojiPicker(val => !val)}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-xl"
                    title="Add emoji"
                    disabled={connectionStatus !== 'connected'}
                  >
                    😊
                  </button>
                </div>
                <button
                  onClick={sendMessage}
                  disabled={connectionStatus !== 'connected'}
                  className={`px-4 py-2 text-white rounded-r ${
                    connectionStatus === 'connected'
                      ? 'bg-blue-500 hover:bg-blue-600'
                      : 'bg-gray-400 cursor-not-allowed'
                  }`}
                >
                  Send
                </button>
              </div>
              {showEmojiPicker && (
                <div className="absolute bottom-20 right-4 z-10">
                  <Picker
                    onEmojiClick={onEmojiClick}
                    pickerStyle={{ width: '100%', boxShadow: 'none' }}
                  />
                </div>
              )}
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            {showSidebar ? 'Select a chat to start messaging' : 'Loading chat...'}
          </div>
        )}
      </div>
    </div>
  );
};

export default ChatInterface;