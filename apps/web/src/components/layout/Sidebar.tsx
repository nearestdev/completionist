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
  UsersThreeIcon,
  ShieldCheckIcon
} from "@phosphor-icons/react/dist/ssr";

export default function Sidebar() {
  const { user, canAccess } = useAuth();
  const { isCollapsed, toggleSidebar } = useSidebar();

  const sidebarWidth = isCollapsed ? "w-[80px]" : "w-[260px]";

  const desktopLinks = [
    { href: "/", icon: <HouseIcon size={24} />, label: "Home Feed" },
    { href: "/my-list", icon: <ListChecksIcon size={24} />, label: "My Lists" },
    { href: "/search", icon: <CompassIcon size={24} />, label: "Explore" },
    { href: "/messages", icon: <ChatCircleDotsIcon size={24} />, label: "Messages" },
    { href: "/rooms", icon: <UsersThreeIcon size={24} />, label: "Rooms" },
    { href: "/challenges", icon: <TrophyIcon size={24} />, label: "Challenges" },
    { href: "/wishlist", icon: <BooksIcon size={24} />, label: "Wishlist" },
    { href: "/settings", icon: <GearIcon size={24} />, label: "Settings" },
  ].filter((link) => canAccess(link.href));

  const mobileLinks = [
    { href: "/", icon: <HouseIcon size={24} /> },
    { href: "/my-list", icon: <ListChecksIcon size={24} /> },
    { href: "/search", icon: <CompassIcon size={24} /> },
    { href: "/challenges", icon: <TrophyIcon size={24} /> },
    { href: "/wishlist", icon: <BooksIcon size={24} /> },
  ].filter((link) => canAccess(link.href));

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
          {desktopLinks.map((link) => (
            <NavLink
              key={link.href}
              href={link.href}
              icon={link.icon}
              label={link.label}
              collapsed={isCollapsed}
            />
          ))}
          {user?.role === "admin" && (
            <NavLink href="/admin" icon={<ShieldCheckIcon size={24} />} label="Admin" collapsed={isCollapsed} />
          )}
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
        {mobileLinks.map((link) => (
          <Link
            key={link.href}
            href={link.href}
            className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center"
          >
            {link.icon}
          </Link>
        ))}
        {user?.role === "admin" && (
          <Link href="/admin" className="p-2 text-muted hover:text-primary transition-colors flex flex-col items-center">
            <ShieldCheckIcon size={24} />
          </Link>
        )}
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
