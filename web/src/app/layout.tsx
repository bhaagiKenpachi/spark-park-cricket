import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';
import { ReduxProvider } from '@/providers/ReduxProvider';
import { AuthProvider } from '@/components/auth/AuthProvider';
import { GraphQLProvider } from '@/providers/GraphQLProvider';
import { AdSenseScript } from '@/components/ads/AdSenseScript';

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
        <ReduxProvider>
          <GraphQLProvider>
            <AuthProvider>{children}</AuthProvider>
          </GraphQLProvider>
        </ReduxProvider>
      </body>
    </html>
  );
}
