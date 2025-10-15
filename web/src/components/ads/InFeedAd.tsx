'use client';

import { useEffect, useRef } from 'react';

declare global {
    interface Window {
        adsbygoogle: unknown[];
    }
}

interface InFeedAdProps {
    adSlot: string;
    adLayout?: 'in-article';
    className?: string;
}

export function InFeedAd({
    adSlot,
    adLayout = 'in-article',
    className = '',
}: InFeedAdProps): React.JSX.Element {
    const clientId = process.env.NEXT_PUBLIC_ADSENSE_CLIENT_ID;
    const isAdPushed = useRef(false);

    useEffect(() => {
        // Prevent double-pushing in React strict mode
        if (isAdPushed.current) return;

        try {
            if (typeof window !== 'undefined') {
                (window.adsbygoogle = window.adsbygoogle || []).push({});
                isAdPushed.current = true;
            }
        } catch (error) {
            console.error('AdSense error:', error);
        }
    }, []);

    if (!clientId) {
        return <></>;
    }

    return (
        <div className={`ad-container my-4 ${className}`}>
            <ins
                className="adsbygoogle"
                style={{ display: 'block', textAlign: 'center' }}
                data-ad-client={clientId}
                data-ad-slot={adSlot}
                data-ad-layout={adLayout}
                data-ad-format="fluid"
            />
        </div>
    );
}

