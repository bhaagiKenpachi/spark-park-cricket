'use client';

import { useState } from 'react';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { logout } from '@/store/reducers/authSlice';
import { User, LogOut, Settings } from 'lucide-react';

export function UserMenu() {
  const dispatch = useAppDispatch();
  const { user, isLoading } = useAppSelector(state => state.auth);
  const [isOpen, setIsOpen] = useState(false);

  const handleLogout = async () => {
    try {
      await dispatch(logout()).unwrap();
      setIsOpen(false);
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  if (!user) {
    return null;
  }

  return (
    <div className="relative">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 p-2"
      >
        {user.picture ? (
          <Image
            src={user.picture}
            alt={user.name}
            width={24}
            height={24}
            className="h-6 w-6 rounded-full"
          />
        ) : (
          <User className="h-6 w-6" />
        )}
        <span className="hidden sm:inline text-sm font-medium">
          {user.name}
        </span>
      </Button>

      {isOpen && (
        <>
          {/* Backdrop */}
          <div
            className="fixed inset-0 z-10"
            onClick={() => setIsOpen(false)}
          />

          {/* Menu */}
          <Card className="absolute right-0 top-full mt-2 w-64 z-20 p-2 shadow-lg">
            <div className="space-y-2">
              {/* User Info */}
              <div className="px-3 py-2 border-b">
                <div className="flex items-center gap-3">
                  {user.picture ? (
                    <Image
                      src={user.picture}
                      alt={user.name}
                      width={40}
                      height={40}
                      className="h-10 w-10 rounded-full"
                    />
                  ) : (
                    <div className="h-10 w-10 rounded-full bg-gray-200 flex items-center justify-center">
                      <User className="h-5 w-5 text-gray-600" />
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-900 truncate">
                      {user.name}
                    </p>
                    <p className="text-xs text-gray-500 truncate">
                      {user.email}
                    </p>
                  </div>
                </div>
              </div>

              {/* Menu Items */}
              <div className="space-y-1">
                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full justify-start gap-2"
                  onClick={() => setIsOpen(false)}
                >
                  <Settings className="h-4 w-4" />
                  Settings
                </Button>

                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full justify-start gap-2 text-red-600 hover:text-red-700 hover:bg-red-50"
                  onClick={handleLogout}
                  disabled={isLoading}
                >
                  <LogOut className="h-4 w-4" />
                  {isLoading ? 'Signing out...' : 'Sign out'}
                </Button>
              </div>
            </div>
          </Card>
        </>
      )}
    </div>
  );
}
