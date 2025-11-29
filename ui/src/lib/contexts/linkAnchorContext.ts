import { createContext, RefObject } from "react";
import { Side } from "../types/general";
import { Directory } from "../types/directory";
import { Playlist } from "../types/playlist";

export type LinkType = "directory" | "playlist"
export type LinkAnchorMap = Record<string, { el: HTMLElement | null, side: Side, directory?: Pick<Directory, "id">, playlist?: Pick<Playlist, "id"> }>;
export type LinkConnection = { from: string; to: string };

interface LinkAnchorContextType {
  registerAnchor: (id: string, anchor: { el: HTMLElement | null, side: Side, directory?: Pick<Directory, "id">, playlist?: Pick<Playlist, "id"> }) => void;
  startConnection: (id: string) => void;
  finishConnection: (id: string) => void;
  removeConnection: (from: string, to: string) => void;
  connections: LinkConnection[];
  draggingFrom: string | null;
  tempPos: { x: number; y: number } | null;
  notifyLayoutChange: () => void;
  anchorsRef: RefObject<LinkAnchorMap>;
  hoveredConnection: LinkConnection | null;
  setHoveredConnection: (connection: LinkConnection | null) => void;
  resetConnections: () => void;
  saveConnections: () => Promise<void>;
}

export const LinkAnchorContext = createContext<LinkAnchorContextType | undefined>(undefined);
