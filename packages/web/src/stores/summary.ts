import { request } from "@/lib/request";
import { Money } from "@/lib/schemas";
import { queryOptions } from "@tanstack/react-query";
import z from "zod";
import { EXPENSES } from "./expense";

const SUMMARY = "summary";

export const ExpenseSummarySchema = z.object({
  total: Money,
  categories: z.record(z.string(), Money),
});

export type ExpenseSummary = z.infer<typeof ExpenseSummarySchema>;

async function fetchExpenseSummary(start: number, end: number) {
  const params = new URLSearchParams({
    start: start.toString(),
    end: end.toString(),
  });
  return request(`/summary/expenses?${params.toString()}`).get(
    ExpenseSummarySchema
  );
}

export const expensesSummaryQuery = (start: Date, end: Date) =>
  queryOptions({
    queryKey: [
      EXPENSES,
      SUMMARY,
      { start: start.getTime(), end: end.getTime() },
    ] as const,
    queryFn: ({ queryKey }) => {
      const [_, __, { start, end }] = queryKey;
      return fetchExpenseSummary(start, end);
    },
  });
