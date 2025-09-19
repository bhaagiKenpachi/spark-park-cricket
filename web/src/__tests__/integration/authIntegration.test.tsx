import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { BrowserRouter } from 'react-router-dom';
import authReducer from '@/store/reducers/authSlice';
import { ProtectedRoute } from '@/components/auth/ProtectedRoute';
import { LoginButton } from '@/components/auth/LoginButton';
import { UserMenu } from '@/components/auth/UserMenu';
import { AuthProvider } from '@/components/auth/AuthProvider';
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
        initiateGoogleLogin: jest.fn(),
    },
}));

import { authService } from '@/services/authService';

const mockAuthService = authService as jest.Mocked<typeof authService>;

// Mock window.location
const mockLocation = {
    href: '',
    assign: jest.fn(),
    replace: jest.fn(),
    reload: jest.fn(),
};

// Delete the existing location property first
delete (window as unknown as { location?: Location }).location;
(window as unknown as { location: Location }).location = mockLocation;

describe('Authentication Integration Tests', () => {
    let store: ReturnType<typeof configureStore>;

    beforeEach(() => {
        store = configureStore({
            reducer: {
                auth: authReducer,
            },
        });
        jest.clearAllMocks();
        mockLocation.href = '';
    });

    const renderWithProviders = (component: React.ReactElement) => {
        return render(
            <Provider store={store}>
                <BrowserRouter>{component}</BrowserRouter>
            </Provider>
        );
    };

    const renderWithAuthProvider = (component: React.ReactElement) => {
        return render(
            <Provider store={store}>
                <BrowserRouter>
                    <AuthProvider>{component}</AuthProvider>
                </BrowserRouter>
            </Provider>
        );
    };

    describe('ProtectedRoute Component', () => {
        it('should show loading state when auth is not initialized', () => {
            renderWithProviders(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>
            );

            expect(screen.getByText('Loading...')).toBeInTheDocument();
            expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
        });

        it('should show login prompt when user is not authenticated', () => {
            // Set initial state to initialized but not authenticated
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: false, user: null },
            });

            renderWithProviders(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>
            );

            expect(screen.getByText('Authentication Required')).toBeInTheDocument();
            expect(
                screen.getByText('Please sign in to access this feature.')
            ).toBeInTheDocument();
            expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
        });

        it('should show protected content when user is authenticated', () => {
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

            // Set initial state to authenticated
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: true, user },
            });

            renderWithProviders(
                <ProtectedRoute>
                    <div>Protected Content</div>
                </ProtectedRoute>
            );

            expect(screen.getByText('Protected Content')).toBeInTheDocument();
            expect(
                screen.queryByText('Authentication Required')
            ).not.toBeInTheDocument();
        });

        it('should show custom fallback when provided', () => {
            // Set initial state to initialized but not authenticated
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: false, user: null },
            });

            renderWithProviders(
                <ProtectedRoute fallback={<div>Custom Fallback</div>}>
                    <div>Protected Content</div>
                </ProtectedRoute>
            );

            expect(screen.getByText('Custom Fallback')).toBeInTheDocument();
            expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
            expect(
                screen.queryByText('Authentication Required')
            ).not.toBeInTheDocument();
        });
    });

    describe('LoginButton Component', () => {
        it('should render login button', () => {
            renderWithProviders(<LoginButton />);

            expect(screen.getByText('Sign in with Google')).toBeInTheDocument();
        });

        it('should initiate Google login when clicked', () => {
            renderWithProviders(<LoginButton />);

            const loginButton = screen.getByText('Sign in with Google');
            fireEvent.click(loginButton);

            expect(mockAuthService.initiateGoogleLogin).toHaveBeenCalled();
        });

        it('should show loading state when clicked', () => {
            renderWithProviders(<LoginButton />);

            const loginButton = screen.getByText('Sign in with Google');
            fireEvent.click(loginButton);

            // The component should show loading state
            expect(screen.getByText('Signing in...')).toBeInTheDocument();
        });
    });

    describe('UserMenu Component', () => {
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

        it('should render user menu when authenticated', () => {
            // Set authenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: true, user: mockUser },
            });

            renderWithProviders(<UserMenu />);

            expect(screen.getByText('Test User')).toBeInTheDocument();
            expect(screen.getByText('test@example.com')).toBeInTheDocument();
        });

        it('should not render when not authenticated', () => {
            // Set unauthenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: false, user: null },
            });

            renderWithProviders(<UserMenu />);

            expect(screen.queryByText('Test User')).not.toBeInTheDocument();
        });

        it('should show dropdown menu when clicked', () => {
            // Set authenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: true, user: mockUser },
            });

            renderWithProviders(<UserMenu />);

            const userButton = screen.getByText('Test User');
            fireEvent.click(userButton);

            expect(screen.getByText('Profile')).toBeInTheDocument();
            expect(screen.getByText('Sign out')).toBeInTheDocument();
        });

        it('should handle logout when sign out is clicked', async () => {
            mockAuthService.logout.mockResolvedValue({
                data: { message: 'Logged out successfully' },
                success: true,
                message: 'Success',
            });

            // Set authenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: true, user: mockUser },
            });

            renderWithProviders(<UserMenu />);

            const userButton = screen.getByText('Test User');
            fireEvent.click(userButton);

            const signOutButton = screen.getByText('Sign out');
            fireEvent.click(signOutButton);

            await waitFor(() => {
                expect(mockAuthService.logout).toHaveBeenCalled();
            });
        });
    });

    describe('AuthProvider Integration', () => {
        it('should initialize authentication on mount', async () => {
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

            mockAuthService.getStoredUser.mockReturnValue(null);
            mockAuthService.isAuthenticated.mockReturnValue(false);
            mockAuthService.getAuthStatus.mockResolvedValue({
                data: mockAuthStatus,
                success: true,
                message: 'Success',
            });

            renderWithAuthProvider(
                <div>
                    <ProtectedRoute>
                        <div>Protected Content</div>
                    </ProtectedRoute>
                </div>
            );

            // Should show loading initially
            expect(screen.getByText('Loading...')).toBeInTheDocument();

            // Wait for initialization to complete
            await waitFor(() => {
                expect(screen.getByText('Protected Content')).toBeInTheDocument();
            });

            expect(mockAuthService.getAuthStatus).toHaveBeenCalled();
        });

        it('should handle authentication initialization error', async () => {
            mockAuthService.getStoredUser.mockReturnValue(null);
            mockAuthService.isAuthenticated.mockReturnValue(false);
            mockAuthService.getAuthStatus.mockRejectedValue(
                new Error('Network error')
            );

            renderWithAuthProvider(
                <div>
                    <ProtectedRoute>
                        <div>Protected Content</div>
                    </ProtectedRoute>
                </div>
            );

            // Should show loading initially
            expect(screen.getByText('Loading...')).toBeInTheDocument();

            // Wait for initialization to complete with error
            await waitFor(() => {
                expect(screen.getByText('Authentication Required')).toBeInTheDocument();
            });

            expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
        });
    });

    describe('Authentication Flow Integration', () => {
        it('should handle complete authentication flow', async () => {
            // Start with unauthenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: false, user: null },
            });

            renderWithProviders(
                <div>
                    <ProtectedRoute>
                        <div>Protected Content</div>
                    </ProtectedRoute>
                    <LoginButton />
                </div>
            );

            // Should show login prompt
            expect(screen.getByText('Authentication Required')).toBeInTheDocument();
            expect(screen.getByText('Sign in with Google')).toBeInTheDocument();

            // Simulate successful authentication
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

            store.dispatch({
                type: 'auth/checkAuthStatus/fulfilled',
                payload: { authenticated: true, user: mockUser },
            });

            // Should now show protected content
            await waitFor(() => {
                expect(screen.getByText('Protected Content')).toBeInTheDocument();
            });

            expect(
                screen.queryByText('Authentication Required')
            ).not.toBeInTheDocument();
        });

        it('should handle logout flow', async () => {
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

            // Start with authenticated state
            store.dispatch({
                type: 'auth/initializeAuth/fulfilled',
                payload: { authenticated: true, user: mockUser },
            });

            mockAuthService.logout.mockResolvedValue({
                data: { message: 'Logged out successfully' },
                success: true,
                message: 'Success',
            });

            renderWithProviders(
                <div>
                    <ProtectedRoute>
                        <div>Protected Content</div>
                    </ProtectedRoute>
                    <UserMenu />
                </div>
            );

            // Should show protected content and user menu
            expect(screen.getByText('Protected Content')).toBeInTheDocument();
            expect(screen.getByText('Test User')).toBeInTheDocument();

            // Click on user menu and logout
            const userButton = screen.getByText('Test User');
            fireEvent.click(userButton);

            const signOutButton = screen.getByText('Sign out');
            fireEvent.click(signOutButton);

            // Wait for logout to complete
            await waitFor(() => {
                expect(screen.getByText('Authentication Required')).toBeInTheDocument();
            });

            expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
            expect(screen.queryByText('Test User')).not.toBeInTheDocument();
        });
    });
});
