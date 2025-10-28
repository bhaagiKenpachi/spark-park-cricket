/**
 * Server-side API utility for fetching teams data
 * Used for server-side rendering and metadata generation
 */

import { VoteTeamWithPlayers } from '@/types/team';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://spark-park.dojima.foundation/api/v1';

export async function fetchTeamsByVoteId(voteId: string): Promise<VoteTeamWithPlayers[]> {
    try {
        const url = `${API_BASE_URL}/votes/${voteId}/teams`;
        const response = await fetch(url, {
            headers: {
                'Content-Type': 'application/json',
                Accept: 'application/json',
            },
            cache: 'no-store', // Always fetch fresh data for metadata
        });

        if (!response.ok) {
            throw new Error(`Failed to fetch teams: ${response.status}`);
        }

        const data = await response.json();
        // Handle nested response structure
        const teams = (data.data?.data || data.data || data) as VoteTeamWithPlayers[];
        return Array.isArray(teams) ? teams : [];
    } catch (error) {
        console.error('Error fetching teams:', error);
        return [];
    }
}

