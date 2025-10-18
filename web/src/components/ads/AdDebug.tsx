'use client';

/**
 * Debug component to verify AdSense environment variables
 * Remove this component once you've verified ads are working/disabled correctly
 */
export function AdDebug(): React.JSX.Element {
    const clientId = process.env.NEXT_PUBLIC_ADSENSE_CLIENT_ID;
    const enableAdsValue = process.env.NEXT_PUBLIC_ENABLE_ADS;
    const adsEnabled = enableAdsValue !== 'false';

    if (typeof window !== 'undefined') {
        console.log('🔍 Ad Debug Info:', {
            NEXT_PUBLIC_ADSENSE_CLIENT_ID: clientId ? 'Set' : 'Not Set',
            NEXT_PUBLIC_ENABLE_ADS: enableAdsValue || 'undefined',
            adsEnabled: adsEnabled,
            expectedBehavior: adsEnabled ? '✅ Ads ENABLED' : '❌ Ads DISABLED',
        });
    }

    return (
        <div style={{
            position: 'fixed',
            bottom: 10,
            right: 10,
            background: 'rgba(0, 0, 0, 0.8)',
            color: 'white',
            padding: '10px',
            borderRadius: '5px',
            fontSize: '12px',
            zIndex: 9999,
            fontFamily: 'monospace',
        }}>
            <div style={{ fontWeight: 'bold', marginBottom: '5px' }}>🔍 Ad Debug</div>
            <div>Client ID: {clientId ? '✅ Set' : '❌ Not Set'}</div>
            <div>ENABLE_ADS: {enableAdsValue || 'undefined'}</div>
            <div style={{ 
                color: adsEnabled ? '#ff6b6b' : '#51cf66',
                fontWeight: 'bold' 
            }}>
                Status: {adsEnabled ? '✅ ENABLED' : '❌ DISABLED'}
            </div>
        </div>
    );
}

