import Link from 'next/link';
import { Button } from '@/components/ui/button';

export default function NotFound() {
    return (
        <div className="min-h-screen bg-gray-50 flex items-center justify-center px-4">
            <div className="text-center max-w-md">
                <h1 className="text-4xl font-bold text-gray-900 mb-4">Teams Not Found</h1>
                <p className="text-gray-600 mb-6">
                    The teams for this vote could not be found. They may not have been created yet.
                </p>
                <Link href="/votes">
                    <Button>Go to Votes</Button>
                </Link>
            </div>
        </div>
    );
}

