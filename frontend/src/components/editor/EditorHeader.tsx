import { Button } from "@/components/ui/button";
import { Play } from "lucide-react";

interface EditorHeaderProps {
    runText: "Run" | "Run Selected";
    loading: boolean;
    onRun: () => void;
}

export const EditorHeader = ({ runText, loading, onRun }: EditorHeaderProps) => (
    <div className="flex items-center justify-between border-b border-border bg-card px-4 py-2">
        <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            SQL Editor
        </span>
        <Button
            size="sm"
            onClick={onRun}
            disabled={loading}
            className="h-7 gap-1.5  px-3 text-xs"
        >
            <Play className="h-3 w-3 fill-current" />
            {loading ? "Running…" : runText} {" "} <p className="text-muted-foreground/80">(Ctrl+Enter)</p>
        </Button>
    </div>
);

