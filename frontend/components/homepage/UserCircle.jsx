import React from 'react';
import Image from 'next/image';
import { profileAPI } from '../../lib/api';


const UserCircle = ({ avatar, name, active = false, highlight = 'var(--primary-accent)' }) => {
  return (
    <div className="flex flex-col items-center gap-1 min-w-[60px]">
      <div className="rounded-full p-[2px]" style={{ background: highlight }}>
        <div className="rounded-full p-[2px]" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <Image src={profileAPI.fetchProfileImage(avatar)} alt="Profile" width={48} height={48} className="w-12 h-12 rounded-full object-cover" />
        </div>
      </div>
      <span className="text-xs" style={{ color: 'var(--primary-text)' }}>{name}</span>
    </div>
  );
};
export default UserCircle;