import { createRoute, Link } from "@tanstack/react-router";
import { root } from "./root";
import { useSuspenseQueries } from "@tanstack/react-query";
import { expensesQuery } from "@/stores/expense";
import { profilesQuery, type Profile } from "@/stores/profile";
import { categoriesQuery, type Category } from "@/stores/category";
import { useMemo } from "react";

function Home() {
  const [expenses, profiles, categories] = useSuspenseQueries({
    queries: [
      expensesQuery(new Date(2000, 0, 1), new Date(2100, 0, 1)),
      profilesQuery,
      categoriesQuery,
    ],
  });

  const profilesMap = useMemo(
    () =>
      profiles.data.reduce(
        (all, p) => ({
          ...all,
          [p.sub]: p,
        }),
        {} as Record<string, Profile>
      ),
    [profiles.data]
  );

  const categoriesMap = useMemo(
    () =>
      categories.data.reduce(
        (all, c) => ({
          ...all,
          [c.id]: c,
        }),
        {} as Record<string, Category>
      ),
    [categories.data]
  );

  const exp = useMemo(
    () =>
      expenses.data.map((e) => ({
        id: e.id,
        label: e.label,
        amount: e.amount,
        timestamp: e.timestamp,
        category: categoriesMap[e.categoryId],
        spender: profilesMap[e.spenderId],
        author: profilesMap[e.authorId],
      })),
    [profilesMap, categoriesMap, expenses.data]
  );

  return (
    <>
      <p>log in or do something idk</p>
      <Link to="/new">go to new</Link>

      <pre>{JSON.stringify(exp, null, 2)}</pre>
    </>
  );
}

export const home = createRoute({
  getParentRoute: () => root,
  path: "/",
  component: Home,
});
