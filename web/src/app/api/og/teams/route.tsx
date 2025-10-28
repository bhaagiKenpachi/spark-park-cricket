import { ImageResponse } from 'next/og';

export const runtime = 'edge';

export async function GET(request: Request) {
    try {
        const { searchParams } = new URL(request.url);
        const voteId = searchParams.get('voteId');

        if (!voteId) {
            return new Response('Missing voteId parameter', { status: 400 });
        }

        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'https://spark-park.dojima.foundation/api/v1';
        const response = await fetch(`${apiUrl}/votes/${voteId}/teams`);
        const data = await response.json();
        const teams = data.data || [];

        const teamA = teams.find((t: any) => t.team_letter === 'A');
        const teamB = teams.find((t: any) => t.team_letter === 'B');

        const teamAName: string = teamA?.team_name || 'Team A';
        const teamBName: string = teamB?.team_name || 'Team B';
        const teamAPlayers: string[] = (teamA?.player_names || []) as string[];
        const teamBPlayers: string[] = (teamB?.player_names || []) as string[];

        // Build two columns of text, padded to align visually using monospace font
        const leftLines = [`Team A: ${teamAName}`, ...teamAPlayers];
        const rightLines = [`Team B: ${teamBName}`, ...teamBPlayers];

        const maxRows = Math.max(leftLines.length, rightLines.length);
        while (leftLines.length < maxRows) leftLines.push('');
        while (rightLines.length < maxRows) rightLines.push('');

        const leftWidth = Math.min(
            48, // cap to prevent overflow
            Math.max(...leftLines.map((s) => (s ? s.length : 0))) + 2,
        );

        const rows: string[] = [];
        for (let i = 0; i < maxRows; i++) {
            const left = (leftLines[i] || '').padEnd(leftWidth, ' ');
            const right = rightLines[i] || '';
            rows.push(`${left} | ${right}`);
        }

        const fullText = rows.join('\n');

        const width = 1200;
        const height = 630;

        return new ImageResponse(
            (
                <div
                    style={{
                        display: 'flex',
                        width: '100%',
                        height: '100%',
                        backgroundColor: '#000000',
                        color: '#ffffff',
                        alignItems: 'center',
                        justifyContent: 'center',
                        padding: '40px',
                        boxSizing: 'border-box',
                        fontFamily: 'Menlo, Consolas, Monaco, Liberation Mono, Courier New, monospace',
                        whiteSpace: 'pre',
                        textAlign: 'left',
                        lineHeight: 1.5,
                        fontSize: 32,
                    }}
                >
                    {fullText}
                </div>
            ),
            { width, height },
        );
    } catch (error) {
        console.error('Error generating OG image:', error);
        return new Response('Failed to generate image', { status: 500 });
    }
}