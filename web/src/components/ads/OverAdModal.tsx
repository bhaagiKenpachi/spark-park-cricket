'use client';

import { useState, useEffect } from 'react';
import { ResponsiveAd } from './ResponsiveAd';
import { X } from 'lucide-react';

interface OverAdModalProps {
    onClose: () => void;
    adSlot: string;
    overNumber: number;
}

export function OverAdModal({
    onClose,
    adSlot,
    overNumber,
}: OverAdModalProps): React.JSX.Element {
    const [countdown, setCountdown] = useState(5);

    useEffect(() => {
        // Countdown timer
        const timer = setInterval(() => {
            setCountdown(prev => {
                if (prev <= 1) {
                    clearInterval(timer);
                    // Use setTimeout to avoid calling onClose during render
                    setTimeout(() => onClose(), 0);
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(timer);
    }, [onClose]);

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
            <div className="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
                {/* Header */}
                <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-2">
                        <div className="text-sm font-semibold text-gray-700">
                            Over {overNumber} - Advertisement
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="text-gray-400 hover:text-gray-600 transition-colors"
                        aria-label="Close ad"
                    >
                        <X className="h-5 w-5" />
                    </button>
                </div>

                {/* Ad Content */}
                <div className="mb-4">
                    <ResponsiveAd adSlot={adSlot} adFormat="rectangle" />
                </div>

                {/* Countdown */}
                <div className="text-center">
                    <div className="text-xs text-gray-500">
                        Ad will close in{' '}
                        <span className="font-bold text-blue-600">{countdown}</span>{' '}
                        second{countdown !== 1 ? 's' : ''}
                    </div>
                    <div className="mt-2 w-full bg-gray-200 rounded-full h-1.5">
                        <div
                            className="bg-blue-600 h-1.5 rounded-full transition-all duration-1000"
                            style={{ width: `${((5 - countdown) / 5) * 100}%` }}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}

