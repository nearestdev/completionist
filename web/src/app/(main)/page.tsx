import { ReactNode } from "react";
import Sidebar from "@/components/layout/Sidebar";
import Widgets from "@/components/layout/Widgets";
import Link from "next/link";
import { 
  HouseIcon, 
  ListChecksIcon, 
  CompassIcon, 
  TrophyIcon, 
  BooksIcon 
} from "@phosphor-icons/react/dist/ssr";

export default function MainLayout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col lg:flex-row">
      <Sidebar />

      <div className="flex-1 flex justify-center lg:ml-[260px]">
        
        <main className="w-full max-w-[680px] px-4 py-8 mb-20 lg:mb-0 min-h-screen">
          {children}
        </main>

        <Widgets />
      </div>

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
    </div>
  );
}