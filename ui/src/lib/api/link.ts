import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { apiGet, apiPost } from "./query"
import { convertLinks, LinkSchema } from "../types/link"
import { STALE_TIME } from "../types/staletime"

const ENDPOINT = "link"

export const useLinkGetAll = () => {
  return useQuery({
    queryKey: ["link"],
    queryFn: async () => (await apiGet(ENDPOINT, convertLinks)).data,
    staleTime: STALE_TIME.MIN_30,
    throwOnError: true,
  })
}

export const useLinkSync = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (links: LinkSchema[]) => apiPost(`${ENDPOINT}/sync`, links),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["link"] }),
    throwOnError: true,
  })
}
