import { dict } from "@/lib/utils";
import type { Category } from "@/stores/category";
import type { Expense } from "@/stores/expense";
import type { Profile } from "@/stores/profile";
import { useMemo } from "react";

type Stores = {
  expenses: Array<Expense>;
  profiles: Array<Profile>;
  categories: Array<Category>;
};

export function useExpenses(stores: Stores) {
  const profiles = useMemo(
    () => dict(stores.profiles, "sub"),
    [stores.profiles]
  );
  const categories = useMemo(
    () => dict(stores.categories, "id"),
    [stores.categories]
  );

  return useMemo(
    () =>
      stores.expenses.map((expense) => ({
        id: expense.id,
        label: expense.label,
        amount: expense.amount,
        timestamp: expense.timestamp,
        category: categories[expense.categoryId],
        spender: profiles[expense.spenderId],
        author: profiles[expense.authorId],
      })),
    [stores.expenses, profiles, categories]
  );
}
