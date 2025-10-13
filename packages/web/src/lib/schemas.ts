import z from "zod";

export const Money = z.number().transform((num) => Math.round(num / 100));
