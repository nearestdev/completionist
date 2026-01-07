"use client";

import { ReactNode } from "react";
import { SidebarProvider, useSidebar } from "@/contexts/SidebarContext";
import Sidebar from "@/components/layout/Sidebar";
import Footer from "@/components/layout/Footer";
import Header from "@/components/layout/Header";

function LayoutContent({ children }: { children: ReactNode }) {
  const { isCollapsed } = useSidebar();
  
  return (
    <div className="flex flex-col min-h-screen bg-background text-foreground transition-colors duration-300">
      <div className={`${isCollapsed ? "lg:ml-[80px]" : "lg:ml-[260px]"} transition-all duration-300`}>
         <Header />
      </div>

      <div className="flex flex-1">
        <Sidebar />
        
        <main 
          className={`
            flex-grow flex flex-col 
            ${isCollapsed ? "lg:ml-[80px]" : "lg:ml-[260px]"}
            transition-all duration-300
            w-full
          `}
        >
          <div className="container mx-auto px-4 py-6 flex-grow">
            {children}
          </div>
          <Footer />
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