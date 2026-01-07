"use client";

import Widgets from "@/components/layout/Widgets";
import { 
  MagnifyingGlassIcon, 
  ListBulletsIcon, 
  TargetIcon,
  ChartLineIcon,
  TrendUpIcon,
  SparkleIcon,
  ArrowRightIcon,
  RocketLaunchIcon
} from "@phosphor-icons/react/dist/ssr";
import TrendingSection from "@/components/features/trending/TrendingSection";
import Link from "next/link";

const QuickActionCard = ({ 
  href, 
  icon: Icon, 
  title, 
  description, 
  gradient 
}: { 
  href: string; 
  icon: any; 
  title: string; 
  description: string; 
  gradient: string;
}) => (
  <Link 
    href={href}
    className="group relative overflow-hidden bg-card rounded-2xl border border-border p-6 shadow-sm hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
  >
    <div className="absolute top-0 right-0 -mt-8 -mr-8 w-32 h-32 bg-primary/5 rounded-full blur-2xl group-hover:bg-primary/10 transition-all duration-500"></div>
    
    <div className="relative z-10">
      <div className={`w-14 h-14 rounded-xl ${gradient} flex items-center justify-center text-white mb-4 shadow-lg group-hover:scale-110 transition-transform duration-300`}>
        <Icon size={28} weight="duotone" />
      </div>
      
      <h3 className="font-heading font-bold text-lg mb-2 text-foreground group-hover:text-primary transition-colors">
        {title}
      </h3>
      <p className="text-sm text-muted leading-relaxed mb-4">
        {description}
      </p>
      
      <div className="flex items-center gap-2 text-sm font-bold text-primary opacity-0 group-hover:opacity-100 transition-opacity">
        <span>Explore</span>
        <ArrowRightIcon size={16} weight="bold" className="group-hover:translate-x-1 transition-transform" />
      </div>
    </div>
  </Link>
);

export default function HomePage() {
  return (
    <div className="flex flex-col lg:flex-row justify-center w-full">
      <main className="w-full max-w-[1200px] px-4 py-8 mb-20 lg:mb-0 min-h-screen">
        
        <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 text-white shadow-2xl mb-12 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div className="absolute inset-0 bg-[url('/patterns/grid.svg')] opacity-10"></div>
          <div className="absolute top-0 right-0 p-12 opacity-10 transform rotate-12">
            <RocketLaunchIcon size={300} weight="duotone" />
          </div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-white/5 rounded-full blur-3xl"></div>
          <div className="absolute top-1/2 right-1/4 w-40 h-40 bg-pink-300/10 rounded-full blur-2xl"></div>
          
          <div className="relative z-10 p-8 md:p-12 lg:p-16">
            <div className="max-w-3xl">
              
              <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold font-heading tracking-tight mb-4 drop-shadow-lg">
                Welcome back, <span className="text-amber-300">Traveler</span>
              </h1>
              
              <p className="text-lg md:text-xl text-indigo-100 mb-8 leading-relaxed max-w-2xl">
                Track your progress across movies, games, books, and more. Complete challenges, 
                climb the leaderboards, and conquer your backlog one item at a time.
              </p>
              
            </div>
          </div>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12 animate-in fade-in slide-in-from-bottom-8 duration-700">
          <QuickActionCard
            href="/search"
            icon={MagnifyingGlassIcon}
            title="Discover Media"
            description="Browse movies, games, books, anime, and music"
            gradient="bg-gradient-to-br from-blue-500 to-cyan-500"
          />
          <QuickActionCard
            href="/my-list"
            icon={ListBulletsIcon}
            title="My Lists"
            description="Manage your personal collection and track progress"
            gradient="bg-gradient-to-br from-emerald-500 to-teal-500"
          />
          <QuickActionCard
            href="/challenges"
            icon={TargetIcon}
            title="Challenges"
            description="Complete challenges and earn rewards"
            gradient="bg-gradient-to-br from-amber-500 to-orange-500"
          />
          <QuickActionCard
            href="/profile"
            icon={ChartLineIcon}
            title="Your Stats"
            description="View your progress and achievements"
            gradient="bg-gradient-to-br from-purple-500 to-pink-500"
          />
        </div>

        <div className="animate-in fade-in slide-in-from-bottom-12 duration-1000">
          <TrendingSection />
        </div>
      </main>

      <Widgets />
    </div>
  );
}