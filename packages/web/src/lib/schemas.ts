import z from "zod";

export const Money = z.number().transform((num) => Math.round(num / 100));
export const Timestamp = z.number().transform((num) => new Date(num));
