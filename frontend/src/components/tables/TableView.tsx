import React, { useState } from "react";
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";

import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { RowDetailsSheet } from "./RowDetailsSheet";
import { Input } from "@/components/ui/input";
import { ChevronRight } from "lucide-react";
import { useIsMobile } from "@/hooks/use-mobile";
import { useRowContext } from "@/hooks/useRows";

interface TableViewProps {
  tableName: string;
}

export type RowData = { hash: string } & Record<string, unknown>;

export const TableView = ({ tableName }: TableViewProps) => {
  const isMobile = useIsMobile();
  const { setSheetOpen, setSheetData, isLoading, data } = useRowContext();
  const [globalFilter, setGlobalFilter] = useState("");

  const columns: ColumnDef<RowData>[] = React.useMemo(
    () => [
      ...(data?.cols?.map((col) => ({
        header: col.columnName,
        accessorKey: col.columnName,
        size: 200,
        maxSize: 200,
      })) || []),

      {
        id: "action",
        size: 0,
        cell: () => (
          <ChevronRight className="h-4 w-4 text-muted-foreground group-hover:text-foreground" />
        ),
      },
    ],
    [data?.cols],
  );

  const flattenedRows = React.useMemo<RowData[]>(
    () =>
      data?.rows?.map((row) => ({
        hash: row.hash,
        ...Object.fromEntries(
          (row.columns || []).map((col) => [col.columnName, col.value]),
        ),
      })) || [],
    [data?.rows],
  );

  const handleRowClick = (row: RowData) => {
    setSheetData({ row: row, tableName });
    setSheetOpen(true);
  };
  const table = useReactTable({
    data: flattenedRows,
    columns: columns,
    state: {
      globalFilter,
      columnPinning: {
        right: ["action"],
      },
    },
    enableColumnPinning: true,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
  });

  return (
    <div className="rounded-md overflow-auto relative">
      <Input
        placeholder="Filter results…"
        value={globalFilter}
        onChange={(e) => setGlobalFilter(e.target.value)}
        className="max-w-sm"
      />

      <RowDetailsSheet />
      <Table>
        <TableCaption>click the row to view complete data</TableCaption>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell
                colSpan={columns.length}
                className="text-center border border-border"
              >
                Loading...
              </TableCell>
            </TableRow>
          ) : (
            table.getRowModel().rows.map((row) => (
              <TableRow
                onClick={() => handleRowClick(row.original)}
                key={row.id}
                className="group cursor-pointer border-l-3 border-b-0 hover:border-primary"
              >
                {row.getVisibleCells().map((cell) => {
                  const isPinned = cell.column.getIsPinned();
                  return (
                    <TableCell
                      key={cell.id}
                      className="text-foreground/80"
                      style={{
                        position: isPinned ? "sticky" : undefined,
                        right: isPinned === "right" ? 0 : undefined,
                        zIndex: isPinned ? 1 : 0,
                        background: isMobile ? "var(--background)" : undefined,

                        maxWidth: `${cell.column.getSize()}px`,
                        overflow: "hidden",
                        textOverflow: "ellipsis",
                        whiteSpace: "nowrap",
                      }}
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  );
                })}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
};
