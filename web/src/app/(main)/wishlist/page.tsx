"use client";

import { useEffect, useState } from "react";
import listService from "@/services/listService";
import { WishlistItem } from "@/types/list";
import Link from "next/link";
import { 
  HeartIcon, 
  TrashIcon, 
  TagIcon, 
  SortAscendingIcon,
  ShoppingCartIcon
} from "@phosphor-icons/react/dist/ssr";
import Button from "@/components/ui/Button";

export default function WishlistPage() {
  const [items, setItems] = useState<WishlistItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchWishlistItems = async () => {
      try {
        setLoading(true);
        const data = await listService.getMyWishlist();
        setItems(data);
      } catch (err) {
        setError("Failed to load your wishlist.");
      } finally {
        setLoading(false);
      }
    };
    fetchWishlistItems();
  }, []);

  const handleRemove = async (id: number) => {
    try {
      await listService.removeFromWishlist(id);
      setItems((prev) => prev.filter((item) => item.id !== id));
    } catch (error) {
      console.error("Failed to remove item", error);
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority.toLowerCase()) {
      case "high": return "bg-red-500/10 text-red-600 border-red-200 dark:border-red-900";
      case "medium": return "bg-orange-500/10 text-orange-600 border-orange-200 dark:border-orange-900";
      case "low": return "bg-blue-500/10 text-blue-600 border-blue-200 dark:border-blue-900";
      default: return "bg-gray-500/10 text-gray-600 border-gray-200";
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center min-h-[60vh]">
      <div className="w-12 h-12 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>
  );

  if (error) return (
     <div className="flex items-center justify-center min-h-[40vh] text-red-500 font-medium">
      {error}
    </div>
  );

  return (
    <div className="container mx-auto px-4 py-8 max-w-5xl animate-in fade-in duration-500">
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="p-3 bg-pink-500/10 rounded-xl text-pink-600">
            <HeartIcon size={32} weight="fill" />
          </div>
          <div>
            <h1 className="text-3xl font-bold font-heading">My Wishlist</h1>
            <p className="text-muted">{items.length} items saved for later</p>
          </div>
        </div>
        
        {items.length > 0 && (
           <Button variant="outline" className="gap-2 hidden md:flex">
             <SortAscendingIcon size={18} />
             Sort by Priority
           </Button>
        )}
      </div>

      {items.length > 0 ? (
        <div className="grid gap-4">
          {items.map((item) => (
            <div
              key={item.id}
              className="group bg-card border border-border hover:border-primary/50 rounded-xl p-5 transition-all hover:shadow-md flex flex-col md:flex-row items-start md:items-center justify-between gap-4"
            >
              <div className="flex items-start gap-4">
                 <div className="w-16 h-16 bg-muted/30 rounded-lg flex items-center justify-center text-muted-foreground">
                    <TagIcon size={24} weight="duotone" />
                 </div>
                 
                 <div>
                   <h2 className="text-lg font-bold text-foreground mb-1 group-hover:text-primary transition-colors">
                     Media Item #{item.mediaItemId.substring(0, 8)}...
                   </h2>
                   <div className="flex items-center gap-2">
                     <span className={`text-xs font-bold uppercase tracking-wider px-2 py-0.5 rounded-full border ${getPriorityColor(item.priority)}`}>
                       {item.priority}
                     </span>
                     {item.priceCents && (
                       <span className="text-sm font-mono text-muted">
                         ${(item.priceCents / 100).toFixed(2)}
                       </span>
                     )}
                   </div>
                 </div>
              </div>

              <div className="flex items-center gap-2 w-full md:w-auto mt-2 md:mt-0">
                <Button className="flex-1 md:flex-none gap-2" size="sm">
                  <ShoppingCartIcon size={16} weight="bold" />
                  View Details
                </Button>
                <button 
                  onClick={() => handleRemove(item.id)}
                  className="p-2 text-muted-foreground hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/10 rounded-lg transition-colors"
                  title="Remove from wishlist"
                >
                  <TrashIcon size={20} />
                </button>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-20 bg-muted/5 rounded-3xl border border-dashed border-border">
          <HeartIcon size={64} className="mx-auto text-muted mb-4 opacity-20" weight="duotone" />
          <h2 className="text-2xl font-bold text-foreground mb-2">Your wishlist is empty</h2>
          <p className="text-muted mb-6 max-w-sm mx-auto">
            Found something you like? Add it to your wishlist to keep track of it.
          </p>
          <Link href="/search">
            <Button size="lg" className="shadow-lg shadow-primary/20">
              Explore Media
            </Button>
          </Link>
        </div>
      )}
    </div>
  );
}