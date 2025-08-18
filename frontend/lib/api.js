import { fetchComments, fetchPostsPaginated } from "./auth";

const API_BASE = process.env.NEXT_PUBLIC_API_URL;

const apiCall = async (endpoint, options = {}) => {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }

  return response.json();
};

export const postAPI = {
  createPost: async (formData) => {
    try {
      const response = await fetch(`${API_BASE}/posts`, {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
      }

      const data = await response.json();
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to create post",
      };
    }
  },
  createComment: async (postId, formData) => {
    try {
      const response = await fetch(`${API_BASE}/posts/${postId}/comments`, {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
      }

      const data = await response.json();
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to create comment",
      };
    }
  },
  fetchPostsPaginated: async (page, limit) => {
    try {
      const data = await apiCall(`/posts?page=${page}&limit=${limit}`);
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to fetch posts",
      };
    }
  },
  updatePost: async (postId, content, image) => {
    try {
      const formData = new FormData();
      formData.append("content", content);

      if (image) {
        formData.append("image", image);
      }
      const response = await fetch(`${API_BASE}/posts/${postId}`, {
        method: "PUT",
        credentials: "include",
        body: formData,
      });
       if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
      }

      const data = await response.json();
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to update post",
      };
    }
  },
  deletePost: async (postId) => {
    try {
      const data = await apiCall(`/posts/${postId}`, { method: "DELETE" });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to delete post",
      };
    }
  },
  fetchComments: async (postId) => {
    try {
      const data = await apiCall(`/posts/${postId}/comments`);
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to fetch comments",
      };
    }
  },
  updateComment: async (postId, commentId, content, image) => {
    try {
      const formData = new FormData();
      formData.append("content", content);

      if (image) {
        formData.append("image", image);
      }

      const response = await fetch(`${API_BASE}/posts/${postId}/comments/${commentId}`, {
        method: "PUT",
        credentials: "include",
        body: formData,
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
      }

      const data = await response.json();
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to update comment",
      };
    }
  },
  deleteComment: async (postId, commentId) => {
    try {
      const data = await apiCall(`/posts/${postId}/comments/${commentId}`, { method: "DELETE" });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to delete post",
      };
    }
  },
  reactToPost: async (postId, reaction) => {
    try {
      const data = await apiCall(`/posts/${postId}/reaction`, {
        method: "POST",
        body: JSON.stringify({ reaction }),
      });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to react to post",
      };
    }
  },
  unreactToPost: async (postId) => {
    try {
      const data = await apiCall(`/posts/${postId}/reaction`, { method: "DELETE" });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to unreact to post",
      };
    }
  },
  reactToComment: async (commentId, reaction) => {
    try {
      const data = await apiCall(`/comments/${commentId}/reaction`, {
        method: "POST",
        body: JSON.stringify({ reaction }),
      });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to react to comment",
      };
    }
  },
  unreactToComment: async (commentId) => {
    try {
      const data = await apiCall(`/comments/${commentId}/reaction`, { method: "DELETE" });
      return { success: true, data };
    } catch (error) {
      return {
        success: false,
        error: error.message || "Failed to unreact to comment",
      };
    }
  },
};

export const chatAPI = {
  // Get private message history
  getPrivateMessages: (userId, limit = 50, offset = 0) =>
    apiCall(
      `/api/messages/private?user=${userId}&limit=${limit}&offset=${offset}`
    ),

  // Get group message history
  getGroupMessages: (groupId, limit = 50, offset = 0) =>
    apiCall(
      `/api/messages/group?group=${groupId}&limit=${limit}&offset=${offset}`
    ),

  // Send group invitation
  sendGroupInvite: (groupId, userId, groupName) =>
    apiCall("/api/groups/invite", {
      method: "POST",
      body: JSON.stringify({
        group_id: groupId,
        user_id: userId,
        group_name: groupName,
      }),
    }),

  // Get notifications
  getNotifications: (limit = 20, offset = 0) =>
    apiCall(`/api/notifications?limit=${limit}&offset=${offset}`),

  // Mark notifications as read
  markNotificationsRead: () =>
    apiCall('/api/notifications/read', { method: 'POST' }),

  // Get unread chats
  getUnreadChats: () =>
    apiCall('/api/chats/unread'),

  // Get unread chat count
  getUnreadChatCount: () =>
    apiCall('/api/chats/unread/count'),
  
  // Get users the current user can message
  getMessageableUsers: () =>
    apiCall('/api/users/messageable'),

  getGroups: () =>
    apiCall('/api/chats/groups'),
};

const fallbackAvatar =
  "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNGM0Y0RjYiLz4KPGNpcmNsZSBjeD0iMjAiIGN5PSIxNiIgcj0iNiIgZmlsbD0iIzlDQTNBRiIvPgo8cGF0aCBkPSJNMzIgMzJDMzIgMjYuNDc3MiAyNy41MjI4IDIyIDIyIDIySDE4QzEyLjQ3NzIgMjIgOCAyNi40NzcyIDggMzJWMzJIMzJWMzJaIiBmaWxsPSIjOUNBM0FGIi8+Cjwvc3ZnPgo=";

export function fetchProfileImage(avatar) {
  if (avatar == "no profile photo") return fallbackAvatar;
  return `${API_BASE}/avatar?avatar=${encodeURIComponent(avatar)}`;
}

function fetchFollowers() {
  // TODO: Fetch followers count from backend
  return null;
}

function fetchFollowing() {
  // TODO: Fetch following count from backend
  return null;
}

function fetchProfileStatus() {
  // TODO: Fetch profile status from backend
  return null;
}

function fetchCommunities() {
  // TODO: Fetch communities list from backend
  return null;
}

async function fetchPendingFollowRequests() {
  return apiCall("/pending-follow-requests", { method: "GET" });
}

const acceptFollowRequest = (requestId, status) =>
  apiCall(`/follow-request/${requestId}/request`, {
    method: "POST",
    body: JSON.stringify({ status }),
  });

const declineFollowRequest = (requestId) =>
  apiCall(`/follow-request/${requestId}/cancel`, { method: "POST" });

const getFollowRequestId = (followeeId) =>
  apiCall(`/follow-request-id/${followeeId}`, { method: "GET" });

const cancelFollowRequest = (requestId) =>
  apiCall(`/follow-request/${requestId}/cancel`, { method: "DELETE" });

export async function updateProfile(profileData) {
  var response = await fetch(`${API_BASE}/EditProfile`, {
    method: "PUT",
    credentials: "include",
    body: profileData,
  });

  // if (!response.ok) {
  //   throw new Error(response.statusText);

  // }

  return {
    status: response.status,
    message: response.json(),
  };
}

export const profileAPI = {
  getProfile: (userId) => apiCall(`/profile/${userId}`),
  getFollowers: async (userId) => {
    try {
      const data = await apiCall(`/profile/${userId}/followers`);
      return { success: true, data };
    } catch (error) {
      return { success: false, error: error.message || 'Failed to fetch followers' };
    }
  },
  getFollowing: async (userId) => {
    try {
      const data = await apiCall(`/profile/${userId}/followees`);
      return { success: true, data };
    } catch (error) {
      return { success: false, error: error.message || 'Failed to fetch following' };
    }
  },
  follow: (followeeid) =>
    apiCall("/follow", {
      method: "POST",
      body: JSON.stringify({ followeeid }),
    }),
  unfollow: (followeeid) =>
    apiCall("/unfollow", {
      method: "DELETE",
      body: JSON.stringify({ followeeid }),
    }),
  fetchProfileImage,
  fetchFollowers,
  fetchFollowing,
  fetchProfileStatus,
  fetchCommunities,
  fetchPendingFollowRequests,
  acceptFollowRequest,
  declineFollowRequest,
  getFollowRequestId,
  cancelFollowRequest,
};

export function fetchGroupImage(avatar) {
  if (!avatar) return fallbackAvatar;
  return `${API_BASE}/group-avatar?avatar=${encodeURIComponent(avatar)}`;
}
