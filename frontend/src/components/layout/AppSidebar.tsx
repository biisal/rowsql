import { ArrowRight, Clock, Plus, SquareCode, Table as TableIcon } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import type { ListTablesRow } from '@/client/types.gen';
import {
	Sidebar as ShadcnSidebar,
	SidebarContent,
	SidebarFooter,
	SidebarGroup,
	SidebarGroupContent,
	SidebarGroupLabel,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
} from '@/components/ui/sidebar';
import { Skeleton } from '@/components/ui/skeleton';

interface SidebarProps extends React.ComponentProps<typeof ShadcnSidebar> {
	tables: ListTablesRow[];
	refreshing: boolean;
	isAppending: boolean;
}

export function AppSidebar({
	tables,
	refreshing: loading,
	isAppending,
	...props
}: SidebarProps) {
	const location = useLocation();

	const isActiveTable = (tableName: string) => {
		return location.pathname === `/tables/${tableName}`;
	};

	const isHome = location.pathname === '/';
	const isHistory = location.pathname === '/history';

	return (
		<ShadcnSidebar {...props}>
			<SidebarHeader>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton size="lg" asChild isActive={isHome}>
							<Link to="/" className="flex items-center gap-2">
								<img src="/logo.png" alt="Logo" className="w-12 h-12" />
								<span className="text-lg font-bold uppercase tracking-widest">
									RowSQL
								</span>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarHeader>

			<SidebarMenu className='p-2'>
				<SidebarMenuItem>
					<SidebarMenuButton
						asChild
						className="bg-primary hover:bg-primary/90 active:bg-primary/50"
					>
						<Link to="/editor">
							<SquareCode className="h-4 w-4" />
							<span>SQL Editor</span>
						</Link>
					</SidebarMenuButton>
				</SidebarMenuItem>
			</SidebarMenu>

			<SidebarContent>
				<SidebarGroup>
					<SidebarGroupLabel>Tables</SidebarGroupLabel>
					<SidebarGroupContent>
						<SidebarMenu>
							{isAppending && <SkeletonTables count={1} />}
							{loading ? (
								<SkeletonTables count={3} />
							) : (
								Array.isArray(tables) &&
								tables.map((table) => (
									<SidebarMenuItem key={table.tableName}>
										<SidebarMenuButton
											asChild
											isActive={isActiveTable(table.tableName)}
										>
											<Link to={`/tables/${table.tableName}`}>
												<TableIcon />
												<span>{table.tableName}</span>
											</Link>
										</SidebarMenuButton>
									</SidebarMenuItem>
								))
							)}
						</SidebarMenu>
					</SidebarGroupContent>
				</SidebarGroup>
			</SidebarContent>

			<SidebarFooter>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton asChild isActive={isHistory} className='bg-muted-foreground/5'>
							<Link to="/history" className="flex items-center justify-between w-full">
								<div className="flex items-center gap-2">
									<Clock className="h-4 w-4" />
									<span>Recent Activity</span>
								</div>
								<ArrowRight />
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
					<SidebarMenuItem>
						<SidebarMenuButton
							asChild
							className="flex justify-center border border-border bg-primary hover:bg-primary/90 "
						>
							<Link to="/new-table">
								<Plus />
								<span>Add Table</span>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>

				</SidebarMenu>
			</SidebarFooter>
		</ShadcnSidebar >
	);
}

function SkeletonTables({ count }: { count: number }) {
	return new Array(count).fill(null).map((_, index) => (
		<SidebarMenuItem key={index}>
			<SidebarMenuButton asChild>
				<div>
					<TableIcon />
					<Skeleton className="w-full h-4 rounded-md" />
				</div>
			</SidebarMenuButton>
		</SidebarMenuItem>
	));
}
