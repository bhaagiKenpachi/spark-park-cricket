import type { NextConfig } from 'next';

// Get CSP domains from environment variables
const getCSPDomains = () => {
  const domains = {
    api: process.env.NEXT_PUBLIC_API_URL || '',
    supabase: process.env.NEXT_PUBLIC_SUPABASE_URL || '',
    allowedOrigins: process.env.NEXT_PUBLIC_ALLOWED_ORIGINS || '',
  };

  // Build connect-src directive with environment-specific domains
  const connectSrcDomains = [
    "'self'",
    'http://localhost:*',
    'https://localhost:*',
    'http://127.0.0.1:*',
    'https://127.0.0.1:*',
    'https://api.whatsapp.com',
    'https://api.iconify.design',
    'https://pagead2.googlesyndication.com',
    'https://*.google.com',
    'https://*.google-analytics.com',
    'https://*.adtrafficquality.google',
    'https://cloudflareinsights.com',
  ];

  // Add environment-specific domains if provided
  if (domains.api) {
    const apiDomain = new URL(domains.api).origin;
    connectSrcDomains.push(apiDomain);
  }
  if (domains.supabase) {
    connectSrcDomains.push(domains.supabase);
  }
  if (domains.allowedOrigins) {
    domains.allowedOrigins.split(',').forEach(origin => {
      const trimmed = origin.trim();
      if (trimmed) connectSrcDomains.push(trimmed);
    });
  }

  return connectSrcDomains.join(' ');
};

const nextConfig: NextConfig = {
  eslint: {
    // Warning: This allows production builds to successfully complete even if
    // your project has ESLint errors.
    ignoreDuringBuilds: true,
  },
  typescript: {
    // Warning: This allows production builds to successfully complete even if
    // your project has type errors.
    ignoreBuildErrors: false,
  },
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
              `connect-src ${getCSPDomains()}; ` +
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
