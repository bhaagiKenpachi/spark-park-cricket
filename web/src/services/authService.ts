import { ApiResponse, ApiError } from './api';

export interface User {
  id: string;
  google_id: string;
  email: string;
  name: string;
  picture?: string;
  email_verified: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}

export interface AuthStatus {
  authenticated: boolean;
  user?: User;
}

export interface AuthResponse {
  user?: User;
  message: string;
}

class AuthService {
  private baseURL: string;

  constructor(
    baseURL: string = process.env.NEXT_PUBLIC_API_URL!
  ) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;

    const defaultHeaders = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    };

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
      mode: 'cors',
      credentials: 'include', // Include cookies for session management
    };

    try {
      const response = await fetch(url, config);

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new ApiError(
          errorData.message || `HTTP error! status: ${response.status}`,
          response.status,
          errorData
        );
      }

      const data = await response.json();
      return {
        data: data.data, // Extract the nested data
        success: true,
        message: data.message,
      };
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        error instanceof Error ? error.message : 'Network error',
        0,
        error
      );
    }
  }

  /**
   * Initiate Google OAuth login
   * @param redirectUrl - Optional URL to redirect to after successful login
   */
  async initiateGoogleLogin(redirectUrl?: string): Promise<void> {
    // Get the current URL if no redirect URL is provided
    const redirect = redirectUrl || (typeof window !== 'undefined' ? window.location.href : '');

    // Redirect to backend Google OAuth endpoint with redirect URL
    const url = redirect
      ? `${this.baseURL}/auth/google?redirect_url=${encodeURIComponent(redirect)}`
      : `${this.baseURL}/auth/google`;

    window.location.href = url;
  }

  /**
   * Check authentication status
   */
  async getAuthStatus(): Promise<ApiResponse<AuthStatus>> {
    const result = await this.request<AuthStatus>('/auth/status');
    return result;
  }

  /**
   * Get current user information
   */
  async getCurrentUser(): Promise<ApiResponse<AuthResponse>> {
    return this.request<AuthResponse>('/auth/me');
  }

  /**
   * Logout current user
   */
  async logout(): Promise<ApiResponse<AuthResponse>> {
    return this.request<AuthResponse>('/auth/logout', {
      method: 'POST',
    });
  }

  /**
   * Update user name
   */
  async updateUserName(name: string): Promise<ApiResponse<User>> {
    return this.request<User>('/users/me', {
      method: 'PUT',
      body: JSON.stringify({ name }),
    });
  }

  /**
   * Check if user is authenticated (client-side check)
   */
  isAuthenticated(): boolean {
    // This is a simple client-side check
    // The actual authentication state should be managed by Redux
    return (
      typeof window !== 'undefined' &&
      localStorage.getItem('auth_authenticated') === 'true'
    );
  }

  /**
   * Store authentication state in localStorage
   */
  setAuthState(authenticated: boolean | undefined, user?: User): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(
        'auth_authenticated',
        String(authenticated === true)
      );
      if (user) {
        localStorage.setItem('auth_user', JSON.stringify(user));
      } else {
        localStorage.removeItem('auth_user');
      }
    }
  }

  /**
   * Get stored user from localStorage
   */
  getStoredUser(): User | null {
    if (typeof window !== 'undefined') {
      const userStr = localStorage.getItem('auth_user');
      if (userStr) {
        try {
          return JSON.parse(userStr);
        } catch {
          return null;
        }
      }
    }
    return null;
  }

  /**
   * Clear authentication state
   */
  clearAuthState(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_authenticated');
      localStorage.removeItem('auth_user');
    }
  }
}

export const authService = new AuthService();
export default authService;
