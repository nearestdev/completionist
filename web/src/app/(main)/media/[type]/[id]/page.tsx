"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import api from "@/services/api";
import { MediaItem } from "@/types/list";
import Image from "next/image";
import Link from "next/link";
import {
  ArrowLeftIcon,
  FilmStripIcon,
  CalendarIcon,
  StarIcon,
  TagIcon,
  PlusIcon,
  HeartIcon
} from "@phosphor-icons/react/dist/ssr";
import Button from "@/components/ui/Button";

export default function MediaDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const mediaType = params.type as string;
  const mediaId = params.id as string;

  const [media, setMedia] = useState<MediaItem | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMediaDetails = async () => {
      try {
        setLoading(true);
        const response = await api.get<MediaItem>(`/media/${mediaId}`);
        setMedia(response.data);
      } catch (err: any) {
        console.error("Failed to fetch media details:", err);
        setError(err.response?.data?.error || "Failed to load media details");
      } finally {
        setLoading(false);
      }
    };

    if (mediaId) {
      fetchMediaDetails();
    }
  }, [mediaId]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[70vh]">
        <div className="w-16 h-16 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
      </div>
    );
  }

  if (error || !media) {
    return (
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        <div className="text-center py-16 bg-card border border-border rounded-3xl">
          <div className="w-20 h-20 bg-red-500/10 rounded-full flex items-center justify-center mx-auto mb-6">
            <FilmStripIcon size={40} className="text-red-500" weight="duotone" />
          </div>
          <h2 className="text-2xl font-bold text-foreground mb-2">Media Not Found</h2>
          <p className="text-muted mb-8">{error || "The requested media could not be found."}</p>
          <Link href="/search">
            <Button>Back to Search</Button>
          </Link>
        </div>
      </div>
    );
  }

  const imageUrl = media.coverImageUrl || "/placeholder.svg";
  const releaseYear = media.releaseDate ? new Date(media.releaseDate).getFullYear() : null;

  const getItemTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      movie: "Movie",
      series: "TV Series",
      anime: "Anime",
      manga: "Manga",
      game: "Game",
      book: "Book",
      music: "Music"
    };
    return labels[type] || type;
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl animate-in fade-in duration-500">
      <button
        onClick={() => router.back()}
        className="flex items-center gap-2 text-muted hover:text-foreground transition-colors mb-6 group"
      >
        <ArrowLeftIcon size={20} className="group-hover:-translate-x-1 transition-transform" />
        <span>Back</span>
      </button>

      <div className="grid md:grid-cols-[300px_1fr] gap-8">
        <div className="relative aspect-[2/3] md:aspect-auto md:h-[450px] rounded-2xl overflow-hidden shadow-2xl group">
          <Image
            src={imageUrl}
            alt={media.title}
            fill
            sizes="(max-width: 768px) 100vw, 300px"
            className="object-cover group-hover:scale-105 transition-transform duration-500"
            priority
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>

        <div className="flex flex-col gap-6">
          <div>
            <div className="flex items-center gap-3 mb-3">
              <span className="px-3 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-primary/10 text-primary border border-primary/20">
                {getItemTypeLabel(media.itemType)}
              </span>
              {releaseYear && (
                <span className="flex items-center gap-1.5 text-sm text-muted">
                  <CalendarIcon size={16} weight="duotone" />
                  {releaseYear}
                </span>
              )}
            </div>

            <h1 className="text-4xl font-bold font-heading mb-4 bg-gradient-to-br from-foreground to-foreground/70 bg-clip-text">
              {media.title}
            </h1>

            {media.description && (
              <div 
                className="text-muted leading-relaxed mb-6 prose prose-sm max-w-none
                  prose-p:my-2 prose-b:text-foreground prose-b:font-semibold 
                  prose-i:text-muted-foreground prose-br:my-1"
                dangerouslySetInnerHTML={{ __html: media.description }}
              />
            )}

            {media.genres && media.genres.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-6">
                <TagIcon size={20} className="text-muted" weight="duotone" />
                {media.genres.map((genre, index) => (
                  <span
                    key={index}
                    className="px-3 py-1 rounded-full text-xs font-medium bg-muted/10 text-muted-foreground border border-border hover:border-primary/50 transition-colors"
                  >
                    {genre}
                  </span>
                ))}
              </div>
            )}

            <div className="flex flex-wrap gap-6 mb-8">
              {media.criticRatingValue && (
                <div className="flex flex-col">
                  <span className="text-xs text-muted mb-1">Critic Score</span>
                  <div className="flex items-center gap-1.5">
                    <StarIcon size={20} weight="fill" className="text-amber-500" />
                    <span className="text-2xl font-bold text-foreground">{media.criticRatingValue.toFixed(1)}</span>
                    {media.criticRatingCount && (
                      <span className="text-sm text-muted">({media.criticRatingCount.toLocaleString()})</span>
                    )}
                  </div>
                </div>
              )}

              {media.userRatingExternal && (
                <div className="flex flex-col">
                  <span className="text-xs text-muted mb-1">User Score</span>
                  <div className="flex items-center gap-1.5">
                    <StarIcon size={20} weight="fill" className="text-blue-500" />
                    <span className="text-2xl font-bold text-foreground">{media.userRatingExternal.toFixed(1)}</span>
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="flex flex-wrap gap-3 mt-auto">
            <Button className="gap-2 shadow-lg shadow-primary/20" size="lg">
              <PlusIcon size={20} weight="bold" />
              Add to My List
            </Button>
            <Button variant="outline" className="gap-2" size="lg">
              <HeartIcon size={20} weight="duotone" />
              Add to Wishlist
            </Button>
          </div>

          <div className="pt-6 border-t border-border">
            <div className="grid grid-cols-2 gap-4 text-sm">
              {media.source && (
                <div>
                  <span className="text-muted">Source</span>
                  <p className="font-medium text-foreground">{media.source}</p>
                </div>
              )}
              {media.externalId && (
                <div>
                  <span className="text-muted">External ID</span>
                  <p className="font-medium text-foreground">{media.externalId}</p>
                </div>
              )}
              <div>
                <span className="text-muted">Added</span>
                <p className="font-medium text-foreground">{new Date(media.createdAt).toLocaleDateString()}</p>
              </div>
              <div>
                <span className="text-muted">Updated</span>
                <p className="font-medium text-foreground">{new Date(media.updatedAt).toLocaleDateString()}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
