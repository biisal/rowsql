import { create } from "zustand";
import type { ErrorResponse, Input, Table } from "../types";
import api from "../axios";
import { toast } from "sonner";
import axios, { AxiosError } from "axios";

interface TablesStore {
  tables: Table[];
  refreshTables: () => Promise<void>;
  tablesRefreshing: boolean;
  tableDeleting: boolean;
  tableCreating: boolean;
  deleteTable: (
    tableName: string,
    verificationQuery: string,
  ) => Promise<boolean>;
  createTable: (tableName: string, inputs: Input[]) => Promise<boolean>;
}

const useTableStore = create<TablesStore>((set, get) => ({
  tables: [],
  tablesRefreshing: false,
  tableDeleting: false,
  tableCreating: false,

  deleteTable: async (tableName: string, verificationQuery: string) => {
    set({ tableDeleting: true });
    try {
      const res = await api.delete(`/table`, {
        data: {
          verificationQuery: verificationQuery,
          tableName: tableName,
        },
      });
      if (res.status === 204) {
        toast.success("Table deleted successfully");
        get().refreshTables();
        return true;
      }
      toast.error("Failed to delete table");
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const axiosError = err as AxiosError<ErrorResponse>;
        const errorMessage =
          axiosError.response?.data?.error ||
          axiosError.message ||
          "Somethihng went wrong!";
        toast.error(errorMessage);
        return false;
      }
      toast.error(
        err instanceof Error ? err.message : "An unknown error occurred",
      );
    } finally {
      set({ tableDeleting: false });
    }
    return false;
  },
  createTable: async (tableName: string, inputs: Input[]) => {
    set({ tableCreating: true });
    try {
      const res = await api.post("/table/form/new", { tableName, inputs });
      if (res.status === 201) {
        toast.success("Table created successfully");
        get().refreshTables();
        return true;
      }
      toast.error("Failed to create table!Something went really wrong");
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const axiosError = err as AxiosError<ErrorResponse>;
        const errorMessage =
          axiosError.response?.data?.error ||
          axiosError.message ||
          "Somethihng went wrong!";
        toast.error(errorMessage);
        return false;
      }
      toast.error(
        err instanceof Error ? err.message : "An unknown error occurred",
      );
    } finally {
      set({ tableCreating: false });
    }
    return false;
  },
  refreshTables: async () => {
    set({ tablesRefreshing: true });
    try {
      const response = await api.get("/tables");
      if (response.data.success && Array.isArray(response.data.data)) {
        const tables = response.data.data;
        set({ tables: tables });
        return;
      }
      set({ tables: [] });
    } catch (error) {
      console.error("Failed to fetch tables:", error);
      set({ tables: [] });
    } finally {
      set({ tablesRefreshing: false });
    }
  },
}));

export default useTableStore;
