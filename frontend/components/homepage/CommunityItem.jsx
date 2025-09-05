import React from 'react';
import Image from 'next/image';

const CommunityItem = ({ icon, name, memberCount, onClick }) => {
  return (
    <div className="flex items-center gap-3 cursor-pointer hover:opacity-80 transition-opacity" onClick={onClick}>
      <div
        className="w-10 h-10 rounded-full overflow-hidden flex items-center justify-center"
        style={{ backgroundColor: 'var(--tertiary-text)' }}
      >
        {icon ? (
          <Image src={icon} alt={name} width={40} height={40} className="w-full h-full object-cover" />
        ) : (
          <span className="text-xs font-bold" style={{ color: 'var(--primary-text)' }}>
            {name.charAt(0).toUpperCase()}
          </span>
        )}
      </div>
      <div className="flex-1">
        <p className="text-sm font-medium" style={{ color: 'var(--primary-text)' }}>{name}</p>
        {memberCount > 0 && (
          <p className="text-xs" style={{ color: 'var(--primary-accent)' }}>
            • {memberCount} members
          </p>
        )}
      </div>
    </div>
  );
};

export default CommunityItem;