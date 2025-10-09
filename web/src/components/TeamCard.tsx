'use client';

import React, { useState } from 'react';
import { VoteTeamWithPlayers, User } from '@/types/team';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Crown, UserPlus, UserMinus, Trash2 } from 'lucide-react';

interface TeamCardProps {
    team: VoteTeamWithPlayers;
    voters: User[];
    allTeams: VoteTeamWithPlayers[]; // All teams to check if a player is already assigned
    currentUserId: string | undefined;
    isAuthenticated: boolean;
    onAddPlayers: (teamId: string, userIds: string[]) => void;
    onRemovePlayer: (teamId: string, playerId: string) => void;
    onUpdateCaptain: (teamId: string, captainId: string) => void;
    onDeleteTeam: (teamId: string) => void;
}

export default function TeamCard({
    team,
    voters,
    allTeams,
    currentUserId,
    isAuthenticated,
    onAddPlayers,
    onRemovePlayer,
    onUpdateCaptain,
    onDeleteTeam,
}: TeamCardProps) {
    const [showPlayerSelection, setShowPlayerSelection] = useState(false);
    const [selectedPlayerIds, setSelectedPlayerIds] = useState<string[]>([]);

    // Get all player IDs from all teams
    const allPlayerIds = new Set<string>();
    allTeams.forEach(t => {
        t.players?.forEach(p => {
            allPlayerIds.add(p.id);
        });
    });

    // Get voters who are not in any team yet
    const availableVoters = voters.filter(
        voter => !allPlayerIds.has(voter.id)
    );

    const handleTogglePlayer = (userId: string) => {
        setSelectedPlayerIds(prev =>
            prev.includes(userId)
                ? prev.filter(id => id !== userId)
                : [...prev, userId]
        );
    };

    const handleAddPlayers = () => {
        if (selectedPlayerIds.length > 0) {
            onAddPlayers(team.id, selectedPlayerIds);
            setSelectedPlayerIds([]);
            setShowPlayerSelection(false);
        }
    };

    const teamColor = team.team_letter === 'A' ? 'blue' : 'red';
    const bgColor = team.team_letter === 'A' ? 'bg-blue-50' : 'bg-red-50';
    const borderColor = team.team_letter === 'A' ? 'border-blue-200' : 'border-red-200';
    const textColor = team.team_letter === 'A' ? 'text-blue-700' : 'text-red-700';

    return (
        <Card className={`${bgColor} ${borderColor} border-2`}>
            <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                    <CardTitle className={`text-lg font-bold ${textColor}`}>
                        Team {team.team_letter}: {team.team_name}
                    </CardTitle>
                    {isAuthenticated && (
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => onDeleteTeam(team.id)}
                            className="h-8 w-8 p-0"
                        >
                            <Trash2 className="h-4 w-4 text-gray-500" />
                        </Button>
                    )}
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Crown className="h-4 w-4" />
                    <span>Captain: {team.captain_name || 'Unknown'}</span>
                </div>
            </CardHeader>

            <CardContent className="space-y-3">
                {/* Player Count */}
                <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Users className="h-4 w-4" />
                    <span>{team.player_count} / 20 players</span>
                </div>

                {/* Player List */}
                {team.players && team.players.length > 0 && (
                    <div className="space-y-2">
                        <p className="text-xs font-semibold text-gray-700">Players:</p>
                        <div className="space-y-1">
                            {team.players.map((player) => (
                                <div
                                    key={player.id}
                                    className="flex items-center justify-between bg-white rounded px-2 py-1.5 text-sm"
                                >
                                    <div className="flex items-center gap-2">
                                        {player.id === team.captain_id && (
                                            <Crown className="h-3 w-3 text-yellow-500" />
                                        )}
                                        <span className="text-gray-800">{player.name}</span>
                                    </div>
                                    {isAuthenticated && player.id !== team.captain_id && (
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => onRemovePlayer(team.id, player.id)}
                                            className="h-6 w-6 p-0"
                                        >
                                            <UserMinus className="h-3 w-3 text-red-500" />
                                        </Button>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* Add Players Button */}
                {isAuthenticated && team.player_count < 20 && (
                    <div className="space-y-2">
                        {!showPlayerSelection ? (
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={() => setShowPlayerSelection(true)}
                                className="w-full"
                                disabled={availableVoters.length === 0}
                            >
                                <UserPlus className="h-4 w-4 mr-2" />
                                Add Players
                            </Button>
                        ) : (
                            <div className="space-y-2">
                                <p className="text-xs font-semibold text-gray-700">
                                    Select players to add:
                                </p>
                                <div className="max-h-32 overflow-y-auto space-y-1">
                                    {availableVoters.map((voter) => (
                                        <label
                                            key={voter.id}
                                            className="flex items-center gap-2 bg-white rounded px-2 py-1.5 cursor-pointer hover:bg-gray-50"
                                        >
                                            <input
                                                type="checkbox"
                                                checked={selectedPlayerIds.includes(voter.id)}
                                                onChange={() => handleTogglePlayer(voter.id)}
                                                className="h-4 w-4"
                                            />
                                            <span className="text-sm">{voter.name}</span>
                                        </label>
                                    ))}
                                </div>
                                <div className="flex gap-2">
                                    <Button
                                        size="sm"
                                        onClick={handleAddPlayers}
                                        disabled={selectedPlayerIds.length === 0}
                                        className="flex-1"
                                    >
                                        Add ({selectedPlayerIds.length})
                                    </Button>
                                    <Button
                                        size="sm"
                                        variant="outline"
                                        onClick={() => {
                                            setShowPlayerSelection(false);
                                            setSelectedPlayerIds([]);
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

                {/* No players available message */}
                {isAuthenticated && availableVoters.length === 0 && team.player_count < 20 && (
                    <p className="text-xs text-gray-500 italic">
                        All voters have been assigned to teams
                    </p>
                )}

                {/* Team full message */}
                {team.player_count >= 20 && (
                    <p className="text-xs text-gray-500 italic">
                        Team is full (20/20 players)
                    </p>
                )}
            </CardContent>
        </Card>
    );
}

