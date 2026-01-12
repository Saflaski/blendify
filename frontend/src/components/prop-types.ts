import { z } from "zod";

export type ControlPanelProps = {
  blendid: string;
  setMode: (val: string) => void;
  setUsers: (val: string[]) => void;
  setUserATopItemsData: (val: CatalogueTopItemsResponse) => void;
  setUserBTopItemsData: (val: CatalogueTopItemsResponse) => void;
  setBlendPercent: (num: number) => void;
  userATopItemApiResponse: CatalogueTopItemsResponse;
  userBTopItemApiResponse: CatalogueTopItemsResponse;
  blendApiResponse: CardApiResponse;
  downloadTopItems: (
    duration: string,
    category: string,
    username: string,
    blendId: string,
    setData: (val: CatalogueTopItemsResponse) => void,
  ) => Promise<void>;
};

// export type BlendApiResponse = {
//   usernames: string[];
//   overallBlendNum: number;
//   ArtistBlend: TypeBlend;
//   AlbumBlend: TypeBlend;
//   TrackBlend: TypeBlend;
// };
// export type MetricKey = keyof BlendApiResponse;

// export type TypeBlend = {
//   OneMonth: number;
//   ThreeMonth: number;
//   OneYear: number;
// };

const TypeBlendSchema = z.object({
  OneMonth: z.number(),
  ThreeMonth: z.number(),
  OneYear: z.number(),
});

export const CardApiResponseSchema = z.object({
  Usernames: z.array(z.string()),
  OverallBlendNum: z.number(),
  ArtistBlend: TypeBlendSchema,
  AlbumBlend: TypeBlendSchema,
  TrackBlend: TypeBlendSchema,
});
export type CardApiResponse = z.infer<typeof CardApiResponseSchema>;

export const CatalogueBlendSchema = z.object({
  Name: z.string(),
  ImageUrl: z.url(),
  Artist: z.string().optional(),
  ArtistUrl: z.url().optional(),
  ArtistImageUrl: z.url().optional(),
  EntryUrl: z.url().optional(),
  Playcounts: z.array(z.number()),
});

export type CatalogueBlendResponse = z.infer<typeof CatalogueBlendSchema>;

export const CatalogueTopItemsSchema = z.object({
  Items: z.array(z.string()),
});

export type CatalogueTopItemsResponse = z.infer<typeof CatalogueTopItemsSchema>;

// `json:"Name"`
// URL            string `json:"EntryUrl,omitempty"`
// ImageURL       string `json:"ImageUrl,omitempty"`
// ArtistName     string `json:"Artist,omitempty"` //Not needed when it's an artist
// ArtistURL      string `json:"ArtistUrl,omitempty"`
// ArtistImageURL string `json:"ArtistImageUrl,omitempty"`
// Playcounts     []int  `json:"Playcounts"`
