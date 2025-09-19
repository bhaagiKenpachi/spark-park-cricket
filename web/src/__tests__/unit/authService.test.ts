import { authService, User, AuthStatus } from '@/services/authService';

// Mock fetch
global.fetch = jest.fn();

describe('AuthService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Clear localStorage
    localStorage.clear();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('getAuthStatus', () => {
    it('should return authentication status successfully', async () => {
      const mockAuthStatus: AuthStatus = {
        authenticated: true,
        user: {
          id: 'user-123',
          google_id: 'google-123',
          email: 'test@example.com',
          name: 'Test User',
          picture: 'https://example.com/picture.jpg',
          email_verified: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      };

      const mockResponse = {
        ok: true,
        status: 200,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          data: mockAuthStatus,
          success: true,
          message: 'Success',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      const result = await authService.getAuthStatus();

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/status'),
        expect.objectContaining({
          credentials: 'include',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Accept: 'application/json',
          }),
          mode: 'cors',
        })
      );

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockAuthStatus);
    });

    it('should handle authentication status check failure', async () => {
      const mockResponse = {
        ok: false,
        status: 401,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          error: 'UNAUTHORIZED',
          message: 'Not authenticated',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      await expect(authService.getAuthStatus()).rejects.toThrow();
    });

    it('should handle network errors', async () => {
      (fetch as jest.Mock).mockRejectedValue(new Error('Network error'));

      await expect(authService.getAuthStatus()).rejects.toThrow(
        'Network error'
      );
    });
  });

  describe('getCurrentUser', () => {
    it('should return current user successfully', async () => {
      const mockUser: User = {
        id: 'user-123',
        google_id: 'google-123',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email_verified: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      const mockResponse = {
        ok: true,
        status: 200,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          data: { user: mockUser },
          success: true,
          message: 'Success',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      const result = await authService.getCurrentUser();

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/me'),
        expect.objectContaining({
          credentials: 'include',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            Accept: 'application/json',
          }),
          mode: 'cors',
        })
      );

      expect(result.success).toBe(true);
      expect(result.data.user).toEqual(mockUser);
    });

    it('should handle unauthorized access', async () => {
      const mockResponse = {
        ok: false,
        status: 401,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          error: 'UNAUTHORIZED',
          message: 'Not authenticated',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      await expect(authService.getCurrentUser()).rejects.toThrow();
    });
  });

  describe('logout', () => {
    it('should logout successfully', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          data: { message: 'Logged out successfully' },
          success: true,
          message: 'Success',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      const result = await authService.logout();

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('/auth/logout'),
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
        })
      );

      expect(result.success).toBe(true);
      expect(result.data.message).toBe('Logged out successfully');
    });

    it('should handle logout failure', async () => {
      const mockResponse = {
        ok: false,
        status: 500,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          error: 'LOGOUT_ERROR',
          message: 'Failed to logout',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      await expect(authService.logout()).rejects.toThrow();
    });
  });

  describe('initiateGoogleLogin', () => {
    it('should call initiateGoogleLogin without errors', () => {
      // This test verifies the method can be called without throwing errors
      // The actual navigation behavior is tested in e2e tests
      expect(() => {
        authService.initiateGoogleLogin();
      }).not.toThrow();
    });
  });

  describe('isAuthenticated', () => {
    it('should return true when user is authenticated in localStorage', () => {
      localStorage.setItem('auth_authenticated', 'true');

      const result = authService.isAuthenticated();

      expect(result).toBe(true);
    });

    it('should return false when user is not authenticated in localStorage', () => {
      localStorage.setItem('auth_authenticated', 'false');

      const result = authService.isAuthenticated();

      expect(result).toBe(false);
    });

    it('should return false when localStorage is empty', () => {
      const result = authService.isAuthenticated();

      expect(result).toBe(false);
    });
  });

  describe('setAuthState', () => {
    it('should set authentication state in localStorage', () => {
      const user: User = {
        id: 'user-123',
        google_id: 'google-123',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email_verified: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      authService.setAuthState(true, user);

      expect(localStorage.getItem('auth_authenticated')).toBe('true');
      expect(JSON.parse(localStorage.getItem('auth_user') || '{}')).toEqual(
        user
      );
    });

    it('should set authentication state without user', () => {
      authService.setAuthState(true);

      expect(localStorage.getItem('auth_authenticated')).toBe('true');
      expect(localStorage.getItem('auth_user')).toBeNull();
    });

    it('should clear user when setting authenticated to false', () => {
      const user: User = {
        id: 'user-123',
        google_id: 'google-123',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email_verified: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      // First set authenticated with user
      authService.setAuthState(true, user);
      expect(localStorage.getItem('auth_user')).not.toBeNull();

      // Then set authenticated to false
      authService.setAuthState(false);
      expect(localStorage.getItem('auth_authenticated')).toBe('false');
      expect(localStorage.getItem('auth_user')).toBeNull();
    });
  });

  describe('getStoredUser', () => {
    it('should return stored user from localStorage', () => {
      const user: User = {
        id: 'user-123',
        google_id: 'google-123',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email_verified: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      localStorage.setItem('auth_user', JSON.stringify(user));

      const result = authService.getStoredUser();

      expect(result).toEqual(user);
    });

    it('should return null when no user is stored', () => {
      const result = authService.getStoredUser();

      expect(result).toBeNull();
    });

    it('should return null when stored user is invalid JSON', () => {
      localStorage.setItem('auth_user', 'invalid-json');

      const result = authService.getStoredUser();

      expect(result).toBeNull();
    });
  });

  describe('clearAuthState', () => {
    it('should clear authentication state from localStorage', () => {
      const user: User = {
        id: 'user-123',
        google_id: 'google-123',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg',
        email_verified: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      };

      localStorage.setItem('auth_authenticated', 'true');
      localStorage.setItem('auth_user', JSON.stringify(user));

      authService.clearAuthState();

      expect(localStorage.getItem('auth_authenticated')).toBeNull();
      expect(localStorage.getItem('auth_user')).toBeNull();
    });
  });

  describe('request method', () => {
    it('should include credentials in requests', async () => {
      const mockResponse = {
        ok: true,
        status: 200,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          data: { authenticated: false },
          success: true,
          message: 'Success',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      await authService.getAuthStatus();

      expect(fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          credentials: 'include',
          mode: 'cors',
        })
      );
    });

    it('should handle response data extraction', async () => {
      const mockData = { authenticated: true };
      const mockResponse = {
        ok: true,
        status: 200,
        headers: new Headers(),
        json: jest.fn().mockResolvedValue({
          data: mockData,
          success: true,
          message: 'Success',
        }),
      };

      (fetch as jest.Mock).mockResolvedValue(mockResponse);

      const result = await authService.getAuthStatus();

      expect(result.data).toEqual(mockData);
      expect(result.success).toBe(true);
      expect(result.message).toBe('Success');
    });
  });
});
