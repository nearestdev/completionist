"use client";

import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { useSidebar } from "@/contexts/SidebarContext";
import { 
  HouseIcon, 
  ListChecksIcon, 
  CompassIcon, 
  TrophyIcon, 
  BooksIcon, 
  GearIcon,
  CaretLeftIcon,
  ChatCircleDotsIcon,
  UsersThreeIcon
} from "@phosphor-icons/react/dist/ssr";

export default function Sidebar() {
  const { user } = useAuth();
  const { isCollapsed, toggleSidebar } = useSidebar();

  const sidebarWidth = isCollapsed ? "w-[80px]" : "w-[260px]";

  return (
    <>
      {/* Desktop Sidebar */}
      <aside 
        className={`
          fixed top-0 left-0 h-screen ${sidebarWidth} 
          border-r border-border bg-card transition-all duration-300
          hidden lg:flex flex-col z-50
        `}
        style={{
          boxShadow: 'var(--shadow-md)'
        }}
      >
        <div className={`flex items-center ${isCollapsed ? "justify-center" : "gap-4"} px-6 h-[80px] border-b border-border`}>
          {!isCollapsed && (
            <div className="text-2xl font-bold font-heading text-foreground transition-opacity duration-300 flex items-center">
              <span>Completionist</span>
            </div>
          )}
          <button 
            onClick={toggleSidebar}
            className={`
              p-2 rounded-full hover:bg-primary/10 text-muted hover:text-primary transition-all
              ${!isCollapsed ? "ml-auto" : ""}
            `}
            aria-label={isCollapsed ? "Expand Sidebar" : "Collapse Sidebar"}
          >
            <CaretLeftIcon 
              size={20} 
              weight="bold" 
              className={`transition-transform duration-300 ${isCollapsed ? "rotate-180" : ""}`} 
            />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex flex-col gap-2 px-4">
          <NavLink href="/" icon={<HouseIcon size={24} />} label="Home Feed" collapsed={isCollapsed} />
          <NavLink href="/my-list" icon={<ListChecksIcon size={24} />} label="My Lists" collapsed={isCollapsed} />
          <NavLink href="/search" icon={<CompassIcon size={24} />} label="Explore" collapsed={isCollapsed} />
          <NavLink href="/messages" icon={<ChatCircleDotsIcon size={24} />} label="Messages" collapsed={isCollapsed} />
          <NavLink href="/rooms" icon={<UsersThreeIcon size={24} />} label="Rooms" collapsed={isCollapsed} />
          <NavLink href="/challenges" icon={<TrophyIcon size={24} />} label="Challenges" collapsed={isCollapsed} />
          <NavLink href="/wishlist" icon={<BooksIcon size={24} />} label="Wishlist" collapsed={isCollapsed} />
          <NavLink href="/settings" icon={<GearIcon size={24} />} label="Settings" collapsed={isCollapsed} />
        </nav>

        <div 
          className="mt-auto p-4 flex flex-col gap-4 border-t border-border" 
          style={{
            background: 'linear-gradient(to bottom, transparent, rgba(99, 102, 241, 0.03))'
          }}
        >
          
          
          <div className={`flex items-center ${isCollapsed ? "justify-center" : "gap-3"}`}>
            <div className="w-10 h-10 rounded-full border-2 border-primary overflow-hidden relative bg-gradient-to-br from-primary/20 to-accent/20 flex-shrink-0">
               <div className="w-full h-full flex items-center justify-center text-primary font-bold bg-primary/10">
                 {user?.username?.charAt(0).toUpperCase() || "G"}
               </div>
            </div>
            
            {!isCollapsed && (
              <div className="flex flex-col overflow-hidden transition-all duration-300">
                <h4 className="text-sm font-semibold text-foreground truncate">
                  {user?.username || "Guest"}
                </h4>
                {user && (
                  <span className="text-xs text-accent font-medium truncate">
                    Joined {new Date(user.createdAt).getFullYear()}
                  </span>
                )}
              </div>
            )}
          </div>
        </div>
      </aside>

      {/* Mobile Navigation */}
      <nav className="lg:hidden fixed bottom-0 left-0 w-full bg-card border-t border-border flex justify-around py-3 px-2 z-50 safe-area-bottom">
        <Link href="/" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
          <HouseIcon size={24} />
        </Link>
        <Link href="/my-list" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
          <ListChecksIcon size={24} />
        </Link>
        <Link href="/search" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
          <CompassIcon size={24} />
        </Link>
        <Link href="/challenges" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
          <TrophyIcon size={24} />
        </Link>
        <Link href="/wishlist" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
          <BooksIcon size={24} />
        </Link>
      </nav>
    </>
  );
}

function NavLink({ href, icon, label, collapsed }: { href: string; icon: React.ReactNode; label: string; collapsed: boolean }) {
  return (
    <Link
      href={href}
      className={`
        flex items-center gap-3 px-3 py-3 rounded-xl 
        text-muted font-medium 
        hover:bg-primary/10 hover:text-primary 
        active:bg-primary/20
        transition-all duration-200
        ${collapsed ? "justify-center" : ""}
      `}
      title={collapsed ? label : undefined}
    >
      <div className={`text-current flex-shrink-0 ${collapsed ? "w-6 h-6" : ""}`}>{icon}</div>
      {!collapsed && (
        <span className="whitespace-nowrap overflow-hidden transition-all duration-300">
          {label}
        </span>
      )}
    </Link>
  );
}