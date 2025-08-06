import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./root";

function Home() {
  return (
    <>
      <p>log in or do something idk</p>
      <Link to="/new">go to new</Link>
    </>
  );
}

export const home = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Home,
});
