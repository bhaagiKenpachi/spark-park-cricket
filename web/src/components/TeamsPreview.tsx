'use client';

import { VoteTeamWithPlayers } from '@/types/team';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Crown, Share2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useEffect, useState } from 'react';

interface TeamsPreviewProps {
    voteId: string;
    teams: VoteTeamWithPlayers[];
}

export function TeamsPreview({ voteId, teams }: TeamsPreviewProps) {
    const [shareUrl, setShareUrl] = useState('');
    const teamA = teams.find(t => t.team_letter === 'A');
    const teamB = teams.find(t => t.team_letter === 'B');

    useEffect(() => {
        if (typeof window !== 'undefined') {
            setShareUrl(`${window.location.origin}/share/vote/${voteId}/teams`);
        }
    }, [voteId]);

    const handleCopyLink = async () => {
        if (!shareUrl) return;

        try {
            if (navigator.clipboard && navigator.clipboard.writeText) {
                await navigator.clipboard.writeText(shareUrl);
                alert('Link copied to clipboard!');
                return;
            }
        } catch (err) {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = shareUrl;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
            alert('Link copied to clipboard!');
        }
    };

    const handleShare = async () => {
        if (!shareUrl) return;
        await handleCopyLink();
    };

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Header */}
            <div className="bg-white border-b shadow-sm">
                <div className="w-full max-w-4xl mx-auto px-4 py-4">
                    <div className="flex items-center justify-between">
                        <h1 className="text-2xl font-bold text-gray-900">
                            {teamA && teamB
                                ? `${teamA.team_name} vs ${teamB.team_name}`
                                : 'Team Matchup'
                            }
                        </h1>
                        {shareUrl && (
                            <Button
                                onClick={handleShare}
                                variant="outline"
                                size="icon"
                                className="h-9 w-9"
                                title="Share or copy link"
                            >
                                <Share2 className="h-4 w-4" />
                            </Button>
                        )}
                    </div>
                </div>
            </div>

            {/* Main Content */}
            <main className="w-full max-w-4xl mx-auto px-4 py-6">
                {teamA && teamB ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                        <TeamCard
                            team={teamA}
                            teamLabel="Team A"
                            bgColor="bg-blue-50"
                            borderColor="border-blue-200"
                            textColor="text-blue-700"
                        />
                        <TeamCard
                            team={teamB}
                            teamLabel="Team B"
                            bgColor="bg-red-50"
                            borderColor="border-red-200"
                            textColor="text-red-700"
                        />
                    </div>
                ) : (
                    <div className="max-w-2xl mx-auto">
                        {teamA && (
                            <TeamCard
                                team={teamA}
                                teamLabel="Team A"
                                bgColor="bg-blue-50"
                                borderColor="border-blue-200"
                                textColor="text-blue-700"
                            />
                        )}
                        {teamB && (
                            <TeamCard
                                team={teamB}
                                teamLabel="Team B"
                                bgColor="bg-red-50"
                                borderColor="border-red-200"
                                textColor="text-red-700"
                            />
                        )}
                    </div>
                )}
            </main>
        </div>
    );
}

interface TeamCardProps {
    team: VoteTeamWithPlayers;
    teamLabel: string;
    bgColor: string;
    borderColor: string;
    textColor: string;
}

function TeamCard({ team, teamLabel, bgColor, borderColor, textColor }: TeamCardProps) {
    return (
        <Card className={`${bgColor} ${borderColor} border-2`}>
            <CardHeader className="pb-3">
                <CardTitle className={`text-xl font-bold ${textColor}`}>
                    {teamLabel}: {team.team_name}
                </CardTitle>
                <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Crown className="h-4 w-4" />
                    <span>Captain: {team.captain_name || 'Unknown'}</span>
                </div>
            </CardHeader>

            <CardContent className="space-y-3">
                <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Users className="h-4 w-4" />
                    <span>{team.player_count} / 20 players</span>
                </div>

                {team.players && team.players.length > 0 && (
                    <div className="space-y-2">
                        <p className="text-xs font-semibold text-gray-700">Players:</p>
                        <div className="space-y-1 max-h-64 overflow-y-auto">
                            {team.players.map((player) => (
                                <div
                                    key={player.id}
                                    className="flex items-center gap-2 bg-white rounded px-2 py-1.5 text-sm"
                                >
                                    {player.id === team.captain_id && (
                                        <Crown className="h-3 w-3 text-yellow-500 flex-shrink-0" />
                                    )}
                                    <span className="text-gray-800">{player.name}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}