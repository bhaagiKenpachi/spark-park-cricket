'use client';

/**
 * TeamCard Component - Team Management with Creator Authorization
 * 
 * AUTHORIZATION SYSTEM:
 * ====================
 * This component implements creator-based authorization for team management:
 * 
 * 1. TEAM CREATOR PERMISSIONS:
 *    - Can delete the entire team
 *    - Can add players to the team
 *    - Can remove players from the team (except captain)
 *    - Sees "✓ You created this team" indicator
 * 
 * 2. NON-CREATOR PERMISSIONS:
 *    - Can only view team information
 *    - Cannot see management buttons (delete, add, remove)
 *    - No visual indicators about management capabilities
 * 
 * 3. AUTHORIZATION CHECK:
 *    - Uses: team.created_by === currentUserId
 *    - All management operations are gated behind isTeamCreator check
 *    - Prevents unauthorized team modifications
 * 
 * 4. CAPTAIN PROTECTION:
 *    - Captain cannot be removed (player.id !== team.captain_id)
 *    - This is separate from creator authorization
 * 
 * IMPORTANT: Do not remove or modify the isTeamCreator checks without
 * understanding the security implications. All team management operations
 * must be restricted to the team creator only.
 */

import React, { useState } from 'react';
import { VoteTeamWithPlayers, User } from '@/types/team';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Users, Crown, UserPlus, UserMinus, Trash2 } from 'lucide-react';
import Checkbox from '@mui/material/Checkbox';
import TextField from '@mui/material/TextField';
import Autocomplete from '@mui/material/Autocomplete';
import CheckBoxOutlineBlankIcon from '@mui/icons-material/CheckBoxOutlineBlank';
import CheckBoxIcon from '@mui/icons-material/CheckBox';

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
    onUpdateCaptain: _onUpdateCaptain,
    onDeleteTeam,
}: TeamCardProps) {
    // AUTHORIZATION: Check if current user is the team creator
    // Only team creators can delete teams, add/remove players
    // This prevents unauthorized team modifications
    const isTeamCreator = currentUserId && team.created_by === currentUserId;
    const [showPlayerSelection, setShowPlayerSelection] = useState(false);
    const [selectedPlayers, setSelectedPlayers] = useState<User[]>([]);
    const icon = <CheckBoxOutlineBlankIcon fontSize="small" />;
    const checkedIcon = <CheckBoxIcon fontSize="small" />;

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

    
    const selectedPlayerIds = selectedPlayers.map(p => p.id);
    const handleAddPlayers = () => {
        if (selectedPlayerIds.length > 0) {
            onAddPlayers(team.id, selectedPlayerIds);
            setSelectedPlayers([]); // Reset the state
            setShowPlayerSelection(false);
        }
    };

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
                    {/* DELETE TEAM BUTTON: Only visible to team creator */}
                    {/* CRITICAL: This prevents unauthorized team deletion */}
                    {isAuthenticated && isTeamCreator && (
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => onDeleteTeam(team.id)}
                            className="h-8 w-8 p-0"
                            title="Delete team (only team creator can delete)"
                        >
                            <Trash2 className="h-4 w-4 text-gray-500" />
                        </Button>
                    )}
                </div>
                <div className="flex items-center gap-2 text-sm text-gray-600">
                    <Crown className="h-4 w-4" />
                    <span>Captain: {team.captain_name || 'Unknown'}</span>
                </div>
                {/* TEAM CREATOR INDICATOR: Shows who has management permissions */}
                {/* This helps users understand why they can/cannot manage the team */}
                {isTeamCreator && (
                    <div className="flex items-center gap-2 text-xs text-blue-600 font-medium">
                        <span>✓ You created this team</span>
                    </div>
                )}
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
                                    {/* REMOVE PLAYER BUTTON: Only visible to team creator */}
                                    {/* CRITICAL: Prevents unauthorized player removal */}
                                    {/* Captain cannot be removed (player.id !== team.captain_id) */}
                                    {isAuthenticated && isTeamCreator && player.id !== team.captain_id && (
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => onRemovePlayer(team.id, player.id)}
                                            className="h-6 w-6 p-0"
                                            title="Remove player (only team creator can remove players)"
                                        >
                                            <UserMinus className="h-3 w-3 text-red-500" />
                                        </Button>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {/* ADD PLAYERS BUTTON: Only visible to team creator */}
                {/* CRITICAL: Prevents unauthorized player addition */}
                {/* Also checks team capacity (player_count < 20) */}
                {isAuthenticated && isTeamCreator && team.player_count < 20 && (
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
                                    <Autocomplete
                                        multiple
                                        id="checkboxes-tags-demo"
                                        options={availableVoters}
                                        disableCloseOnSelect
                                        getOptionLabel={(option) => option.name}
                                        // FIX: Control the component with the `selectedPlayers` state
                                        value={selectedPlayers}
                                        // FIX: Use the top-level onChange to update state
                                        onChange={(event, newValue) => {
                                            setSelectedPlayers(newValue);
                                        }}
                                        renderOption={(props, option, { selected }) => (
                                            <li {...props} key={option.id}>
                                                <Checkbox
                                                    icon={icon}
                                                    checkedIcon={checkedIcon}
                                                    style={{ marginRight: 8 }}
                                                    // FIX: Use the 'selected' prop provided by renderOption
                                                    checked={selected}
                                                />
                                                {option.name}
                                            </li>
                                        )}
                                        // The style prop was causing width issues; it's often better to control width via parent containers.
                                        // style={{ width: 500 }}
                                        renderInput={(params) => (
                                            <TextField
                                                {...params}
                                                // By spreading params first, we can override specific props
                                                label="Select Players"
                                                placeholder="Voters"
                                                size="small"
                                                variant="outlined"
                                                InputLabelProps={{
                                                    ...params.InputLabelProps,
                                                    className: params.InputLabelProps?.className ?? '',
                                                    style: params.InputLabelProps?.style ?? {},
                                                }}
                                            />
                                        )}
                                        />
 
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
                                            setSelectedPlayers([]);
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

                {/* NO PLAYERS AVAILABLE MESSAGE: Only shown to team creator */}
                {/* This message only makes sense for users who can add players */}
                {isAuthenticated && isTeamCreator && availableVoters.length === 0 && team.player_count < 20 && (
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

