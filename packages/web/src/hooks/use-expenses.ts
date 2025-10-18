import { dict } from "@/lib/utils";
import { categoriesQuery } from "@/stores/category";
import { expensesQuery } from "@/stores/expense";
import { profilesQuery } from "@/stores/profile";
import { useSuspenseQueries } from "@tanstack/react-query";
import { useMemo } from "react";

type ExpensesParams = {
  start: Date;
  end: Date;
};

export function useExpenses({ start, end }: ExpensesParams) {
  const [expensesResult, profilesResult, categoriesResult] = useSuspenseQueries(
    { queries: [expensesQuery(start, end), profilesQuery, categoriesQuery] },
  );

  const profiles = useMemo(
    () => dict(profilesResult.data, "sub"),
    [profilesResult.data],
  );
  const categories = useMemo(
    () => dict(categoriesResult.data, "id"),
    [categoriesResult.data],
  );

  return useMemo(
    () =>
      expensesResult.data.map((expense) => ({
        id: expense.id,
        label: expense.label,
        amount: expense.amount,
        timestamp: expense.timestamp,
        category: categories[expense.categoryId],
        spender: profiles[expense.spenderId],
        author: profiles[expense.authorId],
      })),
    [expensesResult.data, profiles, categories],
  );
}
