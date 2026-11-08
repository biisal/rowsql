import { RowUpsertForm } from "@/components/rows/RowUpsertForm";
import { useRowContext } from "@/hooks/useRows";
import { DataView } from "@/components/tables/DataView";
import { RecordSheet } from "@/components/shared/RecordSheet";

export const RowDetailsSheet = () => {
  const {
    viewState,
    rowDetailsSheetData,
    setRowDetailsSheetOpen,
    rowDetailsSheetOpen,
    tableName,
    closeRowDetails,
  } = useRowContext();
  return (
    <RecordSheet
      open={rowDetailsSheetOpen}
      onOpenChange={(open) =>
        !open ? closeRowDetails() : setRowDetailsSheetOpen(true)
      }
      title={rowDetailsSheetData?.tableName || tableName}
      description={
        viewState === "edit"
          ? "Editing record details."
          : "Viewing record details."
      }
    >
      {viewState === "view" && <DataView />}
      {viewState === "edit" && (
        <RowUpsertForm key={rowDetailsSheetData?.row?.key} />
      )}
    </RecordSheet>
  );
};
