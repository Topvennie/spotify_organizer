import { useContext } from "react"
import { LinkAnchorContext } from "../contexts/linkAnchorContext";

export const useLinkAnchor = () => {
  const context = useContext(LinkAnchorContext)
  if (!context) {
    throw new Error("useLinkAnchor must be used within an LinkAnchorProvider")
  }

  return context;
}
