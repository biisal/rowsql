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
import type { RowData } from "@/feature/tables/components/Rows";
import { Link } from "react-router-dom";

interface RowDataSheetProps {
  data: { row: RowData; tableName: string } | null;
  setOpenChange: React.Dispatch<React.SetStateAction<boolean>>;
  open: boolean;
  deleteRow: (hash: string) => void;
}
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
export const RowDetailsSheet = ({
  data,
  setOpenChange,
  open,
  deleteRow,
}: RowDataSheetProps) => {
  return (
    <Sheet open={open} onOpenChange={setOpenChange}>
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-150">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{data?.tableName}</SheetTitle>
          <SheetDescription>
            Viewing{" "}
            <span className="bg-muted-foreground/15 rounded-md p-0.5 px-2">
              {data?.tableName}
            </span>{" "}
            details.
          </SheetDescription>
        </SheetHeader>

        <ScrollArea type="always" className="min-h-0 flex-1">
          <div className="flex flex-col divide-y p-4">
            {Object.entries(data?.row || {}).map(([key, value]) => (
              <div key={key} className="flex flex-col gap-1 py-3">
                <label className="text-xs font-medium text-muted-foreground capitalize">
                  {key}
                </label>
                <p className="font-mono text-sm break-all whitespace-pre-wrap">
                  {formatValue(value)}
                </p>
              </div>
            ))}
          </div>
        </ScrollArea>

        <SheetFooter>
          <div className="grid grid-cols-2 gap-4">
            {data?.row.hash && (
              <>
                <Button className="" asChild>
                  <Link
                    to={`/tables/${data?.tableName}/rows?hash=${data.row.hash}`}
                  >
                    Edit
                  </Link>
                </Button>

                <Button
                  variant={"danger"}
                  onClick={() => deleteRow(data.row.hash)}
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
