import { request } from "@/lib/request";
import { queryOptions, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import z from "zod";

const USER_PROFILE = "user-profile";

const ProfileSchema = z.object({
  firstName: z.string(),
  lastName: z.string(),
  email: z.email(),
  picture: z.string().nullable(),
});

export type Profile = z.infer<typeof ProfileSchema>;

const fetchUserProfile = () => request("/profiles/me").get(ProfileSchema);
const disconnect = () => request("/oauth2/revoke", "none").post();

export const profileQuery = queryOptions({
  queryKey: [USER_PROFILE],
  queryFn: fetchUserProfile,
});

export const useDisconnect = () => {
  const client = useQueryClient();

  return useCallback(async () => {
    await disconnect();
    client.invalidateQueries({ queryKey: [USER_PROFILE] });
  }, [client]);
};
