import React, { useState, useRef, useEffect } from 'react';
import { CameraIcon, EditIcon, UserPlusIcon, X, Save } from 'lucide-react';
import { profileAPI, updateProfile } from '../../lib/api';
import { wsService } from '../../lib/websocket';
import Image from 'next/image';

const ProfileCover = ({ user, currentUser, isOwnProfile, refreshProfile }) => {
  const profileDetails = user?.profile_details || {};
  const [showEditForm, setShowEditForm] = useState(false);
  const [formData, setFormData] = useState({
    firstname: profileDetails.firstname || '',
    lastname: profileDetails.lastname || '',
    nickname: profileDetails.nickname || '',
    email: profileDetails.email || '',
    aboutme: profileDetails.aboutme || '',
    dateofbirth: profileDetails.dateofbirth || '',
    is_private: profileDetails.is_private || false,
    profilePicture: null,
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [previewImage, setPreviewImage] = useState(null);
  const [errorMessage, setErrorMessage] = useState('');
  const fileInputRef = useRef(null);

  // Listen for follow request rejection notifications to refresh profile
  useEffect(() => {
    if (!isOwnProfile && currentUser?.id) {
      const handleFollowRejection = (notification) => {
        // Check if this notification is about a rejection from the current profile user
        if (notification.subtype === 'follow_request_rejected' &&
            notification.user_id === profileDetails.id) {
          // Refresh the profile to update button status
          refreshProfile();
        }
      };

      wsService.onMessage('notification', handleFollowRejection);

      return () => {
        wsService.removeHandler('notification', handleFollowRejection);
      };
    }
  }, [isOwnProfile, currentUser?.id, profileDetails.id, refreshProfile]);

  const handleFollow = async () => {
    try {
      if (profileDetails.followbtnstatus === 'following') {
        await profileAPI.unfollow(profileDetails.id);
      } else if (profileDetails.followbtnstatus === 'pending') {
        // Cancel the pending follow request
        const requestResponse = await profileAPI.getFollowRequestId(profileDetails.id);
        if (requestResponse && requestResponse.request_id) {
          await profileAPI.cancelFollowRequest(requestResponse.request_id);
        }
      } else {
        await profileAPI.follow(profileDetails.id);
      }
      refreshProfile();
    } catch (error) {
      console.error('Error following/unfollowing user:', error);
    }
  };

  const handleEditProfile = () => {
    setFormData({
      firstname: profileDetails.firstname || '',
      lastname: profileDetails.lastname || '',
      nickname: profileDetails.nickname || '',
      email: profileDetails.email || '',
      aboutme: profileDetails.aboutme || '',
      dateofbirth: profileDetails.dateofbirth || '',
      is_private: profileDetails.is_private || false,
      profilePicture: null,
    });
    setPreviewImage(null);
    setShowEditForm(true);
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setFormData(prev => ({
        ...prev,
        profilePicture: file
      }));
      
      // Create preview URL
      const reader = new FileReader();
      reader.onload = (e) => {
        setPreviewImage(e.target.result);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleImageClick = () => {
    fileInputRef.current?.click();
  };

  const handleFormChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrorMessage(''); // Clear any previous errors
    
    try {
      // Create FormData object to handle file uploads
      const submitData = new FormData();
      
      // Append all form fields
      submitData.append('firstname', formData.firstname);
      submitData.append('lastname', formData.lastname);
      submitData.append('nickname', formData.nickname);
      submitData.append('email', formData.email);
      submitData.append('aboutme', formData.aboutme);
      submitData.append('dateofbirth', formData.dateofbirth);
      submitData.append('is_private', formData.is_private);
      
      // Append profile picture if selected
      if (formData.profilePicture) {
        submitData.append('profilePicture', formData.profilePicture);
      }
      
      
      let response=await updateProfile(submitData);
       if(response.status===200){ 
      setShowEditForm(false);
      setPreviewImage(null);
      setErrorMessage('');
      refreshProfile();
       }else if (response.status === 409){

         setErrorMessage("This email already exists for another user");
       }else{
        setErrorMessage("Failed to update profile. Please try again.");
       }
     
    } catch (error) {
      console.error('Error updating profile:', error);
      
      // Set user-friendly error message
      let errorMsg = 'Failed to update profile. Please try again.';
      
      if (error.response) {
        // Server responded with error status
        if (error.response.status === 400) {
          errorMsg = 'Invalid data provided. Please check your inputs.';
        } else if (error.response.status === 401) {
          errorMsg = 'You are not authorized to perform this action.';
        } else if (error.response.status === 413) {
          errorMsg = 'File size too large. Please choose a smaller image.';
        } else if (error.response.data && error.response.data.message) {
          errorMsg = error.response.data.message;
        }else if (error.response.status === 500) {
          errorMsg = 'Internal Server Error';
        }else if(error.response.status === 409){
          errorMsg= "This email already exists for another user";
        }
      } else if (error.message) {
        errorMsg = error.message;
      }
      
      setErrorMessage(errorMsg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCloseForm = () => {
    setShowEditForm(false);
    setErrorMessage('');
    setPreviewImage(null);
  };

  const renderFollowButton = () => {
    if (profileDetails.followbtnstatus === 'hide') {
      return null;
    }

    let buttonText = 'Follow';
    if (profileDetails.followbtnstatus === 'pending') {
      buttonText = 'Cancel Request';
    } else if (profileDetails.followbtnstatus === 'following') {
      buttonText = 'Following';
    }

    return (
      <button
        className="flex items-center gap-2 px-6 py-2 rounded-lg font-medium border cursor-pointer"
        style={{
          backgroundColor: 'transparent',
          borderColor: 'var(--border-color)',
          color: 'var(--primary-text)',
        }}
        onClick={handleFollow}
      >
        <UserPlusIcon className="w-4 h-4" />
        {buttonText}
      </button>
    );
  };
  

  return (
    <div className="w-full">
      {/* Cover Photo */}
      <div className="relative h-80 w-full overflow-hidden rounded-lg -z-50">
        <div
          className="absolute inset-0"
          style={{
            backgroundColor: 'var(--secondary-background)',
            backgroundImage: `url(${profileAPI.fetchProfileImage(profileDetails.avatar || '')})`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundRepeat: 'no-repeat',
            backgroundBlendMode: 'overlay'
          }}
        />
        
        {/* Cover Photo Edit Button - Only show for own profile */}
        {isOwnProfile && (
          <button className="absolute top-4 right-4 p-2 rounded-lg bg-black bg-opacity-50 hover:bg-opacity-70 transition-all">
            <CameraIcon className="w-5 h-5 text-white" />
          </button>
        )}
      </div>

      {/* Profile Info Section */}
      <div className="relative -mt-20 px-6 pb-4">
        <div className="flex flex-col md:flex-row md:items-end md:justify-between">
          {/* Profile Picture and Basic Info */}
          <div className="flex flex-col md:flex-row md:items-end gap-4">
            {/* Profile Picture */}
            <div className="relative">
              <div
                className={`w-40 h-40 rounded-full p-1 ${isOwnProfile ? 'cursor-pointer hover:opacity-80 transition-opacity' : ''}`}
                style={{ backgroundColor: 'var(--primary-accent)' }}
                onClick={isOwnProfile ? handleImageClick : undefined}
              >
                <Image
                  src={profileAPI.fetchProfileImage(profileDetails.avatar || '')}
                  alt="Profile"
                  width={160}
                  height={160}
                  className="w-full h-full rounded-full object-cover"
                />
              </div>
              {/* Profile Picture Edit Button - Only show for own profile */}
              {isOwnProfile && (
                <button
                  className="absolute bottom-2 right-2 p-2 rounded-full hover:opacity-80 transition-opacity"
                  style={{ backgroundColor: 'var(--tertiary-text)' }}
                  onClick={handleImageClick}
                >
                  <CameraIcon className="w-4 h-4 text-white" />
                </button>
              )}
            </div>

            {/* User Info */}
            <div className="pb-4 mt-4">
              <h1 className="text-3xl font-bold text-white">
                {profileDetails.firstname} {profileDetails.lastname}
              </h1>
              <p className="text-lg" style={{ color: 'var(--secondary-text)' }}>
                @{profileDetails.nickname}
              </p>
          <span className="inline-block px-3 py-1 rounded-full text-xs font-semibold" style={{ backgroundColor: 'var(--tertiary-text)', color: 'var(--primary-text)' }}>
                 {profileDetails.profile ? 'Public Account' : 'Private Account'}
           </span>
          
           
              
              {/* Stats */}
              <div className="flex gap-6 mt-3">
                <div className="text-center">
                  <div className="font-bold text-white">{profileDetails.numberofposts}</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Posts</div>
                </div>
                <div className="text-center">
                  <div className="font-bold text-white">{profileDetails.numberoffollowers}</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Followers</div>
                </div>
                <div className="text-center">
                  <div className="font-bold text-white">{profileDetails.numberoffollowees}</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Following</div>
                </div>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 mt-4 md:mt-0">
            {/* Edit Profile Button - Only show for own profile */}
            {isOwnProfile && (
              <button
                className="px-6 py-2 rounded-lg font-medium flex items-center gap-2 cursor-pointer hover:opacity-80 transition-opacity"
                style={{
                  backgroundColor: 'var(--tertiary-text)',
                  color: 'var(--primary-text)',
                }}
                onClick={handleEditProfile}
              >
                <EditIcon className="w-4 h-4" />
                Edit Profile
              </button>
            )}
            {renderFollowButton()}
          </div>
        </div>
      </div>

      {/* Edit Profile Modal */}
      {showEditForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div
            className="w-full max-w-2xl max-h-[90vh] overflow-y-auto rounded-lg p-6"
            style={{ backgroundColor: 'var(--primary-background)' }}
          >
            {/* Modal Header */}
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-2xl font-bold text-white">Edit Profile</h2>
              <button
                onClick={handleCloseForm}
                className="p-2 rounded-lg hover:bg-opacity-20 hover:bg-white transition-all"
              >
                <X className="w-6 h-6 text-white" />
              </button>
            </div>

            {/* Edit Form */}
            <form onSubmit={handleFormSubmit} className="space-y-6">
              {/* Error Message */}
              {errorMessage && (
                <div className="p-4 rounded-lg border-l-4 border-red-500 bg-red-50 bg-opacity-10">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <X className="h-5 w-5 text-red-400" />
                    </div>
                    <div className="ml-3">
                      <p className="text-sm text-red-300">{errorMessage}</p>
                    </div>
                    <div className="ml-auto pl-3">
                      <div className="-mx-1.5 -my-1.5">
                        <button
                          type="button"
                          onClick={() => setErrorMessage('')}
                          className="inline-flex rounded-md p-1.5 text-red-400 hover:text-red-300 focus:outline-none"
                        >
                          <X className="h-4 w-4" />
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              )}
              {/* Profile Picture Section */}
              <div className="flex items-center gap-4 mb-6">
                <div className="relative">
                  <Image
                    src={previewImage || profileAPI.fetchProfileImage(profileDetails.avatar || '')}
                    alt="Profile"
                    width={80}
                    height={80}
                    className="w-20 h-20 rounded-full object-cover cursor-pointer hover:opacity-80 transition-opacity"
                    onClick={handleImageClick}
                  />
                  <button
                    type="button"
                    className="absolute bottom-0 right-0 p-1 rounded-full bg-blue-500 hover:bg-blue-600 transition-colors"
                    onClick={handleImageClick}
                  >
                    <CameraIcon className="w-4 h-4 text-white" />
                  </button>
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-white">Profile Picture</h3>
                  <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
                    Click the image or camera icon to change your profile picture
                  </p>
                </div>
              </div>

              {/* Hidden File Input */}
              <input
                type="file"
                ref={fileInputRef}
                onChange={handleImageChange}
                accept="image/*"
                className="hidden"
              />

              {/* Form Fields */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-white mb-2">
                    First Name
                  </label>
                  <input
                    type="text"
                    name="firstname"
                    value={formData.firstname}
                    onChange={handleFormChange}
                    className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    style={{
                      backgroundColor: 'var(--secondary-background)',
                      borderColor: 'var(--border-color)',
                      color: 'var(--primary-text)',
                    }}
                    required
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-white mb-2">
                    Last Name
                  </label>
                  <input
                    type="text"
                    name="lastname"
                    value={formData.lastname}
                    onChange={handleFormChange}
                    className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    style={{
                      backgroundColor: 'var(--secondary-background)',
                      borderColor: 'var(--border-color)',
                      color: 'var(--primary-text)',
                    }}
                    required
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  Nickname
                </label>
                <input
                  type="text"
                  name="nickname"
                  value={formData.nickname}
                  onChange={handleFormChange}
                  className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  style={{
                    backgroundColor: 'var(--secondary-background)',
                    borderColor: 'var(--border-color)',
                    color: 'var(--primary-text)',
                  }}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  Email
                </label>
                <input
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleFormChange}
                  className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  style={{
                    backgroundColor: 'var(--secondary-background)',
                    borderColor: 'var(--border-color)',
                    color: 'var(--primary-text)',
                  }}
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  Date of Birth
                </label>
                <input
                  type="date"
                  name="dateofbirth"
                  value={formData.dateofbirth}
                  onChange={handleFormChange}
                  className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  style={{
                    backgroundColor: 'var(--secondary-background)',
                    borderColor: 'var(--border-color)',
                    color: 'var(--primary-text)',
                  }}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-2">
                  About Me
                </label>
                <textarea
                  name="aboutme"
                  value={formData.aboutme}
                  onChange={handleFormChange}
                  rows={4}
                  className="w-full px-3 py-2 rounded-lg border focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                  style={{
                    backgroundColor: 'var(--secondary-background)',
                    borderColor: 'var(--border-color)',
                    color: 'var(--primary-text)',
                  }}
                  placeholder="Tell us about yourself..."
                />
              </div>

              {/* Privacy Setting */}
              <div className="flex items-center gap-3">
                <input
                  type="checkbox"
                  id="is_private"
                  name="is_private"
                  checked={formData.is_private}
                  onChange={handleFormChange}
                  className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                />
                <label htmlFor="is_private" className="text-sm font-medium text-white">
                  Make my account private
                </label>
              </div>

              {/* Form Actions */}
              <div className="flex gap-3 pt-4">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="flex-1 px-6 py-3 rounded-lg font-medium flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-white"
                >
                  <Save className="w-4 h-4" />
                  {isSubmitting ? 'Saving...' : 'Save Changes'}
                </button>
                <button
                  type="button"
                  onClick={handleCloseForm}
                  className="px-6 py-3 rounded-lg font-medium border hover:bg-opacity-10 hover:bg-white transition-colors"
                  style={{
                    borderColor: 'var(--border-color)',
                    color: 'var(--primary-text)',
                  }}
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProfileCover;