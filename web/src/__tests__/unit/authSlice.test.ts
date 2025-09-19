import { configureStore } from '@reduxjs/toolkit';
import authReducer, {
    checkAuthStatus,
    getCurrentUser,
    logout,
    initializeAuth,
    setUser,
    clearUser,
    setError,
    clearError,
    setLoading,
    AuthState,
} from '@/store/reducers/authSlice';
import { User } from '@/services/authService';

// Mock the authService
jest.mock('@/services/authService', () => ({
    authService: {
        getAuthStatus: jest.fn(),
        getCurrentUser: jest.fn(),
        logout: jest.fn(),
        getStoredUser: jest.fn(),
        isAuthenticated: jest.fn(),
        setAuthState: jest.fn(),
        clearAuthState: jest.fn(),
    },
}));

import { authService } from '@/services/authService';

const mockAuthService = authService as jest.Mocked<typeof authService>;

describe('AuthSlice', () => {
    let store: ReturnType<typeof configureStore<{ auth: AuthState }>>;

    beforeEach(() => {
        store = configureStore({
            reducer: {
                auth: authReducer,
            },
        });
        jest.clearAllMocks();
    });

    describe('initial state', () => {
        it('should have correct initial state', () => {
            const state = store.getState().auth;

            expect(state.user).toBeNull();
            expect(state.isAuthenticated).toBe(false);
            expect(state.isLoading).toBe(false);
            expect(state.error).toBeNull();
            expect(state.isInitialized).toBe(false);
        });
    });

    describe('synchronous actions', () => {
        it('should set user correctly', () => {
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

            store.dispatch(setUser(user));

            const state = store.getState().auth;
            expect(state.user).toEqual(user);
            expect(state.isAuthenticated).toBe(true);
            expect(state.error).toBeNull();
            expect(mockAuthService.setAuthState).toHaveBeenCalledWith(true, user);
        });

        it('should clear user correctly', () => {
            // First set a user
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

            store.dispatch(setUser(user));

            // Then clear user
            store.dispatch(clearUser());

            const state = store.getState().auth;
            expect(state.user).toBeNull();
            expect(state.isAuthenticated).toBe(false);
            expect(state.error).toBeNull();
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });

        it('should set error correctly', () => {
            const errorMessage = 'Test error message';

            store.dispatch(setError(errorMessage));

            const state = store.getState().auth;
            expect(state.error).toBe(errorMessage);
        });

        it('should clear error correctly', () => {
            // First set an error
            store.dispatch(setError('Test error'));

            // Then clear error
            store.dispatch(clearError());

            const state = store.getState().auth;
            expect(state.error).toBeNull();
        });

        it('should set loading correctly', () => {
            store.dispatch(setLoading(true));

            const state = store.getState().auth;
            expect(state.isLoading).toBe(true);

            store.dispatch(setLoading(false));

            const newState = store.getState().auth as AuthState;
            expect(newState.isLoading).toBe(false);
        });
    });

    describe('checkAuthStatus async thunk', () => {
        it('should handle successful authentication status check', async () => {
            const mockAuthStatus = {
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

            mockAuthService.getAuthStatus.mockResolvedValue({
                data: mockAuthStatus,
                success: true,
                message: 'Success',
            });

            await store.dispatch(checkAuthStatus());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(true);
            expect(state.user).toEqual(mockAuthStatus.user);
            expect(state.error).toBeNull();
            expect(mockAuthService.setAuthState).toHaveBeenCalledWith(
                true,
                mockAuthStatus.user
            );
        });

        it('should handle failed authentication status check', async () => {
            const errorMessage = 'Authentication failed';

            mockAuthService.getAuthStatus.mockRejectedValue(new Error(errorMessage));

            await store.dispatch(checkAuthStatus());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBe(errorMessage);
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });

        it('should set loading state during authentication status check', async () => {
            mockAuthService.getAuthStatus.mockImplementation(
                () => new Promise(resolve => setTimeout(resolve, 100))
            );

            const promise = store.dispatch(checkAuthStatus());

            // Check loading state is true
            let state = store.getState().auth as AuthState;
            expect(state.isLoading).toBe(true);
            expect(state.error).toBeNull();

            await promise;

            // Check loading state is false after completion
            state = store.getState().auth as AuthState;
            expect(state.isLoading).toBe(false);
        });
    });

    describe('getCurrentUser async thunk', () => {
        it('should handle successful get current user', async () => {
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

            mockAuthService.getCurrentUser.mockResolvedValue({
                data: { user: mockUser, message: 'Success' },
                success: true,
                message: 'Success',
            });

            await store.dispatch(getCurrentUser());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(true);
            expect(state.user).toEqual(mockUser);
            expect(state.error).toBeNull();
            expect(mockAuthService.setAuthState).toHaveBeenCalledWith(true, mockUser);
        });

        it('should handle failed get current user', async () => {
            const errorMessage = 'Failed to get current user';

            mockAuthService.getCurrentUser.mockRejectedValue(new Error(errorMessage));

            await store.dispatch(getCurrentUser());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBe(errorMessage);
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });
    });

    describe('logout async thunk', () => {
        it('should handle successful logout', async () => {
            // First set a user
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

            store.dispatch(setUser(user));

            mockAuthService.logout.mockResolvedValue({
                data: { message: 'Logged out successfully' },
                success: true,
                message: 'Success',
            });

            await store.dispatch(logout());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBeNull();
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });

        it('should handle failed logout', async () => {
            const errorMessage = 'Failed to logout';

            mockAuthService.logout.mockRejectedValue(new Error(errorMessage));

            await store.dispatch(logout());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBe(errorMessage);
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });
    });

    describe('initializeAuth async thunk', () => {
        it('should initialize with stored authenticated state', async () => {
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

            const mockAuthStatus = {
                authenticated: true,
                user: mockUser,
            };

            mockAuthService.getStoredUser.mockReturnValue(mockUser);
            mockAuthService.isAuthenticated.mockReturnValue(true);
            mockAuthService.getAuthStatus.mockResolvedValue({
                data: mockAuthStatus,
                success: true,
                message: 'Success',
            });

            await store.dispatch(initializeAuth());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(true);
            expect(state.user).toEqual(mockUser);
            expect(state.error).toBeNull();
            expect(state.isInitialized).toBe(true);
            expect(mockAuthService.setAuthState).toHaveBeenCalledWith(true, mockUser);
        });

        it('should clear local state when server verification fails', async () => {
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

            const mockAuthStatus = {
                authenticated: false,
                user: undefined,
            };

            mockAuthService.getStoredUser.mockReturnValue(mockUser);
            mockAuthService.isAuthenticated.mockReturnValue(true);
            mockAuthService.getAuthStatus.mockResolvedValue({
                data: mockAuthStatus,
                success: true,
                message: 'Success',
            });

            await store.dispatch(initializeAuth());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBeNull();
            expect(state.isInitialized).toBe(true);
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });

        it('should check server when no stored state', async () => {
            const mockAuthStatus = {
                authenticated: false,
                user: undefined,
            };

            mockAuthService.getStoredUser.mockReturnValue(null);
            mockAuthService.isAuthenticated.mockReturnValue(false);
            mockAuthService.getAuthStatus.mockResolvedValue({
                data: mockAuthStatus,
                success: true,
                message: 'Success',
            });

            await store.dispatch(initializeAuth());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBeNull();
            expect(state.isInitialized).toBe(true);
            expect(mockAuthService.getAuthStatus).toHaveBeenCalled();
        });

        it('should handle initialization error', async () => {
            const errorMessage = 'Initialization failed';

            mockAuthService.getStoredUser.mockReturnValue(null);
            mockAuthService.isAuthenticated.mockReturnValue(false);
            mockAuthService.getAuthStatus.mockRejectedValue(new Error(errorMessage));

            await store.dispatch(initializeAuth());

            const state = store.getState().auth;
            expect(state.isLoading).toBe(false);
            expect(state.isAuthenticated).toBe(false);
            expect(state.user).toBeNull();
            expect(state.error).toBeNull(); // Error is caught and handled gracefully
            expect(state.isInitialized).toBe(true);
            expect(mockAuthService.clearAuthState).toHaveBeenCalled();
        });
    });
});
