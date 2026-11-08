import { Button } from "@/components/ui/button";
import { SheetClose, SheetFooter } from "@/components/ui/sheet";
import { RecordFields } from "@/components/shared/RecordFields";
import { RecordSheet } from "@/components/shared/RecordSheet";

interface QueryRowSheetProps {
  columns: string[];
  row: unknown[] | null;
  rowNumber: number;
  onClose: () => void;
}

export function QueryRowSheet({
  columns,
  row,
  rowNumber,
  onClose,
}: QueryRowSheetProps) {
  const fields = columns.map((name, index) => ({ name, value: row?.[index] }));
  return (
    <RecordSheet
      open={row !== null}
      onOpenChange={(open) => !open && onClose()}
      title={`Result row ${rowNumber}`}
      description="Viewing record details."
    >
      <RecordFields fields={fields} />
      <SheetFooter>
        <SheetClose asChild>
          <Button className="w-full" variant="outline">
            Close
          </Button>
        </SheetClose>
      </SheetFooter>
    </RecordSheet>
  );
}
