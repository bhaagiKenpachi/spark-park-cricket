'use client';

import { Button } from '@/components/ui/button';
import { LogIn } from 'lucide-react';

export function LoginButton() {
  const handleGoogleLogin = () => {
    // Redirect directly to backend OAuth endpoint
    // The backend will handle the OAuth flow and redirect back
    const apiBaseUrl =
      process.env.NEXT_PUBLIC_API_URL || 'https://spark-park.dojima.foundation/api/v1';
    const oauthUrl = `${apiBaseUrl}/auth/google`;

    window.location.href = oauthUrl;
  };

  return (
    <Button
      onClick={handleGoogleLogin}
      className="flex items-center gap-1.5 bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 h-8 px-2.5"
      data-cy="login-button"
      size="sm"
    >
      <LogIn className="h-3.5 w-3.5" />
      <span className="text-xs font-medium">Sign in</span>
    </Button>
  );
}
