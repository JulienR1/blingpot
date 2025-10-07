import { request } from "@/lib/request";
import { useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";
import z from "zod";

const EXPENSES = "expenses";

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
    [q],
  );
};
