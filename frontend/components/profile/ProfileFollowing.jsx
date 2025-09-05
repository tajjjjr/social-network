import React, { useState, useEffect } from 'react';
import { profileAPI } from '../../lib/api';
import { useRouter } from 'next/navigation';
import Image from 'next/image';

const ProfileFollowing = ({ user, currentUser, isOwnProfile }) => {
  const [following, setFollowing] = useState([]);

  const router = useRouter();

  useEffect(() => {
    const fetchFollowing = async () => {
      if (user && user.profile_details && user.profile_details.id) {
        try {
          const data = await profileAPI.getFollowing(user.profile_details.id);
          setFollowing(data.user || []);
        } catch (error) {
          console.error('Error fetching following:', error);
        }
      }
    };

    fetchFollowing();
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
        <h3 className="text-xl font-bold mb-4 text-white">Following ({following.length})</h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {following && following.length > 0 ? (
            following.map((followedUser) => (
              <div
                key={followedUser.follower_id}
                className="rounded-lg p-4 text-center cursor-pointer hover:opacity-80 transition-opacity"
                style={{ backgroundColor: 'var(--secondary-background)' }}
                onClick={() => handleViewProfile(followedUser.follower_id)}
              >
                <Image
                  src={profileAPI.fetchProfileImage(followedUser.avatar)}
                  alt={followedUser.firstname}
                  width={80}
                  height={80}
                  className="w-20 h-20 rounded-full mx-auto mb-3"
                />
                <h4 className="font-medium text-white text-sm">
                  {followedUser.firstname} {followedUser.lastname}
                </h4>
              </div>
            ))
          ) : (
            <p className="text-white">Not following anyone yet. Explore and connect with other users!</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileFollowing;