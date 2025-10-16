import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'lh3.googleusercontent.com',
        port: '',
        pathname: '/**',
      },
    ],
  },
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'Content-Security-Policy',
            value:
              "default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob:; " +
              "script-src 'self' 'unsafe-inline' 'unsafe-eval' https://pagead2.googlesyndication.com https://*.google.com https://*.adtrafficquality.google https://www.googletagmanager.com https://static.cloudflareinsights.com https://challenges.cloudflare.com; " +
              "connect-src 'self' http://localhost:* https://localhost:* http://127.0.0.1:* https://127.0.0.1:* https://cricket.dojima.foundation https://cricket-dev.dojima.foundation https://spark-park.dojima.foundation https://spark-park-dev.dojima.foundation https://ochhmsslirapqqzcgvek.supabase.co https://api.whatsapp.com https://api.iconify.design https://pagead2.googlesyndication.com https://*.google.com https://*.google-analytics.com https://*.adtrafficquality.google https://cloudflareinsights.com; " +
              "img-src 'self' data: https: blob:; " +
              "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
              "font-src 'self' data: https://fonts.gstatic.com; " +
              "frame-src 'self' https://googleads.g.doubleclick.net https://tpc.googlesyndication.com https://*.adtrafficquality.google https://www.google.com;",
          },
        ],
      },
    ];
  },
};

export default nextConfig;
