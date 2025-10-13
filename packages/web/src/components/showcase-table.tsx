import { useExpenses } from "@/hooks/use-expenses";
import { Suspense } from "react";

export function ShowcaseTable() {
  return (
    <Suspense fallback={"loading :)"}>
      <ShowcaseTableContents />
    </Suspense>
  );
}

function ShowcaseTableContents() {
  const expenses = useExpenses({
    start: new Date(2000, 0, 1),
    end: new Date(2100, 0, 1),
  });

  return <pre>{JSON.stringify(expenses, null, 2)}</pre>;
}
