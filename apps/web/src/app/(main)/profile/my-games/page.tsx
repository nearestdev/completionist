"use client";

import { useEffect, useState } from "react";
import steamService from "@/services/steamService";
import { SteamOwnedGames } from "@/types/steam";
import Image from "next/image";

export default function MyGamesPage() {
  const [games, setGames] = useState<SteamOwnedGames | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchOwnedGames = async () => {
      try {
        setLoading(true);
        const ownedGames = await steamService.getMySteamOwnedGames();
        setGames(ownedGames);
      } catch (err) {
        setError("Failed to load your Steam games.");
      } finally {
        setLoading(false);
      }
    };
    fetchOwnedGames();
  }, []);

  if (loading) return <div>Loading your games...</div>;
  if (error) return <div className="text-red-500">{error}</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">My Steam Games</h1>
      {games && games.response.games.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {games.response.games.map((game) => (
            <div
              key={game.appid}
              className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden"
            >
              <div className="relative h-48">
                <Image
                  src={`https://steamcdn-a.akamaihd.net/steam/apps/${game.appid}/header.jpg`}
                  alt={game.name}
                  fill
                  sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
                  className="object-cover"
                />
              </div>
              <div className="p-4">
                <h3 className="font-bold text-lg truncate">{game.name}</h3>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  Playtime: {Math.round(game.playtime_forever / 60)} hours
                </p>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p>You have no games on your Steam account or your profile is private.</p>
      )}
    </div>
  );
}