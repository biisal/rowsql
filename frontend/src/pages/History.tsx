import { useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Loader2, Search } from "lucide-react";
import { listHistoryOptions } from "@/client/@tanstack/react-query.gen";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";

type ActionKind = "create" | "insert" | "update" | "delete" | "sql" | "other";
type Scope = "all" | "tables" | "rows" | "sql";

interface HistoryEvent {
  id: number;
  message: string;
  time: string;
  kind: ActionKind;
  tableName: string | null;
  sqlText: string | null;
}

interface ActionStyle {
  label: string;
  className: string;
}

const ALL_TABLES = "__all_tables__";

const actionStyles: Record<ActionKind, ActionStyle> = {
  create: {
    label: "Create",
    className: "bg-green-500/10 text-green-400 border-green-500/20",
  },
  insert: {
    label: "Insert",
    className: "bg-blue-500/10 text-blue-400 border-blue-500/20",
  },
  update: {
    label: "Update",
    className: "bg-amber-500/10 text-amber-400 border-amber-500/20",
  },
  delete: {
    label: "Delete",
    className: "bg-red-500/10 text-red-400 border-red-500/20",
  },
  sql: {
    label: "SQL",
    className: "bg-muted text-muted-foreground border-border",
  },
  other: {
    label: "Action",
    className: "bg-muted text-muted-foreground border-border",
  },
};

const scopeTabs: { value: Scope; label: string }[] = [
  { value: "all", label: "All" },
  { value: "tables", label: "Tables" },
  { value: "rows", label: "Rows" },
  { value: "sql", label: "SQL" },
];

const tableNamePattern = /'([^']+)'/;
const sqlPrefix = "Ran SQL: ";

function parseEvent(item: {
  id: number;
  message: string;
  time: string;
}): HistoryEvent {
  const lowerMessage = item.message.toLowerCase();
  const tableName = item.message.match(tableNamePattern)?.[1] ?? null;
  const base = { ...item, tableName, sqlText: null };

  if (item.message.startsWith(sqlPrefix)) {
    return {
      ...base,
      kind: "sql",
      sqlText: item.message.slice(sqlPrefix.length),
    };
  }
  if (lowerMessage.startsWith("created table")) {
    return { ...base, kind: "create" };
  }
  if (lowerMessage.startsWith("dropped table")) {
    return { ...base, kind: "delete" };
  }
  if (lowerMessage.startsWith("inserted")) {
    return { ...base, kind: "insert" };
  }
  if (lowerMessage.startsWith("updated")) {
    return { ...base, kind: "update" };
  }
  if (lowerMessage.startsWith("deleted")) {
    return { ...base, kind: "delete" };
  }
  return { ...base, kind: "other" };
}

function matchesScope(event: HistoryEvent, scope: Scope) {
  switch (scope) {
    case "all":
      return true;
    case "sql":
      return event.kind === "sql";
    case "rows":
      return (
        event.kind !== "sql" && event.message.toLowerCase().includes(" row ")
      );
    case "tables":
      return (
        event.kind !== "sql" && event.message.toLowerCase().includes(" table ")
      );
  }
}

function startOfDay(date: Date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

function dayLabel(date: Date) {
  const dayDifference = Math.round(
    (startOfDay(new Date()).getTime() - startOfDay(date).getTime()) / 86400000,
  );
  if (dayDifference === 0) return "Today";
  if (dayDifference === 1) return "Yesterday";
  return date.toLocaleDateString(undefined, {
    weekday: "long",
    month: "short",
    day: "numeric",
  });
}

function relativeTime(timeString: string) {
  const elapsedMs = Date.now() - new Date(timeString).getTime();
  const minutes = Math.floor(elapsedMs / 60000);
  const hours = Math.floor(elapsedMs / 3600000);
  const days = Math.floor(elapsedMs / 86400000);

  if (minutes < 1) return "Just now";
  if (minutes < 60) return `${minutes} min ago`;
  if (hours < 24) return `${hours} hr ago`;
  if (days < 7) return `${days} day${days > 1 ? "s" : ""} ago`;
  return new Date(timeString).toLocaleDateString();
}

function groupByDay(events: HistoryEvent[]) {
  const groups: { label: string; events: HistoryEvent[] }[] = [];
  for (const event of events) {
    const label = dayLabel(new Date(event.time));
    const lastGroup = groups[groups.length - 1];
    if (lastGroup && lastGroup.label === label) {
      lastGroup.events.push(event);
    } else {
      groups.push({ label, events: [event] });
    }
  }
  return groups;
}

function EventDescription({ event }: { event: HistoryEvent }) {
  if (event.kind === "sql") {
    return (
      <div className="min-w-0 flex-1">
        <p className="text-sm">Ran a SQL query</p>
        <code className="block truncate font-mono text-[12.5px] text-muted-foreground">
          {event.sqlText}
        </code>
      </div>
    );
  }

  if (!event.tableName) {
    return <p className="min-w-0 flex-1 text-sm">{event.message}</p>;
  }

  const [beforeTable, afterTable] = event.message.split(`'${event.tableName}'`);
  return (
    <p className="min-w-0 flex-1 break-words text-sm">
      {beforeTable}
      <code className="rounded border border-border bg-muted/40 px-1.5 py-px font-mono text-[12.5px]">
        {event.tableName}
      </code>
      {afterTable}
    </p>
  );
}

function EventRow({ event }: { event: HistoryEvent }) {
  const style = actionStyles[event.kind];
  return (
    <div className="flex flex-wrap items-center gap-x-4 gap-y-2 border-b border-border/50 px-4 py-3.5 last:border-b-0">
      <div className="w-[76px] shrink-0">
        <span
          className={cn(
            "inline-flex h-[22px] items-center rounded border px-2 text-xs font-medium",
            style.className,
          )}
        >
          {style.label}
        </span>
      </div>
      <EventDescription event={event} />
      <div className="flex shrink-0 items-center gap-3">
        <span className="text-[13px] text-muted-foreground">
          {relativeTime(event.time)}
        </span>
        <span className="min-w-6 text-right font-mono text-xs text-muted-foreground/70">
          #{event.id}
        </span>
      </div>
    </div>
  );
}

export function History() {
  const [searchParams, setSearchParams] = useSearchParams();
  const currentPage = parseInt(searchParams.get("page") || "1");

  const [searchText, setSearchText] = useState("");
  const [scope, setScope] = useState<Scope>("all");
  const [selectedTable, setSelectedTable] = useState(ALL_TABLES);

  const {
    data: history,
    isLoading,
    error: queryError,
  } = useQuery(listHistoryOptions({ query: { page: currentPage } }));

  const events = useMemo(() => (history ?? []).map(parseEvent), [history]);

  const tableNames = useMemo(
    () =>
      Array.from(
        new Set(events.map((event) => event.tableName).filter(Boolean)),
      ).sort() as string[],
    [events],
  );

  const visibleEvents = useMemo(() => {
    const query = searchText.trim().toLowerCase();
    return events.filter((event) => {
      if (!matchesScope(event, scope)) return false;
      if (selectedTable !== ALL_TABLES && event.tableName !== selectedTable) {
        return false;
      }
      return !query || event.message.toLowerCase().includes(query);
    });
  }, [events, scope, selectedTable, searchText]);

  const dayGroups = useMemo(() => groupByDay(visibleEvents), [visibleEvents]);

  const goToPage = (page: number) => {
    if (page < 1) return;
    setSearchParams({ page: page.toString() });
  };

  const hasNextPage = events.length > 0;
  const hasPreviousPage = currentPage > 1;

  return (
    <div className="flex-1 overflow-y-auto">
      <div className="mx-auto flex max-w-[1040px] flex-col gap-5 px-6 pb-12 pt-8">
        <div className="flex flex-col gap-1">
          <h1 className="text-[22px] font-semibold tracking-tight">History</h1>
          <p className="text-sm text-muted-foreground">
            Every change made to this database, newest first.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <label className="relative min-w-60 flex-1">
            <span className="sr-only">Search history</span>
            <Search className="pointer-events-none absolute left-3 top-1/2 h-[15px] w-[15px] -translate-y-1/2 text-muted-foreground" />
            <input
              value={searchText}
              onChange={(event) => setSearchText(event.target.value)}
              placeholder="Search by table or SQL"
              className="h-9 w-full rounded-md border border-border bg-background pl-[34px] pr-3 text-[13.5px] outline-none placeholder:text-muted-foreground focus-visible:ring-1 focus-visible:ring-ring"
            />
          </label>

          <div className="flex gap-0.5 rounded-md border border-border bg-background p-[3px]">
            {scopeTabs.map((tab) => (
              <button
                key={tab.value}
                type="button"
                onClick={() => setScope(tab.value)}
                className={cn(
                  "h-[22px] rounded px-3 text-[13px] font-medium leading-none transition-colors",
                  scope === tab.value
                    ? "bg-muted text-foreground"
                    : "text-muted-foreground hover:text-foreground",
                )}
              >
                {tab.label}
              </button>
            ))}
          </div>

          <Select value={selectedTable} onValueChange={setSelectedTable}>
            <SelectTrigger className="h-9 w-[155px] text-[13.5px]">
              <SelectValue placeholder="All tables" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={ALL_TABLES}>All tables</SelectItem>
              {tableNames.map((name) => (
                <SelectItem key={name} value={name}>
                  {name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <section className="overflow-hidden rounded-lg border border-border bg-card/40">
          {isLoading ? (
            <div className="flex h-48 items-center justify-center">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
            </div>
          ) : queryError ? (
            <p className="p-6 text-sm text-destructive">
              {queryError.detail || "Failed to fetch history"}
            </p>
          ) : dayGroups.length === 0 ? (
            <p className="p-10 text-center text-sm text-muted-foreground">
              No history matches your filters on this page.
            </p>
          ) : (
            dayGroups.map((group) => (
              <div key={group.label}>
                <div className="flex items-center justify-between border-b border-border bg-muted/20 px-4 py-2.5">
                  <span className="text-xs font-medium text-muted-foreground">
                    {group.label}
                  </span>
                  <span className="text-xs text-muted-foreground/80">
                    {group.events.length} event
                    {group.events.length === 1 ? "" : "s"}
                  </span>
                </div>
                {group.events.map((event) => (
                  <EventRow key={event.id} event={event} />
                ))}
              </div>
            ))
          )}

          <div className="flex flex-wrap items-center justify-between gap-3 px-4 py-2.5">
            <span className="text-[12.5px] text-muted-foreground">
              Page {currentPage}
            </span>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={!hasPreviousPage}
                onClick={() => goToPage(currentPage - 1)}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={!hasNextPage}
                onClick={() => goToPage(currentPage + 1)}
              >
                Next
              </Button>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
