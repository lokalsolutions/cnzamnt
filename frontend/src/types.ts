export interface User {
  id: number;
  display_name: string;
  handle: string;
  cnz_balance: number;
  created_at: string;
}

export interface Comment {
  id: number;
  artwork_id: number;
  user_id: number;
  rating_cnz: number;
  body: string;
  created_at: string;
  author?: User;
  artist_cnz: number;
}

export interface Artwork {
  id: number;
  user_id: number;
  title: string;
  image_url: string;
  caption: string;
  created_at: string;
  artist?: User;
  comments?: Comment[];
  comment_count?: number;
}

export interface CreateArtworkInput {
  title: string;
  image_url: string;
  caption: string;
}
