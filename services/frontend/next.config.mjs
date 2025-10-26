/** @type {import('next').NextConfig} */
const nextConfig = {
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  output: 'standalone',
  turbopack: false,
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://perfumist:8000/:path*',
      },
    ]
  },
}

export default nextConfig
