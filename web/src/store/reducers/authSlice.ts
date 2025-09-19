import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { User } from '@/services/authService';
import { authService } from '@/services/authService';

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  isInitialized: boolean;
}

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
  isInitialized: false,
};

// Async thunks
export const checkAuthStatus = createAsyncThunk(
  'auth/checkAuthStatus',
  async (_, { rejectWithValue }) => {
    console.log('=== REDUX: checkAuthStatus START ===');
    console.log('Document cookies:', document.cookie);
    console.log(
      'LocalStorage auth state:',
      localStorage.getItem('auth_authenticated')
    );

    try {
      const response = await authService.getAuthStatus();
      console.log('=== REDUX: checkAuthStatus SUCCESS ===');
      console.log('Response:', response);
      console.log('Response.data:', response.data);
      console.log('Response.data.authenticated:', response.data?.authenticated);
      console.log('Response.data.user:', response.data?.user);
      return response.data;
    } catch (error: unknown) {
      console.error('=== REDUX: checkAuthStatus ERROR ===');
      console.error('Auth status check failed:', error);
      const errorMessage =
        error instanceof Error
          ? error.message
          : 'Failed to check authentication status';
      return rejectWithValue(errorMessage);
    }
  }
);

export const getCurrentUser = createAsyncThunk(
  'auth/getCurrentUser',
  async (_, { rejectWithValue }) => {
    try {
      const response = await authService.getCurrentUser();
      return response.data;
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : 'Failed to get current user';
      return rejectWithValue(errorMessage);
    }
  }
);

export const logout = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      const response = await authService.logout();
      return response.data;
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : 'Failed to logout';
      return rejectWithValue(errorMessage);
    }
  }
);

export const initializeAuth = createAsyncThunk(
  'auth/initializeAuth',
  async () => {
    console.log('=== INITIALIZE AUTH START ===');
    console.log('Document cookies:', document.cookie);
    console.log(
      'LocalStorage auth state:',
      localStorage.getItem('auth_authenticated')
    );
    console.log('LocalStorage user:', localStorage.getItem('auth_user'));

    try {
      // First check if we have stored auth state
      const storedUser = authService.getStoredUser();
      const isStoredAuthenticated = authService.isAuthenticated();
      console.log('Stored auth state:', { storedUser, isStoredAuthenticated });

      if (isStoredAuthenticated && storedUser) {
        // Verify with server
        console.log('Verifying stored auth with server...');
        const response = await authService.getAuthStatus();
        console.log('Server response:', response);
        if (response.data.authenticated && response.data.user) {
          console.log('Server verification successful');
          return response.data;
        } else {
          // Server says not authenticated, clear local state
          console.log('Server verification failed, clearing local state');
          authService.clearAuthState();
          return { authenticated: false, user: null };
        }
      }

      // No stored state, check with server
      console.log('No stored state, checking with server...');
      const response = await authService.getAuthStatus();
      console.log('Server auth status:', response);
      return response.data;
    } catch (error) {
      console.error('=== INITIALIZE AUTH ERROR ===', error);
      // If there's an error, assume not authenticated
      authService.clearAuthState();
      return { authenticated: false, user: null };
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
      state.isAuthenticated = true;
      state.error = null;
      authService.setAuthState(true, action.payload);
    },
    clearUser: state => {
      state.user = null;
      state.isAuthenticated = false;
      state.error = null;
      authService.clearAuthState();
    },
    setError: (state, action: PayloadAction<string>) => {
      state.error = action.payload;
    },
    clearError: state => {
      state.error = null;
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
  },
  extraReducers: builder => {
    builder
      // Check Auth Status
      .addCase(checkAuthStatus.pending, state => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(checkAuthStatus.fulfilled, (state, action) => {
        console.log('=== REDUX: checkAuthStatus.fulfilled ===');
        console.log('Action payload:', action.payload);
        console.log(
          'Action payload.authenticated:',
          action.payload?.authenticated
        );
        console.log('Action payload.user:', action.payload?.user);

        state.isLoading = false;
        state.isAuthenticated = action.payload.authenticated;
        state.user = action.payload.user || null;
        state.error = null;
        authService.setAuthState(
          action.payload.authenticated,
          action.payload.user || undefined
        );

        console.log('Updated state:', {
          isAuthenticated: state.isAuthenticated,
          user: state.user,
          isLoading: state.isLoading,
          error: state.error,
        });
      })
      .addCase(checkAuthStatus.rejected, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = action.payload as string;
        authService.clearAuthState();
      })
      // Get Current User
      .addCase(getCurrentUser.pending, state => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(getCurrentUser.fulfilled, (state, action) => {
        state.isLoading = false;
        if (action.payload.user) {
          state.user = action.payload.user;
          state.isAuthenticated = true;
          authService.setAuthState(true, action.payload.user);
        }
        state.error = null;
      })
      .addCase(getCurrentUser.rejected, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = action.payload as string;
        authService.clearAuthState();
      })
      // Logout
      .addCase(logout.pending, state => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(logout.fulfilled, state => {
        state.isLoading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = null;
        authService.clearAuthState();
      })
      .addCase(logout.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
        // Still clear local state even if logout failed
        state.isAuthenticated = false;
        state.user = null;
        authService.clearAuthState();
      })
      // Initialize Auth
      .addCase(initializeAuth.pending, state => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(initializeAuth.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = action.payload.authenticated;
        state.user = action.payload.user || null;
        state.error = null;
        state.isInitialized = true;
        authService.setAuthState(
          action.payload.authenticated,
          action.payload.user || undefined
        );
      })
      .addCase(initializeAuth.rejected, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = false;
        state.user = null;
        state.error = action.payload as string;
        state.isInitialized = true;
        authService.clearAuthState();
      });
  },
});

export const { setUser, clearUser, setError, clearError, setLoading } =
  authSlice.actions;
export default authSlice.reducer;
