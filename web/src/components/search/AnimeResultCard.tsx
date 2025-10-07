import { JikanAnime } from "@/types/jikan";
import Image from "next/image";

interface AnimeResultCardProps {
  anime: JikanAnime;
}

export default function AnimeResultCard({ anime }: AnimeResultCardProps) {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div className="relative h-64">
        <Image
          src={anime.images.webp.image_url}
          alt={anime.title}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="font-bold text-lg truncate">{anime.title}</h3>
        <p className="text-sm text-gray-600 dark:text-gray-400">{anime.type}</p>
      </div>
    </div>
  );
}