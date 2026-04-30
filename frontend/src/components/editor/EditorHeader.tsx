import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

interface EditorHeaderProps {
    runText: "Run" | "Run Selected";
    loading: boolean;
    onRun: () => void;
}

export const EditorHeader = ({ runText, loading, onRun }: EditorHeaderProps) => {

    const isMac = navigator.userAgent.includes("Macintosh");
    return (< div className="flex items-center justify-between border-b border-border bg-card px-4 py-2" >
        <span className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            SQL Editor
        </span>
        <Button
            size="sm"
            onClick={onRun}
            disabled={loading}
            className=" gap-1.5  px-3 text-xs"
        >
            {runText}
            <Separator orientation="vertical" className="h-full bg-muted-foreground/50" />
            <span className="">
                {isMac ? "Cmd" : "Ctrl"}
            </span>
            <span>+</span>
            Enter
        </Button>
    </div >

    )

} 