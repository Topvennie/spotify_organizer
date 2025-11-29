import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor";
import { Side } from "@/lib/types/general";
import { Playlist } from "@/lib/types/playlist";
import { useState } from "react";
import { PlaylistCover } from "../playlist/PlaylistCover";
import { getLinkPlaylistId } from "./util";

type Props = {
  playlist: Playlist;
  side: Side;
}

export const LinkPlaylist = ({ playlist, side }: Props) => {
  const { registerAnchor, startConnection, finishConnection, hoveredConnection } = useLinkAnchor()

  const id = getLinkPlaylistId(playlist, side)

  const isHoveredConnection = side === "left" ? hoveredConnection?.from === id : hoveredConnection?.to === id
  const [isHovered, setIsHovered] = useState(false)

  return (
    <div
      onMouseDown={(e) => { e.stopPropagation(); startConnection(id) }}
      onMouseUp={(e) => { e.stopPropagation(); finishConnection(id) }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className="relative flex items-center justify-self-center gap-2 cursor-pointer"
    >
      <div className="w-8">
        <PlaylistCover playlist={playlist} />
      </div>
      <span className="whitespace-nowrap truncate">{playlist.name}</span>
      <span className="text-muted text-sm">{playlist.tracks}</span>
      <div
        id={id}
        ref={el => registerAnchor(id, { el, side, playlist })}
        className={`absolute ${side === "left" ? "-right-1" : "-left-5"}`}
      >
        <div
          className={`w-3 h-3 rounded-full bg-blue-500 cursor-pointer ${isHoveredConnection ? "bg-red-500" : isHovered ? "bg-blue-700" : ""}`}
        />
      </div>
    </div>
  )
}
