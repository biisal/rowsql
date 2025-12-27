import * as React from "react";
import { useEffect, useState } from "react";
import { useParams, Link, useSearchParams } from "react-router-dom";
import { Plus, Loader2, MoreHorizontal } from "lucide-react";
import axios from "axios";
import api from "@/lib/axios";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AppPagination } from "@/components/app-pagination";
import { cn } from "@/lib/utils";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { DeletAlert } from "@/components/delete-alert";

interface Column {
  columnName: string;
  dataType: string;
}

interface Table {
  table_name: string;
  table_schema: string;
}

type CellValue = string | number | boolean | null;
type RowData = CellValue[];

interface TableData {
  Page: number;
  Tables: Table[];
  Cols: Column[];
  ActiveTable: string;
  Rows: RowData[];
}

export function TableRows() {
  const { tableName } = useParams<{ tableName: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const [refresh, setRefesh] = useState(0);
  const [data, setData] = useState<TableData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedRows, setSelectedRows] = useState<Record<number, boolean>>({});

  const page = parseInt(searchParams.get("page") || "1");

  const fetchData = React.useCallback(async () => {
    if (!tableName) return;

    setLoading(true);
    setError(null);
    try {
      const response = await api.get(`/table/${tableName}?page=${page}`);
      if (response.data.success) {
        setData(response.data.data);
        setSelectedRows({}); // Reset selection on page change
      }
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.message);
      } else {
        setError(
          err instanceof Error ? err.message : "An unknown error occurred",
        );
      }
    } finally {
      setLoading(false);
    }
  }, [tableName, page, refresh]); // eslint-disable-line

  const deleteRow = async (hash: string) => {
    try {
      const res = await api.delete(`/table/${tableName}/row/${hash}`);
      if (res.data.success) {
        toast.success("Row deleted successfully");
        setRefesh((r) => r + 1);
        return;
      }
      toast.error("Failed to delete row");
    } catch (err) {
      console.error(err);
      if (axios.isAxiosError(err)) {
        toast.error(err.message);
      } else {
        toast.error(
          err instanceof Error ? err.message : "An unknown error occurred",
        );
      }
    }
  };

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const toggleRowSelection = (index: number) => {
    setSelectedRows((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const toggleAllSelection = () => {
    if (!data) return;
    const allSelected =
      data.Rows.length > 0 && data.Rows.every((_, idx) => selectedRows[idx]);
    if (allSelected) {
      setSelectedRows({});
    } else {
      const newSelection: Record<number, boolean> = {};
      data.Rows.forEach((_, idx) => {
        newSelection[idx] = true;
      });
      setSelectedRows(newSelection);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-destructive font-medium">Error: {error}</div>
      </div>
    );
  }

  if (!data) {
    return <div className="p-8">No data found.</div>;
  }

  const isAllSelected =
    data.Rows.length > 0 && data.Rows.every((_, idx) => selectedRows[idx]);
  const isSomeSelected =
    data.Rows.some((_, idx) => selectedRows[idx]) && !isAllSelected;

  return (
    <div className="flex-1 p-4 md:p-8 overflow-y-auto w-full max-w-full">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold tracking-tight">
            {data.ActiveTable}
          </h1>
          <div className="flex items-center justify-center gap-1">
            <DeletAlert tableName={data.ActiveTable} />
            <Link to={`/table/${data.ActiveTable}/form`}>
              <Button className="shadow-lg shadow-primary/20">
                <Plus className="mr-2 h-4 w-4" /> Insert Record
              </Button>
            </Link>
          </div>
        </div>

        <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
          <CardHeader className="px-6 py-4 border-b border-border/50">
            <CardTitle className="text-lg font-medium">Table Data</CardTitle>
          </CardHeader>
          <CardContent className="p-4">
            <div className="rounded-md border overflow-auto  relative">
              <div
                className="grid min-w-full"
                style={{
                  gridTemplateColumns: `40px ${data.Cols.map(() => "minmax(150px, 1fr)").join(" ")} 80px`,
                }}
              >
                {/* Checkbox Header */}
                <div className="h-10 px-2 flex items-center justify-center border-b bg-muted/50 sticky top-0 z-20">
                  <Checkbox
                    checked={
                      isAllSelected
                        ? true
                        : isSomeSelected
                          ? "indeterminate"
                          : false
                    }
                    onCheckedChange={toggleAllSelection}
                    aria-label="Select all"
                  />
                </div>

                {/* Data Columns Headers */}
                {data.Cols.map((col) => (
                  <div
                    key={col.columnName}
                    className="h-10 px-2 text-left align-middle font-medium text-muted-foreground flex items-center border-b bg-muted/50 sticky top-0 z-20"
                  >
                    {col.columnName}
                  </div>
                ))}

                {/* Actions Header */}
                <div className="h-10 px-2 text-right align-middle font-medium text-muted-foreground flex items-center justify-end border-b bg-muted/50 sticky top-0 right-0 z-30 shadow-[-1px_0_0_0_var(--border)]">
                  Actions
                </div>

                {/* Body Rows */}
                {data.Rows.length > 0 ? (
                  data.Rows.map((row, rowIndex) => {
                    const hash = String(row[0] ?? "");
                    const isSelected = !!selectedRows[rowIndex];

                    return (
                      <React.Fragment key={rowIndex}>
                        {/* Checkbox Cell */}
                        <div
                          className={cn(
                            "p-2 flex items-center justify-center border-b bg-card group-hover:bg-muted/50 transition-colors",
                            isSelected && "bg-muted",
                          )}
                        >
                          <Checkbox
                            checked={isSelected}
                            onCheckedChange={() => toggleRowSelection(rowIndex)}
                            aria-label="Select row"
                          />
                        </div>

                        {row
                          .slice(1)
                          .map((cell: CellValue, cellIndex: number) => (
                            <div
                              key={cellIndex}
                              className={cn(
                                "p-2 align-middle flex items-center border-b bg-card group-hover:bg-muted/50 transition-colors lowercase",
                                isSelected && "bg-muted",
                              )}
                            >
                              {cell === null ? "NULL" : String(cell)}
                            </div>
                          ))}

                        {/* Actions Cell */}
                        <div
                          className={cn(
                            "p-2 align-middle flex items-center justify-end border-b bg-card group-hover:bg-muted/50 transition-colors sticky right-0 z-10 shadow-[-1px_0_0_0_var(--border)]",
                            isSelected && "bg-muted",
                          )}
                        >
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" className="h-8 w-8 p-0">
                                <span className="sr-only">Open menu</span>
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent
                              align="end"
                              className="px-4 py-2"
                            >
                              <DropdownMenuLabel>Actions</DropdownMenuLabel>
                              <DropdownMenuItem>
                                <Link
                                  to={`/table/${data.ActiveTable}/form?hash=${hash}&page=${data.Page}`}
                                >
                                  <DropdownMenuItem>Edit Row</DropdownMenuItem>
                                </Link>
                              </DropdownMenuItem>
                              <Button
                                variant="danger"
                                className="w-full"
                                onClick={() => deleteRow(hash)}
                              >
                                Delete Row
                              </Button>
                              <DropdownMenuSeparator />
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </React.Fragment>
                    );
                  })
                ) : (
                  <div className="col-span-full h-24 flex items-center justify-center">
                    No results.
                  </div>
                )}
              </div>
            </div>
            <div className="flex items-center justify-between py-4">
              <div className="text-muted-foreground flex-1 text-sm">
                {Object.values(selectedRows).filter(Boolean).length} of{" "}
                {data.Rows.length} row(s) selected.
              </div>
              <AppPagination
                currentPage={page}
                onPageChange={(newPage) =>
                  setSearchParams({ page: String(newPage) })
                }
                hasNextPage={data.Rows.length > 0}
                hasPreviousPage={page > 1}
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
