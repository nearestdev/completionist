"use client";
import { useEffect, useState } from "react";
import steamService from "@/services/steamService";
import { SteamAccount } from "@/types/steam";
import Button from "@/components/ui/Button";
import { useAuth } from "@/hooks/useAuth";
export default function SettingsPage() {
  const { user } = useAuth();
  const [steamAccount, setSteamAccount] = useState<SteamAccount | null>(null);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchSteamAccount = async () => {
      try {
        const account = await steamService.getMySteamAccount();
        setSteamAccount(account);
      } catch (error) {
        console.log("No Steam account linked.");
      } finally {
        setLoading(false);
      }
    };
    if (user) {
      fetchSteamAccount();
    }
  }, [user]);
  const handleLinkSteam = () => {
    steamService.steamLogin();
  };
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Settings</h1>
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Account Integrations</h2>
        {loading ? (
          <p>Loading...</p>
        ) : steamAccount ? (
          <div>
            <p>
              <strong>Steam Account Linked:</strong> {steamAccount.persona}
            </p>
          </div>
        ) : (
          <div>
            <p>No Steam account linked.</p>
            <Button onClick={handleLinkSteam} className="mt-2 max-w-xs">
              Link Steam Account
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}