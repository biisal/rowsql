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
import { RowInsertOrUpdateForm } from "@/feature/rows/components/row-insert-or-update-form";
import { useRowContext } from "@/hooks/useRows";

const formatValue = (value: unknown) => {
  console.log({ value });
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
  const { sheetData, setSheetOpen, sheetOpen, deleteRow } = useRowContext();
  console.log({ sheetData });

  return (
    <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-xl">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{sheetData?.tableName}</SheetTitle>
          <SheetDescription>Viewing record details.</SheetDescription>
        </SheetHeader>

        <ScrollArea type="always" className="min-h-0 flex-1">
          <div className="flex flex-col divide-y p-4">
            {Object.entries(sheetData?.row || {}).map(([key, value]) => (
              <div key={key} className="flex flex-col gap-1 py-3">
                <label className="text-xs font-medium text-muted-foreground capitalize">
                  {key}
                </label>
                <p className="font-mono text-sm break-all whitespace-pre-wrap">
                  {(() => {
                    switch (typeof value) {
                      case "undefined":
                        return <p className="text-muted-foreground">-</p>;
                      case null:
                        return <p className="text-muted-foreground">null</p>;
                      default:
                        return formatValue(value);
                    }
                  })()}
                </p>
              </div>
            ))}
          </div>
        </ScrollArea>

        <SheetFooter>
          <div className="grid grid-cols-2 gap-4">
            {sheetData?.row.hash && (
              <>
                <RowInsertOrUpdateForm hash={sheetData.row.hash}>
                  <Button className="">Edit</Button>
                </RowInsertOrUpdateForm>

                <Button
                  variant={"danger"}
                  onClick={() => deleteRow(sheetData.row.hash)}
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
