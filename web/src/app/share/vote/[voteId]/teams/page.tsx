import { Metadata } from 'next';
import { TeamsPreview } from '@/components/TeamsPreview';
import { fetchTeamsByVoteId } from '@/lib/teams-api';
import { notFound } from 'next/navigation';

interface PageProps {
    params: Promise<{
        voteId: string;
    }>;
}

async function getSiteUrl(): Promise<string> {
    // For local testing with ngrok/localtunnel, set NEXT_PUBLIC_NGROK_URL
    return process.env.NEXT_PUBLIC_NGROK_URL ||
        process.env.NEXT_PUBLIC_SITE_URL ||
        (process.env.VERCEL_URL ? `https://${process.env.VERCEL_URL}` : '') ||
        'https://spark-park-cricket.vercel.app';
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
    const { voteId } = await params;
    const teams = await fetchTeamsByVoteId(voteId);

    const teamA = teams.find(t => t.team_letter === 'A');
    const teamB = teams.find(t => t.team_letter === 'B');

    const siteUrl = await getSiteUrl();
    const url = `${siteUrl}/share/vote/${voteId}/teams`;

    // Generate title and description
    let title = 'Team Matchup';
    let description = 'Check out the teams!';

    if (teamA && teamB) {
        title = `Team ${teamA.team_name} vs Team ${teamB.team_name}`;
        const teamAPlayers = teamA.player_count || 0;
        const teamBPlayers = teamB.player_count || 0;
        description = `Team A: ${teamA.team_name} (${teamAPlayers} players) vs Team B: ${teamB.team_name} (${teamBPlayers} players)`;
    } else if (teamA) {
        title = `Team ${teamA.team_name} (Team A)`;
        description = `Team A: ${teamA.team_name} with ${teamA.player_count || 0} players`;
    } else if (teamB) {
        title = `Team ${teamB.team_name} (Team B)`;
        description = `Team B: ${teamB.team_name} with ${teamB.player_count || 0} players`;
    }

    // Generate preview text for WhatsApp
    // WhatsApp preview shows ~2-3 lines, so keep it concise (no newlines, single line)
    const previewText = teamA && teamB
        ? `🏏 ${teamA.team_name} vs ${teamB.team_name} | Team A: ${teamA.player_count || 0}/20 players | Team B: ${teamB.player_count || 0}/20 players`
        : description;

    // Use our working custom OG image that shows both teams side by side
    const ogImageUrl = `${siteUrl}/api/og/teams?voteId=${voteId}`;

    return {
        title,
        description,
        openGraph: {
            title,
            description: previewText,
            url,
            siteName: 'Spark Park Cricket',
            locale: 'en_US',
            type: 'website',
            images: [
                {
                    url: ogImageUrl,
                    width: 1200,
                    height: 630,
                    alt: title,
                },
            ],
        },
        twitter: {
            card: 'summary_large_image',
            title,
            description: previewText,
            images: [ogImageUrl],
        },
        alternates: {
            canonical: url,
        },
        other: {
            ...(process.env.NEXT_PUBLIC_FACEBOOK_APP_ID && { 'fb:app_id': process.env.NEXT_PUBLIC_FACEBOOK_APP_ID }),
            'article:author': 'Spark Park Cricket',
            'article:section': 'Cricket Tournament',
            'og:updated_time': new Date().toISOString(),
        },
    };
}

export default async function ShareTeamsPage({ params }: PageProps) {
    const { voteId } = await params;
    const teams = await fetchTeamsByVoteId(voteId);

    // Check if we have at least one team
    if (!teams || teams.length === 0) {
        notFound();
    }

    return <TeamsPreview voteId={voteId} teams={teams} />;
}

