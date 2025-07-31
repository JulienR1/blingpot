import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./Root";

function Index() {
  return (
    <>
      <p>log in or do something idk</p>
      <Link to="/new">go to new</Link>
    </>
  );
}

export const index = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Index,
});
