import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./root";
import { useSuspenseQueries } from "@tanstack/react-query";
import { expensesQuery } from "@/stores/expense";
import { profilesQuery } from "@/stores/profile";
import { categoriesQuery } from "@/stores/category";
import { useExpenses } from "@/hooks/use-expenses";

function Home() {
  const queries = useSuspenseQueries({
    queries: [
      expensesQuery(new Date(2000, 0, 1), new Date(2100, 0, 1)),
      profilesQuery,
      categoriesQuery,
    ],
  });

  const expenses = useExpenses({
    expenses: queries[0].data,
    profiles: queries[1].data,
    categories: queries[2].data,
  });

  return (
    <>
      <p>log in or do something idk</p>
      <Link to="/new">go to new</Link>

      <pre>{JSON.stringify(expenses, null, 2)}</pre>
    </>
  );
}

export const home = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Home,
});
