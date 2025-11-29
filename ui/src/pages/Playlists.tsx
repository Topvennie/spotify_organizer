import { PlaylistTableView } from "@/components/playlist/PlaylistTableView"
import { PlaylistTreeView } from "@/components/playlist/PlaylistTreeView"
import { Switch, Title } from "@mantine/core"
import { useState } from "react"

type View = "table" | "tree"

const storageKey = "music-playlist-view"

export const Playlists = () => {
  const [view, setView] = useState<View>(localStorage.getItem(storageKey) as View ?? "table")

  const handleCheckToggle = () => {
    let newView: View

    if (view === "table") newView = "tree"
    else newView = "table"

    localStorage.setItem(storageKey, newView)
    setView(newView)
  }

  return (
    <div className="flex flex-col gap-8">
      <Title order={1} className="text-center">Playlists</Title>
      <div className="self-end">
        <Switch
          checked={view == "tree"}
          onChange={handleCheckToggle}
          label="Tree view"
        />
      </div>
      {view === "tree"
        ? <PlaylistTreeView />
        : <PlaylistTableView />
      }
    </div>
  )
}

