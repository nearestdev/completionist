"use client";
import { useEffect, useState } from "react";
import steamService from "@/services/steamService";
import lastfmService from "@/services/lastfmService";
import { SteamAccount } from "@/types/steam";
import { LastFMAccount } from "@/types/lastfm";
import Button from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";

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
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Settings</h1>
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Account Integrations</h2>
        {loading ? (
          <p>Loading...</p>
        ) : (
          <div className="space-y-6">
            <div>
              <h3 className="font-semibold text-lg">Steam</h3>
              {steamAccount ? (
                <div className="mt-2">
                  <p className="mb-2">
                    <strong>Account Linked:</strong> {steamAccount.persona}
                  </p>
                  <Link href="/profile/steam">
                    <Button className="max-w-xs">View Steam Details</Button>
                  </Link>
                </div>
              ) : (
                <>
                  <p className="mt-2">No Steam account linked.</p>
                  <Button onClick={handleLinkSteam} className="mt-2 max-w-xs">
                    Link Steam Account
                  </Button>
                </>
              )}
            </div>
            <div className="border-t pt-6">
              <h3 className="font-semibold text-lg">Last.fm</h3>
              {lastfmAccount ? (
                <div className="mt-2">
                  <p className="mb-2">
                    <strong>Account Linked:</strong> {lastfmAccount.username}
                  </p>
                  <Link href="/profile/lastfm">
                    <Button className="max-w-xs">View Last.fm Details</Button>
                  </Link>
                </div>
              ) : (
                <>
                  <p className="mt-2">No Last.fm account linked.</p>
                  <Button onClick={handleLinkLastFM} className="mt-2 max-w-xs">
                    Link Last.fm Account
                  </Button>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}