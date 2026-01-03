"use client";
import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { 
  HouseIcon, 
  ListChecksIcon, 
  CompassIcon, 
  TrophyIcon, 
  BooksIcon, 
  GearIcon
} from "@phosphor-icons/react/dist/ssr";

export default function Sidebar() {
  const { user } = useAuth();

  return (
    <aside className="fixed top-0 left-0 h-screen w-[260px] p-6 border-r border-border flex flex-col bg-background hidden lg:flex z-50">
      <div className="mb-10 text-2xl font-bold font-heading text-foreground">
        <span>Completionist</span>
      </div>

      <nav className="flex flex-col gap-2">
        <NavLink href="/" icon={<HouseIcon size={24} />} label="Home Feed" />
        <NavLink href="/my-list" icon={<ListChecksIcon size={24} />} label="My Lists" />
        <NavLink href="/search" icon={<CompassIcon size={24} />} label="Explore" />
        <NavLink href="/challenges" icon={<TrophyIcon size={24} />} label="Challenges" />
        <NavLink href="/wishlist" icon={<BooksIcon size={24} />} label="Wishlist" />
        <NavLink href="/settings" icon={<GearIcon size={24} />} label="Settings" />
      </nav>

      <div className="mt-auto pt-5 border-t border-border flex items-center gap-3">
        <div className="w-10 h-10 rounded-full border-2 border-primary overflow-hidden relative bg-gray-200 flex-shrink-0">
           <div className="w-full h-full flex items-center justify-center text-primary font-bold bg-indigo-100">
             {user?.username?.charAt(0).toUpperCase() || "G"}
           </div>
        </div>
        <div className="flex flex-col overflow-hidden">
          <h4 className="text-sm font-semibold text-foreground truncate">
            {user?.username || "Guest"}
          </h4>
          {user && (
            <span className="text-xs text-muted font-medium truncate">
              Joined {new Date(user.createdAt).getFullYear()}
            </span>
          )}
        </div>
      </div>
    </aside>
  );
}

function NavLink({ href, icon, label }: { href: string; icon: React.ReactNode; label: string }) {
  return (
    <Link
      href={href}
      className="flex items-center gap-3 px-4 py-3 rounded-xl text-muted font-medium hover:bg-indigo-50 hover:text-primary transition-colors dark:hover:bg-gray-800"
    >
      <div className="text-current">{icon}</div>
      <span>{label}</span>
    </Link>
  );
}