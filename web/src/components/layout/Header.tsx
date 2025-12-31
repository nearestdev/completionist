"use client";
import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";
export default function Header() {
  const { user, logout } = useAuth();
  return (
    <header className="bg-white shadow dark:bg-gray-800">
      <nav className="container mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="text-xl font-semibold text-gray-700 dark:text-gray-100">
            <Link href="/" className="text-gray-800 dark:text-gray-100 hover:text-gray-700 dark:hover:text-gray-300">
              Completionist
            </Link>
          </div>
          <div className="flex items-center space-x-4">
            <Link href="/posts" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
              Feed
            </Link>
            <Link href="/search" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
              Search
            </Link>
            <Link href="/users" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
              Find Users
            </Link>
            {user ? (
              <>
                <Link href="/my-list" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
                  My List
                </Link>
                <Link href="/wishlist" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
                  Wishlist
                </Link>
                <Link href={`/profile/${user.username}`} className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
                  My Profile
                </Link>
                <Link href="/settings" className="text-gray-600 dark:text-gray-300 hover:text-gray-800">
                  Settings
                </Link>
                <button
                  onClick={logout}
                  className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="px-4 py-2 text-sm font-medium text-gray-700 rounded-md hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700">
                  Login
                </Link>
                <Link href="/register" className="ml-2 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </nav>
    </header>
  );
}