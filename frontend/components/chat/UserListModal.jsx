import React, { useState, useEffect } from 'react';
import { chatAPI, profileAPI } from '../../lib/api';
import { useRouter } from 'next/navigation';
import { X } from 'lucide-react';

const UserListModal = ({ user, onClose }) => {
  const [followingUsers, setFollowingUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const router = useRouter();

  useEffect(() => {
    const fetchMessageableUsers = async () => {
      try {
        setLoading(true);
        const response = await chatAPI.getMessageableUsers();
        if (response) {
          // Sort followers by nickname in ascending order
          const sortedUsers = response.sort((a, b) => a.nickname.localeCompare(b.nickname));
          setFollowingUsers(sortedUsers);
        } else {
          setError('Failed to fetch users');
        }
      } catch (err) {
        setError('Network error: ' + err.message);
      } finally {
        setLoading(false);
      }
    };

    if (user?.id) {
      fetchMessageableUsers();
    }
  }, [user]);

  const handleUserClick = (selectedUser) => {
    router.push(`/chats/${selectedUser.id}?nickname=${encodeURIComponent(selectedUser.nickname)}`);
    onClose(); // Close the modal after navigating
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50">
        <div className="bg-white p-6 rounded-lg shadow-lg w-96 text-center">
          <p>Loading users...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="fixed inset-0 bg-gray-600 bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50">
        <div className="bg-white p-6 rounded-lg shadow-lg w-96 text-center">
          <p className="text-red-500">Error: {error}</p>
          <button onClick={onClose} className="mt-4 px-4 py-2 bg-blue-500 text-white rounded">Close</button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-96 max-h-[80vh] flex flex-col"
        style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold">Start a Conversation</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <X className="w-6 h-6" />
          </button>
        </div>
        <div className="flex-1 overflow-y-auto">
          {followingUsers.length === 0 ? (
            <p className="text-center text-gray-500">You are not following any users yet.</p>
          ) : (
            <ul>
              {followingUsers.map((followedUser) => (
                <li key={followedUser.id} className="mb-2">
                  <button
                    onClick={() => handleUserClick(followedUser)}
                    className="w-full text-left p-3 rounded-lg hover:bg-gray-100 transition-colors flex items-center gap-3"
                    style={{ backgroundColor: 'var(--primary-background)', color: 'var(--primary-text)' }}
                  >
                    <img
                      src={profileAPI.fetchProfileImage(followedUser.avatar || '')}
                      alt={`${followedUser.nickname}'s avatar`}
                      className="w-8 h-8 rounded-full object-cover"
                    />
                    <span>{followedUser.nickname}</span>
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
};

export default UserListModal;
