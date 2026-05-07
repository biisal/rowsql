import { useEffect, useState } from "react";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { type SetURLSearchParams } from "react-router-dom";
import type { ColValue } from "@/client";

interface InitialVal {
  col: string | null;
  order: "ASC" | "DESC";
}

interface RowOrderFormProps {
  cols: ColValue[];
  setUrlParams: SetURLSearchParams;
  initialValue: InitialVal;
}

export const RowOrderForm = ({
  cols: columns,
  setUrlParams,
  initialValue,
}: RowOrderFormProps) => {
  const [selectedCol, setSelectedCol] = useState<string | null>(
    initialValue.col ??
      (columns && columns.length > 0 ? columns[0].columnName : null),
  );
  const [selectedOrder, setSelectedOrder] = useState<string>(
    initialValue.order ?? "asc",
  );

  useEffect(() => {
    setUrlParams((prev) => {
      const newParams = new URLSearchParams(prev);
      if (selectedCol) {
        newParams.set("col", selectedCol);
      }
      if (selectedOrder) {
        newParams.set("order", selectedOrder);
      }
      return newParams;
    });
  }, []);

  const handleColChange = (value: string) => {
    setSelectedCol(value);
    setUrlParams((prev) => {
      const newParams = new URLSearchParams(prev);
      newParams.set("col", value || "");
      return newParams;
    });
  };

  const handleOrderChange = (value: string) => {
    setSelectedOrder(value);
    setUrlParams((prev) => {
      const newParams = new URLSearchParams(prev);
      newParams.set("order", value);
      return newParams;
    });
  };
  return (
    <div className="flex flex-col gap-1">
      <p className="text-sm font-medium text-muted-foreground">Order by:</p>
      <div className="flex items-center justify-center gap-4 flex-wrap">
        <Select
          value={selectedCol || undefined}
          onValueChange={handleColChange}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select a column" />
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Columns</SelectLabel>
                {columns.map((col, idx) => (
                  <SelectItem key={idx} value={col.columnName || null!}>
                    {col.columnName || (
                      <div className="h-6 w-full bg-primary-foreground"></div>
                    )}
                  </SelectItem>
                ))}
              </SelectGroup>
            </SelectContent>
          </SelectTrigger>
        </Select>

        <Select value={selectedOrder} onValueChange={handleOrderChange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Select order" />
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Order</SelectLabel>
                <SelectItem value={"ASC"}>Ascending</SelectItem>
                <SelectItem value={"DESC"}>Descending</SelectItem>
              </SelectGroup>
            </SelectContent>
          </SelectTrigger>
        </Select>
      </div>
    </div>
  );
};
