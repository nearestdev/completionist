import { useEffect, useState } from "react";
import { mediaService } from "@/services/api";
import Image from "next/image";
import Link from "next/link";
import { FireIcon, ArrowRightIcon } from "@phosphor-icons/react";

interface MediaItem {
  id: string;
  title: string;
  coverImageUrl?: string;
  itemType: string;
}

export default function TrendingSection() {
  const [trending, setTrending] = useState<MediaItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTrending = async () => {
      try {
        const res = await mediaService.getTrendingMedia();
        if (res.data) {
          console.log("Trending media data:", res.data);
          setTrending(res.data);
        }
      } catch (e) {
        console.error("Failed to fetch trending", e);
      } finally {
        setLoading(false);
      }
    };
    fetchTrending();
  }, []);

  if (loading) return <div className="animate-pulse h-80 bg-card/50 rounded-2xl border border-border"></div>;
  if (!trending.length) return null;

  return (
    <section>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-gradient-to-br from-orange-500 to-red-500 rounded-xl shadow-lg">
            <FireIcon size={24} className="text-white" weight="fill" />
          </div>
          <div>
            <h2 className="font-heading font-bold text-2xl text-foreground">Trending Now</h2>
            <p className="text-sm text-muted">Popular media across all categories</p>
          </div>
        </div>
        <Link 
          href="/search" 
          className="flex items-center gap-2 text-sm font-bold text-primary hover:text-primary-hover transition-colors group"
        >
          <span>View All</span>
          <ArrowRightIcon size={16} weight="bold" className="group-hover:translate-x-1 transition-transform" />
        </Link>
      </div>
      
      <div className="relative">
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
          {trending.slice(0, 10).map((item) => (
            <Link
              key={item.id}
              href={`/media/${item.itemType}/${item.id}`}
              className="group relative bg-card rounded-xl overflow-hidden border border-border shadow-sm hover:shadow-xl transition-all duration-300 hover:-translate-y-1"
            >
              <div className="aspect-[2/3] relative overflow-hidden bg-muted">
                {item.coverImageUrl ? (
                  <Image 
                    src={item.coverImageUrl} 
                    alt={item.title} 
                    fill 
                    unoptimized
                    className="object-cover group-hover:scale-110 transition-transform duration-500"
                    onError={(e) => console.error("Image failed to load:", item.coverImageUrl, e)}
                  />
                ) : (
                  <div className="w-full h-full flex items-center justify-center text-muted-foreground text-xs">
                    No Image
                  </div>
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                  <div className="absolute bottom-0 left-0 right-0 p-4 translate-y-2 group-hover:translate-y-0 transition-transform duration-300">
                    <div className="flex items-center gap-2 text-white/90 text-xs mb-2">
                      <div className="w-1.5 h-1.5 rounded-full bg-orange-500"></div>
                      <span className="uppercase font-bold tracking-wide">{item.itemType}</span>
                    </div>
                  </div>
                </div>
              </div>
              <div className="p-3">
                <h3 className="font-bold text-sm truncate text-foreground group-hover:text-primary transition-colors" title={item.title}>
                  {item.title}
                </h3>
                <p className="text-xs text-muted capitalize mt-1">{item.itemType}</p>
              </div>
            </Link>
          ))}
        </div>
      </div>
    </section>
  );
}
