import { categoriesQuery, type Category } from "@/stores/category";
import { expensesSummaryQuery } from "@/stores/summary";
import { useSuspenseQueries } from "@tanstack/react-query";
import { Suspense, useMemo } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { dict } from "@/lib/utils";
import { moneyFormatter } from "@/lib/formatters";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";

export function SummaryTable() {
  return (
    <Suspense fallback="loading">
      <SummaryTableContents />
    </Suspense>
  );
}

const column = createColumnHelper<{ category: Category; subtotal: number }>();
const columns = [
  column.accessor("category", {
    header: "Catégorie",
    cell: ({ getValue }) => <CategoryCell {...getValue()} />,
  }),
  column.accessor("subtotal", {
    header: "Sous-total",
    cell: ({ getValue }) => (
      <div className="text-center">{moneyFormatter.format(getValue())}</div>
    ),
  }),
];

function CategoryCell({ label, iconName }: Category) {
  return (
    <Tooltip>
      <TooltipTrigger>
        <div className="flex items-center gap-2 max-w-52 w-fit">
          <span className="material-symbols-outlined">
            <span className="block h-full text-base">{iconName}</span>
          </span>
          <span className="text-ellipsis overflow-hidden text-sm">{label}</span>
        </div>
      </TooltipTrigger>
      <TooltipContent>{label}</TooltipContent>
    </Tooltip>
  );
}

function SummaryTableContents() {
  const [summary, categories] = useSuspenseQueries({
    queries: [
      expensesSummaryQuery(new Date(2000, 0, 1), new Date(2100, 0, 1)),
      categoriesQuery,
    ],
  });

  const data = useMemo(() => {
    const categoriesMap = dict(categories.data, "id");
    return Object.entries(summary.data?.categories ?? {})
      .map(([categoryId, subtotal]) => ({
        category: categoriesMap[parseInt(categoryId)],
        subtotal,
      }))
      .sort((a, b) => a.category.order - b.category.order);
  }, [summary.data, categories.data]);

  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Table className="max-w-80 mx-auto">
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <TableHead
                key={header.id}
                colSpan={header.colSpan}
                className="font-bold text-lg last:text-center"
              >
                {flexRender(
                  header.column.columnDef.header,
                  header.getContext()
                )}
              </TableHead>
            ))}
          </TableRow>
        ))}
      </TableHeader>
      <TableBody>
        <TableRow>
          <TableCell>
            <strong>Total</strong>
          </TableCell>
          <TableCell>
            {moneyFormatter.format(summary.data?.total ?? 0)}
          </TableCell>
        </TableRow>
        {table.getRowModel().rows.map((row) => (
          <TableRow key={row.id}>
            {row.getVisibleCells().map((cell) => (
              <TableCell key={cell.id} className="py-1">
                {flexRender(cell.column.columnDef.cell, cell.getContext())}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
