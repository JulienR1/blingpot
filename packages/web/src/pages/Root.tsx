import { useSuspenseQuery, type QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { query } from "@/stores/profile";

export const root = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  component: Root,
  loader: ({ context }) => context.queryClient.ensureQueryData(query),
});

function Root() {
  const { data: profile } = useSuspenseQuery(query);

  if (profile == null) {
    return <a href="/oauth2/authenticate">connect!</a>;
  }
  return <Outlet />;
}
