import { TMDBSearchItem } from "@/types/tmdb";
import Image from "next/image";
import { mediaService } from "@/services/api";
import { useRouter } from "next/navigation";

interface TvResultCardProps {
  show: TMDBSearchItem;
}

export default function TvResultCard({ show }: TvResultCardProps) {
  const router = useRouter();
  const imageUrl = show.poster_path
    ? `https://image.tmdb.org/t/p/w500${show.poster_path}`
    : "/placeholder.svg";

  const handleClick = async () => {
    try {
      const response = await mediaService.getMediaDetails("TMDB", String(show.id), "series");
      router.push(`/media/series/${response.data.id}`);
    } catch (e) {
      console.error("Failed to fetch media details", e);
    }
  };

  return (
    <div 
      className="bg-card rounded-xl border border-border overflow-hidden transition-all duration-300 hover:shadow-lg group cursor-pointer"
      style={{ boxShadow: 'var(--shadow-sm)' }}
      onClick={handleClick}
    >
      <div className="relative h-64 overflow-hidden">
        <Image
          src={imageUrl}
          alt={show.name || "TV show poster"}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </div>
      <div className="p-4">
        <span className="text-[10px] uppercase tracking-wider font-bold text-primary bg-primary/10 px-2 py-0.5 rounded mb-2 inline-block">
          TV Show
        </span>
        <h3 className="font-heading font-bold text-foreground text-lg truncate mb-1">{show.name}</h3>
        <p className="text-sm text-muted">
          {show.first_air_date?.substring(0, 4)}
        </p>
      </div>
    </div>
  );
}