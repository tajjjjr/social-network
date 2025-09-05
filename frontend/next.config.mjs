/** @type {import('next').NextConfig} */
const nextConfig = {
  outputFileTracingRoot: '/home/jamos/Documents/social-network/frontend',
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9000',
        pathname: '/avatar/**',
      },
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '9000',
        pathname: '/group-avatar/**',
      },
    ],
  },
};

export default nextConfig;
