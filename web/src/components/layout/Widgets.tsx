"use client";
import { useEffect, useState } from "react";
import Link from "next/link";
import { UserPlusIcon, UsersIcon } from "@phosphor-icons/react/dist/ssr";
import socialService from "@/services/socialService";
import { UserFollowResponse } from "@/types/social";
import { useAuth } from "@/hooks/useAuth";
import { getUserRank } from "@/services/rankService";
import { UserRank } from "@/types/rank";
import RankBadge from "@/components/rank/RankBadge";

export default function Widgets() {
  const { user } = useAuth();
  const [suggestions, setSuggestions] = useState<UserFollowResponse[]>([]);
  const [rank, setRank] = useState<UserRank | null>(null);
  const [loading, setLoading] = useState(false);
  const [rankLoading, setRankLoading] = useState(false);

  useEffect(() => {
    if (!user) return;

    const fetchSuggestions = async () => {
      try {
        setLoading(true);
        const data = await socialService.getFollowSuggestions();
        setSuggestions(data);
      } catch (error) {
        console.error("Failed to load suggestions", error);
      } finally {
        setLoading(false);
      }
    };

    const fetchRank = async () => {
      try {
        setRankLoading(true);
        const data = await getUserRank();
        setRank(data);
      } catch (error) {
        console.error("Failed to load rank", error);
      } finally {
        setRankLoading(false);
      }
    };

    fetchSuggestions();
    fetchRank();
  }, [user]);

  if (!user) {
    return null;
  }

  if (!loading && suggestions.length === 0) {
    return (
      <aside className="hidden xl:flex flex-col w-[320px] p-6 space-y-6 h-screen sticky top-0 border-l border-border bg-background">
        <div className="text-sm text-muted">No suggestions available.</div>
      </aside>
    );
  }

  return (
    <aside className="hidden xl:flex flex-col w-[320px] p-6 space-y-6 h-screen sticky top-0 overflow-y-auto border-l border-border bg-background">
      
      {/* Rank Card */}
      {rank && (
        <div className="w-full">
           <RankBadge rank={rank} />
        </div>
      )}

      <div className="bg-card rounded-xl p-5 border border-border shadow-sm">
        <div className="flex items-center gap-2 mb-4 text-foreground">
          <UsersIcon size={20} weight="bold" className="text-primary" />
          <h3 className="font-heading font-semibold text-lg">Who to Follow</h3>
        </div>
        {loading ? (
          <div className="text-sm text-muted animate-pulse">Loading suggestions...</div>
        ) : (
          <ul className="space-y-4">
            {suggestions.map((user) => (
              <li key={user.userId} className="flex items-center justify-between gap-3">
                <div className="flex items-center gap-3 overflow-hidden">
                  <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-primary font-bold text-xs flex-shrink-0">
                    {user.username.charAt(0).toUpperCase()}
                  </div>
                  <div className="flex flex-col min-w-0">
                    <Link 
                      href={`/profile/${user.username}`}
                      className="font-bold text-sm text-foreground hover:underline truncate"
                    >
                      {user.username}
                    </Link>
                    <span className="text-xs text-muted truncate">
                      Joined {new Date(user.followedAt).getFullYear()}
                    </span>
                  </div>
                </div>
                <Link href={`/profile/${user.username}`}>
                    <button 
                      className="p-1.5 text-primary hover:bg-indigo-50 rounded-full transition-colors"
                      title="View Profile"
                    >
                      <UserPlusIcon size={18} />
                    </button>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </aside>
  );
}