"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import socialService from "@/services/socialService";
import postService from "@/services/postService";
import steamService from "@/services/steamService";
import lastfmService from "@/services/lastfmService";
import { EnhancedUserProfile } from "@/types/social";
import { PostWithDetails } from "@/types/posts";
import { SteamAccount } from "@/types/steam";
import { LastFMAccount } from "@/types/lastfm";
import { useAuth } from "@/hooks/useAuth";
import Button from "@/components/ui/Button";
import PostCard from "@/components/posts/PostCard";
import Link from "next/link";
import Image from "next/image";
import { 
  UsersIcon, 
  GameControllerIcon, 
  MusicNotesIcon, 
  GearIcon,
  UserCircleIcon,
  ChatCircleTextIcon 
} from "@phosphor-icons/react/dist/ssr";

export default function ProfilePage() {
  const { username } = useParams<{ username: string }>();
  const { user: currentUser } = useAuth();
  const [profile, setProfile] = useState<EnhancedUserProfile | null>(null);
  const [posts, setPosts] = useState<PostWithDetails[]>([]);
  const [steamAccount, setSteamAccount] = useState<SteamAccount | null>(null);
  const [lastfmAccount, setLastfmAccount] = useState<LastFMAccount | null>(
    null
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const isOwnProfile = currentUser?.username === username;

  useEffect(() => {
    const fetchProfileData = async () => {
      if (!username) return;
      try {
        setLoading(true);
        const profileData = await socialService.getUserProfile(username);
        setProfile(profileData);
        const postsData = await postService.getUserPosts(username);
        setPosts(postsData);
        if (isOwnProfile) {
          try {
            const steamData = await steamService.getMySteamAccount();
            setSteamAccount(steamData);
          } catch (steamError) {
            console.log("No Steam account linked for this user.");
          }
          try {
            const lastfmData = await lastfmService.getMyLastFMAccount();
            setLastfmAccount(lastfmData);
          } catch (lastfmError) {
            console.log("No Last.fm account linked for this user.");
          }
        }
      } catch (err) {
        setError("Failed to load profile.");
      } finally {
        setLoading(false);
      }
    };
    fetchProfileData();
  }, [username, isOwnProfile]);

  const handleFollow = async () => {
    if (!profile) return;
    try {
      await socialService.followUser(profile.username);
      setProfile((prev) =>
        prev
          ? {
              ...prev,
              isFollowing: true,
              stats: {
                ...prev.stats,
                followersCount: prev.stats.followersCount + 1,
              },
            }
          : null
      );
    } catch (error) {
      console.error("Failed to follow user", error);
    }
  };

  const handleUnfollow = async () => {
    if (!profile) return;
    try {
      await socialService.unfollowUser(profile.username);
      setProfile((prev) =>
        prev
          ? {
              ...prev,
              isFollowing: false,
              stats: {
                ...prev.stats,
                followersCount: prev.stats.followersCount - 1,
              },
            }
          : null
      );
    } catch (error) {
      console.error("Failed to unfollow user", error);
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>
  );

  if (error) return (
    <div className="flex items-center justify-center min-h-[40vh] text-red-500 font-medium">
      {error}
    </div>
  );

  if (!profile) return <div>Profile not found.</div>;

  return (
    <div className="container mx-auto px-4 py-8 max-w-5xl animate-in fade-in duration-500">
      
      {/* Profile Header Card */}
      <div className="relative overflow-hidden bg-card rounded-3xl shadow-xl mb-10 border border-border group">
        {/* Background Pattern */}
        <div className="absolute inset-0 bg-gradient-to-r from-indigo-500/10 to-purple-500/10 opacity-50" />
        <div className="absolute top-0 right-0 p-12 opacity-5 transform rotate-12">
          <UserCircleIcon size={300} weight="duotone" />
        </div>

        <div className="relative z-10 p-8 flex flex-col md:flex-row items-center md:items-start gap-8">
          {/* Avatar */}
          <div className="relative shrink-0">
             <div className="w-32 h-32 rounded-full ring-4 ring-white dark:ring-gray-700 shadow-2xl overflow-hidden bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
               {steamAccount?.avatar ? (
                 <Image
                   src={steamAccount.avatar}
                   alt={`${profile.username}'s avatar`}
                   width={128}
                   height={128}
                   className="object-cover"
                 />
               ) : (
                 <UserCircleIcon size={80} className="text-gray-400" weight="fill" />
               )}
             </div>
             {/* Integration Badges */}
             <div className="absolute -bottom-2 -right-2 flex gap-1">
               {steamAccount && (
                 <div className="bg-[#171a21] text-white p-1.5 rounded-full shadow-md border border-gray-700" title="Steam Linked">
                   <GameControllerIcon size={16} weight="fill" />
                 </div>
               )}
               {lastfmAccount && (
                 <div className="bg-[#ba0000] text-white p-1.5 rounded-full shadow-md border border-red-900" title="Last.fm Linked">
                   <MusicNotesIcon size={16} weight="fill" />
                 </div>
               )}
             </div>
          </div>

          {/* Info */}
          <div className="flex-1 text-center md:text-left space-y-4">
            <div>
              <h1 className="text-4xl font-heading font-bold text-foreground tracking-tight mb-2">
                {profile.username}
              </h1>
              <p className="text-muted text-lg max-w-2xl mx-auto md:mx-0 leading-relaxed">
                {profile.bio || "This user prefers to keep their bio mysterious."}
              </p>
            </div>

            {/* Stats */}
            <div className="flex items-center justify-center md:justify-start gap-8 pt-2">
              <Link href={`/profile/${profile.username}/followers`} className="group flex items-center gap-2 hover:text-primary transition-colors">
                <div className="bg-primary/10 p-2 rounded-lg text-primary group-hover:bg-primary/20 transition-colors">
                  <UsersIcon size={20} weight="bold" />
                </div>
                <div className="text-left">
                  <span className="block font-bold text-lg leading-none">{profile.stats.followersCount}</span>
                  <span className="text-xs text-muted font-medium uppercase tracking-wider">Followers</span>
                </div>
              </Link>
              
              <Link href={`/profile/${profile.username}/following`} className="group flex items-center gap-2 hover:text-primary transition-colors">
                 <div className="bg-primary/10 p-2 rounded-lg text-primary group-hover:bg-primary/20 transition-colors">
                  <UsersIcon size={20} weight="duotone" />
                </div>
                <div className="text-left">
                  <span className="block font-bold text-lg leading-none">{profile.stats.followingCount}</span>
                  <span className="text-xs text-muted font-medium uppercase tracking-wider">Following</span>
                </div>
              </Link>
            </div>
            
            {/* Actions */}
            <div className="pt-4 flex flex-wrap justify-center md:justify-start gap-3">
              {!isOwnProfile && currentUser && (
                 profile.isFollowing ? (
                  <Button onClick={handleUnfollow} variant="outline" className="border-red-200 text-red-600 hover:bg-red-50 hover:border-red-300 dark:border-red-900/30 dark:text-red-400 dark:hover:bg-red-900/10">
                    Unfollow
                  </Button>
                ) : (
                  <Button onClick={handleFollow} className="bg-primary hover:bg-primary-hover text-white shadow-lg shadow-primary/20">
                    Follow
                  </Button>
                )
              )}
              
              {isOwnProfile && (
                <Link href="/settings">
                  <Button variant="outline" className="gap-2">
                    <GearIcon size={18} />
                    Edit Profile
                  </Button>
                </Link>
              )}
            </div>
          </div>
        </div>

        {/* Integration Details (Own Profile Only) */}
        {isOwnProfile && (steamAccount || lastfmAccount) && (
           <div className="border-t border-border bg-muted/5 p-4 flex flex-wrap justify-center md:justify-start gap-4">
              {steamAccount && (
                <Link href="/profile/steam">
                   <div className="flex items-center gap-2 px-4 py-2 bg-[#171a21]/5 hover:bg-[#171a21]/10 dark:bg-[#171a21]/40 dark:hover:bg-[#171a21]/60 rounded-lg text-sm font-medium transition-colors cursor-pointer">
                      <GameControllerIcon size={18} />
                      Manage Steam
                   </div>
                </Link>
              )}
              {lastfmAccount && (
                <Link href="/profile/lastfm">
                   <div className="flex items-center gap-2 px-4 py-2 bg-[#ba0000]/5 hover:bg-[#ba0000]/10 dark:bg-[#ba0000]/20 dark:hover:bg-[#ba0000]/30 rounded-lg text-sm font-medium transition-colors cursor-pointer text-red-700 dark:text-red-400">
                      <MusicNotesIcon size={18} />
                      Manage Last.fm
                   </div>
                </Link>
              )}
           </div>
        )}
      </div>

      {/* Posts Section */}
      <div className="space-y-6">
        <div className="flex items-center gap-3 mb-6 border-b border-border pb-4">
           <ChatCircleTextIcon size={28} weight="duotone" className="text-primary" />
           <h2 className="text-2xl font-bold font-heading">Recent Activity</h2>
        </div>

        <div className="grid gap-6">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
          {posts.length === 0 && (
            <div className="text-center py-16 bg-card rounded-2xl border border-dashed border-border text-muted">
              <ChatCircleTextIcon size={48} className="mx-auto mb-4 opacity-20" />
              <p className="text-lg font-medium">No posts yet</p>
              <p className="text-sm opacity-70">This user hasn't shared any updates.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}