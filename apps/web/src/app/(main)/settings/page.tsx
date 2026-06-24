"use client";

import { useEffect, useState } from "react";
import steamService from "@/services/steamService";
import lastfmService from "@/services/lastfmService";
import { SteamAccount } from "@/types/steam";
import { LastFMAccount } from "@/types/lastfm";
import Button from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";
import Image from "next/image";
import { 
  GearIcon, 
  GameControllerIcon, 
  MusicNotesIcon, 
  CheckCircleIcon, 
  LinkBreakIcon,
  PlugsIcon
} from "@phosphor-icons/react/dist/ssr";

export default function SettingsPage() {
  const { user } = useAuth();
  const [steamAccount, setSteamAccount] = useState<SteamAccount | null>(null);
  const [lastfmAccount, setLastfmAccount] = useState<LastFMAccount | null>(
    null
  );
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAccounts = async () => {
      if (!user) return;
      try {
        const steam = steamService.getMySteamAccount();
        const lastfm = lastfmService.getMyLastFMAccount();
        const [steamRes, lastfmRes] = await Promise.allSettled([steam, lastfm]);

        if (steamRes.status === "fulfilled") {
          setSteamAccount(steamRes.value);
        }
        if (lastfmRes.status === "fulfilled") {
          setLastfmAccount(lastfmRes.value);
        }
      } catch (error) {
        console.error("Failed to fetch accounts", error);
      } finally {
        setLoading(false);
      }
    };
    fetchAccounts();
  }, [user]);

  const handleLinkSteam = () => {
    steamService.steamLogin();
  };

  const handleLinkLastFM = () => {
    lastfmService.lastfmAuth();
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl animate-in fade-in duration-500">
      <div className="flex items-center gap-3 mb-8">
        <div className="p-3 bg-primary/10 rounded-xl text-primary">
          <GearIcon size={32} weight="duotone" />
        </div>
        <div>
          <h1 className="text-3xl font-bold font-heading">Settings</h1>
          <p className="text-muted">Manage your account connections and preferences</p>
        </div>
      </div>

      <div className="space-y-8">
        <section>
          <div className="flex items-center gap-2 mb-4 text-foreground/80 font-semibold">
            <PlugsIcon size={20} weight="bold" />
            <h2>Integrations</h2>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <div className="bg-card border border-border rounded-xl p-6 relative overflow-hidden group hover:border-primary/50 transition-colors">
              <div className="absolute top-0 right-0 p-4 opacity-5 group-hover:opacity-10 transition-opacity">
                <GameControllerIcon size={120} weight="fill" />
              </div>
              
              <div className="relative z-10 flex flex-col h-full">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="p-2.5 bg-[#171a21] text-white rounded-lg shadow-sm">
                      <GameControllerIcon size={24} weight="fill" />
                    </div>
                    <div>
                      <h3 className="font-bold text-lg">Steam</h3>
                      <p className="text-xs text-muted">Game Library</p>
                    </div>
                  </div>
                  {steamAccount && (
                    <div className="flex items-center gap-1 text-emerald-500 text-xs font-bold uppercase tracking-wide bg-emerald-500/10 px-2 py-1 rounded-full">
                      <CheckCircleIcon size={14} weight="fill" />
                      Connected
                    </div>
                  )}
                </div>

                <div className="flex-1">
                  {loading ? (
                    <div className="h-10 bg-muted/20 rounded animate-pulse" />
                  ) : steamAccount ? (
                    <div className="bg-muted/30 p-3 rounded-lg border border-border mb-4">
                      <div className="flex items-center gap-3">
                        <Image
                          src={steamAccount.avatar}
                          alt="Avatar"
                          width={40}
                          height={40}
                          className="w-10 h-10 rounded-md"
                        />
                        <div>
                          <div className="font-bold text-sm">{steamAccount.persona}</div>
                          <div className="text-xs text-muted">ID: {steamAccount.steamId}</div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <p className="text-sm text-muted mb-4">
                      Link your Steam account to import your games, achievements, and playtime.
                    </p>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-border">
                  {steamAccount ? (
                    <Link href="/profile/steam" className="block">
                      <Button variant="outline" className="w-full">
                        Manage Integration
                      </Button>
                    </Link>
                  ) : (
                    <Button onClick={handleLinkSteam} className="w-full bg-[#171a21] hover:bg-[#2a475e] text-white border-transparent">
                      Link Steam Account
                    </Button>
                  )}
                </div>
              </div>
            </div>

            <div className="bg-card border border-border rounded-xl p-6 relative overflow-hidden group hover:border-red-500/30 transition-colors">
              <div className="absolute top-0 right-0 p-4 opacity-5 group-hover:opacity-10 transition-opacity text-red-500">
                <MusicNotesIcon size={120} weight="fill" />
              </div>
              
              <div className="relative z-10 flex flex-col h-full">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <div className="p-2.5 bg-[#ba0000] text-white rounded-lg shadow-sm">
                      <MusicNotesIcon size={24} weight="fill" />
                    </div>
                    <div>
                      <h3 className="font-bold text-lg">Last.fm</h3>
                      <p className="text-xs text-muted">Music History</p>
                    </div>
                  </div>
                   {lastfmAccount && (
                    <div className="flex items-center gap-1 text-emerald-500 text-xs font-bold uppercase tracking-wide bg-emerald-500/10 px-2 py-1 rounded-full">
                      <CheckCircleIcon size={14} weight="fill" />
                      Connected
                    </div>
                  )}
                </div>

                <div className="flex-1">
                  {loading ? (
                    <div className="h-10 bg-muted/20 rounded animate-pulse" />
                  ) : lastfmAccount ? (
                    <div className="bg-muted/30 p-3 rounded-lg border border-border mb-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-md bg-red-100 dark:bg-red-900/20 flex items-center justify-center text-red-600 dark:text-red-400">
                          <MusicNotesIcon size={20} />
                        </div>
                        <div>
                          <div className="font-bold text-sm">{lastfmAccount.username}</div>
                          <div className="text-xs text-muted">Scrobbles tracked</div>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <p className="text-sm text-muted mb-4">
                      Connect Last.fm to track your listening history and generate music charts.
                    </p>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-border">
                  {lastfmAccount ? (
                    <Link href="/profile/lastfm" className="block">
                      <Button variant="outline" className="w-full">
                        Manage Integration
                      </Button>
                    </Link>
                  ) : (
                    <Button onClick={handleLinkLastFM} className="w-full bg-[#ba0000] hover:bg-[#d51007] text-white border-transparent">
                      Link Last.fm Account
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
