import { createRouter } from "@tanstack/react-router";
import { root, index, newTransaction } from "./pages";

const routeTree = root.addChildren([index, newTransaction]);
export const router = createRouter({
  routeTree,
  context: { queryClient: undefined! },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
