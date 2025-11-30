import z from "zod";
import { API } from "./api";
import { convertUser, User } from "./user";
import { JSONBody } from "./general";

export interface Playlist {
  id: number;
  spotifyId: string;
  owner?: User;
  name: string;
  description?: string;
  public: boolean;
  trackAmount: number;
  collaborative: boolean;
  hasCover: boolean;
}

export const convertPlaylist = (playlist: API.Playlist): Playlist => {
  return {
    id: playlist.id,
    spotifyId: playlist.spotify_id,
    owner: playlist.owner ? convertUser(playlist.owner) : undefined,
    name: playlist.name,
    description: playlist.description,
    public: playlist.public,
    trackAmount: playlist.track_amount,
    collaborative: playlist.collaborative,
    hasCover: playlist.has_cover,
  }
}

export const convertPlaylists = (playlists: API.Playlist[]): Playlist[] => {
  return playlists.map(convertPlaylist)
}

export const convertPlaylistsSchema = (playlists: Playlist[]): PlaylistSchema[] => {
  return playlists.map(p => ({
    id: p.id,
    name: p.name,
    trackAmount: p.trackAmount,
    hasCover: p.hasCover,
  }))
}

export const playlistSchema = z.object({
  id: z.number(),
  name: z.string(),
  trackAmount: z.number(),
  hasCover: z.boolean(),
})
export type PlaylistSchema = z.infer<typeof playlistSchema> & JSONBody;
