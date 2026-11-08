import { Search, X } from "lucide-react";
import { SidebarInput } from "@/components/ui/sidebar";

interface SidebarTableSearchProps {
  searchText: string;
  onSearchTextChange: (searchText: string) => void;
  onClear: () => void;
}

export function SidebarTableSearch({
  searchText,
  onSearchTextChange,
  onClear,
}: SidebarTableSearchProps) {
  return (
    <div
      role="search"
      className="relative group-data-[collapsible=icon]:hidden"
    >
      <Search className="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
      <SidebarInput
        type="search"
        value={searchText}
        onChange={(event) => onSearchTextChange(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Escape") onClear();
        }}
        placeholder="Search tables…"
        aria-label="Search tables"
        className="pl-8 pr-8 [&::-webkit-search-cancel-button]:hidden"
      />
      {searchText && (
        <button
          type="button"
          onClick={onClear}
          aria-label="Clear table search"
          className="absolute right-1.5 top-1/2 -translate-y-1/2 rounded-sm p-1 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
        >
          <X className="size-3.5" />
        </button>
      )}
    </div>
  );
}
