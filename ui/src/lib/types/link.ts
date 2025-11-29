import { z } from "zod";
import { API } from "./api";
import { JSONBody } from "./general";

export interface Link {
  id: number;
  sourceDirectoryId?: number;
  sourcePlaylistId?: number;
  targetDirectoryId?: number;
  targetPlaylistId?: number;
}

export const convertLink = (l: API.Link): Link => {
  return {
    id: l.id,
    sourceDirectoryId: l.source_directory_id,
    sourcePlaylistId: l.source_playlist_id,
    targetDirectoryId: l.target_directory_id,
    targetPlaylistId: l.target_playlist_id,
  }
}

export const convertLinks = (l: API.Link[]): Link[] => {
  return l.map(convertLink)
}

export const convertLinkSchema = (links: Link[]): LinkSchema[] => {
  return links.map(l => ({
    id: l.id,
    sourceDirectoryId: l.sourceDirectoryId,
    sourcePlaylistId: l.sourcePlaylistId,
    targetDirectoryId: l.targetDirectoryId,
    targetPlaylistId: l.targetPlaylistId,
  }))
}

export const linkSchema = z.object({
  id: z.number().optional(),
  sourceDirectoryId: z.number().optional(),
  sourcePlaylistId: z.number().optional(),
  targetDirectoryId: z.number().optional(),
  targetPlaylistId: z.number().optional(),
}).superRefine((data, ctx) => {
  const { sourceDirectoryId, sourcePlaylistId, targetDirectoryId, targetPlaylistId } = data;

  const sourceDirSet = sourceDirectoryId !== undefined;
  const sourcePlSet = sourcePlaylistId !== undefined;

  if (sourceDirSet === sourcePlSet) {
    ctx.addIssue({
      code: "custom",
      message: "Exactly one of sourceDirectoryId or sourcePlaylistId must be provided",
      path: ["sourceDirectoryId"],
    });
    ctx.addIssue({
      code: "custom",
      message: "Exactly one of sourceDirectoryId or sourcePlaylistId must be provided",
      path: ["sourcePlaylistId"],
    });
  }

  const targetDirSet = targetDirectoryId !== undefined;
  const targetPlSet = targetPlaylistId !== undefined;

  if (targetDirSet === targetPlSet) {
    ctx.addIssue({
      code: "custom",
      message: "Exactly one of targetDirectoryId or targetPlaylistId must be provided",
      path: ["targetDirectoryId"],
    });
    ctx.addIssue({
      code: "custom",
      message: "Exactly one of targetDirectoryId or targetPlaylistId must be provided",
      path: ["targetPlaylistId"],
    });
  }
});
export type LinkSchema = z.infer<typeof linkSchema> & JSONBody
