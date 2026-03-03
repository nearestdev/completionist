"use client";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import socialService from "@/services/socialService";
import { UserFollowResponse } from "@/types/social";
import UserCard from "@/components/users/UserCard";
export default function FollowingPage() {
  const { username } = useParams<{ username: string }>();
  const [following, setFollowing] = useState<UserFollowResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    const fetchFollowing = async () => {
      if (!username) return;
      try {
        setLoading(true);
        const data = await socialService.getUserFollowing(username);
        setFollowing(data);
      } catch (err) {
        setError("Failed to load following list.");
      } finally {
        setLoading(false);
      }
    };
    fetchFollowing();
  }, [username]);
  if (loading) return <div>Loading...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Following</h1>
      {following.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {following.map((user) => (
            <UserCard key={user.userId} user={user} />
          ))}
        </div>
      ) : (
        <p>This user is not following anyone yet.</p>
      )}
    </div>
  );
}