import { JikanManga } from "@/types/jikan";
import Image from "next/image";

interface MangaResultCardProps {
  manga: JikanManga;
}

export default function MangaResultCard({ manga }: MangaResultCardProps) {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
      <div className="relative h-64">
        <Image
          src={manga.images.webp.image_url}
          alt={manga.title}
          fill
          sizes="(max-width: 640px) 100vw, (max-width: 768px) 50vw, (max-width: 1024px) 33vw, 25vw"
          className="object-cover"
        />
      </div>
      <div className="p-4">
        <h3 className="font-bold text-lg truncate">{manga.title}</h3>
      </div>
    </div>
  );
}