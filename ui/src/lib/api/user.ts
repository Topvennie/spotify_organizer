import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { convertUserToModel } from "../types/user";
import { apiGet, apiPost } from "./query";

const ENDPOINT_AUTH = "auth"
const ENDPOINT_USER = "user"

const STALE_30_MIN = 30 * 60 * 1000;

export const useUser = () => {
  return useQuery({
    queryKey: ["user"],
    queryFn: async () => (await apiGet(`${ENDPOINT_USER}/me`, convertUserToModel)).data,
    retry: 0,
    staleTime: STALE_30_MIN,
  })
}

export const useUserLogin = () => {
  window.location.href = `/api/${ENDPOINT_AUTH}/login/spotify`
}

export const useUserLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => (await apiPost(`${ENDPOINT_AUTH}/logout`)).data,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["user"] })
  })
}

