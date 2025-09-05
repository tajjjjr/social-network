'use client';
import React, { useState, useEffect } from 'react';
import { SearchIcon, PlusIcon } from 'lucide-react';
import CommunityItem from '../homepage/CommunityItem';
import { profileAPI, api, fetchGroupImage } from '../../lib/api';
import { useRouter } from 'next/navigation';
import Image from 'next/image';

const ProfileSidebar = ({ profile, connectionStatus = 'disconnected' }) => {
  const router = useRouter();
  const [userGroups, setUserGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  const handleProfileClick = () => {
    router.push('/profile');
  };

  useEffect(() => {
    const fetchUserGroups = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/my-groups`, {
          credentials: 'include'
        });
        if (response.ok) {
          const data = await response.json();
          setUserGroups(data || []);
        } else {
          setUserGroups([]);
        }
      } catch (error) {
        console.error('Error fetching user groups:', error);
        setUserGroups([]);
      } finally {
        setLoading(false);
      }
    };

    fetchUserGroups();

    // Listen for group creation and join events
    const handleGroupCreated = () => {
      fetchUserGroups();
    };
    const handleGroupJoined = () => {
      fetchUserGroups();
    };
    window.addEventListener('groupCreated', handleGroupCreated);
    window.addEventListener('groupJoined', handleGroupJoined);

    return () => {
      window.removeEventListener('groupCreated', handleGroupCreated);
      window.removeEventListener('groupJoined', handleGroupJoined);
    };
  }, []);
  
  return <div className="w-72 flex flex-col gap-6">
      {/* Profile Card */}
      <div
        className="rounded-xl p-4 flex flex-col items-center"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <div className="relative">
          <div
            className="w-24 h-24 rounded-full flex items-center justify-center"
            style={{ backgroundColor: 'var(--primary-accent)' }}
          >
            <Image
              src={profileAPI.fetchProfileImage(profile.avatar? profile.avatar : '')}
              alt="Profile"
              width={80}
              height={80}
              className="w-20 h-20 rounded-full"
            />
    
          </div>

           <div className="mt-2 text-center">
          <span className="inline-block px-3 py-1 rounded-full text-xs font-semibold" style={{ backgroundColor: 'var(--tertiary-text)', color: 'var(--primary-text)' }}>
                 {profile.profile ? 'Public Account' : 'Private Account'}
           </span>
           </div>


          {/* WebSocket Connection Status Indicator */}
          <div
            className={`absolute top-0 right-0 w-4 h-4 rounded-full border-2 border-white ${
              connectionStatus === 'connected' ? 'bg-green-500' :
              connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
            }`}
            title={`WebSocket ${connectionStatus}`}
          />

          <button
            className="absolute bottom-0 right-0 p-1 rounded-full"
            style={{ backgroundColor: 'var(--tertiary-text)' }}
          >
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
              <path
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                stroke="white"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
        </div>

        <div className="flex justify-between w-full mt-4">
          <div className="text-center">
            <div className="text-xl font-bold">{profile.numberoffollowers}</div>
            <div className="text-xs" style={{ color: 'var(--secondary-text)' }}>
              Followers
            </div>
          </div>
          <div className="text-center">
            <div className="text-xl font-bold">{profile.numberoffollowees}</div>
            <div className="text-xs" style={{ color: 'var(--secondary-text)' }}>
              Following
            </div>
          </div>
        </div>

        <div className="mt-4 text-center">
          <h2 className="text-lg font-bold">{profile.firstname} {profile.lastname}</h2>
          <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
            @{profile.nickname}
          </p>
        </div>

        <div className="mt-4 text-center text-sm">
          <p>✨ {profile.about_me} ✨</p>
          {/* TODO: Replace with profile status from backend */}
        </div>

        <button
          className="mt-4 w-full py-2 rounded-lg text-sm border cursor-pointer"
          style={{
            backgroundColor: 'var(--primary-background)',
            borderColor: 'var(--tertiary-text)',
            color: 'var(--primary-text)',
          }}
          onClick={handleProfileClick}
        >
          My Profile
        </button>
      </div>

      {/* Communities Card */}
      <div
        className="rounded-xl p-4"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <div className="flex justify-between items-center mb-3">
          <h3 className="font-bold">Groups</h3>
          <div className="flex gap-2">
            <button
              className="p-1.5 rounded-full"
              style={{ backgroundColor: 'transparent' }}
              onMouseOver={e => (e.currentTarget.style.backgroundColor = 'var(--tertiary-text)')}
              onMouseOut={e => (e.currentTarget.style.backgroundColor = 'transparent')}
              onClick={() => router.push('/groups')}
            >
              <SearchIcon className="w-4 h-4 cursor-pointer" />
            </button>
            <button
              className="p-1.5 rounded-full"
              style={{ backgroundColor: 'transparent' }}
              onMouseOver={e => (e.currentTarget.style.backgroundColor = 'var(--tertiary-text)')}
              onMouseOut={e => (e.currentTarget.style.backgroundColor = 'transparent')}
              onClick={() => router.push('/groups')}
            >
              <PlusIcon className="w-4 h-4 cursor-pointer" />
            </button>
          </div>
        </div>

        <div className="flex flex-col gap-3">
          {loading ? (
            <div className="text-center text-sm" style={{ color: 'var(--secondary-text)' }}>
              Loading groups...
            </div>
          ) : userGroups.length > 0 ? (
            userGroups.map(group => (
              <CommunityItem 
                key={group.id}
                icon={group.avatar && group.avatar.trim() !== '' ? fetchGroupImage(group.avatar) : null} 
                name={group.title} 
                memberCount={0}
                onClick={() => router.push(`/groups/${group.id}`)}
              />
            ))
          ) : (
            <div className="text-center text-sm" style={{ color: 'var(--secondary-text)' }}>
              You have not joined any groups yet
            </div>
          )}
        </div>
      </div>
    </div>;
};

export default ProfileSidebar;
