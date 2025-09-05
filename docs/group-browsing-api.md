# Group Browsing API

This document describes the group browsing capabilities implemented in the social network application.

## Endpoints

### GET /groups
**Description:** Retrieve all public groups  
**Authentication:** Required  
**Response:**
```json
[
  {
    "id": 1,
    "title": "Photography Club",
    "description": "A group for photography enthusiasts",
    "creator_id": 123,
    "privacy": "public",
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

### GET /groups/search?query={searchTerm}
**Description:** Search public groups by title  
**Authentication:** Required  
**Parameters:**
- `query` (string): Search term to match against group titles
**Response:** Same as GET /groups

### GET /groups/my-groups
**Description:** Retrieve groups that the authenticated user is a member of  
**Authentication:** Required  
**Response:** Same as GET /groups

### POST /groups
**Description:** Create a new group  
**Authentication:** Required  
**Request Body:**
```json
{
  "title": "New Group",
  "description": "Group description",
  "privacy": "public"
}
```
**Response:**
```json
{
  "id": 2,
  "title": "New Group",
  "description": "Group description",
  "creator_id": 123,
  "privacy": "public",
  "created_at": "2024-01-15T11:00:00Z"
}
```

### POST /groups/{groupID}/join-request
**Description:** Send a join request to a group  
**Authentication:** Required  
**Response:**
```json
{
  "request": {
    "id": 1,
    "group_id": 2,
    "user_id": 123,
    "status": "pending"
  },
  "message": "Join request approved"
}
```

## Frontend Components

### GroupBrowser
Located at: `frontend/components/groups/GroupBrowser.jsx`

Features:
- Browse all public groups
- Search groups by title
- View user's joined groups
- Send join requests
- Tabbed interface for different views

### GroupCreation
Located at: `frontend/components/groups/GroupCreation.jsx`

Features:
- Create new groups
- Set group privacy (public/private)
- Form validation

### Groups Page
Located at: `frontend/src/app/groups/page.js`

Main page that combines group creation and browsing functionality.

## Usage

1. Navigate to `/groups` in the frontend
2. Use the "Browse All Groups" tab to see all public groups
3. Use the "My Groups" tab to see groups you've joined
4. Use the search bar to find specific groups
5. Click "Join Group" to send a join request
6. Use the "Create New Group" form to create your own group

## Database Schema

The group browsing functionality relies on these tables:
- `groups`: Stores group information
- `group_members`: Tracks group membership
- `group_requests`: Handles join requests

## Security

- All endpoints require authentication
- Users can only see public groups when browsing
- Private groups are only visible to members
- Join requests are required for group membership