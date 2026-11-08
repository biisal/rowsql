import { ScrollArea } from "@/components/ui/scroll-area";

export interface RecordField {
  name: string;
  value: unknown;
}

const formatValue = (value: unknown) => {
  if (typeof value === "string") {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  }
  if (typeof value === "object") return JSON.stringify(value, null, 2);
  return String(value);
};

export const RecordFields = ({ fields }: { fields: RecordField[] }) => (
  <ScrollArea className="min-h-0 flex-1">
    <div className="flex flex-col divide-y p-4">
      {fields.map((field) => (
        <div key={field.name} className="flex flex-col gap-1 py-3">
          <label className="text-xs font-medium text-muted-foreground capitalize">
            {field.name}
          </label>
          <div className="font-mono text-sm break-all whitespace-pre-wrap">
            {field.value === undefined ? (
              <p className="text-muted-foreground">-</p>
            ) : field.value === null ? (
              <p className="text-muted-foreground">null</p>
            ) : (
              formatValue(field.value)
            )}
          </div>
        </div>
      ))}
    </div>
  </ScrollArea>
);
