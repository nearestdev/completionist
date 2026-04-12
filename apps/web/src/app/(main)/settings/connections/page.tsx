"use client";

import { useEffect, useState } from "react";
import { Link as LinkIcon, ArrowsClockwise, Trash } from "@phosphor-icons/react";
import importService from "@/services/importService";
import { ConnectedAccount, ImportJob } from "@/types/connected_account";
import Button from "@/components/ui/Button";

const PROVIDERS = ["myanimelist", "anilist", "trakt", "letterboxd", "goodreads"];

export default function ConnectionsPage() {
  const [accounts, setAccounts] = useState<ConnectedAccount[]>([]);
  const [jobs, setJobs] = useState<ImportJob[]>([]);

  useEffect(() => {
    importService.getConnections().then((data) => setAccounts(data || [])).catch(() => {});
    importService.getImportJobs().then((data) => setJobs(data || [])).catch(() => {});
  }, []);

  const connected = new Map(accounts.map((a) => [a.provider, a]));

  const handleConnect = async (provider: string) => {
    try {
      const { url } = await importService.connectProvider(provider);
      window.location.href = url;
    } catch {}
  };

  const handleSync = async (provider: string) => {
    try {
      const job = await importService.triggerSync(provider);
      setJobs([job, ...jobs]);
    } catch {}
  };

  const handleDisconnect = async (provider: string) => {
    await importService.disconnect(provider);
    setAccounts(accounts.filter((a) => a.provider !== provider));
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold text-foreground mb-8">Connected Accounts</h1>

      <div className="space-y-4 mb-8">
        {PROVIDERS.map((provider) => {
          const account = connected.get(provider);
          return (
            <div key={provider} className="bg-card rounded-xl border border-border p-4 flex items-center justify-between">
              <div>
                <h3 className="font-semibold text-foreground capitalize">{provider}</h3>
                {account ? (
                  <p className="text-sm text-muted">
                    Connected as {account.providerUsername || account.providerUserId}
                    {account.lastSyncedAt && ` · Last synced ${new Date(account.lastSyncedAt).toLocaleDateString()}`}
                  </p>
                ) : (
                  <p className="text-sm text-muted">Not connected</p>
                )}
              </div>
              <div className="flex gap-2">
                {account ? (
                  <>
                    <Button size="sm" variant="outline" onClick={() => handleSync(provider)}>
                      <ArrowsClockwise size={14} className="mr-1" /> Sync
                    </Button>
                    <Button size="sm" variant="ghost" onClick={() => handleDisconnect(provider)}>
                      <Trash size={14} />
                    </Button>
                  </>
                ) : (
                  <Button size="sm" variant="primary" onClick={() => handleConnect(provider)}>
                    <LinkIcon size={14} className="mr-1" /> Connect
                  </Button>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {jobs.length > 0 && (
        <>
          <h2 className="text-lg font-semibold text-foreground mb-4">Import History</h2>
          <div className="space-y-3">
            {jobs.map((job) => (
              <div key={job.id} className="bg-card rounded-lg border border-border p-3">
                <div className="flex justify-between text-sm">
                  <span className="font-medium text-foreground capitalize">{job.provider}</span>
                  <span className={`text-xs px-2 py-0.5 rounded-full ${
                    job.status === "completed" ? "bg-accent/10 text-accent" :
                    job.status === "running" ? "bg-primary/10 text-primary" :
                    job.status === "failed" ? "bg-red-500/10 text-red-500" : "bg-muted/10 text-muted"
                  }`}>{job.status}</span>
                </div>
                <p className="text-xs text-muted mt-1">
                  {job.importedItems} imported · {job.skippedItems} skipped · {job.failedItems} failed
                </p>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
