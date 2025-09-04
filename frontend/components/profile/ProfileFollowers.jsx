import React, { useState, useEffect } from 'react';
import { profileAPI } from '../../lib/api';
import { useRouter } from 'next/navigation';
import Image from 'next/image';

const ProfileFollowers = ({ user, currentUser, isOwnProfile }) => {
  const [followers, setFollowers] = useState([]);

  const router = useRouter();

  useEffect(() => {
    const fetchFollowers = async () => {
      if (user && user.profile_details && user.profile_details.id) {
        try {
          const data = await profileAPI.getFollowers(user.profile_details.id);
          setFollowers(data.user || []);
        } catch (error) {
          console.error('Error fetching followers:', error);
        }
      }
    };

    
    fetchFollowers();
  }, [user]);

  // Handle view profile
  const handleViewProfile = (userId) => {
    router.push(`/profile/${userId}`);
  };


  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Followers ({followers.length})</h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {followers && followers.length > 0 ? (
            followers.map((follower) => (
              <div
                key={follower.follower_id}
                className="rounded-lg p-4 text-center cursor-pointer hover:opacity-80 transition-opacity"
                style={{ backgroundColor: 'var(--secondary-background)' }}
                onClick={() => handleViewProfile(follower.follower_id)}
              >
                <Image
                  src={profileAPI.fetchProfileImage(follower.avatar)}
                  alt={follower.firstname}
                  width={80}
                  height={80}
                  className="w-20 h-20 rounded-full mx-auto mb-3"
                />
                <h4 className="font-medium text-white text-sm">
                  {follower.firstname} {follower.lastname}
                </h4>
                <p className="text-xs mt-1" style={{ color: 'var(--secondary-text)' }}>
                  {follower.mutualFollowers} mutual followers
                </p>
              </div>
            ))
          ) : (
            <p className="text-white">No followers yet. Start sharing your content to attract followers!</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileFollowers;