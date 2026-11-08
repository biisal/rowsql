import { Button } from "@/components/ui/button";
import { SheetClose, SheetFooter } from "@/components/ui/sheet";
import { useRowContext } from "@/hooks/useRows";
import { RecordFields } from "@/components/shared/RecordFields";

export const DataView = () => {
  const {
    rowDetailsSheetData: rowDetailssheetData,
    deleteRow,
    setViewState,
  } = useRowContext();
  if (!rowDetailssheetData) return null;
  const row = rowDetailssheetData.row.columns;

  const handleEdit = () => {
    setViewState("edit");
  };
  return (
    <>
      <RecordFields
        fields={(row || []).map((col) => ({
          name: col.columnName,
          value: col.value,
        }))}
      />

      <SheetFooter>
        <div className="grid grid-cols-2 gap-4">
          {rowDetailssheetData?.row.key && (
            <>
              <Button
                variant={"danger"}
                onClick={() => deleteRow(rowDetailssheetData.row.key)}
              >
                Delete
              </Button>
            </>
          )}
          <Button onClick={handleEdit}>Edit</Button>
        </div>
        <SheetClose asChild>
          <Button className="w-full" variant={"outline"}>
            Close
          </Button>
        </SheetClose>
      </SheetFooter>
    </>
  );
};
