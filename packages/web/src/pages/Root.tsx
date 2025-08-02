import { useSuspenseQuery, type QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { profileQuery } from "@/stores/profile";

export const root = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  component: Root,
  loader: ({ context }) => context.queryClient.ensureQueryData(profileQuery),
});

function Root() {
  const { data: profile } = useSuspenseQuery(profileQuery);
  const authenticated = profile != null;

  return authenticated ? (
    <Outlet />
  ) : (
    <a href="/oauth2/authenticate">connect!</a>
  );
}
