import { TMDBSearchItem } from "@/types/tmdb";
import Image from "next/image";

interface MovieResultCardProps {
  movie: TMDBSearchItem;
}

export default function MovieResultCard({ movie }: MovieResultCardProps) {
  const imageUrl = movie.poster_path
    ? `https://image.tmdb.org/t/p/w500${movie.poster_path}`
    : "/placeholder.svg";

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div className="relative h-64">
        <Image
          src={imageUrl}
          alt={movie.title || "Movie poster"}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="font-bold text-lg truncate">{movie.title}</h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">{movie.release_date?.substring(0, 4)}</p>
      </div>
    </div>
  );
}