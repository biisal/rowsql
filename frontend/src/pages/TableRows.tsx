import { useParams, useSearchParams } from "react-router-dom";
import { AppPagination } from "@/components/shared/AppPagination";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DeleteAlert, RowOrderForm, Rows } from "@/feature/tables";
import { RowProvider, useRowContext } from "@/hooks/useRows";
import { RowInsertOrUpdateForm } from "@/feature/rows/components/row-insert-or-update-form";

function TablePageContent() {
  const { tableName } = useParams<{ tableName: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const { page, rowFetchError, data } = useRowContext();

  if (rowFetchError) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-destructive font-medium">
          Error: {rowFetchError?.detail || "An unknown error occurred"}
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
            <RowInsertOrUpdateForm />
          </div>
        </div>

        <Card className="border-border/50 bg-card/50 backdrop-blur-sm">
          <CardHeader className="px-6 py-4 flex justify-between items-center flex-wrap border-b border-border/50">
            <CardTitle className="text-lg font-medium">Table Data</CardTitle>
            <RowOrderForm
              cols={data.cols || []}
              initialValue={{
                col: searchParams.get("col"),
                order: searchParams.get("order")?.toUpperCase() as
                  | "ASC"
                  | "DESC",
              }}
              setUrlParams={setSearchParams}
            />
          </CardHeader>
          <CardContent className="p-4">
            <Rows tableName={tableName || ""} />
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

export function TablePage() {
  const { tableName } = useParams<{ tableName: string }>();

  return (
    <RowProvider tableName={tableName!}>
      <TablePageContent />
    </RowProvider>
  );
}
