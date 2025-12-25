import { z } from "zod";

export type ControlPanelProps = {
  setMode: (val: string) => void;
  setUsers: (val: string[]) => void;
  setBlendPercent: (num: number) => void;
  blendApiResponse: BlendApiResponse;
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

export const BlendApiResponseSchema = z.object({
  Usernames: z.array(z.string()),
  OverallBlendNum: z.number(),
  ArtistBlend: TypeBlendSchema,
  AlbumBlend: TypeBlendSchema,
  TrackBlend: TypeBlendSchema,
});
export type BlendApiResponse = z.infer<typeof BlendApiResponseSchema>;
