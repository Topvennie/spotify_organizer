import { Directory } from "@/lib/types/directory";
import { LinkNode } from "./LinkNode";
import { Side } from "@/lib/types/general";

type Props = {
  directory: Directory;
  side: Side;
}

export const LinkTree = ({ directory, side }: Props) => {
  return <LinkNode
    directory={directory}
    side={side}
    level={0}
  />
}


