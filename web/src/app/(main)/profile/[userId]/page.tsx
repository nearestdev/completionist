"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import socialService from "@/services/socialService";
import postService from "@/services/postService";
import steamService from "@/services/steamService";
import { EnhancedUserProfile } from "@/types/social";
import { PostWithDetails } from "@/types/posts";
import { SteamAccount } from "@/types/steam";
import { useAuth } from "@/hooks/useAuth";
import Button from "@/components/ui/Button";
import PostCard from "@/components/posts/PostCard";
import Link from "next/link";
import Image from "next/image";

export default function ProfilePage() {
  const { userId } = useParams();
  const { user: currentUser } = useAuth();
  const [profile, setProfile] = useState<EnhancedUserProfile | null>(null);
  const [posts, setPosts] = useState<PostWithDetails[]>([]);
  const [steamAccount, setSteamAccount] = useState<SteamAccount | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProfileData = async () => {
      if (!userId) return;
      try {
        setLoading(true);
        const profileData = await socialService.getUserProfile(Number(userId));
        setProfile(profileData);
        const postsData = await postService.getUserPosts(Number(userId));
        setPosts(postsData);
        try {
          const steamData = await steamService.getUserSteamAccount(
            Number(userId)
          );
          setSteamAccount(steamData);
        } catch (steamError) {
          console.log("No Steam account linked for this user.");
        }
      } catch (err) {
        setError("Failed to load profile.");
      } finally {
        setLoading(false);
      }
    };
    fetchProfileData();
  }, [userId]);

  const handleFollow = async () => {
    if (!profile) return;
    try {
      await socialService.followUser(profile.id);
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
      await socialService.unfollowUser(profile.id);
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

  const isOwnProfile = currentUser?.id === profile.id;

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
            href={`/profile/${profile.id}/followers`}
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
            href={`/profile/${profile.id}/following`}
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
        {isOwnProfile && steamAccount && (
          <div className="mt-6 text-center">
            <Link href="/profile/my-games">
              <Button>My Games</Button>
            </Link>
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