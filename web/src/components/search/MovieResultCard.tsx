import { TMDBSearchItem } from "@/types/tmdb";
import Image from "next/image";
import { mediaService } from "@/services/api";
import { useRouter } from "next/navigation";

interface MovieResultCardProps {
  movie: TMDBSearchItem;
}

export default function MovieResultCard({ movie }: MovieResultCardProps) {
  const router = useRouter();
  const imageUrl = movie.poster_path
    ? `https://image.tmdb.org/t/p/w500${movie.poster_path}`
    : "/placeholder.svg";

  const handleClick = async () => {
    try {
      const response = await mediaService.getMediaDetails("TMDB", String(movie.id), "movie");
      console.log("Media details fetched:", response.data);
      router.push(`/media/movie/${response.data.id}`);
    } catch (e: any) {
      console.error("Failed to fetch media details:", {
        message: e.message,
        response: e.response?.data,
        status: e.response?.status,
        error: e
      });
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
          alt={movie.title || "Movie poster"}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover transition-transform duration-300 group-hover:scale-105"
        />
      </div>
      <div className="p-4">
        <span className="text-[10px] uppercase tracking-wider font-bold text-primary bg-primary/10 px-2 py-0.5 rounded mb-2 inline-block">
          Movie
        </span>
        <h3 className="font-heading font-bold text-foreground text-lg truncate mb-1">{movie.title}</h3>
        <p className="text-sm text-muted">
          {movie.release_date?.substring(0, 4)}
        </p>
      </div>
    </div>
  );
}