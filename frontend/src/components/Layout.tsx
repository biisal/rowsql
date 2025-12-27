import { useEffect } from "react";
import { Outlet } from "react-router-dom";
import {
  SidebarProvider,
  SidebarTrigger,
  SidebarInset,
} from "@/components/ui/sidebar";
import { Sidebar } from "./Sidebar";
import { Separator } from "@/components/ui/separator";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Toaster } from "sonner";
import useTableStore from "@/lib/store/use-table";

export function Layout() {
  const { tables, refreshTables, tablesRefreshing, tableAppending } =
    useTableStore();

  useEffect(() => {
    refreshTables();
  }, [refreshTables]);

  return (
    <SidebarProvider>
      <Sidebar
        refreshing={tablesRefreshing}
        isAppending={tableAppending}
        tables={tables}
      />
      <SidebarInset className="min-w-0 overflow-hidden">
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbPage>RowSQL</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4 min-w-0 overflow-hidden">
          <Outlet />
        </div>
        <Toaster />
      </SidebarInset>
    </SidebarProvider>
  );
}
