import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'Content-Security-Policy',
            value:
              "connect-src 'self' http://localhost:* https://localhost:* http://127.0.0.1:* https://127.0.0.1:* https://cricket.dojima.foundation https://cricket-dev.dojima.foundation https://ochhmsslirapqqzcgvek.supabase.co https://api.whatsapp.com wss://.supabase.co https://api.iconify.design; default-src 'self' 'unsafe-inline' 'unsafe-eval' data: blob:;",
          },
        ],
      },
    ];
  },
};

export default nextConfig;
