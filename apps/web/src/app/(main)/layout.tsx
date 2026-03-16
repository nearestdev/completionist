"use client";

import { ReactNode, useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";
import { SidebarProvider, useSidebar } from "@/contexts/SidebarContext";
import Sidebar from "@/components/layout/Sidebar";
import Footer from "@/components/layout/Footer";
import Header from "@/components/layout/Header";
import { useAuth } from "@/hooks/useAuth";

function LayoutContent({ children }: { children: ReactNode }) {
  const { isCollapsed } = useSidebar();
  const { loading, canAccess } = useAuth();
  const pathname = usePathname();
  const router = useRouter();
  
  const isFullScreenPage = /^\/rooms\/\d+$/.test(pathname) || pathname === '/messages';
  const allowedForCurrentRole = canAccess(pathname);

  useEffect(() => {
    if (loading) {
      return;
    }

    if (allowedForCurrentRole) {
      return;
    }

    const next = encodeURIComponent(pathname);
    router.replace(`/login?next=${next}`);
  }, [loading, allowedForCurrentRole, pathname, router]);

  if (loading || !allowedForCurrentRole) {
    return null;
  }
  
  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground transition-colors duration-300">
      <div className={`${isCollapsed ? "lg:ml-[80px]" : "lg:ml-[260px]"} transition-all duration-300 pointer-events-none fixed top-0 w-full z-50`}> 
         <div className="pointer-events-auto w-full">
            <Header />
         </div>
      </div>
      {/* Spacer for fixed header */}
      <div className="h-[73px]" /> 

      <div className="flex flex-1">
        <Sidebar />
        
        <main 
          className={`
            flex-grow flex flex-col 
            ${isCollapsed ? "lg:ml-[80px]" : "lg:ml-[260px]"}
            transition-all duration-300
            w-full
            ${isFullScreenPage ? 'h-[calc(100vh-73px)] overflow-hidden' : ''}
          `}
        >
          {children}
          {!isFullScreenPage && <Footer />}
        </main>
      </div>
    </div>
  );
}

export default function MainLayout({ children }: { children: ReactNode }) {
  return (
    <SidebarProvider>
      <LayoutContent>{children}</LayoutContent>
    </SidebarProvider>
  );
}
