import { useExpenses } from "@/hooks/use-expenses";
import { Suspense } from "react";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import type { Expense } from "@/stores/expense";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./ui/table";
import type { Category } from "@/stores/category";
import { moneyFormatter, timestampFormatter } from "@/lib/formatters";

export function ShowcaseTable() {
  return (
    <Suspense fallback={"attends un peu ça charge"}>
      <ShowcaseTableContents />
    </Suspense>
  );
}

const column = createColumnHelper<Expense>();
const columns = [
  column.accessor("timestamp", {
    header: "Date",
    cell: ({ getValue }) => timestampFormatter.format(getValue()),
  }),
  column.accessor("label", { header: "Description" }),
  column.accessor("amount", {
    header: "Montant",
    cell: ({ getValue }) => moneyFormatter.format(getValue()),
  }),
  column.accessor("category", {
    header: "Catégorie",
    cell: ({ getValue }) => <CategoryCell category={getValue()} />,
  }),
  column.accessor("spender", {
    header: "Source",
    cell: ({ getValue }) => `${getValue().firstName} ${getValue().lastName}`,
  }),
];

function CategoryCell({ category }: { category: Category }) {
  return (
    <div>
      <span className="material-symbols-outlined">{category.iconName}</span>
      <span>{category.label}</span>
    </div>
  );
}

function ShowcaseTableContents() {
  const expenses = useExpenses({
    start: new Date(2000, 0, 1),
    end: new Date(2100, 0, 1),
  });

  const table = useReactTable({
    columns,
    data: expenses,
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
