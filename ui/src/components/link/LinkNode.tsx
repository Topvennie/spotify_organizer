import { LinkConnection } from "@/lib/contexts/linkAnchorContext";
import { useLinkAnchor } from "@/lib/hooks/useLinkAnchor";
import { Directory } from "@/lib/types/directory";
import { Playlist } from "@/lib/types/playlist";
import { useMemo, useState } from "react";
import { FaRegFolder, FaRegFolderOpen } from "react-icons/fa6";
import { LinkPlaylist } from "./LinkPlaylist";
import { getLinkDirectoryId, getLinkPlaylistId } from "./util";
import { Side } from "@/lib/types/general";

type Props = {
  directory: Directory;
  side: Side;
  level: number;
}

const getPlaylists = (directory: Directory): Playlist[] => {
  return [
    ...directory.playlists,
    ...(directory.children?.flatMap(getPlaylists) ?? [])
  ]
}

const getDirectories = (directories: Directory[]): Directory[] => {
  return [...directories, ...(directories.flatMap(d => getDirectories(d?.children ?? [])))]
}

export const LinkNode = ({ directory, side, level }: Props) => {
  const [expanded, setExpanded] = useState(true)
  const { registerAnchor, startConnection, finishConnection, connections, notifyLayoutChange, hoveredConnection } = useLinkAnchor()

  const id = getLinkDirectoryId(directory, side)

  const isHoveredConnection = side === "left" ? hoveredConnection?.from === id : hoveredConnection?.to === id
  const [isHovered, setIsHovered] = useState(false)

  const handleExpandToggle = () => {
    setExpanded(prev => !prev)
    setTimeout(() => notifyLayoutChange(), 0) // With a timeout so the DOM can change first
  }

  const hiddenConnections = useMemo(() => {
    if (expanded) return 0

    const playlists = getPlaylists(directory)
    const directories = getDirectories(directory.children ?? [])

    const ids = [...playlists.map(p => getLinkPlaylistId(p, side)), ...directories.map(d => getLinkDirectoryId(d, side))]

    let cons: LinkConnection[]
    if (side === "left") cons = connections.filter(c => ids.includes(c.from))
    else cons = connections.filter(c => ids.includes(c.to))

    return cons.length
  }, [expanded, connections, directory, side])

  return (
    <div className="flex flex-col gap-1">
      <div
        onClick={handleExpandToggle}
        style={{ marginLeft: level * 16 }}
        onMouseDown={(e) => { e.stopPropagation(); startConnection(id) }}
        onMouseUp={(e) => { e.stopPropagation(); finishConnection(id) }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        className="relative flex items-center gap-2 rounded-md bg-gray-200 p-4 cursor-pointer hover:bg-gray-100"
      >
        {expanded
          ? <FaRegFolderOpen />
          : <FaRegFolder />
        }

        <span className="font-semibold">{directory.name}</span>
        <span className="text-muted">{(directory.children?.length ?? 0) + directory.playlists.length}</span>
        {hiddenConnections > 0 && <span className="ml-auto mr-2 text-red-500">{hiddenConnections}</span>}

        <div
          id={id}
          ref={el => registerAnchor(id, { el, side, directory })}
          className={`absolute ${side === "left" ? "-right-1" : "-left-1"}`}
        >
          <div
            className={`w-3 h-3 rounded-full bg-blue-500 cursor-pointer ${isHoveredConnection ? "bg-red-500" : isHovered ? "bg-blue-700" : ""}`}
          />
        </div>
      </div>

      {expanded && (
        <>
          <div style={{ marginLeft: (level + (side === "left" ? 1 : 2)) * 16 }} className="flex flex-col gap-1">
            {directory.playlists?.map(p => (
              <LinkPlaylist
                key={p.id}
                playlist={p}
                side={side}
              />
            ))}
          </div>

          <div className="flex flex-col gap-1">
            {directory.children?.map(child => (
              <LinkNode
                key={child.id}
                directory={child}
                side={side}
                level={level + 1}
              />
            ))}
          </div>
        </>
      )}
    </div>
  )
}
