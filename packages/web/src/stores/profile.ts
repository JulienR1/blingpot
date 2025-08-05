import { request } from "@/lib/request";
import { queryOptions, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import z from "zod";

const PROFILES = "profiles";

const ProfileSchema = z.object({
  sub: z.string(),
  firstName: z.string(),
  lastName: z.string(),
  email: z.email(),
  picture: z.string().nullable(),
});

export type Profile = z.infer<typeof ProfileSchema>;

const fetchAllProfiles = () =>
  request("/profiles")
    .get(z.array(ProfileSchema))
    .then((profiles) => profiles ?? []);
const fetchUserProfile = () => request("/profiles/me").get(ProfileSchema);
const disconnect = () => request("/oauth2/revoke", "none").post();

export const profilesQuery = queryOptions({
  queryKey: [PROFILES],
  queryFn: fetchAllProfiles,
});

export const profileQuery = queryOptions({
  queryKey: [PROFILES, "me"],
  queryFn: fetchUserProfile,
});

export const useDisconnect = () => {
  const client = useQueryClient();

  return useCallback(async () => {
    await disconnect();
    client.invalidateQueries({ queryKey: [PROFILES, "me"] });
  }, [client]);
};
