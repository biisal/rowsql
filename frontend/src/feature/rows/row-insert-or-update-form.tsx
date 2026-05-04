import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";

export const IowInsertOrUpdateForm = () => {
  return (
    <Sheet>
      <SheetContent className="flex min-w-[90%] flex-col md:min-w-150">
        <SheetHeader className="bg-muted-foreground/10">
          <SheetTitle>{"Sheet Title"}</SheetTitle>
          <SheetDescription>Viewing record details.</SheetDescription>
        </SheetHeader>
      </SheetContent>
    </Sheet>
  );
};
