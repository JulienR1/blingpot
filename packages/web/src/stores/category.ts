import { request } from "@/lib/request";
import { queryOptions } from "@tanstack/react-query";
import z from "zod";

const CATEGORIES = "categories";

const CategorySchema = z.object({
  id: z.number(),
  label: z.string(),
  color: z.object({
    foreground: z.string(),
    background: z.string(),
  }),
  iconName: z.string(),
  order: z.number(),
});

export type Category = z.infer<typeof CategorySchema>;

const fetchAllCategories = () =>
  request("/categories")
    .get(z.array(CategorySchema))
    .then((categories) => categories ?? []);

export const categoriesQuery = queryOptions({
  queryKey: [CATEGORIES],
  queryFn: fetchAllCategories,
});
