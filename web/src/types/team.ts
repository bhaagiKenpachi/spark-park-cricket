// Team voting types
export interface VoteTeam {
    id: string;
    vote_id: string;
    team_name: string;
    team_letter: 'A' | 'B';
    captain_id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
}

export interface TeamPlayer {
    id: string;
    team_id: string;
    user_id: string;
    created_at: string;
}

export interface User {
    id: string;
    name: string;
    email: string;
    picture?: string;
}

export interface VoteTeamWithPlayers extends VoteTeam {
    captain?: User;
    players?: User[];
    player_count: number;
    captain_name?: string;
    player_names?: string[];
}

// Request types
export interface CreateVoteTeamRequest {
    vote_id?: string; // Set from URL param
    team_name: string;
    team_letter: 'A' | 'B';
    captain_id: string;
}

export interface UpdateVoteTeamRequest {
    team_name?: string;
    captain_id?: string;
}

export interface AddPlayerRequest {
    user_id: string;
}

export interface TeamAssignmentRequest {
    user_ids: string[];
}


