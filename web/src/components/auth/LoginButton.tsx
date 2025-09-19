'use client';

import { Button } from '@/components/ui/button';
import { LogIn } from 'lucide-react';

export function LoginButton() {
  const handleGoogleLogin = () => {
    console.log('=== FRONTEND OAUTH LOGIN START ===');
    console.log('Current URL:', window.location.href);
    console.log('Current domain:', window.location.hostname);
    console.log('Current port:', window.location.port);
    console.log('Document cookies:', document.cookie);
    console.log(
      'LocalStorage auth state:',
      localStorage.getItem('auth_authenticated')
    );
    console.log('LocalStorage user:', localStorage.getItem('auth_user'));

    // Redirect directly to backend OAuth endpoint
    // The backend will handle the OAuth flow and redirect back
    const oauthUrl = 'http://localhost:8080/api/v1/auth/google';
    console.log('Redirecting to OAuth URL:', oauthUrl);

    window.location.href = oauthUrl;
  };

  return (
    <Button
      onClick={handleGoogleLogin}
      className="flex items-center gap-2 bg-white text-gray-700 border border-gray-300 hover:bg-gray-50"
      data-cy="login-button"
    >
      <LogIn className="h-4 w-4" />
      Sign in with Google
    </Button>
  );
}
