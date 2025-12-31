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
  if (loading) return <div>Loading profile...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!profile) return <div>Profile not found.</div>;
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="bg-white dark:bg-gray-800 shadow-md rounded-lg p-6 mb-8">
        <div className="flex items-center space-x-4">
          {steamAccount && (
            <Image
              src={steamAccount.avatar}
              alt={`${profile.username}'s Steam avatar`}
              width={80}
              height={80}
              className="rounded-full"
            />
          )}
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              {profile.username}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              {profile.bio || "No bio yet."}
            </p>
          </div>
        </div>
        <div className="mt-6 flex justify-around text-center">
          <Link
            href={`/profile/${profile.username}/followers`}
            className="hover:underline"
          >
            <div>
              <p className="font-semibold text-lg">
                {profile.stats.followersCount}
              </p>
              <p className="text-gray-600 dark:text-gray-400">Followers</p>
            </div>
          </Link>
          <Link
            href={`/profile/${profile.username}/following`}
            className="hover:underline"
          >
            <div>
              <p className="font-semibold text-lg">
                {profile.stats.followingCount}
              </p>
              <p className="text-gray-600 dark:text-gray-400">Following</p>
            </div>
          </Link>
        </div>
        {!isOwnProfile && currentUser && (
          <div className="mt-6 flex justify-center">
            {profile.isFollowing ? (
              <Button onClick={handleUnfollow}>Unfollow</Button>
            ) : (
              <Button onClick={handleFollow}>Follow</Button>
            )}
          </div>
        )}
        {isOwnProfile && (steamAccount || lastfmAccount) && (
          <div className="mt-6 border-t pt-4">
            <h3 className="font-semibold text-lg mb-2 text-center">
              My Integrations
            </h3>
            <div className="flex gap-4 justify-center">
              {steamAccount && (
                <Link href="/profile/steam">
                  <Button>Steam Details</Button>
                </Link>
              )}
              {lastfmAccount && (
                <Link href="/profile/lastfm">
                  <Button>Last.fm Details</Button>
                </Link>
              )}
            </div>
          </div>
        )}
      </div>
      <div>
        <h2 className="text-2xl font-bold mb-4">Posts</h2>
        <div className="space-y-6">
          {posts.map((post) => (
            <PostCard key={post.id} post={post} />
          ))}
          {posts.length === 0 && <p>This user has not made any posts yet.</p>}
        </div>
      </div>
    </div>
  );
}