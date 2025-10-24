import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import { ReduxProvider } from '@/providers/ReduxProvider';
import { AuthProvider } from '@/components/auth/AuthProvider';
import { GraphQLProvider } from '@/providers/GraphQLProvider';
import { PostHogProvider } from '@/providers/PostHogProvider';
import { AdSenseScript } from '@/components/ads/AdSenseScript';
import ErrorBoundary from '@/components/ErrorBoundary';

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin'],
});

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'Spark Park Cricket',
  description: 'Cricket Tournament Management System',
  other: {
    'google-adsense-account': 'ca-pub-5474524579770573',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AdSenseScript />
        <ErrorBoundary>
          <PostHogProvider>
            <ReduxProvider>
              <GraphQLProvider>
                <AuthProvider>{children}</AuthProvider>
              </GraphQLProvider>
            </ReduxProvider>
          </PostHogProvider>
        </ErrorBoundary>
      </body>
    </html>
  );
}
