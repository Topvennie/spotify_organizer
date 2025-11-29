import { getLinkDirectoryId, getLinkPlaylistId } from "@/components/link/util";
import { LoadingSpinner } from "@/components/molecules/LoadingSpinner";
import { notifications } from "@mantine/notifications";
import { PropsWithChildren, useCallback, useEffect, useMemo, useRef, useState } from "react";
import z from "zod";
import { useLinkGetAll, useLinkSync } from "../api/link";
import { LinkAnchorContext, LinkAnchorMap, LinkConnection } from "../contexts/linkAnchorContext";
import { Directory } from "../types/directory";
import { Side } from "../types/general";
import { Link, linkSchema, LinkSchema } from "../types/link";
import { Playlist } from "../types/playlist";
import { debounce } from "../utils";

const linkToConnection = (link: Link): LinkConnection => {
  let from = ""
  let to = ""

  if (link.sourceDirectoryId) from = getLinkDirectoryId({ id: link.sourceDirectoryId }, "left")
  else if (link.sourcePlaylistId) from = getLinkPlaylistId({ id: link.sourcePlaylistId }, "left")

  if (link.targetDirectoryId) to = getLinkDirectoryId({ id: link.targetDirectoryId }, "right")
  else if (link.targetPlaylistId) to = getLinkPlaylistId({ id: link.targetPlaylistId }, "right")

  return { from, to }
}

export const LinkAnchorProvider = ({ children }: PropsWithChildren) => {
  const { data: links, isLoading } = useLinkGetAll()
  const saveLinks = useLinkSync()

  const anchorsRef = useRef<LinkAnchorMap>({})
  const allAnchorsRef = useRef<LinkAnchorMap>({})
  const [connections, setConnections] = useState<LinkConnection[]>([])
  const [draggingFrom, setDraggingFrom] = useState<string | null>(null)
  const [tempPos, setTempPos] = useState<{ x: number; y: number } | null>(null)
  const [, setLayoutVersion] = useState(0)
  const [hoveredConnection, setHoveredConnection] = useState<LinkConnection | null>(null)

  const observers = useRef<Record<string, ResizeObserver | null>>({})

  const debouncedLayoutChange = useRef(
    debounce(() => setLayoutVersion(v => v + 1), 40)
  ).current

  const notifyLayoutChange = useCallback(() => {
    debouncedLayoutChange()
  }, [debouncedLayoutChange])

  const registerAnchor = useCallback((id: string, anchor: { el: HTMLElement | null, side: Side, directory?: Pick<Directory, "id">, playlist?: Pick<Playlist, "id"> }) => {
    if (anchor.el) {
      anchorsRef.current[id] = anchor

      if (!observers.current[id]) {
        const observer = new ResizeObserver(() => notifyLayoutChange())
        observer.observe(anchor.el)
        observers.current[id] = observer
      }

      if (!(id in allAnchorsRef.current)) allAnchorsRef.current[id] = anchor
    } else {
      delete anchorsRef.current[id]

      observers.current[id]?.disconnect()
      delete observers.current[id]
    }

    notifyLayoutChange()
  }, [notifyLayoutChange])

  const startConnection = useCallback((id: string) => {
    setDraggingFrom(id);
    document.body.style.userSelect = "none"
  }, [])

  const addConnection = useCallback((from: string, to: string) => {
    setConnections(prev => {
      if (prev.some(c => c.from === from && c.to === to)) {
        return prev
      }
      return [...prev, { from, to }]
    })
  }, [])

  const finishConnection = useCallback((id: string) => {
    if (draggingFrom) {
      const from = anchorsRef.current[draggingFrom]
      const to = anchorsRef.current[id]

      // Both elements need to exist
      // And a directory cannot point to itself (on the other side)
      // And a playlist cannot point to itself (on the other side)
      if (from && to &&
        !(from.directory && from.directory === to.directory) &&
        !(from.playlist && from.playlist === to.playlist)
      ) {
        if (from.side === "left" && to.side === "right") addConnection(draggingFrom, id)
        else if (from.side === "right" && to.side === "left") addConnection(id, draggingFrom)
      }
    }

    setDraggingFrom(null)
    setTempPos(null)
    document.body.style.userSelect = "auto"
  }, [draggingFrom, addConnection])

  const removeConnection = useCallback((from: string, to: string) => {
    setConnections(prev => prev.filter(p => p.from !== from || p.to !== to))
  }, [])

  useEffect(() => {
    const onMove = (e: MouseEvent) => {
      if (!draggingFrom) return

      setTempPos({ x: e.clientX, y: e.clientY })
    };

    const onUp = () => {
      if (!draggingFrom) return

      setDraggingFrom(null)
      setTempPos(null)
    };

    window.addEventListener("mousemove", onMove)
    window.addEventListener("mouseup", onUp)

    return () => {
      window.removeEventListener("mousemove", onMove)
      window.removeEventListener("mouseup", onUp)
    }
  }, [draggingFrom])

  useEffect(() => {
    window.addEventListener("resize", notifyLayoutChange)
    window.addEventListener("scroll", notifyLayoutChange)

    return () => {
      window.removeEventListener("resize", notifyLayoutChange)
      window.removeEventListener("scroll", notifyLayoutChange)
    }
  }, [notifyLayoutChange])

  const resetConnections = useCallback(() => setConnections(links?.map(linkToConnection) ?? []), [links])

  const saveConnections = useCallback(async () => {
    const links: LinkSchema[] = connections.map(c => {
      const from = allAnchorsRef.current[c.from]
      const to = allAnchorsRef.current[c.to]

      return {
        sourceDirectoryId: from.directory?.id,
        sourcePlaylistId: from.playlist?.id,
        targetDirectoryId: to.directory?.id,
        targetPlaylistId: to.playlist?.id,
      }
    })

    const linkSchemas = z.array(linkSchema)
    const result = linkSchemas.safeParse(links)
    if (!result.success) {
      notifications.show({ variant: "error", title: "Validation error", message: result.error.message })
      return
    }

    await saveLinks.mutateAsync(result.data, {
      onSuccess: () => notifications.show({ variant: "success", message: "Links synced" })
    })
  }, [connections, saveLinks])

  useEffect(() => {
    if (!links) return

    setConnections(links.map(linkToConnection))
  }, [links])

  const value = useMemo(() => ({
    registerAnchor,
    startConnection,
    finishConnection,
    removeConnection,
    connections,
    draggingFrom,
    tempPos,
    notifyLayoutChange,
    anchorsRef,
    hoveredConnection,
    setHoveredConnection,
    resetConnections,
    saveConnections,
  }), [registerAnchor, startConnection, finishConnection, removeConnection, connections, draggingFrom, tempPos, notifyLayoutChange, anchorsRef, hoveredConnection, setHoveredConnection, resetConnections, saveConnections])

  if (isLoading) return <LoadingSpinner />

  return <LinkAnchorContext value={value}>{children}</LinkAnchorContext>
}
