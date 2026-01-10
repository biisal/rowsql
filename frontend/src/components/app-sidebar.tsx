import * as React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Table as TableIcon, Plus, Clock } from 'lucide-react';
import {
	Sidebar as ShadcnSidebar,
	SidebarContent,
	SidebarGroup,
	SidebarGroupContent,
	SidebarGroupLabel,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarFooter,
} from '@/components/ui/sidebar';
import { Skeleton } from './ui/skeleton';

interface Table {
	tableName: string;
	tableSchema: string;
}

interface SidebarProps extends React.ComponentProps<typeof ShadcnSidebar> {
	tables: Table[];
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
						<SidebarMenuButton asChild isActive={isHistory}>
							<Link to="/history">
								<Clock />
								<span>Recent Activity</span>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
					<SidebarMenuItem>
						<SidebarMenuButton
							asChild
							className="bg-primary text-primary-foreground hover:bg-primary/90 justify-center shadow-lg shadow-primary/20"
						>
							<Link to="/new-table">
								<Plus />
								<span>Add Table</span>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarFooter>
		</ShadcnSidebar>
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
