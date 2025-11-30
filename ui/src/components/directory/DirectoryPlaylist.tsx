import { PlaylistSchema } from "@/lib/types/playlist";
import { cn } from "@/lib/utils";
import { ActionIcon } from "@mantine/core";
import { ComponentProps } from "react";
import { FaTrashCan } from "react-icons/fa6";
import { PlaylistCover } from "../playlist/PlaylistCover";

type Props = {
  playlist: PlaylistSchema;
  onDelete: (playlist: PlaylistSchema) => void;
} & ComponentProps<"div">

export const DirectoryPlaylist = ({ playlist, onDelete, className, ...props }: Props) => {
  return (
    <div className={cn("flex justify-between rounded-md bg-white p-4", className)} {...props}>
      <div className="flex items-center gap-2">
        <PlaylistCover playlist={playlist} />
        <span className="whitespace-nowrap">{playlist.name}</span>
        <span className="text-muted text-sm">{playlist.trackAmount}</span>
      </div>
      <ActionIcon onClick={() => onDelete(playlist)} color="red" variant="subtle">
        <FaTrashCan />
      </ActionIcon>
    </div>
  )
}
