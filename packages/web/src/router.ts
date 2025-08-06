import { createRouter } from "@tanstack/react-router";
import { root, home, newExpense } from "./pages";

const routeTree = root.addChildren([home, newExpense]);
export const router = createRouter({
  routeTree,
  context: { queryClient: undefined! },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}
