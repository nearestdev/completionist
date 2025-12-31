"use client";
import { useEffect, useState } from "react";
import steamService from "@/services/steamService";
import { SteamAccount, SteamOwnedGames } from "@/types/steam";
import Image from "next/image";
import Link from "next/link";
import Button from "@/components/ui/Button";

export default function MySteamPage() {
  const [account, setAccount] = useState<SteamAccount | null>(null);
  const [games, setGames] = useState<SteamOwnedGames | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSteamData = async () => {
      try {
        setLoading(true);
        const accountData = await steamService.getMySteamAccount();
        setAccount(accountData);
        const ownedGames = await steamService.getMySteamOwnedGames();
        setGames(ownedGames);
      } catch (err) {
        setError("Failed to load your Steam information. Your profile may be private.");
      } finally {
        setLoading(false);
      }
    };
    fetchSteamData();
  }, []);

  if (loading) return <div>Loading your Steam information...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!account) return <div>Could not load Steam account.</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-8 flex items-center gap-6">
        <Image
          src={account.avatar}
          alt={`${account.persona}'s avatar`}
          width={128}
          height={128}
          className="rounded-lg"
        />
        <div>
          <h1 className="text-3xl font-bold">{account.persona}</h1>
          <a
            href={`https://steamcommunity.com/profiles/${account.steamId}`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-indigo-500 hover:underline"
          >
            View on Steam
          </a>
        </div>
      </div>

      <h2 className="text-2xl font-bold mb-6">
        My Games ({games?.response.game_count || 0})
      </h2>
      {games && games.response.games.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          {games.response.games
            .sort((a, b) => b.playtime_forever - a.playtime_forever)
            .map((game) => (
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
        <p>
          You have no games on your Steam account or your game details are private.
        </p>
      )}
    </div>
  );
}