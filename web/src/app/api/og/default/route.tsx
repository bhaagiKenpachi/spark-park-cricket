import { ImageResponse } from 'next/og';

export const runtime = 'edge';

export async function GET() {
    return new ImageResponse(
        (
            <div
                style={{
                    height: '100%',
                    width: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    backgroundColor: '#f3f4f6',
                    backgroundImage: 'linear-gradient(45deg, #3b82f6 0%, #ef4444 100%)',
                }}
            >
                <div
                    style={{
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                        backgroundColor: 'white',
                        padding: '40px',
                        borderRadius: '20px',
                        boxShadow: '0 20px 25px -5px rgba(0, 0, 0, 0.1)',
                        maxWidth: '800px',
                    }}
                >
                    <div
                        style={{
                            fontSize: '48px',
                            fontWeight: 'bold',
                            color: '#1f2937',
                            marginBottom: '20px',
                            textAlign: 'center',
                        }}
                    >
                        🏏 Cricket Match
                    </div>

                    <div
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '40px',
                            marginBottom: '20px',
                        }}
                    >
                        <div
                            style={{
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'center',
                                backgroundColor: '#dbeafe',
                                padding: '20px',
                                borderRadius: '12px',
                                minWidth: '200px',
                            }}
                        >
                            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#1e40af' }}>
                                Team A
                            </div>
                            <div style={{ fontSize: '20px', color: '#1e40af', marginTop: '8px' }}>
                                vs
                            </div>
                        </div>

                        <div
                            style={{
                                fontSize: '32px',
                                fontWeight: 'bold',
                                color: '#6b7280',
                            }}
                        >
                            VS
                        </div>

                        <div
                            style={{
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'center',
                                backgroundColor: '#fef2f2',
                                padding: '20px',
                                borderRadius: '12px',
                                minWidth: '200px',
                            }}
                        >
                            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#dc2626' }}>
                                Team B
                            </div>
                            <div style={{ fontSize: '20px', color: '#dc2626', marginTop: '8px' }}>
                                vs
                            </div>
                        </div>
                    </div>

                    <div
                        style={{
                            fontSize: '18px',
                            color: '#6b7280',
                            marginTop: '20px',
                        }}
                    >
                        Spark Park Cricket
                    </div>
                </div>
            </div>
        ),
        {
            width: 1200,
            height: 630,
        }
    );
}
