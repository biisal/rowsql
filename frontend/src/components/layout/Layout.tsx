import { Outlet } from 'react-router-dom';
import {
	SidebarProvider,
	SidebarTrigger,
	SidebarInset,
} from '@/components/ui/sidebar';
import { AppSidebar } from './AppSidebar';
import { Separator } from '@/components/ui/separator';
import {
	Breadcrumb,
	BreadcrumbItem,
	BreadcrumbList,
	BreadcrumbPage,
} from '@/components/ui/breadcrumb';
import { useQuery, useMutationState } from '@tanstack/react-query';
import { listTablesOptions } from '@/client/@tanstack/react-query.gen';
import { Toaster } from 'sonner';
import { GitStarButton } from './GitStarButton';

export function Layout() {
	const { data: tables, isLoading: tablesRefreshing } = useQuery(listTablesOptions());

	const pendingMutations = useMutationState({
		filters: { status: 'pending', mutationKey: ['createNewTable'] },
		select: (mutation) => mutation.state.status,
	});

	const tableAppending = pendingMutations.length > 0;

	return (
		<SidebarProvider>
			<AppSidebar
				refreshing={tablesRefreshing}
				isAppending={tableAppending}
				tables={tables || []}
			/>
			<SidebarInset className="min-w-0 overflow-hidden">
				<header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
					<SidebarTrigger className="-ml-1" />
					<Separator orientation="vertical" className="mr-2 h-4" />
					<Breadcrumb className="w-full flex justify-between">
						<BreadcrumbList>
							<BreadcrumbItem>
								<BreadcrumbPage>RowSQL</BreadcrumbPage>
							</BreadcrumbItem>
						</BreadcrumbList>
						<GitStarButton />
					</Breadcrumb>
				</header>
				<div className="flex flex-1 flex-col gap-4 min-w-0 overflow-hidden">
					<Outlet />
				</div>
				<Toaster />
			</SidebarInset>
		</SidebarProvider>
	);
}
