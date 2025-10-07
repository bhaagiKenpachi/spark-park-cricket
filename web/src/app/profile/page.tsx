'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { useAppSelector } from '@/store/hooks';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import EditProfileForm from '@/components/auth/EditProfileForm';
import {
    Home,
    User,
    Mail,
    Calendar,
    CheckCircle,
    Vote,
    Shield,
    XCircle,
    Edit
} from 'lucide-react';

export default function ProfilePage(): React.JSX.Element {
    const router = useRouter();
    const { user, isAuthenticated, isLoading } = useAppSelector(state => state.auth);
    const [isEditing, setIsEditing] = useState(false);

    useEffect(() => {
        if (!isLoading && !isAuthenticated) {
            router.push('/');
        }
    }, [isAuthenticated, isLoading, router]);

    if (isLoading) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Loading profile...</p>
                </div>
            </div>
        );
    }

    if (!user) {
        return (
            <div className="min-h-screen bg-gray-50 flex items-center justify-center">
                <div className="text-center">
                    <p className="text-gray-600">Please sign in to view your profile</p>
                </div>
            </div>
        );
    }

    const formatDate = (dateString?: string) => {
        if (!dateString) return 'N/A';
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Top Navigation Bar */}
            <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
                <div className="w-full max-w-md mx-auto px-3 py-2">
                    <div className="flex items-center justify-between">
                        <Link href="/" className="flex items-center gap-2">
                            <div className="bg-blue-600 p-1.5 rounded-lg">
                                <Home className="h-4 w-4 text-white" />
                            </div>
                            <span className="text-sm font-semibold text-gray-900">Spark Park</span>
                        </Link>

                        <Link href="/votes">
                            <Button variant="ghost" size="sm" className="h-8 px-2.5 flex items-center gap-1.5">
                                <Vote className="h-3.5 w-3.5" />
                                <span className="text-xs font-medium">Votes</span>
                            </Button>
                        </Link>
                    </div>
                </div>
            </nav>

            {/* Sub Navigation Bar */}
            <div className="bg-white border-b">
                <div className="w-full max-w-md mx-auto px-3 py-3">
                    <div className="flex items-center gap-2">
                        <div className="bg-purple-100 p-1.5 rounded-lg">
                            <User className="h-4 w-4 text-purple-600" />
                        </div>
                        <div>
                            <h1 className="text-base font-bold text-gray-900">Profile</h1>
                            <p className="text-xs text-gray-500">Your account information</p>
                        </div>
                    </div>
                </div>
            </div>

            <main className="w-full max-w-md mx-auto px-4 py-4">
                {/* Edit Profile Form */}
                {isEditing ? (
                    <div className="mb-4">
                        <EditProfileForm
                            onCancel={() => setIsEditing(false)}
                            onSuccess={() => setIsEditing(false)}
                        />
                    </div>
                ) : (
                    <>
                        {/* Profile Picture and Name */}
                        <div className="bg-white rounded-lg shadow-sm border mb-4">
                            <div className="p-6 text-center">
                                {user.picture ? (
                                    <Image
                                        src={user.picture}
                                        alt={user.name}
                                        width={80}
                                        height={80}
                                        className="h-20 w-20 rounded-full mx-auto mb-3"
                                    />
                                ) : (
                                    <div className="h-20 w-20 rounded-full bg-gray-200 flex items-center justify-center mx-auto mb-3">
                                        <User className="h-10 w-10 text-gray-600" />
                                    </div>
                                )}
                                <h2 className="text-lg font-bold text-gray-900">{user.name}</h2>
                                <p className="text-sm text-gray-500 mt-1">{user.email}</p>

                                {user.email_verified && (
                                    <div className="mt-3 inline-flex items-center gap-1 px-3 py-1 bg-green-50 border border-green-200 rounded-full">
                                        <CheckCircle className="h-3 w-3 text-green-600" />
                                        <span className="text-xs text-green-700 font-medium">Verified</span>
                                    </div>
                                )}

                                {/* Edit Button */}
                                <Button
                                    onClick={() => setIsEditing(true)}
                                    variant="outline"
                                    size="sm"
                                    className="mt-4 flex items-center gap-2 mx-auto"
                                >
                                    <Edit className="h-3.5 w-3.5" />
                                    Edit Profile
                                </Button>
                            </div>
                        </div>
                    </>
                )}

                {/* Account Details */}
                <div className="bg-white rounded-lg shadow-sm border">
                    <div className="p-3 border-b">
                        <div className="flex items-center gap-2">
                            <Shield className="h-4 w-4 text-gray-600" />
                            <span className="text-sm font-medium text-gray-700">Account Details</span>
                        </div>
                    </div>
                    <div className="p-4 space-y-4">
                        <div className="space-y-1">
                            <label className="text-xs font-medium text-gray-500">User ID</label>
                            <p className="text-sm text-gray-900 font-mono bg-gray-50 px-2 py-1.5 rounded border text-xs break-all">
                                {user.id}
                            </p>
                        </div>

                        <div className="space-y-1">
                            <label className="text-xs font-medium text-gray-500">Email</label>
                            <div className="flex items-center gap-2">
                                <Mail className="h-4 w-4 text-gray-400" />
                                <p className="text-sm text-gray-900">{user.email}</p>
                            </div>
                        </div>

                        <div className="space-y-1">
                            <label className="text-xs font-medium text-gray-500">Google ID</label>
                            <p className="text-sm text-gray-900 font-mono bg-gray-50 px-2 py-1.5 rounded border text-xs break-all">
                                {user.google_id}
                            </p>
                        </div>

                        <div className="space-y-1">
                            <label className="text-xs font-medium text-gray-500">Account Created</label>
                            <div className="flex items-center gap-2">
                                <Calendar className="h-4 w-4 text-gray-400" />
                                <p className="text-sm text-gray-900">{formatDate(user.created_at)}</p>
                            </div>
                        </div>

                        {user.last_login_at && (
                            <div className="space-y-1">
                                <label className="text-xs font-medium text-gray-500">Last Login</label>
                                <div className="flex items-center gap-2">
                                    <Calendar className="h-4 w-4 text-gray-400" />
                                    <p className="text-sm text-gray-900">{formatDate(user.last_login_at)}</p>
                                </div>
                            </div>
                        )}

                        <div className="space-y-1">
                            <label className="text-xs font-medium text-gray-500">Email Verification</label>
                            <div className="flex items-center gap-2">
                                {user.email_verified ? (
                                    <>
                                        <CheckCircle className="h-4 w-4 text-green-600" />
                                        <span className="text-sm text-green-700">Verified</span>
                                    </>
                                ) : (
                                    <>
                                        <XCircle className="h-4 w-4 text-red-600" />
                                        <span className="text-sm text-red-700">Not Verified</span>
                                    </>
                                )}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Quick Actions - Only show when not editing */}
                {!isEditing && (
                    <div className="mt-6">
                        <Link href="/">
                            <Button variant="outline" className="w-full flex items-center justify-center gap-2">
                                <Home className="h-4 w-4" />
                                Back to Home
                            </Button>
                        </Link>
                    </div>
                )}
            </main>

            <footer className="bg-white border-t mt-8">
                <div className="w-full max-w-md mx-auto py-3 px-4">
                    <p className="text-center text-gray-500 text-xs">
                        © 2024 Spark Park Cricket. All rights reserved.
                    </p>
                </div>
            </footer>
        </div>
    );
}

