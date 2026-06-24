"use client";
import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";
import ThemeSwitcher from "@/components/common/ThemeSwitcher";
import StreakPill from "@/components/streaks/StreakPill";
import MemberBadge from "@/components/subscription/MemberBadge";

export default function Header() {
  const { user, logout } = useAuth();
  
  return (
    <header 
      className="bg-card border-b border-border sticky top-0 z-40 backdrop-blur-sm bg-card/95"
      style={{ boxShadow: 'var(--shadow-sm)' }}
    >
      <nav className="container mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-6">
          </div>
          
          <div className="flex items-center gap-3">
            <ThemeSwitcher />
            
            {user ? (
              <>
                <StreakPill userId={user.id} />
                {(user.role === "member" || user.role === "admin") && <MemberBadge />}
                <Link
                  href={`/profile/${user.username}`}
                  className="hidden sm:flex px-4 py-2 rounded-full text-sm font-medium text-muted hover:text-foreground hover:bg-primary/10 transition-all"
                >
                  Profile
                </Link>
                <button
                  onClick={logout}
                  className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-red-500 to-red-600 rounded-full hover:from-red-600 hover:to-red-700 transition-all shadow-sm hover:shadow-md"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link 
                  href="/login" 
                  className="px-4 py-2 text-sm font-medium text-muted hover:text-foreground rounded-full hover:bg-primary/10 transition-all"
                >
                  Login
                </Link>
                <Link 
                  href="/register" 
                  className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-primary to-primary-hover rounded-full hover:opacity-90 transition-all shadow-sm hover:shadow-md"
                >
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
