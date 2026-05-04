import {
  createContext,
  type ReactNode,
  useState,
  type Dispatch,
  type SetStateAction,
} from "react";

export type RowData = { hash: string } & Record<string, unknown>;

type SheetData = { row: RowData; tableName: string };

export interface RowContextType {
  sheetData: SheetData | null;
  setSheetData: Dispatch<SetStateAction<SheetData | null>>;
  sheetOpen: boolean;
  setSheetOpen: Dispatch<SetStateAction<boolean>>;
  globalFilter: string;
  setGlobalFilter: Dispatch<SetStateAction<string>>;
}

const RowContext = createContext<RowContextType | null>(null);

export const RowProvider = ({ children }: { children: ReactNode }) => {
  const [sheetData, setSheetData] = useState<SheetData | null>(null);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [globalFilter, setGlobalFilter] = useState("");
  return (
    <RowContext.Provider
      value={{
        sheetData,
        setSheetData,
        sheetOpen,
        globalFilter,
        setGlobalFilter,
        setSheetOpen,
      }}
    >
      {children}
    </RowContext.Provider>
  );
};
