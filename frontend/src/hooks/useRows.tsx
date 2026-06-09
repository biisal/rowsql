import {
  createContext,
  type ReactNode,
  useState,
  useContext,
  type Dispatch,
  type SetStateAction,
} from "react";
import type { ErrorModel, ListRowsResponse, RowSet } from "@/client";
import {
  deleteRowMutation,
  listRowsOptions,
  listRowsQueryKey,
} from "@/client/@tanstack/react-query.gen";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { useSearchParams } from "react-router-dom";

type SheetData = { row: RowSet; tableName: string };

interface RowContextType {
  viewState: ViewState;
  setViewState: Dispatch<SetStateAction<ViewState>>;
  rowDetailsSheetData: SheetData | null;
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
  openRowDetails: (row?: RowSet, mode?: ViewState) => void;
  closeRowDetails: () => void;
}

type ViewState = "view" | "edit";

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
  const [viewState, setViewState] = useState<ViewState>("view");
  const [globalFilter, setGlobalFilter] = useState("");

  const [searchParams] = useSearchParams();
  const page = parseInt(searchParams.get("page") || "1");
  const col = searchParams.get("col");
  const order = searchParams.get("order")?.toUpperCase() as "ASC" | "DESC";

  const queryClient = useQueryClient();

  const openRowDetails = (row?: RowSet, mode: ViewState = "view") => {
    setViewState(mode);
    setRowDetailsSheetData(row ? { row, tableName } : null);
    setRowDetailsSheetOpen(true);
  };

  const closeRowDetails = () => {
    setRowDetailsSheetOpen(false);
    // Reset after animation if needed, but for now simple reset
    setViewState("view");
    setRowDetailsSheetData(null);
  };

  const deleteMutation = useMutation({
    ...deleteRowMutation(),
    onSuccess: () => {
      toast.success("Row deleted successfully");
      queryClient.invalidateQueries({
        queryKey: listRowsQueryKey({ path: { tableName: tableName! } }),
      });
      closeRowDetails();
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
        viewState,
        setViewState,
        rowDetailsSheetData,
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
        openRowDetails,
        closeRowDetails,
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
