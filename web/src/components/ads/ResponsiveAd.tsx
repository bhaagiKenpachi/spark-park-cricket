'use client';

import { useEffect, useRef } from 'react';

declare global {
    interface Window {
        adsbygoogle: unknown[];
    }
}

interface ResponsiveAdProps {
    adSlot: string;
    adFormat?: 'auto' | 'fluid' | 'rectangle';
    className?: string;
}

export function ResponsiveAd({
    adSlot,
    adFormat = 'auto',
    className = '',
}: ResponsiveAdProps): React.JSX.Element {
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
        <div className={`ad-container ${className}`}>
            <ins
                className="adsbygoogle"
                style={{ display: 'block' }}
                data-ad-client={clientId}
                data-ad-slot={adSlot}
                data-ad-format={adFormat}
                data-full-width-responsive="true"
            />
        </div>
    );
}

