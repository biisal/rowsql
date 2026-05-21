import type { RunSqlQueryOutputBody } from "@/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnDef,
  type SortingState,
} from "@tanstack/react-table";
import {
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  ChevronUp,
  Search,
} from "lucide-react";
import { useMemo, useState } from "react";

type SqlRow = NonNullable<NonNullable<RunSqlQueryOutputBody["rows"]>[number]>;

interface QueryTableProps {
  data?: RunSqlQueryOutputBody;
  error?: string;
}

export function QueryTable({ data, error }: QueryTableProps) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [globalFilter, setGlobalFilter] = useState("");

  const columns = useMemo<ColumnDef<SqlRow>[]>(
    () =>
      data?.columns?.map((col, i) => ({
        id: col,
        header: col,
        accessorFn: (row) => row[i],
      })) ?? [],
    [data?.columns],
  );

  // eslint-disable-next-line react-hooks/incompatible-library
  const table = useReactTable({
    data: (data?.rows as SqlRow[]) ?? [],
    columns,
    state: { sorting, globalFilter },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  if (error) return <TableEmptyState message={error} variant="error" />;
  if (!data) return <TableEmptyState message="Run a query to see results" />;
  if (!data.columns?.length)
    return (
      <TableEmptyState
        message="Query executed successfully"
        variant="success"
      />
    );

  const { pageIndex } = table.getState().pagination;

  return (
    <div className="flex h-full flex-col bg-background">
      <div className="flex items-center gap-2 border-b border-border px-3 py-2">
        <Search className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
        <Input
          placeholder="Filter results…"
          value={globalFilter}
          onChange={(e) => setGlobalFilter(e.target.value)}
          className="h-7 max-w-xs border-0 bg-transparent p-0 text-xs shadow-none focus-visible:ring-0"
        />
        <span className="ml-auto text-xs text-muted-foreground">
          {table.getFilteredRowModel().rows.length} row
          {table.getFilteredRowModel().rows.length !== 1 ? "s" : ""}
        </span>
      </div>

      <ScrollArea type="always" className="flex-1 overflow-auto">
        <Table className="w-full text-xs">
          <TableHeader className="sticky top-0 z-10 bg-muted/60 backdrop-blur">
            {table.getHeaderGroups().map((hg) => (
              <TableRow
                key={hg.id}
                className="border-border hover:bg-transparent"
              >
                {hg.headers.map((header) => {
                  const sorted = header.column.getIsSorted();
                  return (
                    <TableHead
                      key={header.id}
                      onClick={header.column.getToggleSortingHandler()}
                      className="h-8 cursor-pointer select-none whitespace-nowrap px-3 font-semibold uppercase tracking-wider text-muted-foreground hover:text-foreground"
                    >
                      <span className="flex items-center gap-1">
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                        {sorted === "asc" && <ChevronUp className="h-3 w-3" />}
                        {sorted === "desc" && (
                          <ChevronDown className="h-3 w-3" />
                        )}
                      </span>
                    </TableHead>
                  );
                })}
              </TableRow>
            ))}
          </TableHeader>

          <TableBody>
            {table.getRowModel().rows.length > 0 ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  className="border-border/50 hover:bg-muted/30"
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className="max-w-[320px] truncate whitespace-nowrap px-3 py-1.5 font-mono"
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-muted-foreground"
                >
                  No results match your filter.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </ScrollArea>

      {table.getPageCount() > 1 && (
        <div className="flex shrink-0 items-center justify-between border-t border-border px-3 py-2">
          <span className="text-xs text-muted-foreground">
            Page {pageIndex + 1} of {table.getPageCount()}
          </span>
          <div className="flex gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => table.previousPage()}
              disabled={!table.getCanPreviousPage()}
            >
              <ChevronLeft className="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-6 w-6"
              onClick={() => table.nextPage()}
              disabled={!table.getCanNextPage()}
            >
              <ChevronRight className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

interface TableEmptyStateProps {
  message: string;
  variant?: "default" | "error" | "success";
}

const variantStyles = {
  default: "text-muted-foreground",
  error: "text-red-400",
  success: "text-emerald-400",
};

const TableEmptyState = ({
  message,
  variant = "default",
}: TableEmptyStateProps) => (
  <div className="flex h-full w-full items-center justify-center">
    <p className={`text-xs ${variantStyles[variant]}`}>{message}</p>
  </div>
);
