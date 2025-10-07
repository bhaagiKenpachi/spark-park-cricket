'use client';

import { useState } from 'react';
import Link from 'next/link';
import { VoteList } from '@/components/VoteList';
import { VoteForm } from '@/components/VoteForm';
import { VoteView } from '@/components/VoteView';
import { Vote } from '@/types/vote';
import { Button } from '@/components/ui/button';
import { ArrowLeft, Plus, Home, Vote as VoteIcon } from 'lucide-react';
import { LoginButton } from '@/components/auth/LoginButton';
import { UserMenu } from '@/components/auth/UserMenu';
import { useAppSelector } from '@/store/hooks';

type ViewMode = 'list' | 'create' | 'edit' | 'view';

export default function VotesPage(): React.JSX.Element {
    const { isAuthenticated } = useAppSelector(state => state.auth);
    const [viewMode, setViewMode] = useState<ViewMode>('list');
    const [selectedVote, setSelectedVote] = useState<Vote | null>(null);
    const [selectedVoteId, setSelectedVoteId] = useState<string | null>(null);

    const handleCreateVote = () => {
        setSelectedVote(null);
        setViewMode('create');
    };

    const handleEditVote = (vote: Vote) => {
        setSelectedVote(vote);
        setViewMode('edit');
    };

    const handleViewVote = (voteId: string) => {
        setSelectedVoteId(voteId);
        setViewMode('view');
    };

    const handleBackToList = () => {
        setViewMode('list');
        setSelectedVote(null);
        setSelectedVoteId(null);
    };

    const handleFormSuccess = () => {
        setViewMode('list');
        setSelectedVote(null);
    };

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Top Navigation Bar */}
            <nav className="bg-white border-b sticky top-0 z-50 shadow-sm">
                <div className="w-full max-w-md mx-auto px-3 py-2">
                    <div className="flex items-center justify-between">
                        {/* Left: Logo/Home */}
                        <Link href="/" className="flex items-center gap-2">
                            <div className="bg-blue-600 p-1.5 rounded-lg">
                                <Home className="h-4 w-4 text-white" />
                            </div>
                            <span className="text-sm font-semibold text-gray-900">Spark Park</span>
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

            {/* Sub Navigation Bar for Votes */}
            <div className="bg-white border-b">
                <div className="w-full max-w-md mx-auto px-3 py-3">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            {viewMode !== 'list' ? (
                                <Button
                                    onClick={handleBackToList}
                                    variant="ghost"
                                    size="sm"
                                    className="p-1.5 h-auto"
                                >
                                    <ArrowLeft className="h-4 w-4" />
                                </Button>
                            ) : (
                                <div className="bg-purple-100 p-1.5 rounded-lg">
                                    <VoteIcon className="h-4 w-4 text-purple-600" />
                                </div>
                            )}
                            <div>
                                <h1 className="text-base font-bold text-gray-900">
                                    {viewMode === 'list' && 'Votes'}
                                    {viewMode === 'create' && 'Create Vote'}
                                    {viewMode === 'edit' && 'Edit Vote'}
                                    {viewMode === 'view' && 'Vote Details'}
                                </h1>
                                {viewMode === 'list' && (
                                    <p className="text-xs text-gray-500">Browse and participate</p>
                                )}
                            </div>
                        </div>

                        {viewMode === 'list' && (
                            <Button
                                onClick={handleCreateVote}
                                size="sm"
                                className="flex items-center gap-1 h-8 px-3 bg-purple-600 hover:bg-purple-700"
                            >
                                <Plus className="h-3.5 w-3.5" />
                                <span className="text-xs font-medium">New</span>
                            </Button>
                        )}
                    </div>
                </div>
            </div>

            <main>
                {viewMode === 'list' && (
                    <VoteList
                        onCreateVote={handleCreateVote}
                        onViewVote={handleViewVote}
                        onEditVote={handleEditVote}
                    />
                )}

                {viewMode === 'create' && (
                    <VoteForm
                        onSuccess={handleFormSuccess}
                        onCancel={handleBackToList}
                    />
                )}

                {viewMode === 'edit' && selectedVote && (
                    <VoteForm
                        vote={selectedVote}
                        onSuccess={handleFormSuccess}
                        onCancel={handleBackToList}
                    />
                )}

                {viewMode === 'view' && selectedVoteId && (
                    <VoteView
                        voteId={selectedVoteId}
                        onBack={handleBackToList}
                    />
                )}
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
