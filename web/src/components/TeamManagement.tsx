'use client';

/**
 * TeamManagement Component - Team Management Interface
 * 
 * AUTHORIZATION INTEGRATION:
 * ==========================
 * This component passes the currentUserId to TeamCard components for
 * creator-based authorization. The TeamCard component uses this to:
 * 
 * 1. Check if current user is team creator (team.created_by === currentUserId)
 * 2. Show/hide management buttons based on creator status
 * 3. Display appropriate visual indicators
 * 
 * IMPORTANT: The currentUserId prop is critical for security.
 * Without it, TeamCard cannot determine authorization levels.
 */

import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '@/store/reducers';
import { VoteWithResults } from '@/types/vote';
import { User, VoteTeamWithPlayers } from '@/types/team';
import TeamCard from './TeamCard';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Shield, RefreshCw, Share2 } from 'lucide-react';
import {
    fetchTeamsRequest,
    createTeamRequest,
    deleteTeamRequest,
    addPlayersRequest,
    removePlayerRequest,
    updateTeamRequest,
} from '@/store/reducers/voteTeamSlice';

interface TeamManagementProps {
    vote: VoteWithResults;
    voters: User[];
    isAuthenticated: boolean;
    currentUserId: string | undefined;
}

export default function TeamManagement({
    vote,
    voters,
    isAuthenticated,
    currentUserId,
}: TeamManagementProps) {
    const dispatch = useDispatch();
    const { teams, loading, error } = useSelector((state: RootState) => state.voteTeam);

    const [showCreateForm, setShowCreateForm] = useState(false);
    const [teamName, setTeamName] = useState('');
    const [teamLetter, setTeamLetter] = useState<'A' | 'B'>('A');
    const [captainId, setCaptainId] = useState('');
    const [errorMessage, setErrorMessage] = useState('');

    // Fetch teams on mount
    useEffect(() => {
        if (vote?.vote?.id) {
            dispatch(fetchTeamsRequest({ voteId: vote.vote.id }));
        }
    }, [dispatch, vote?.vote?.id]);

    const teamA = teams.find((t: VoteTeamWithPlayers) => t.team_letter === 'A');
    const teamB = teams.find((t: VoteTeamWithPlayers) => t.team_letter === 'B');

    const handleCreateTeam = () => {
        if (!teamName.trim() || !captainId || !vote?.vote?.id || errorMessage) {
            return;
        }
        dispatch(createTeamRequest({
            voteId: vote.vote.id,
            teamData: {
                team_name: teamName,
                team_letter: teamLetter,
                captain_id: captainId,
            },
        }));

        // Reset form
        setTeamName('');
        setCaptainId('');
        setShowCreateForm(false);
    };

    const handleAddPlayers = (teamId: string, userIds: string[]) => {
        if (vote?.vote?.id) {
            dispatch(addPlayersRequest({
                teamId,
                voteId: vote.vote.id,
                playerData: { user_ids: userIds },
            }));
        }
    };

    const handleRemovePlayer = (teamId: string, playerId: string) => {
        if (vote?.vote?.id) {
            dispatch(removePlayerRequest({
                teamId,
                voteId: vote.vote.id,
                playerId
            }));
        }
    };

    const handleUpdateCaptain = (teamId: string, newCaptainId: string) => {
        dispatch(updateTeamRequest({
            teamId,
            teamData: { captain_id: newCaptainId },
        }));
    };

    const handleDeleteTeam = (teamId: string) => {
        if (confirm('Are you sure you want to delete this team?')) {
            dispatch(deleteTeamRequest({ teamId }));
        }
    };

    const handleRefresh = () => {
        if (vote?.vote?.id) {
            dispatch(fetchTeamsRequest({ voteId: vote.vote.id }));
        }
    };

    // Determine which team letters are available
    const availableLetters: ('A' | 'B')[] = [];
    if (!teamA) availableLetters.push('A');
    if (!teamB) availableLetters.push('B');

    // Generate share URL
    const shareUrl = typeof window !== 'undefined' && vote?.vote?.id
        ? `${window.location.origin}/share/vote/${vote.vote.id}/teams`
        : '';

    const handleCopyLink = async () => {
        if (!shareUrl) return;

        try {
            // Try modern clipboard API first
            if (navigator.clipboard && navigator.clipboard.writeText) {
                await navigator.clipboard.writeText(shareUrl);
                alert('Link copied to clipboard!');
                return;
            }
        } catch (err) {
            // Fall through to fallback method
        }

        // Fallback: use execCommand for older browsers or when clipboard API fails
        try {
            const textArea = document.createElement('textarea');
            textArea.value = shareUrl;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            const successful = document.execCommand('copy');
            document.body.removeChild(textArea);

            if (successful) {
                alert('Link copied to clipboard!');
            } else {
                throw new Error('Copy command failed');
            }
        } catch (err) {
            // Last resort: show the URL in an alert for manual copying
            prompt('Copy this link:', shareUrl);
        }
    };

    const handleShare = async () => {
        if (!shareUrl) return;
        await handleCopyLink();
    };

    return (
        <div className="space-y-4">
            {/* Header */}
            <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold flex items-center gap-2">
                    <Shield className="h-5 w-5" />
                    Team Division
                </h3>
                <div className="flex items-center gap-2">
                    {/* Share button - only show when both teams exist */}
                    {teamA && teamB && shareUrl && (
                        <Button
                            variant="outline"
                            size="icon"
                            onClick={handleShare}
                            className="h-9 w-9"
                            title="Share or copy link"
                        >
                            <Share2 className="h-4 w-4" />
                        </Button>
                    )}
                    <Button
                        variant="outline"
                        size="icon"
                        onClick={handleRefresh}
                        disabled={loading}
                        className="h-9 w-9"
                        title="Refresh teams"
                    >
                        <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                    </Button>
                </div>
            </div>

            {/* Error Message */}
            {error && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                    <p className="text-sm text-red-600">{error}</p>
                </div>
            )}

            {/* Teams Display */}
            <div className="space-y-3">
                {/* TEAM A: Pass currentUserId for creator authorization */}
                {teamA && (
                    <TeamCard
                        team={teamA}
                        voters={voters}
                        allTeams={teams}
                        currentUserId={currentUserId} // CRITICAL: Required for authorization
                        isAuthenticated={isAuthenticated}
                        onAddPlayers={handleAddPlayers}
                        onRemovePlayer={handleRemovePlayer}
                        onUpdateCaptain={handleUpdateCaptain}
                        onDeleteTeam={handleDeleteTeam}
                    />
                )}

                {/* TEAM B: Pass currentUserId for creator authorization */}
                {teamB && (
                    <TeamCard
                        team={teamB}
                        voters={voters}
                        allTeams={teams}
                        currentUserId={currentUserId} // CRITICAL: Required for authorization
                        isAuthenticated={isAuthenticated}
                        onAddPlayers={handleAddPlayers}
                        onRemovePlayer={handleRemovePlayer}
                        onUpdateCaptain={handleUpdateCaptain}
                        onDeleteTeam={handleDeleteTeam}
                    />
                )}
            </div>

            {/* Create Team Button/Form */}
            {isAuthenticated && availableLetters.length > 0 && (
                <div className="space-y-3">
                    {!showCreateForm ? (
                        <Button
                            onClick={() => setShowCreateForm(true)}
                            className="w-full"
                            variant="outline"
                        >
                            <Shield className="h-4 w-4 mr-2" />
                            Create Team {availableLetters[0]}
                        </Button>
                    ) : (
                        <div className="bg-white rounded-lg border p-4 space-y-3">
                            <div className="space-y-2">
                                <Label htmlFor="teamName">Team Name</Label>
                                <Input
                                    id="teamName"
                                    value={teamName}
                                    onChange={(e) => {
                                        setTeamName(e.target.value);
                                        if (e.target.value.length > 0 && e.target.value.length < 2) {
                                            setErrorMessage('Minimum 2 characters required');
                                        } else {
                                            setErrorMessage('');
                                            }
                                    }}
                                    placeholder="e.g., Blue Warriors"
                                    maxLength={100}
                                />
                                    {errorMessage && 
                                        <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                                            <p className="text-sm text-red-600">{errorMessage}</p>
                                        </div>}
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="teamLetter">Team</Label>
                                <Select
                                    value={teamLetter}
                                    onValueChange={(value) => setTeamLetter(value as 'A' | 'B')}
                                >
                                    <SelectTrigger>
                                        <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {availableLetters.map((letter) => (
                                            <SelectItem key={letter} value={letter}>
                                                Team {letter}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            <div className="space-y-2">
                                <Label htmlFor="captain">Captain (must be a voter)</Label>
                                <Select
                                    value={captainId}
                                    onValueChange={setCaptainId}
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select captain..." />
                                    </SelectTrigger>
                                    <SelectContent>
                                        {voters.map((voter) => (
                                            <SelectItem key={voter.id} value={voter.id}>
                                                {voter.name}
                                            </SelectItem>
                                        ))}
                                    </SelectContent>
                                </Select>
                            </div>

                            <div className="flex gap-2">
                                <Button
                                    onClick={handleCreateTeam}
                                    disabled={loading || !teamName.trim() || !captainId}
                                    className="flex-1"
                                >
                                    {loading ? 'Creating...' : 'Create Team'}
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={() => {
                                        setShowCreateForm(false);
                                        setTeamName('');
                                        setCaptainId('');
                                    }}
                                    className="flex-1"
                                >
                                    Cancel
                                </Button>
                            </div>
                        </div>
                    )}
                </div>
            )}

            {/* No Teams Message */}
            {!loading && teams.length === 0 && (
                <div className="text-center py-6 text-gray-500">
                    <Shield className="h-12 w-12 mx-auto mb-2 text-gray-300" />
                    <p className="text-sm">No teams created yet</p>
                    {isAuthenticated && (
                        <p className="text-xs mt-1">
                            Create teams to divide voters into Team A and Team B
                        </p>
                    )}
                </div>
            )}
        </div>
    );
}

