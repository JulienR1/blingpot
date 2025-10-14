import { request } from "@/lib/request";
import { queryOptions, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import type { Category } from "./category";
import type { Profile } from "./profile";
import z from "zod";
import { Money, Timestamp } from "@/lib/schemas";

export const EXPENSES = "expenses";

const ExpenseSchema = z.object({
  id: z.number(),
  spenderId: z.string(),
  label: z.string(),
  amount: Money,
  timestamp: Timestamp,
  authorId: z.string(),
  categoryId: z.number(),
});

export type ExpenseCore = z.infer<typeof ExpenseSchema>;
export type Expense = Omit<
  ExpenseCore,
  "spenderId" | "categoryId" | "authorId"
> & {
  category: Category;
  spender: Profile;
  author: Profile;
};

export async function fetchExpenses(start: number, end: number) {
  const params = new URLSearchParams({
    start: start.toString(),
    end: end.toString(),
  });
  const expenses = await request(`/expenses?${params.toString()}`).get(
    z.array(ExpenseSchema)
  );
  return expenses ?? [];
}

export const expensesQuery = (start: Date, end: Date) =>
  queryOptions({
    queryKey: [
      EXPENSES,
      { start: start.getTime(), end: end.getTime() },
    ] as const,
    queryFn: ({ queryKey }) => {
      const { start, end } = queryKey[1];
      return fetchExpenses(start, end);
    },
  });

const CreateResponse = z.object({ id: z.number() });

type CreatePayload = {
  label: string;
  amount: string;
  timestamp: Date;
  spenderId: string;
  categoryId: number | null;
};

export const useCreate = () => {
  const q = useQueryClient();
  return useCallback(
    async ({
      label,
      amount,
      timestamp,
      spenderId,
      categoryId,
    }: CreatePayload) => {
      const body = {
        label,
        spenderId,
        amount: Math.floor(100 * parseFloat(amount)),
        timestamp: timestamp.getTime(),
        categoryId,
      };

      const id = await request("/expenses").post(CreateResponse, { body });
      if (id != null) {
        q.invalidateQueries({ queryKey: [EXPENSES, id] });
      }
    },
    [q]
  );
};
