import React from 'react';

const EndOfPosts = () => {
  return (
    <div className="text-center py-8 border-t" style={{ borderColor: 'var(--border-color)' }}>
      <div className="mb-2" style={{ color: 'var(--secondary-text)' }}>
        You&apos;ve reached the end
      </div>
      <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>
        No more posts to show
      </div>
    </div>
  );
};

export default EndOfPosts;