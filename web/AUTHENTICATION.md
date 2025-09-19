# Authentication Integration

This document explains how Google authentication has been integrated into the Spark Park Cricket frontend.

## Overview

The frontend now includes a complete Google OAuth2 authentication system that integrates with the backend authentication service.

## Features

- **Google Sign-In**: Users can sign in using their Google accounts
- **Session Management**: Secure session handling with HTTP-only cookies
- **User Profile**: Display user information and profile pictures
- **Protected Routes**: Components can be protected with authentication requirements
- **Automatic State Management**: Authentication state is managed through Redux

## Components

### AuthProvider

- Wraps the entire application to initialize authentication state
- Automatically checks authentication status on app load
- Handles loading states during authentication initialization

### LoginButton

- Displays a "Sign in with Google" button
- Redirects users to the backend Google OAuth endpoint
- Shows loading state during authentication process

### UserMenu

- Displays user profile information when authenticated
- Shows user avatar (from Google profile picture)
- Provides logout functionality
- Dropdown menu with user details

### ProtectedRoute

- Wrapper component for protecting routes that require authentication
- Shows login prompt for unauthenticated users
- Can accept custom fallback components

## Services

### AuthService

- Handles all authentication-related API calls
- Manages localStorage for client-side authentication state
- Provides methods for login, logout, and status checking

## Redux Integration

### AuthSlice

- Manages authentication state in Redux store
- Provides async thunks for authentication operations
- Handles loading states and error management

### State Structure

```typescript
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  isInitialized: boolean;
}
```

## Usage Examples

### Using Authentication State

```tsx
import { useSelector } from 'react-redux';
import { RootState } from '@/store/reducers';

function MyComponent() {
  const { isAuthenticated, user, isLoading } = useSelector(
    (state: RootState) => state.auth
  );

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      {isAuthenticated ? <p>Welcome, {user?.name}!</p> : <p>Please sign in</p>}
    </div>
  );
}
```

### Protecting Routes

```tsx
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';

function AdminPanel() {
  return (
    <ProtectedRoute>
      <div>Admin content here</div>
    </ProtectedRoute>
  );
}
```

### Custom Fallback for Protected Routes

```tsx
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';
import { LoginButton } from '@/components/auth/LoginButton';

function PremiumContent() {
  return (
    <ProtectedRoute
      fallback={
        <div className="text-center">
          <h2>Premium Content</h2>
          <p>Sign in to access premium features</p>
          <LoginButton />
        </div>
      }
    >
      <div>Premium content here</div>
    </ProtectedRoute>
  );
}
```

## API Integration

The authentication system integrates with the backend API endpoints:

- `GET /auth/status` - Check authentication status
- `GET /auth/me` - Get current user information
- `POST /auth/logout` - Logout current user
- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/google/callback` - Handle Google OAuth callback

## Environment Configuration

The frontend expects the following environment variable:

```env
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

## Security Features

- **HTTP-Only Cookies**: Session cookies are HTTP-only and secure
- **CSRF Protection**: Built-in CSRF protection through secure cookies
- **Automatic Token Refresh**: Handled by the backend session management
- **Secure Redirects**: OAuth redirects are validated by the backend

## Testing

Authentication functionality is covered by unit tests:

- `authService.test.ts` - Tests for authentication service
- Redux slice tests for state management
- Component tests for UI interactions

## Development

To test the authentication system:

1. Start the backend server on port 8081
2. Start the frontend development server: `npm run dev`
3. Navigate to `http://localhost:3000`
4. Click "Sign in with Google" to test the authentication flow

## Troubleshooting

### Common Issues

1. **Redux Thunk Error**: If you see "Actions must be plain objects" error, ensure Redux Thunk middleware is enabled in the store configuration
2. **CORS Errors**: Ensure the backend is configured to allow requests from the frontend domain
3. **Session Not Persisting**: Check that cookies are being set correctly and the domain configuration is correct
4. **Redirect Issues**: Verify the Google OAuth redirect URL is configured correctly in Google Console

### Debug Mode

Enable debug logging by checking the browser's Network tab and Redux DevTools to monitor authentication state changes.
