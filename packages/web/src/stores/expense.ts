import { request } from "@/lib/request";
import { queryOptions, useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import type { Category } from "./category";
import type { Profile } from "./profile";
import z from "zod";

const EXPENSES = "expenses";

const ExpenseSchema = z.object({
  id: z.number(),
  spenderId: z.string(),
  label: z.string(),
  amount: z.number().transform((num) => Math.round(num / 100)),
  timestamp: z.number().transform((num) => new Date(num)),
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

type FetchParams = { queryKey: [string, { start: number; end: number }] };
export async function fetchExpenses({ queryKey }: FetchParams) {
  const [_, opts] = queryKey;
  const params = new URLSearchParams({
    start: opts.start.toString(),
    end: opts.end.toString(),
  });
  const expenses = await request(`/expenses?${params.toString()}`).get(
    z.array(ExpenseSchema)
  );
  return expenses ?? [];
}

export const expensesQuery = (start: Date, end: Date) =>
  queryOptions({
    queryKey: [EXPENSES, { start: start.getTime(), end: end.getTime() }],
    queryFn: fetchExpenses,
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
        timestamp: Math.floor(timestamp.getTime() / 1000),
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
