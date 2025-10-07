# Voting API Examples

This document provides examples of how to use the voting API endpoints.

## Overview

The voting system allows users to create polls with multiple options and enables other users to vote on them. It supports both single-choice and multiple-choice voting.

## Features

- Create votes with descriptions and multiple options
- Support for single-choice and multiple-choice voting
- User authentication required for voting
- One vote per user per poll
- Real-time vote results
- Vote management (close, cancel)

## API Endpoints

### 1. Create Vote

**POST** `/api/v1/votes`

Creates a new vote poll.

**Headers:**
```
Content-Type: application/json
Cookie: user_session=your_session_cookie
```

**Request Body:**
```json
{
  "title": "Best Cricket Player of 2024",
  "description": "Vote for your favorite cricket player from the list below",
  "type": "single",
  "options": [
    "Virat Kohli",
    "Rohit Sharma", 
    "Jasprit Bumrah",
    "Ravindra Jadeja"
  ]
}
```

**Response:**
```json
{
  "data": {
    "message": "Vote created successfully",
    "vote": {
      "id": "vote-uuid",
      "title": "Best Cricket Player of 2024",
      "description": "Vote for your favorite cricket player from the list below",
      "type": "single",
      "status": "active",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T10:00:00Z",
      "updated_at": "2025-01-27T10:00:00Z",
      "closed_at": null
    }
  }
}
```

### 2. Get Vote

**GET** `/api/v1/votes/{id}`

Retrieves a vote with its options.

**Response:**
```json
{
  "data": {
    "vote": {
      "id": "vote-uuid",
      "title": "Best Cricket Player of 2024",
      "description": "Vote for your favorite cricket player from the list below",
      "type": "single",
      "status": "active",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T10:00:00Z",
      "updated_at": "2025-01-27T10:00:00Z",
      "closed_at": null
    },
    "options": [
      {
        "id": "option-uuid-1",
        "vote_id": "vote-uuid",
        "text": "Virat Kohli",
        "description": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T10:00:00Z"
      },
      {
        "id": "option-uuid-2",
        "vote_id": "vote-uuid",
        "text": "Rohit Sharma",
        "description": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T10:00:00Z"
      }
    ]
  }
}
```

### 3. Get Vote Results

**GET** `/api/v1/votes/{id}/results`

Retrieves a vote with results and user's vote (if authenticated).

**Response:**
```json
{
  "data": {
    "vote": {
      "id": "vote-uuid",
      "title": "Best Cricket Player of 2024",
      "description": "Vote for your favorite cricket player from the list below",
      "type": "single",
      "status": "active",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T10:00:00Z",
      "updated_at": "2025-01-27T10:00:00Z",
      "closed_at": null
    },
    "options": [
      {
        "id": "option-uuid-1",
        "vote_id": "vote-uuid",
        "text": "Virat Kohli",
        "description": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T10:00:00Z"
      },
      {
        "id": "option-uuid-2",
        "vote_id": "vote-uuid",
        "text": "Rohit Sharma",
        "description": "",
        "created_at": "2025-01-27T10:00:00Z",
        "updated_at": "2025-01-27T10:00:00Z"
      }
    ],
    "results": {
      "option-uuid-1": 15,
      "option-uuid-2": 8,
      "option-uuid-3": 12,
      "option-uuid-4": 5
    },
    "user_vote": {
      "id": "user-vote-uuid",
      "vote_id": "vote-uuid",
      "user_id": "user-uuid",
      "selected_options": ["option-uuid-1"],
      "voted_at": "2025-01-27T10:30:00Z"
    },
    "total_votes": 40,
    "voted_users": ["user-uuid-1", "user-uuid-2", "user-uuid-3"]
  }
}
```

### 4. Cast Vote

**POST** `/api/v1/votes/{id}/vote`

Casts a vote for a specific poll.

**Headers:**
```
Content-Type: application/json
Cookie: user_session=your_session_cookie
```

**Request Body:**
```json
{
  "selected_options": ["option-uuid-1"]
}
```

**Response:**
```json
{
  "data": {
    "message": "Vote cast successfully"
  }
}
```

### 5. List Votes

**GET** `/api/v1/votes`

Lists votes with optional filters.

**Query Parameters:**
- `status`: Filter by status (active, closed, cancelled)
- `type`: Filter by type (single, multiple)
- `created_by`: Filter by creator user ID
- `limit`: Number of results (default: 20, max: 100)
- `offset`: Number of results to skip (default: 0)

**Example:**
```
GET /api/v1/votes?status=active&type=single&limit=10
```

**Response:**
```json
{
  "data": [
    {
      "id": "vote-uuid-1",
      "title": "Best Cricket Player of 2024",
      "description": "Vote for your favorite cricket player",
      "type": "single",
      "status": "active",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T10:00:00Z",
      "updated_at": "2025-01-27T10:00:00Z",
      "closed_at": null
    },
    {
      "id": "vote-uuid-2",
      "title": "Favorite Cricket Format",
      "description": "Which format do you enjoy most?",
      "type": "multiple",
      "status": "active",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T09:00:00Z",
      "updated_at": "2025-01-27T09:00:00Z",
      "closed_at": null
    }
  ]
}
```

### 6. Update Vote

**PUT** `/api/v1/votes/{id}`

Updates an existing vote (only by creator).

**Headers:**
```
Content-Type: application/json
Cookie: user_session=your_session_cookie
```

**Request Body:**
```json
{
  "title": "Updated Vote Title",
  "description": "Updated description",
  "status": "closed"
}
```

**Response:**
```json
{
  "data": {
    "message": "Vote updated successfully",
    "vote": {
      "id": "vote-uuid",
      "title": "Updated Vote Title",
      "description": "Updated description",
      "type": "single",
      "status": "closed",
      "created_by": "user-uuid",
      "created_at": "2025-01-27T10:00:00Z",
      "updated_at": "2025-01-27T11:00:00Z",
      "closed_at": "2025-01-27T11:00:00Z"
    }
  }
}
```

### 7. Delete Vote

**DELETE** `/api/v1/votes/{id}`

Deletes a vote (only by creator).

**Headers:**
```
Cookie: user_session=your_session_cookie
```

**Response:**
```json
{
  "data": {
    "message": "Vote deleted successfully"
  }
}
```

### 8. Get User Vote

**GET** `/api/v1/votes/{id}/my-vote`

Gets the current user's vote for a specific poll.

**Headers:**
```
Cookie: user_session=your_session_cookie
```

**Response:**
```json
{
  "data": {
    "id": "user-vote-uuid",
    "vote_id": "vote-uuid",
    "user_id": "user-uuid",
    "selected_options": ["option-uuid-1"],
    "voted_at": "2025-01-27T10:30:00Z"
  }
}
```

### 9. Check if User Voted

**GET** `/api/v1/votes/{id}/has-voted`

Checks if the current user has voted on a specific poll.

**Headers:**
```
Cookie: user_session=your_session_cookie
```

**Response:**
```json
{
  "data": {
    "has_voted": true
  }
}
```

### 10. Close Vote

**POST** `/api/v1/votes/{id}/close`

Closes a vote (only by creator).

**Headers:**
```
Cookie: user_session=your_session_cookie
```

**Response:**
```json
{
  "data": {
    "message": "Vote closed successfully"
  }
}
```

### 11. Cancel Vote

**POST** `/api/v1/votes/{id}/cancel`

Cancels a vote (only by creator).

**Headers:**
```
Cookie: user_session=your_session_cookie
```

**Response:**
```json
{
  "data": {
    "message": "Vote cancelled successfully"
  }
}
```

## Error Responses

### Validation Error
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "title is required"
  }
}
```

### Unauthorized
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

### Forbidden
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Only vote creator can update vote"
  }
}
```

### Already Voted
```json
{
  "error": {
    "code": "ALREADY_VOTED",
    "message": "User has already voted on this poll"
  }
}
```

### Vote Closed
```json
{
  "error": {
    "code": "INVALID_STATE",
    "message": "Cannot vote on closed or cancelled vote"
  }
}
```

## cURL Examples

### Create a Vote
```bash
curl -X POST http://localhost:8081/api/v1/votes \
  -H "Content-Type: application/json" \
  -H "Cookie: user_session=your_session_cookie" \
  -d '{
    "title": "Best Cricket Player",
    "description": "Vote for your favorite player",
    "type": "single",
    "options": ["Player A", "Player B", "Player C"]
  }'
```

### Cast a Vote
```bash
curl -X POST http://localhost:8081/api/v1/votes/vote-uuid/vote \
  -H "Content-Type: application/json" \
  -H "Cookie: user_session=your_session_cookie" \
  -d '{
    "selected_options": ["option-uuid-1"]
  }'
```

### Get Vote Results
```bash
curl -X GET http://localhost:8081/api/v1/votes/vote-uuid/results \
  -H "Cookie: user_session=your_session_cookie"
```

### List Active Votes
```bash
curl -X GET "http://localhost:8081/api/v1/votes?status=active&limit=10" \
  -H "Cookie: user_session=your_session_cookie"
```

## Vote Types

### Single Choice
- Users can select only one option
- Use `"type": "single"` when creating a vote
- `selected_options` array should contain exactly one option ID

### Multiple Choice
- Users can select multiple options
- Use `"type": "multiple"` when creating a vote
- `selected_options` array can contain multiple option IDs

## Vote Status

- **active**: Vote is open for voting
- **closed**: Vote is closed, no more votes accepted
- **cancelled**: Vote is cancelled

## Authentication

All voting operations require user authentication. Users must be logged in to:
- Create votes
- Cast votes
- Update/delete their own votes
- View their own vote history

Public access is available for:
- Viewing votes and their options
- Viewing vote results (without user's vote details)
- Listing votes
