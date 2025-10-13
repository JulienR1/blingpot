import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./root";
import { usePrefetchQuery } from "@tanstack/react-query";
import { expensesQuery } from "@/stores/expense";
import { profilesQuery } from "@/stores/profile";
import { categoriesQuery } from "@/stores/category";
import { ShowcaseTable } from "@/components/showcase-table";

function Home() {
  usePrefetchQuery(profilesQuery);
  usePrefetchQuery(categoriesQuery);
  usePrefetchQuery(expensesQuery(new Date(2000, 0, 1), new Date(2100, 0, 1)));

  return (
    <>
      <p>log in or do something idk</p>
      <Link to="/new">go to new</Link>

      <ShowcaseTable />
    </>
  );
}

export const home = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Home,
});
