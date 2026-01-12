"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import api from "@/services/api";
import listService from "@/services/listService";
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
  HeartIcon,
  CheckCircleIcon
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
  const [addingToList, setAddingToList] = useState(false);
  const [addingToWishlist, setAddingToWishlist] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [showPriorityModal, setShowPriorityModal] = useState(false);

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

  const handleAddToList = async () => {
    if (!media) return;

    try {
      setAddingToList(true);
      setErrorMessage(null);
      setSuccessMessage(null);

      await listService.createListItem({
        mediaData: {
          itemType: media.itemType,
          source: media.source,
          externalId: media.externalId,
          title: media.title,
          description: media.description,
          coverImageUrl: media.coverImageUrl,
          releaseDate: media.releaseDate,
          genres: media.genres,
          metadata: media.metadata,
        },
        listData: {
          status: "planning",
        },
      });

      setSuccessMessage("Successfully added to your list!");
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (err: any) {
      console.error("Failed to add to list:", err);
      setErrorMessage(
        err.response?.data?.error || "Failed to add to list. It may already be in your list."
      );
      setTimeout(() => setErrorMessage(null), 5000);
    } finally {
      setAddingToList(false);
    }
  };

  const handleAddToWishlist = async (priority: "low" | "medium" | "high") => {
    if (!media) return;

    try {
      setAddingToWishlist(true);
      setErrorMessage(null);
      setSuccessMessage(null);
      setShowPriorityModal(false);

      await listService.addToWishlist({
        mediaItemId: media.id,
        priority,
      });

      setSuccessMessage("Successfully added to your wishlist!");
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (err: any) {
      console.error("Failed to add to wishlist:", err);
      setErrorMessage(
        err.response?.data?.error || "Failed to add to wishlist. It may already be in your wishlist."
      );
      setTimeout(() => setErrorMessage(null), 5000);
    } finally {
      setAddingToWishlist(false);
    }
  };

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

          {(successMessage || errorMessage) && (
            <div className={`flex items-center gap-2 px-4 py-3 rounded-xl border animate-in fade-in slide-in-from-top-2 ${
              successMessage 
                ? "bg-green-500/10 border-green-500/30 text-green-600" 
                : "bg-red-500/10 border-red-500/30 text-red-600"
            }`}>
              {successMessage && <CheckCircleIcon size={20} weight="fill" />}
              <span className="text-sm font-medium">{successMessage || errorMessage}</span>
            </div>
          )}

          <div className="flex flex-wrap gap-3 mt-auto">
            <Button 
              className="gap-2 shadow-lg shadow-primary/20" 
              size="lg"
              onClick={handleAddToList}
              disabled={addingToList || addingToWishlist}
            >
              {addingToList ? (
                <>
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Adding...
                </>
              ) : (
                <>
                  <PlusIcon size={20} weight="bold" />
                  Add to My List
                </>
              )}
            </Button>
            <Button 
              variant="outline" 
              className="gap-2" 
              size="lg"
              onClick={() => setShowPriorityModal(true)}
              disabled={addingToList || addingToWishlist}
            >
              {addingToWishlist ? (
                <>
                  <div className="w-5 h-5 border-2 border-primary/30 border-t-primary rounded-full animate-spin" />
                  Adding...
                </>
              ) : (
                <>
                  <HeartIcon size={20} weight="duotone" />
                  Add to Wishlist
                </>
              )}
            </Button>
          </div>

          {showPriorityModal && (
            <div 
              className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-in fade-in duration-200"
              onClick={() => setShowPriorityModal(false)}
            >
              <div 
                className="bg-card border border-border rounded-2xl p-6 max-w-sm w-full shadow-2xl animate-in zoom-in-95 duration-200"
                onClick={(e) => e.stopPropagation()}
              >
                <h3 className="text-xl font-bold text-foreground mb-2">Set Priority</h3>
                <p className="text-sm text-muted mb-6">How important is this item to you?</p>
                
                <div className="flex flex-col gap-3">
                  <button
                    onClick={() => handleAddToWishlist("high")}
                    disabled={addingToWishlist}
                    className="w-full px-4 py-3 rounded-xl font-medium transition-all border-2 bg-red-500/10 text-red-600 border-red-500/30 hover:bg-red-500/20 hover:border-red-500/50 disabled:opacity-50"
                  >
                    🔥 High Priority
                  </button>
                  
                  <button
                    onClick={() => handleAddToWishlist("medium")}
                    disabled={addingToWishlist}
                    className="w-full px-4 py-3 rounded-xl font-medium transition-all border-2 bg-orange-500/10 text-orange-600 border-orange-500/30 hover:bg-orange-500/20 hover:border-orange-500/50 disabled:opacity-50"
                  >
                    ⭐ Medium Priority
                  </button>
                  
                  <button
                    onClick={() => handleAddToWishlist("low")}
                    disabled={addingToWishlist}
                    className="w-full px-4 py-3 rounded-xl font-medium transition-all border-2 bg-blue-500/10 text-blue-600 border-blue-500/30 hover:bg-blue-500/20 hover:border-blue-500/50 disabled:opacity-50"
                  >
                    💡 Low Priority
                  </button>
                </div>

                <button
                  onClick={() => setShowPriorityModal(false)}
                  className="w-full mt-4 px-4 py-2 text-sm text-muted hover:text-foreground transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          )}

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
