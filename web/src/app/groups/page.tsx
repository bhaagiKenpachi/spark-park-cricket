'use client';

import { useState } from 'react';
import Link from 'next/link';
import { GroupList } from '@/components/GroupList';
import { GroupView } from '@/components/GroupView';
import { GroupWithCreator } from '@/types/group';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Users } from 'lucide-react';
import { LoginButton } from '@/components/auth/LoginButton';
import { UserMenu } from '@/components/auth/UserMenu';
import { useAppSelector } from '@/store/hooks';

export default function GroupsPage(): React.JSX.Element {
    const { isAuthenticated } = useAppSelector(state => state.auth);
    const [selectedGroup, setSelectedGroup] = useState<GroupWithCreator | null>(null);

    const handleGroupSelect = (group: GroupWithCreator) => {
        setSelectedGroup(group);
    };

    const handleBack = () => {
        setSelectedGroup(null);
    };

    const handleEdit = (group: GroupWithCreator) => {
        setSelectedGroup(group);
    };

    if (selectedGroup) {
        return (
            <GroupView
                group={selectedGroup}
                onBack={handleBack}
                onEdit={handleEdit}
            />
        );
    }

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Top Navigation Bar */}
            <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
                <div className="w-full max-w-md mx-auto px-3 py-2">
                    <div className="flex items-center justify-between">
                        {/* Left: Back to Home */}
                        <Link href="/">
                            <Button
                                variant="ghost"
                                size="sm"
                                className="flex items-center gap-2 p-1.5 h-auto"
                            >
                                <ArrowLeft className="h-4 w-4" />
                                <span className="text-sm font-medium">Back to Home</span>
                            </Button>
                        </Link>

                        {/* Right: Auth */}
                        <div className="flex items-center gap-2">
                            {isAuthenticated ? (
                                <UserMenu />
                            ) : (
                                <LoginButton />
                            )}
                        </div>
                    </div>
                </div>
            </nav>

            {/* Sub Navigation Bar for Groups */}
            <div className="bg-white border-b">
                <div className="w-full max-w-md mx-auto px-3 py-3">
                    <div className="flex items-center gap-2">
                        <div className="bg-blue-100 p-1.5 rounded-lg">
                            <Users className="h-4 w-4 text-blue-600" />
                        </div>
                        <div>
                            <h1 className="text-base font-bold text-gray-900">Groups</h1>
                            <p className="text-xs text-gray-500">Manage voting groups</p>
                        </div>
                    </div>
                </div>
            </div>

            <main className="w-full max-w-md mx-auto">
                <GroupList onGroupSelect={handleGroupSelect} />
            </main>

            <footer className="bg-white border-t">
                <div className="w-full max-w-md mx-auto py-3 px-4">
                    <p className="text-center text-gray-500 text-xs">
                        © 2024 Spark Park Cricket. All rights reserved.
                    </p>
                </div>
            </footer>
        </div>
    );
}
