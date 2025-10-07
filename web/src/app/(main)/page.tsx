"use client";

import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";

export default function HomePage() {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-4">Welcome to Completionist</h1>
      {user ? (
        <p>
          Hello, {user.username}! This is your dashboard. More features coming
          soon.
        </p>
      ) : (
        <p>
          Please{" "}
          <Link href="/login" className="text-indigo-600 hover:underline">
            login
          </Link>{" "}
          or{" "}
          <Link href="/register" className="text-indigo-600 hover:underline">
            register
          </Link>{" "}
          to continue.
        </p>
      )}
    </div>
  );
}