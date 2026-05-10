import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { RowUpsertForm } from "@/components/rows/RowUpsertForm";
import { useRowContext } from "@/hooks/useRows";

const formatValue = (value: unknown) => {
  if (typeof value === "string") {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }
  return value?.toString();
};
export const RowDetailsSheet = () => {
  const {
    rowDetailsSheetData: rowDetailssheetData,
    setRowDetailsSheetOpen,
    rowDetailsSheetOpen,
    deleteRow,
  } = useRowContext();
  if (!rowDetailssheetData) return null;
  const row = rowDetailssheetData.row.columns;
  return (
    <Sheet open={rowDetailsSheetOpen} onOpenChange={setRowDetailsSheetOpen}>
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-xl">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{rowDetailssheetData?.tableName}</SheetTitle>
          <SheetDescription>Viewing record details.</SheetDescription>
        </SheetHeader>

        <ScrollArea type="always" className="min-h-0 flex-1">
          <div className="flex flex-col divide-y p-4">
            {(row || []).map((col) => (
              <div key={col.columnName} className="flex flex-col gap-1 py-3">
                <label className="text-xs font-medium text-muted-foreground capitalize">
                  {col.columnName}
                </label>
                <div className="font-mono text-sm break-all whitespace-pre-wrap">
                  {(() => {
                    const value = col.value;
                    if (value === undefined)
                      return <p className="text-muted-foreground">-</p>;
                    if (value === null)
                      return <p className="text-muted-foreground">null</p>;
                    return formatValue(value);
                  })()}
                </div>
              </div>
            ))}
          </div>
        </ScrollArea>

        <SheetFooter>
          <div className="grid grid-cols-2 gap-4">
            {rowDetailssheetData?.row.hash && (
              <>
                <RowUpsertForm hash={rowDetailssheetData.row.hash}>
                  <Button className="">Edit</Button>
                </RowUpsertForm>

                <Button
                  variant={"danger"}
                  onClick={() => deleteRow(rowDetailssheetData.row.hash)}
                >
                  Delete
                </Button>
              </>
            )}
          </div>
          <SheetClose asChild>
            <Button className="w-full" variant={"outline"}>
              Close
            </Button>
          </SheetClose>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};
