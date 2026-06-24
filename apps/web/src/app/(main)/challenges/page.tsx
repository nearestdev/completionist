"use client";

import { useEffect, useState } from "react";
import challengeService from "@/services/challengeService";
import { getUserRank } from "@/services/rankService";
import { LeaderboardEntry, Season, UserChallenge } from "@/types/challenge";
import { UserRank, RankNames } from "@/types/rank";
import ChallengeCard from "@/components/challenges/ChallengeCard";
import { 
  CrownIcon, 
  CalendarCheckIcon, 
  TargetIcon, 
  TrendUpIcon,
  TrophyIcon,
  MedalIcon,
  ShieldCheckIcon
} from "@phosphor-icons/react/dist/ssr";
import Link from "next/link";

export default function ChallengesPage() {
  const [season, setSeason] = useState<Season | null>(null);
  const [challenges, setChallenges] = useState<UserChallenge[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [userRank, setUserRank] = useState<UserRank | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [seasonRes, challengesRes, leaderboardRes, userRankRes] =
          await Promise.allSettled([
            challengeService.getActiveSeason(),
            challengeService.getMyChallenges(),
            challengeService.getLeaderboard(),
            getUserRank(),
          ]);

        if (seasonRes.status === "fulfilled") {
          setSeason(seasonRes.value);
        }
        if (challengesRes.status === "fulfilled") {
          setChallenges(challengesRes.value || []);
        }
        if (leaderboardRes.status === "fulfilled") {
          setLeaderboard(leaderboardRes.value);
        }
        if (userRankRes.status === "fulfilled") {
          setUserRank(userRankRes.value);
        }
      } catch (error) {
        console.error("Failed to load challenge data", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="flex flex-col items-center gap-4 text-muted">
          <div className="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
          <p className="font-medium animate-pulse">Loading season data...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl animate-in fade-in duration-500">
      
      <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-indigo-900 to-purple-900 text-white shadow-2xl mb-12">
        <div className="absolute inset-0 bg-[url('/patterns/grid.svg')] opacity-10"></div>
        <div className="absolute top-0 right-0 p-12 opacity-10 transform rotate-12">
          <TrophyIcon size={400} weight="duotone" />
        </div>
        
        <div className="relative z-10 p-8 md:p-12">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
            <div className="space-y-4">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 text-xs font-bold uppercase tracking-wider text-indigo-100">
                <CalendarCheckIcon size={14} weight="bold" />
                <span>Current Season</span>
              </div>
              
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold font-heading tracking-tight text-white mb-2 drop-shadow-sm">
                {season ? season.name : "Pre-Season"}
              </h1>
              
              {season && (
                <p className="text-indigo-200 text-lg max-w-md">
                  Complete challenges to earn XP, climb the ranks, and unlock exclusive rewards before {new Date(season.endAt).toLocaleDateString()}.
                </p>
              )}
            </div>

            <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/10 w-full md:w-auto min-w-[280px] shadow-xl">
              <div className="flex items-center gap-4 mb-6">
                <div className="p-3 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl shadow-lg ring-4 ring-white/5">
                  <ShieldCheckIcon size={32} weight="fill" className="text-white" />
                </div>
                <div>
                  <div className="text-xs uppercase font-bold text-indigo-200 tracking-wider">Current Rank</div>
                  <div className="text-2xl font-bold font-heading text-white">
                    {userRank ? RankNames[userRank.currentRank] : "Unranked"}
                  </div>
                </div>
              </div>

              <div className="space-y-3">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-indigo-200 flex items-center gap-2">
                    <TrendUpIcon size={16} />
                    Current Elo
                  </span>
                  <span className="font-bold font-mono text-white">{userRank?.currentElo || 0}</span>
                </div>
                <div className="h-1.5 w-full bg-black/20 rounded-full overflow-hidden">
                  <div 
                    className="h-full bg-gradient-to-r from-amber-300 to-orange-400 rounded-full"
                    style={{ width: `${Math.min(100, ((userRank?.currentElo || 0) / 3000) * 100)}%` }}
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-12 gap-8 items-start">
        
        <div className="lg:col-span-8 space-y-8">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold font-heading flex items-center gap-3 text-foreground">
              <div className="p-2 bg-primary/10 text-primary rounded-lg">
                <TargetIcon size={24} weight="duotone" />
              </div>
              Active Challenges
            </h2>
            <div className="text-sm font-medium text-muted bg-card border border-border px-3 py-1 rounded-full shadow-sm">
              {challenges.length} Available
            </div>
          </div>

          {challenges.length > 0 ? (
            <div className="grid md:grid-cols-2 gap-5">
              {challenges.map((c) => (
                <ChallengeCard key={c.id} challenge={c} />
              ))}
            </div>
          ) : (
            <div className="bg-card rounded-2xl p-12 text-center border border-dashed border-border">
              <div className="w-16 h-16 bg-muted/20 text-muted rounded-full flex items-center justify-center mx-auto mb-4">
                <CalendarCheckIcon size={32} />
              </div>
              <h3 className="text-lg font-bold text-foreground mb-2">All Caught Up!</h3>
              <p className="text-muted max-w-sm mx-auto">
                You&apos;ve engaged with all active challenges. Check back later for new quests or explore the community lists.
              </p>
            </div>
          )}
        </div>

        <div className="lg:col-span-4 space-y-6">
          <section className="bg-card rounded-2xl shadow-sm border border-border overflow-hidden">
            <div className="p-6 border-b border-border bg-gradient-to-r from-card to-muted/5">
              <h2 className="text-xl font-bold font-heading flex items-center gap-2 text-foreground">
                <CrownIcon size={24} weight="duotone" className="text-amber-500" />
                Top Players
              </h2>
            </div>
            
            <div className="divide-y divide-border">
              {leaderboard.map((entry, index) => (
                <div 
                  key={entry.userId} 
                  className={`p-4 flex items-center gap-4 transition-colors hover:bg-muted/5 ${
                    index < 3 ? "bg-gradient-to-r from-amber-50/50 to-transparent dark:from-amber-900/10" : ""
                  }`}
                >
                  <div className={`
                    w-10 h-10 flex items-center justify-center rounded-xl font-bold text-sm shadow-sm
                    ${index === 0 ? "bg-gradient-to-br from-amber-300 to-amber-500 text-white ring-2 ring-amber-200 dark:ring-amber-900" : 
                      index === 1 ? "bg-gradient-to-br from-gray-200 to-gray-400 text-gray-800" : 
                      index === 2 ? "bg-gradient-to-br from-orange-200 to-orange-400 text-orange-900" : 
                      "bg-muted text-muted-foreground"}
                  `}>
                    {index < 3 ? <MedalIcon size={20} weight="fill" /> : entry.rank}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <Link 
                      href={`/profile/${entry.username}`} 
                      className="font-bold text-foreground hover:text-primary truncate block transition-colors"
                    >
                      {entry.username}
                    </Link>
                    <div className="text-xs text-muted flex items-center gap-1">
                      <span className="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                      Active
                    </div>
                  </div>
                  
                  <div className="text-right">
                    <div className="font-bold font-mono text-primary">{entry.totalXp}</div>
                    <div className="text-[10px] uppercase font-bold text-muted">XP</div>
                  </div>
                </div>
              ))}
              
              {leaderboard.length === 0 && (
                <div className="p-8 text-center text-muted">
                  <p>Leaderboard is building...</p>
                </div>
              )}
            </div>
            
            <div className="p-4 bg-muted/5 border-t border-border text-center">
              <Link href="/leaderboard" className="text-sm font-bold text-primary hover:text-primary-hover transition-colors flex items-center justify-center gap-1">
                View Full Rankings
                <TrendUpIcon size={16} weight="bold" />
              </Link>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
