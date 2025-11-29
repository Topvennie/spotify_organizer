import { Directory } from "@/lib/types/directory";
import { Side } from "@/lib/types/general";
import { Playlist } from "@/lib/types/playlist";

export const getLinkDirectoryId = (directory: Pick<Directory, "id">, side: Side) => {
  return `directory-${side}-${directory.id}`
}

export const getLinkPlaylistId = (playlist: Pick<Playlist, "id">, side: Side) => {
  return `playlist-${side}-${playlist.id}`
}
