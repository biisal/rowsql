import type { ReactNode } from "react";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

interface RecordSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  children: ReactNode;
}

export const RecordSheet = ({
  open,
  onOpenChange,
  title,
  description,
  children,
}: RecordSheetProps) => (
  <Sheet open={open} onOpenChange={onOpenChange}>
    <SheetContent className="flex min-w-[90%] flex-col md:min-w-xl">
      <SheetHeader className="bg-muted-foreground/10">
        <SheetTitle>{title}</SheetTitle>
        <SheetDescription>{description}</SheetDescription>
      </SheetHeader>
      {children}
    </SheetContent>
  </Sheet>
);
