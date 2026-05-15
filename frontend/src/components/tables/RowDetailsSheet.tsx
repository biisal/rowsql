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

export const RowDetailsSheet = ({ goBack }: { goBack?: () => void }) => {
  const {
    viewState,
    setViewState,
    rowDetailsSheetData: rowDetailssheetData,
    setRowDetailsSheetOpen,
    rowDetailsSheetOpen,
  } = useRowContext();
  return (
    <Sheet open={rowDetailsSheetOpen} onOpenChange={setRowDetailsSheetOpen}>
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-xl">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{rowDetailssheetData?.tableName}</SheetTitle>
          <SheetDescription>Viewing record details.</SheetDescription>
        </SheetHeader>

        {viewState === "view" && <DataView />}
        {viewState === "edit" && (
          <RowUpsertForm
            hash={rowDetailssheetData?.row.hash}
            goBack={goBack ?? (() => setViewState("view"))}
          />
        )}
      </SheetContent>
    </Sheet>
  );
};
