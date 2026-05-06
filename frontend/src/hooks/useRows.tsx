import {
  createContext,
  type ReactNode,
  useState,
  useContext,
  type Dispatch,
  type SetStateAction,
} from "react";
import type { ErrorModel, ListRowsResponse } from "@/client";
import {
  deleteRowMutation,
  listRowsOptions,
  listRowsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useSearchParams } from "react-router-dom";

export type RowData = { hash: string } & Record<string, unknown>;

type SheetData = { row: RowData; tableName: string };

export interface RowContextType {
  rowDetailsSheetData: SheetData | null;
  setRowDetailsSheetData: Dispatch<SetStateAction<SheetData | null>>;
  rowDetailsSheetOpen: boolean;
  setRowDetailsSheetOpen: Dispatch<SetStateAction<boolean>>;
  globalFilter: string;
  setGlobalFilter: Dispatch<SetStateAction<string>>;
  tableName: string;
  isLoading: boolean;
  deleteRow: (hash: string) => void;
  data?: ListRowsResponse;
  page: number;
  rowFetchError: ErrorModel | null;
}

const RowContext = createContext<RowContextType | null>(null);
interface RowProviderProps {
  children: ReactNode;
  page?: number;
  tableName: string;
}
export const RowProvider = ({ children, tableName }: RowProviderProps) => {
  const [rowDetailsSheetData, setRowDetailsSheetData] =
    useState<SheetData | null>(null);
  const [rowDetailsSheetOpen, setRowDetailsSheetOpen] = useState(false);
  const [globalFilter, setGlobalFilter] = useState("");

  const [searchParams] = useSearchParams();
  const page = parseInt(searchParams.get("page") || "1");
  const col = searchParams.get("col");
  const order = searchParams.get("order")?.toUpperCase() as "ASC" | "DESC";

  const queryClient = useQueryClient();

  const deleteMutation = useMutation({
    ...deleteRowMutation(),
    onSuccess: () => {
      toast.success("Row deleted successfully");
      queryClient.invalidateQueries({
        queryKey: listRowsQueryKey({ path: { tableName: tableName! } }),
      });
      setRowDetailsSheetOpen(false);
      setRowDetailsSheetData(null);
    },
  });

  const deleteRow = async (hash: string) => {
    await deleteMutation.mutateAsync({
      path: { tableName: tableName!, hash },
      query: { page },
    });
  };

  const {
    data,
    isLoading,
    error: rowFetchError,
  } = useQuery(
    listRowsOptions({
      path: { tableName: tableName! },
      query: {
        page,
        column: col || undefined,
        order: order || undefined,
      },
    }),
  );
  return (
    <RowContext.Provider
      value={{
        rowDetailsSheetData,
        setRowDetailsSheetData,
        rowDetailsSheetOpen,
        globalFilter,
        setGlobalFilter,
        setRowDetailsSheetOpen,
        tableName,
        isLoading,
        deleteRow,
        data,
        rowFetchError,
        page,
      }}
    >
      {children}
    </RowContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export function useRowContext() {
  const context = useContext(RowContext);
  if (!context) {
    throw new Error("useRowContext must be used within a RowProvider");
  }
  return context;
}
