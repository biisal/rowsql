import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { RowUpsertForm } from "@/components/rows/RowUpsertForm";
import { useRowContext } from "@/hooks/useRows";
import { DataView } from "@/components/tables/DataView";

export const RowDetailsSheet = () => {
  const {
    viewState,
    rowDetailsSheetData: rowDetailssheetData,
    setRowDetailsSheetOpen,
    rowDetailsSheetOpen,
    tableName,
    closeRowDetails,
  } = useRowContext();
  return (
    <Sheet
      open={rowDetailsSheetOpen}
      onOpenChange={(open) =>
        !open ? closeRowDetails() : setRowDetailsSheetOpen(true)
      }
    >
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-xl">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{rowDetailssheetData?.tableName || tableName}</SheetTitle>
          <SheetDescription>
            {viewState === "edit"
              ? "Editing record details."
              : "Viewing record details."}
          </SheetDescription>
        </SheetHeader>

        {viewState === "view" && <DataView />}
        {viewState === "edit" && (
          <RowUpsertForm hash={rowDetailssheetData?.row?.hash} />
        )}
      </SheetContent>
    </Sheet>
  );
};
