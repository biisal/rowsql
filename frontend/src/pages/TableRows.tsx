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
              tableName={tableName || ""}
              isLoading={isLoading}
              data={data}
              deleteRow={deleteRow}
            />
            <div className="flex items-center flex-col sm:flex-row justify-center md:justify-between  py-4">
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
