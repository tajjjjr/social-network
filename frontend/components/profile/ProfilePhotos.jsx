import React from 'react';
import { profileAPI } from '../../lib/api';
import Image from 'next/image';

const ProfilePhotos = ({ user, isOwnProfile }) => {
  var photos = user?.photos || [];
  console.log('ProfilePhotos rendered with user:', user, 'and photos:', user?.photos);

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Photos ({photos.length})</h3>
        {photos.length > 0 ? (
          <div className="grid grid-cols-3 gap-2">
            {photos.map((photo) => (
              <div
                key={photo.image}
                className="aspect-square rounded-lg overflow-hidden cursor-pointer hover:opacity-80 transition-opacity"
                style={{ backgroundColor: 'var(--secondary-background)' }}
              >
                <Image
                  src={profileAPI.fetchProfileImage(photo.image)}
                  alt={photo.image}
                  width={120}
                  height={120}
                  className="w-full h-full object-cover"
                />
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <p className="text-gray-400">{isOwnProfile ? "You have no photos." : "This user has no photos."}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default ProfilePhotos;