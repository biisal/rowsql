import * as ResizablePrimitive from "react-resizable-panels";

import { cn } from "@/lib/utils";

function ResizablePanelGroup({
  className,
  ...props
}: ResizablePrimitive.GroupProps) {
  return (
    <ResizablePrimitive.Group
      data-slot="resizable-panel-group"
      className={cn(
        "flex h-full w-full aria-[orientation=vertical]:flex-col",
        className,
      )}
      {...props}
    />
  );
}

function ResizablePanel({ ...props }: ResizablePrimitive.PanelProps) {
  return <ResizablePrimitive.Panel data-slot="resizable-panel" {...props} />;
}

function ResizableHandle({
  withHandle,
  className,
  ...props
}: ResizablePrimitive.SeparatorProps & {
  withHandle?: boolean;
}) {
  return (
    <ResizablePrimitive.Separator
      data-slot="resizable-handle"
      className={cn(
        "relative flex items-center justify-center bg-transparent",
        "aria-[orientation=horizontal]:h-1.5 aria-[orientation=horizontal]:w-full",
        "aria-[orientation=vertical]:w-1.5 aria-[orientation=vertical]:h-full",
        "focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring",
        className,
      )}
      {...props}
    >
      {withHandle && (
        <div className="z-10 rounded-full bg-zinc-600 aria-[orientation=horizontal]:h-[3px] aria-[orientation=horizontal]:w-8 [aria-orientation=vertical]&:h-8 [aria-orientation=vertical]&:w-[3px] h-[3px] w-8" />
      )}
    </ResizablePrimitive.Separator>
  );
}

export { ResizableHandle, ResizablePanel, ResizablePanelGroup };
