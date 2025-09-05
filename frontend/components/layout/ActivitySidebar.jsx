import React, { useEffect, useState, useCallback } from 'react';
import ActivityItem from '../homepage/ActivityItem';
import { wsService } from '../../lib/websocket';
import { notificationService } from '../../lib/notificationService';
import { profileAPI, groupAPI } from '../../lib/api';

const ActivitySidebar = () => {
  const [activities, setActivities] = useState([]);

  const formatTimeSince = (timestamp) => {
    const date = new Date(timestamp);
    const now = Date.now();
    const diffMs = now - date.getTime();
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHr = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHr / 24);

    if (diffSec < 60) return 'Just now';
    if (diffMin < 60) return `${diffMin} min ago`;
    if (diffHr < 24) return `${diffHr} hr${diffHr > 1 ? 's' : ''} ago`;
    return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`;
  };

  const getAvatarUrl = (avatar, isGroup) => {
    if (!avatar) return '/default-group-avatar.png';
    return isGroup
      ? groupAPI.getAvatarUrl(avatar)
      : profileAPI.fetchProfileImage(avatar);
  };

  const handleRequest = useCallback((notification) => {
    console.log('Received follow_request notification:', notification);

    const {
      user_name,
      avatar,
      timestamp,
      user_id,
      request_id,
    } = notification;

    if (!request_id) return;

    const activity = {
      image: getAvatarUrl(avatar, false),
      name: user_name || 'Unknown User',
      action: 'sent a follow request',
      time: formatTimeSince(timestamp),
      isGroup: false,
      isPartial: false,
      userId: user_id,
      requestId: request_id,
    };

    setActivities((prev) => {
      if (prev.some((a) => a.requestId === activity.requestId && a.action === activity.action)) {
        return prev;
      }
      return [activity, ...prev];
    });
  }, []);

  const handleFollow = useCallback((notification) => {
    console.log('Received follow notification:', notification);

    const {
      user_name,
      avatar,
      timestamp,
      user_id,
      request_id,
    } = notification;

    if (!request_id) return;

    const activity = {
      image: getAvatarUrl(avatar, false),
      name: user_name || 'Unknown User',
      action: 'followed you',
      time: formatTimeSince(timestamp),
      isGroup: false,
      isPartial: true,
      userId: user_id,
      requestId: request_id,
    };

    setActivities((prev) => {
      if (prev.some((a) => a.requestId === activity.requestId && a.action === activity.action)) {
        return prev;
      }
      return [activity, ...prev];
    });
  }, []);

  const handleNotification = useCallback((notification) => {
    notificationService.handleNotification(notification);
  }, []);

  const handleAccept = (requestId) => {
    profileAPI.acceptFollowRequest(requestId, 'accepted')
      .then(() => {
        setActivities((prev) => prev.filter((a) => a.requestId !== requestId));
      })
      .catch((err) => {
        console.error('Failed to accept follow request:', err);
      });
  };

  const handleDecline = (requestId) => {
    profileAPI.acceptFollowRequest(requestId, 'rejected')
      .then(() => {
        setActivities((prev) => prev.filter((a) => a.requestId !== requestId));
      })
      .catch((err) => {
        console.error('Failed to decline follow request:', err);
      });
  };

  const handleFollowRequestCancelled = (notification) => {
    // Remove the cancelled follow request from activities
    setActivities((prev) => prev.filter((a) =>
      !(a.requestId === notification.request_id && a.userId === notification.user_id)
    ));
  };

  useEffect(() => {
    wsService.onMessage('notification', handleNotification);
    notificationService.onNotification('follow_request', handleRequest);
    notificationService.onNotification('follow', handleFollow);
    notificationService.onNotification('follow_request_cancelled', handleFollowRequestCancelled);

    profileAPI.fetchPendingFollowRequests()
      .then((requests) => {
        const pendingRequests = Array.isArray(requests?.user) ? requests.user : [];

        const formatted = pendingRequests.map((req) => ({
          image: getAvatarUrl(req.avatar, false),
          name: req.firstname || 'Unknown User',
          action: 'sent a follow request',
          time: formatTimeSince(req.requested_at),
          isGroup: false,
          isPartial: false,
          userId: req.follower_id,
          requestId: req.request_id,
        }));

        setActivities((prev) => {
          const uniqueFormatted = formatted.filter(
            (f) => !prev.some((p) => p.requestId === f.requestId && p.action === f.action)
          );
          return [...uniqueFormatted, ...prev];
        });
      })
      .catch((err) => {
        console.error('Failed to load pending follow requests:', err);
      });

    return () => {
      wsService.removeHandler('notification', handleNotification);
      notificationService.removeHandler('follow_request', handleRequest);
      notificationService.removeHandler('follow', handleFollow);
      notificationService.removeHandler('follow_request_cancelled', handleFollowRequestCancelled);
    };
  }, [handleRequest, handleNotification, handleFollow, handleFollowRequestCancelled]);

  return (
    <div className="w-72">
      <div className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
        <h3 className="font-bold mb-4">Recent activity</h3>
        <div className="flex flex-col gap-4 cursor-pointer">
          {activities.length === 0 ? (
            <p className="text-sm text-gray-500">No recent activity</p>
          ) : (
            activities
              .filter((activity) => activity?.requestId != null)
              .map((activity) => (
                <ActivityItem
                  key={`${activity.requestId}-${activity.action}`}
                  {...activity}
                  onAccept={handleAccept}
                  onDecline={handleDecline}
                />
              ))
          )}
        </div>
      </div>
    </div>
  );
};

export default ActivitySidebar;
