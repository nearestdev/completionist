"use client";
import { useEffect, useState } from "react";
import lastfmService from "@/services/lastfmService";
import { LastFMAccount, Track } from "@/types/lastfm";
import Image from "next/image";
import Link from "next/link";

export default function MyLastFMPage() {
  const [account, setAccount] = useState<LastFMAccount | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLastFMData = async () => {
      try {
        setLoading(true);
        const accountData = await lastfmService.getMyLastFMAccount();
        setAccount(accountData);
        const recentTracksData = await lastfmService.getMyRecentTracks(30);
        setTracks(recentTracksData.recenttracks.track);
      } catch (err) {
        setError("Failed to load your Last.fm information.");
      } finally {
        setLoading(false);
      }
    };
    fetchLastFMData();
  }, []);

  if (loading) return <div>Loading your Last.fm information...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!account) return <div>Could not load Last.fm account.</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-8">
        <h1 className="text-3xl font-bold">{account.username}</h1>
        <a
          href={`https://www.last.fm/user/${account.username}`}
          target="_blank"
          rel="noopener noreferrer"
          className="text-indigo-500 hover:underline"
        >
          View on Last.fm
        </a>
      </div>

      <h2 className="text-2xl font-bold mb-6">Recent Tracks</h2>
      {tracks.length > 0 ? (
        <div className="space-y-4">
          {tracks.map((track, index) => {
            const artUrl =
              track.image.find((i) => i.size === "large")?.["#text"] ||
              "/placeholder.svg";
            const isNowPlaying = track["@attr"]?.nowplaying === "true";
            return (
              <a
                key={index}
                href={track.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-4 bg-white dark:bg-gray-800 p-3 rounded-lg shadow-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
              >
                <div className="flex-shrink-0 w-16 h-16 relative">
                  <Image
                    src={artUrl}
                    alt={track.album["#text"]}
                    fill
                    sizes="64px"
                    className="rounded-md object-cover"
                  />
                </div>
                <div className="flex-grow">
                  <p className="font-bold">{track.name}</p>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {track.artist["#text"]}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-500">
                    {track.album["#text"]}
                  </p>
                </div>
                {isNowPlaying && (
                  <span className="text-sm text-green-500 font-semibold animate-pulse">
                    Now Playing
                  </span>
                )}
              </a>
            );
          })}
        </div>
      ) : (
        <p>No recent tracks found.</p>
      )}
    </div>
  );
}