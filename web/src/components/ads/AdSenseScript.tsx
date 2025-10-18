'use client';

import Script from 'next/script';

export function AdSenseScript(): React.JSX.Element {
    const clientId = process.env.NEXT_PUBLIC_ADSENSE_CLIENT_ID;
    const adsEnabled = process.env.NEXT_PUBLIC_ENABLE_ADS !== 'false';

    if (!adsEnabled) {
        console.info('Ads are disabled via NEXT_PUBLIC_ENABLE_ADS environment variable.');
        return <></>;
    }

    if (!clientId) {
        console.warn('AdSense client ID not configured');
        return <></>;
    }

    return (
        <Script
            async
            src={`https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=${clientId}`}
            crossOrigin="anonymous"
            strategy="afterInteractive"
        />
    );
}
