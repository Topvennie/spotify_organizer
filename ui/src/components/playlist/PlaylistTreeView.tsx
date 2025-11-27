import { useDirectoryGetAll } from "@/lib/api/directory";
import { Directory } from "@/lib/types/directory";
import { Group, RenderTreeNodePayload, Tree, TreeNodeData } from "@mantine/core";
import { useMemo } from "react";
import { FaRegCirclePlay, FaRegFolder, FaRegFolderOpen } from "react-icons/fa6";
import { LoadingSpinner } from "../molecules/LoadingSpinner";

interface FileIconProps {
  isFolder: boolean;
  expanded: boolean;
}

const FileIcon = ({ isFolder, expanded }: FileIconProps) => {
  if (!isFolder) {
    return <FaRegCirclePlay />
  }

  if (expanded) {
    return <FaRegFolderOpen />
  }

  return <FaRegFolder />
}

const Leaf = ({ node, expanded, hasChildren, elementProps }: RenderTreeNodePayload) => {
  return (
    <Group gap={5} {...elementProps}>
      <FileIcon isFolder={hasChildren} expanded={expanded} />
      <span>{node.label}</span>
    </Group>
  )
}

const toData = (directory: Directory): TreeNodeData => {
  const playlists = directory.playlists.map(p => ({ label: p.name, value: p.name }))

  return {
    label: directory.name,
    value: directory.name,
    children: [...playlists, ...(directory.children?.map(toData) ?? [])],
  }
}

export const PlaylistTreeView = () => {
  const { data: directories, isLoading } = useDirectoryGetAll()

  const data = useMemo(() => directories?.map(toData) ?? [], [directories])
  console.log(data)

  if (isLoading) return <LoadingSpinner />

  return (
    <Tree
      selectOnClick
      clearSelectionOnOutsideClick
      data={data}
      renderNode={payload => <Leaf {...payload} />}
    />
  )
}
