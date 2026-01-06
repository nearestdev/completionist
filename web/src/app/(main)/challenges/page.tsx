"use client";

import { useEffect, useState } from "react";
import challengeService from "@/services/challengeService";
import { LeaderboardEntry, Season, UserChallenge } from "@/types/challenge";
import ChallengeCard from "@/components/challenges/ChallengeCard";
import { CrownIcon, CalendarBlankIcon } from "@phosphor-icons/react/dist/ssr";
import Link from "next/link";

export default function ChallengesPage() {
  const [season, setSeason] = useState<Season | null>(null);
  const [challenges, setChallenges] = useState<UserChallenge[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [seasonRes, challengesRes, leaderboardRes] =
          await Promise.allSettled([
            challengeService.getActiveSeason(),
            challengeService.getMyChallenges(),
            challengeService.getLeaderboard(),
          ]);

        if (seasonRes.status === "fulfilled") {
          setSeason(seasonRes.value);
        }
        if (challengesRes.status === "fulfilled") {
          setChallenges(challengesRes.value);
        }
        if (leaderboardRes.status === "fulfilled") {
          setLeaderboard(leaderboardRes.value);
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
    return <div className="p-8 text-center text-gray-500">Loading challenges...</div>;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Season Header */}
      <div className="mb-10 bg-gradient-to-r from-indigo-600 to-purple-600 rounded-2xl p-8 text-white shadow-lg">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4">
          <div>
            <div className="flex items-center gap-2 mb-2 opacity-90">
              <CalendarBlankIcon size={20} />
              <span className="text-sm font-medium uppercase tracking-wider">
                Current Season
              </span>
            </div>
            <h1 className="text-3xl md:text-4xl font-bold font-heading">
              {season ? season.name : "Pre-Season"}
            </h1>
            {season && (
              <p className="mt-2 opacity-80">
                Ends on {new Date(season.endAt).toLocaleDateString()}
              </p>
            )}
          </div>
          <div className="bg-white/20 backdrop-blur-sm p-4 rounded-xl text-center min-w-[120px]">
            <div className="text-sm uppercase font-bold opacity-75">Your Rank</div>
            <div className="text-3xl font-bold">#42</div> 
            {/* Note: Rank logic would ideally come from a specific 'getMyRank' endpoint or derived from leaderboard if loaded fully */}
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Main Column: Challenges */}
        <div className="lg:col-span-2 space-y-8">
          <section>
            <h2 className="text-2xl font-bold mb-4 flex items-center gap-2">
              <span className="bg-indigo-100 text-indigo-600 p-1.5 rounded-lg">
                🏆
              </span>{" "}
              Active Challenges
            </h2>
            
            {challenges.length > 0 ? (
              <div className="grid sm:grid-cols-2 gap-4">
                {challenges.map((c) => (
                  <ChallengeCard key={c.id} challenge={c} />
                ))}
              </div>
            ) : (
              <div className="bg-white dark:bg-gray-800 p-8 rounded-xl text-center border border-dashed border-gray-300 dark:border-gray-700">
                <p className="text-gray-500 mb-4">No active challenges found.</p>
                <p className="text-sm text-gray-400">
                  Complete items in your list or write reviews to unlock hidden achievements!
                </p>
              </div>
            )}
          </section>
        </div>

        {/* Sidebar Column: Leaderboard */}
        <div className="space-y-6">
          <section className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-6">
            <h2 className="text-xl font-bold mb-6 flex items-center gap-2">
              <CrownIcon size={24} className="text-amber-500" weight="fill" />
              Leaderboard
            </h2>
            
            <div className="space-y-4">
              {leaderboard.map((entry, index) => (
                <div 
                  key={entry.userId} 
                  className={`flex items-center justify-between p-3 rounded-lg ${
                    index < 3 ? "bg-amber-50 dark:bg-amber-900/10" : ""
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <div className={`
                      w-8 h-8 flex items-center justify-center rounded-full font-bold text-sm
                      ${index === 0 ? "bg-amber-400 text-white" : 
                        index === 1 ? "bg-gray-300 text-gray-700" : 
                        index === 2 ? "bg-orange-300 text-orange-800" : 
                        "bg-gray-100 dark:bg-gray-700 text-gray-500"}
                    `}>
                      {entry.rank}
                    </div>
                    <Link 
                      href={`/profile/${entry.username}`} 
                      className="font-medium hover:text-indigo-600 transition-colors"
                    >
                      {entry.username}
                    </Link>
                  </div>
                  <div className="font-mono text-sm font-semibold text-gray-500">
                    {entry.totalXp} XP
                  </div>
                </div>
              ))}
              
              {leaderboard.length === 0 && (
                <p className="text-center text-gray-500 text-sm">
                  Be the first to claim the throne!
                </p>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}