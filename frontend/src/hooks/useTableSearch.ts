import { useMemo, useState } from "react";
import type { ListTablesRow } from "@/client/types.gen";

export function filterTablesByName(
  tables: ListTablesRow[],
  searchText: string,
): ListTablesRow[] {
  const normalizedSearch = searchText.trim().toLowerCase();
  if (!normalizedSearch) return tables;
  return tables.filter((table) =>
    table.tableName.toLowerCase().includes(normalizedSearch),
  );
}

export function useTableSearch(tables: ListTablesRow[]) {
  const [searchText, setSearchText] = useState("");

  const filteredTables = useMemo(
    () => filterTablesByName(tables, searchText),
    [tables, searchText],
  );

  return {
    searchText,
    setSearchText,
    clearSearch: () => setSearchText(""),
    filteredTables,
    isSearching: searchText.trim() !== "",
  };
}
