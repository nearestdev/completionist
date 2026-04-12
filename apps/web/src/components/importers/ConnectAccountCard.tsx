"use client";

import { ArrowsClockwise, Link as LinkIcon, Trash } from "@phosphor-icons/react";
import { ConnectedAccount } from "@/types/connected_account";
import Button from "@/components/ui/Button";

interface ConnectAccountCardProps {
  provider: string;
  account?: ConnectedAccount;
  onConnect: (provider: string) => void;
  onSync: (provider: string) => void;
  onDisconnect: (provider: string) => void;
}

export default function ConnectAccountCard({ provider, account, onConnect, onSync, onDisconnect }: ConnectAccountCardProps) {
  return (
    <div className="bg-card rounded-xl border border-border p-4 flex items-center justify-between">
      <div>
        <h3 className="font-semibold text-foreground capitalize">{provider}</h3>
        {account ? (
          <p className="text-sm text-muted">
            {account.providerUsername || account.providerUserId || "Connected"}
            {account.lastSyncedAt && ` · Synced ${new Date(account.lastSyncedAt).toLocaleDateString()}`}
          </p>
        ) : (
          <p className="text-sm text-muted">Not connected</p>
        )}
      </div>
      <div className="flex gap-2">
        {account ? (
          <>
            <Button size="sm" variant="outline" onClick={() => onSync(provider)}>
              <ArrowsClockwise size={14} className="mr-1" /> Sync
            </Button>
            <Button size="sm" variant="ghost" onClick={() => onDisconnect(provider)}>
              <Trash size={14} />
            </Button>
          </>
        ) : (
          <Button size="sm" variant="primary" onClick={() => onConnect(provider)}>
            <LinkIcon size={14} className="mr-1" /> Connect
          </Button>
        )}
      </div>
    </div>
  );
}
