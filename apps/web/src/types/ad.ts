export type AdPlacement = "feed" | "sidebar" | "banner" | "profile";

export interface AdCampaign {
  id: number;
  name: string;
  advertiserName: string;
  contentHtml: string;
  imageUrl?: string;
  targetUrl: string;
  placement: AdPlacement;
  status: string;
  createdAt: string;
}
