import React from 'react';
import Image from 'next/image';

const ActivityItem = ({
  image,
  name,
  action,
  time,
  isGroup = false,
  isPartial = false,
  requestId,
  onAccept,
  onDecline,
}) => {
  if (isPartial) {
    return (
      <div className="flex items-center gap-3">
        <Image src={image} alt={name} width={40} height={40} className="w-10 h-10 rounded-full" />
        <div className="flex-1">
          <div className="flex items-center gap-1">
            <span className="font-medium text-sm" style={{ color: 'var(--primary-text)' }}>{name}</span>
            <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>{action}</span>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-start gap-3">
      <div className="relative">
        <Image src={image} alt={name} width={40} height={40} className="w-10 h-10 rounded-full" />
        {isGroup && (
          <div
            className="absolute -top-1 -right-1 w-4 h-4 rounded-full flex items-center justify-center text-[8px]"
            style={{ backgroundColor: 'var(--warning-color)', color: 'white' }}
          >
            3
          </div>
        )}
      </div>
      <div className="flex-1">
        <div className="flex items-center gap-1">
          <span className="font-medium text-sm" style={{ color: 'var(--primary-text)' }}>{name}</span>
          <span className="text-sm" style={{ color: 'var(--secondary-text)' }}>{action}</span>
        </div>
        {time && <p className="text-xs" style={{ color: 'var(--secondary-text)' }}>{time}</p>}
        <div className="flex justify-between items-center mt-2 gap-2">
          <button
            className="px-4 py-1 rounded-full text-xs flex-1"
            style={{ backgroundColor: 'var(--primary-accent)', color: 'black' }}
            onClick={() => onAccept(requestId)} // Pass requestId to onAccept
          >
            Accept
          </button>
          <button
            className="px-4 py-1 rounded-full text-xs flex-1 border"
            style={{
              backgroundColor: 'var(--primary-background)',
              borderColor: 'var(--tertiary-text)',
              color: 'var(--primary-text)',
            }}
            onClick={() => onDecline(requestId)} // Pass requestId to onDecline
          >
            Decline
          </button>
        </div>
      </div>
    </div>
  );
};
export default ActivityItem;