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

export function SummaryTable() {
  return (
    <Suspense fallback="loading">
      <SummaryTableContents />
    </Suspense>
  );
}

const column = createColumnHelper<{ category: Category; subtotal: number }>();
const columns = [
  column.accessor("category", { cell: ({ getValue }) => getValue().label }),
  column.accessor("subtotal", {
    cell: ({ getValue }) => moneyFormatter.format(getValue()),
  }),
];

function SummaryTableContents() {
  const [summary, categories] = useSuspenseQueries({
    queries: [
      expensesSummaryQuery(new Date(2000, 0, 1), new Date(2100, 0, 1)),
      categoriesQuery,
    ],
  });

  const data = useMemo(() => {
    const categoriesMap = dict(categories.data, "id");
    return Object.entries(summary.data?.categories ?? {}).map(
      ([categoryId, subtotal]) => ({
        category: categoriesMap[parseInt(categoryId)],
        subtotal,
      })
    );
  }, [summary.data, categories.data]);

  const table = useReactTable({
    columns,
    data,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Table>
      <TableHeader>
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map((header) => (
              <TableHead key={header.id} colSpan={header.colSpan}>
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
              <TableCell key={cell.id}>
                {flexRender(cell.column.columnDef.cell, cell.getContext())}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
