import { useState } from "react";
import { useParams, Link, useSearchParams } from "react-router-dom";
import { Plus } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  listRowsOptions,
  deleteRowMutation,
  listRowsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { Button } from "@/components/ui/button";
import { AppPagination } from "@/components/shared/AppPagination";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { DeleteAlert, RowOrderForm, Rows } from "@/feature/tables";

export function TablePage() {
  const { tableName } = useParams<{ tableName: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const queryClient = useQueryClient();
  const [selectedRows, setSelectedRows] = useState<Record<number, boolean>>({});
  const page = parseInt(searchParams.get("page") || "1");
  const col = searchParams.get("col");
  const order = searchParams.get("order")?.toUpperCase() as "ASC" | "DESC";

  const { data, isLoading, error } = useQuery(
    listRowsOptions({
      path: { tableName: tableName! },
      query: {
        page,
        column: col || undefined,
        order: order || undefined,
      },
    }),
  );

  const deleteMutation = useMutation({
    ...deleteRowMutation(),
    onSuccess: () => {
      toast.success("Row deleted successfully");
      queryClient.invalidateQueries({
        queryKey: listRowsQueryKey({ path: { tableName: tableName! } }),
      });
    },
    onError: (err) => {
      const errorMessage = err?.detail || "Failed to delete row";
      toast.error(errorMessage);
    },
  });

  const deleteRow = async (hash: string) => {
    await deleteMutation.mutateAsync({
      path: { tableName: tableName!, hash },
      query: { page },
    });
  };

  // Reset search params when table name changes
  // useEffect(() => {
  // 	setSearchParams({}, { replace: true });
  // 	setSelectedRows({});
  // }, [tableName, setSearchParams]);

  const toggleRowSelection = (index: number) => {
    setSelectedRows((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const toggleAllSelection = () => {
    if (!data?.rows) return;
    const allSelected =
      data.rows.length > 0 && data.rows.every((_, idx) => selectedRows[idx]);
    if (allSelected) {
      setSelectedRows({});
    } else {
      const newSelection: Record<number, boolean> = {};
      data.rows.forEach((_, idx) => {
        newSelection[idx] = true;
      });
      setSelectedRows(newSelection);
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading...</div>;
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-destructive font-medium">
          Error: {error?.detail || "An unknown error occurred"}
        </div>
      </div>
    );
  }

  if (!data || !data.rows) {
    return <div className="p-8">No data found.</div>;
  }

  const isAllSelected =
    data.rows.length > 0 && data.rows.every((_, idx) => selectedRows[idx]);
  const isSomeSelected =
    data.rows.some((_, idx) => selectedRows[idx]) && !isAllSelected;

  return (
    <div className="flex-1 p-4 md:p-8 overflow-y-auto w-full max-w-full">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold tracking-tight">
            {data.activeTable}
          </h1>
          <div className="flex items-center justify-center gap-1">
            <DeleteAlert tableName={data.activeTable} />
            <Link to={`/tables/${data.activeTable}/rows/`}>
              <Button className="shadow-lg shadow-primary/20">
                <Plus className="mr-2 h-4 w-4" /> Insert Record
              </Button>
            </Link>
          </div>
        </div>

        <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
          <CardHeader className="px-6 py-4 flex justify-between items-center flex-wrap border-b border-border/50">
            <CardTitle className="text-lg font-medium">Table Data</CardTitle>
            <RowOrderForm
              cols={data.cols || []}
              initialValue={{ col, order }}
              setUrlParams={setSearchParams}
            />
          </CardHeader>
          <CardContent className="p-4">
            <Rows
              data={data}
              selectedRows={selectedRows}
              isAllSelected={isAllSelected}
              isSomeSelected={isSomeSelected}
              toggleAllSelection={toggleAllSelection}
              toggleRowSelection={toggleRowSelection}
              deleteRow={deleteRow}
            />
            <div className="flex items-center flex-col sm:flex-row justify-center md:justify-between  py-4">
              <div className="text-muted-foreground flex-1  text-sm">
                {Object.values(selectedRows).filter(Boolean).length} of{" "}
                {data.rows?.length || 0} row(s) selected.
              </div>
              <AppPagination
                currentPage={page}
                totalPages={data.totalPages || 0}
                onPageChange={(newPage) =>
                  setSearchParams((prev) => {
                    const newParams = new URLSearchParams(prev);
                    newParams.set("page", String(newPage));
                    return newParams;
                  })
                }
                hasNextPage={data.hasNextPage}
                hasPreviousPage={page > 1}
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
