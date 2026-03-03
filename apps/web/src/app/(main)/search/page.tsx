"use client";

import Link from "next/link";
import { 
  TelevisionIcon, 
  FilmStripIcon, 
  GameControllerIcon, 
  BookIcon, 
  BookOpenIcon, 
  MagnifyingGlassIcon,
  MonitorPlayIcon
} from "@phosphor-icons/react/dist/ssr";

export default function SearchHubPage() {
  const searchCategories = [
    { 
      name: "Movies", 
      href: "/search/movies", 
      icon: FilmStripIcon, 
      color: "from-red-500 to-rose-600",
      description: "Discover trending films and cinematic masterpieces"
    },
    { 
      name: "TV Shows", 
      href: "/search/tv", 
      icon: TelevisionIcon, 
      color: "from-blue-500 to-indigo-600",
       description: "Binge-worthy series and episodes"
    },
    { 
      name: "Anime", 
      href: "/search/anime", 
      icon: MonitorPlayIcon, 
      color: "from-pink-500 to-purple-600",
      description: "Japanese animation and seasonal hits"
    },
    { 
      name: "Manga", 
      href: "/search/manga", 
      icon: BookOpenIcon, 
      color: "from-slate-400 to-slate-500",
      description: "Graphic novels and serialized comics"
    },
    { 
      name: "Games", 
      href: "/search/games", 
      icon: GameControllerIcon, 
      color: "from-emerald-500 to-teal-600",
      description: "Top-rated titles across all platforms"
    },
    { 
      name: "Books", 
      href: "/search/books", 
      icon: BookIcon, 
      color: "from-amber-500 to-orange-600",
      description: "Literary classics and bestsellers"
    },
  ];

  return (
    <div className="container mx-auto px-4 py-12 max-w-6xl animate-in fade-in duration-500">
      <div className="text-center mb-16 space-y-4">
        <div className="inline-flex items-center justify-center p-4 bg-primary/10 rounded-full mb-4 text-primary animate-bounce-slow">
           <MagnifyingGlassIcon size={32} weight="bold" />
        </div>
        <h1 className="text-4xl md:text-5xl font-heading font-bold text-foreground tracking-tight">
          Explore the Collection
        </h1>
        <p className="text-xl text-muted max-w-2xl mx-auto">
          Dive into a vast database of entertainment. Select a category to start your journey.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {searchCategories.map((category) => (
          <Link
            key={category.name}
            href={category.href}
            className="group relative overflow-hidden bg-card rounded-2xl border border-border shadow-md hover:shadow-xl transition-all duration-300 hover:-translate-y-1"
          >
            <div className={`absolute top-0 right-0 w-32 h-32 bg-gradient-to-br ${category.color} opacity-5 rounded-full blur-2xl -mr-16 -mt-16 group-hover:opacity-10 transition-opacity duration-500`} />
            
            <div className="p-8">
              <div className={`w-14 h-14 rounded-xl bg-gradient-to-br ${category.color} flex items-center justify-center text-white shadow-lg mb-6 group-hover:scale-110 transition-transform duration-300`}>
                <category.icon size={28} weight="fill" />
              </div>
              
              <h2 className="text-2xl font-heading font-bold text-foreground mb-2 group-hover:text-primary transition-colors">
                {category.name}
              </h2>
              <p className="text-muted text-sm leading-relaxed">
                {category.description}
              </p>
            </div>
            
            <div className="absolute bottom-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-primary/50 to-transparent transform scale-x-0 group-hover:scale-x-100 transition-transform duration-500" />
          </Link>
        ))}
      </div>
    </div>
  );
}