import type { NextConfig } from "next";
import dotenv from "dotenv";
import path from "path";

dotenv.config({ path: path.resolve(__dirname, "../.env") });

const backendBaseURL = process.env.BACKEND_BASE_URL || "http://localhost";
const backendPort = process.env.BACKEND_PORT || "8080";
const frontendPort = process.env.PORT || "3000";

const nextConfig: NextConfig = {
  env: {
    NEXT_PUBLIC_BACKEND_BASE_URL: backendBaseURL,
    NEXT_PUBLIC_BACKEND_PORT: backendPort,
    NEXT_PUBLIC_API_URL: `${backendBaseURL}:${backendPort}/api`,
  },
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "cdn.myanimelist.net",
        port: "",
        pathname: "/images/**",
      },
      {
        protocol: "https",
        hostname: "image.tmdb.org",
        port: "",
        pathname: "/t/p/**",
      },
      {
        protocol: "https",
        hostname: "media.rawg.io",
        port: "",
        pathname: "/media/**",
      },
      {
        protocol: "https",
        hostname: "avatars.steamstatic.com",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "steamcdn-a.akamaihd.net",
        port: "",
        pathname: "/**",
      },
      {
        protocol: "https",
        hostname: "books.google.com",
        port: "",
        pathname: "/books/content/**",
      },
      {
        protocol: "https",
        hostname: "lastfm.freetls.fastly.net",
        port: "",
        pathname: "/i/u/**",
      },
    ],
  },
};

export default nextConfig;