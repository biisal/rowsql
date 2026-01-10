import {
	Database,
	Plus,
	Book,
	Activity,
	Server,
	Table as TableIcon,
	Clock,
	Loader2,
	ChevronRight,
	History as HistoryIcon,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { useEffect, useState } from 'react';
import api from '@/lib/axios';
import useTableStore from '@/lib/store/use-table';

interface HistoryItem {
	id: number;
	message: string;
	time: string;
}

export function Home() {
	const { tables, tablesRefreshing } = useTableStore();
	const [recentHistory, setRecentHistory] = useState<HistoryItem[]>([]);
	const [historyLoading, setHistoryLoading] = useState(true);

	useEffect(() => {
		fetchRecentHistory();
	}, []);

	const fetchRecentHistory = async () => {
		try {
			const response = await api.get('/history/recent');
			setRecentHistory(response.data.data || []);
		} catch (err) {
			console.error('Failed to fetch recent history:', err);
		} finally {
			setHistoryLoading(false);
		}
	};

	const formatTime = (timeString: string) => {
		const date = new Date(timeString);
		const now = new Date();
		const diffMs = now.getTime() - date.getTime();
		const diffMins = Math.floor(diffMs / 60000);
		const diffHours = Math.floor(diffMs / 3600000);
		const diffDays = Math.floor(diffMs / 86400000);

		if (diffMins < 1) return 'Just now';
		if (diffMins < 60) return `${diffMins}m ago`;
		if (diffHours < 24) return `${diffHours}h ago`;
		if (diffDays < 7) return `${diffDays}d ago`;

		return date.toLocaleDateString();
	};

	const getActionBadge = (message: string) => {
		if (message.toLowerCase().includes('created table')) {
			return (
				<Badge
					variant="secondary"
					className="text-xs bg-green-500/10 text-green-500 border-green-500/20"
				>
					Create
				</Badge>
			);
		}
		if (
			message.toLowerCase().includes('dropped table') ||
			message.toLowerCase().includes('deleted')
		) {
			return (
				<Badge
					variant="secondary"
					className="text-xs bg-red-500/10 text-red-500 border-red-500/20"
				>
					Delete
				</Badge>
			);
		}
		if (message.toLowerCase().includes('updated')) {
			return (
				<Badge
					variant="secondary"
					className="text-xs bg-blue-500/10 text-blue-500 border-blue-500/20"
				>
					Update
				</Badge>
			);
		}
		if (message.toLowerCase().includes('inserted')) {
			return (
				<Badge
					variant="secondary"
					className="text-xs bg-purple-500/10 text-purple-500 border-purple-500/20"
				>
					Insert
				</Badge>
			);
		}
		return (
			<Badge variant="secondary" className="text-xs">
				Action
			</Badge>
		);
	};

	return (
		<div className="flex-1 p-8 overflow-y-auto">
			<div className="max-w-6xl mx-auto space-y-8 animate-fade-in-up">
				<div className="relative overflow-hidden rounded-3xl bg-linear-to-br from-primary/10 via-background to-background border border-border/50 p-8 md:p-12">
					<div className="relative z-10 max-w-2xl space-y-6">
						<div className="flex items-center flex-wrap">
							<img src="/logo.png" className="h-20 w-20" />
							<h1 className="text-4xl md:text-5xl font-bold tracking-tight text-foreground">
								Welcome to{' '}
								<span className="bg-linear-to-br  to-primary via-primary from-secondary-foreground bg-clip-text text-transparent">
									RowSQL
								</span>
							</h1>
						</div>
						<p className="text-lg text-muted-foreground leading-relaxed">
							Your modern, powerful interface for database management. Monitor
							performance, manage schemas, and query data with ease.
						</p>
						<div className="flex flex-wrap gap-4 pt-2">
							<Link to="/tables/create/new">
								<button className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg bg-primary text-primary-foreground font-medium hover:bg-primary/90 transition-all shadow-lg shadow-primary/20">
									<Plus className="w-5 h-5" /> Create New Table
								</button>
							</Link>
							<Link to="/docs">
								<button className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg bg-card border border-border hover:bg-muted/50 transition-all text-foreground font-medium">
									<Book className="w-5 h-5" /> Documentation
								</button>
							</Link>
						</div>
					</div>
					<div className="absolute right-0 top-0 -mt-20 -mr-20 w-96 h-96 bg-primary/5 rounded-full blur-3xl"></div>
					<div className="absolute bottom-0 right-20 w-64 h-64 bg-secondary/5 rounded-full blur-3xl"></div>
				</div>

				<div className="grid grid-cols-1 md:grid-cols-3 gap-6">
					<Card className="bg-card/50 backdrop-blur-sm border-border/50 hover:border-primary/50 transition-all">
						<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
							<CardTitle className="text-sm font-medium text-muted-foreground">
								Total Tables
							</CardTitle>
							<TableIcon className="h-4 w-4 text-primary" />
						</CardHeader>
						<CardContent>
							<div className="text-2xl font-bold">
								{tablesRefreshing ? (
									<Loader2 className="w-6 h-6 animate-spin" />
								) : (
									tables.length
								)}
							</div>
							<p className="text-xs text-muted-foreground mt-1">
								Active database tables
							</p>
						</CardContent>
					</Card>
					<Card className="bg-card/50 backdrop-blur-sm border-border/50 hover:border-primary/50 transition-all">
						<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
							<CardTitle className="text-sm font-medium text-muted-foreground">
								Recent Operations
							</CardTitle>
							<Activity className="h-4 w-4 text-primary" />
						</CardHeader>
						<CardContent>
							<div className="text-2xl font-bold">
								{historyLoading ? (
									<Loader2 className="w-6 h-6 animate-spin" />
								) : (
									recentHistory.length
								)}
							</div>
							<p className="text-xs text-muted-foreground mt-1">
								Last 10 database operations
							</p>
						</CardContent>
					</Card>
					<Card className="bg-card/50 backdrop-blur-sm border-border/50 hover:border-primary/50 transition-all">
						<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
							<CardTitle className="text-sm font-medium text-muted-foreground">
								Database Type
							</CardTitle>
							<Server className="h-4 w-4 text-primary" />
						</CardHeader>
						<CardContent>
							<div className="text-2xl font-bold">
								<Badge
									variant="outline"
									className="text-primary border-primary/20 bg-primary/10"
								>
									Active
								</Badge>
							</div>
							<p className="text-xs text-muted-foreground mt-1">
								Multi-database support
							</p>
						</CardContent>
					</Card>
				</div>

				<div className="grid grid-cols-1 md:grid-cols-2 gap-6">
					<Card className="border-border/50">
						<CardHeader>
							<CardTitle>Quick Actions</CardTitle>
							<CardDescription>
								Common tasks you might want to perform.
							</CardDescription>
						</CardHeader>
						<CardContent className="grid gap-4">
							<Link to="/tables/create/new">
								<div className="flex items-center justify-between p-4 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors cursor-pointer group">
									<div className="flex items-center gap-4">
										<div className="p-2 rounded-md bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
											<Plus className="w-5 h-5" />
										</div>
										<div>
											<h4 className="font-medium">Create New Table</h4>
											<p className="text-sm text-muted-foreground">
												Design and create a new database table
											</p>
										</div>
									</div>
									<ChevronRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
								</div>
							</Link>

							<Link to="/history">
								<div className="flex items-center justify-between p-4 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors cursor-pointer group">
									<div className="flex items-center gap-4">
										<div className="p-2 rounded-md bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
											<HistoryIcon className="w-5 h-5" />
										</div>
										<div>
											<h4 className="font-medium">View Full History</h4>
											<p className="text-sm text-muted-foreground">
												Browse complete database operation logs
											</p>
										</div>
									</div>
									<ChevronRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
								</div>
							</Link>

							{!tablesRefreshing && tables.length > 0 && (
								<Link to={`/tables/${tables[0].tableName}`}>
									<div className="flex items-center justify-between p-4 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors cursor-pointer group">
										<div className="flex items-center gap-4">
											<div className="p-2 rounded-md bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors">
												<Database className="w-5 h-5" />
											</div>
											<div>
												<h4 className="font-medium">Browse Tables</h4>
												<p className="text-sm text-muted-foreground">
													View and manage table data
												</p>
											</div>
										</div>
										<ChevronRight className="w-4 h-4 text-muted-foreground group-hover:text-primary transition-colors" />
									</div>
								</Link>
							)}
						</CardContent>
					</Card>

					<Card className="border-border/50">
						<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
							<div>
								<CardTitle className="flex items-center gap-2">
									<Clock className="w-5 h-5 text-primary" />
									Recent Activity
								</CardTitle>
								<CardDescription>
									Latest changes in your database.
								</CardDescription>
							</div>
							<Link to="/history">
								<Button variant="ghost" size="sm" className="text-xs">
									View All
									<ChevronRight className="w-3 h-3 ml-1" />
								</Button>
							</Link>
						</CardHeader>
						<CardContent>
							{historyLoading ? (
								<div className="flex items-center justify-center py-8">
									<Loader2 className="w-6 h-6 animate-spin text-primary" />
								</div>
							) : recentHistory.length === 0 ? (
								<div className="text-center py-8">
									<p className="text-sm text-muted-foreground">
										No recent activity
									</p>
									<p className="text-xs text-muted-foreground/70 mt-1">
										Start performing database operations
									</p>
								</div>
							) : (
								<div className="space-y-3">
									{recentHistory.map((item) => (
										<div
											key={item.id}
											className="flex items-center gap-3 pb-3 border-b border-border/50 last:border-0 last:pb-0"
										>
											<div className="w-2 h-2 rounded-full bg-primary shrink-0"></div>
											<div className="flex-1 min-w-0">
												<p className="text-sm font-medium truncate">
													{item.message}
												</p>
												<p className="text-xs text-muted-foreground">
													{formatTime(item.time)}
												</p>
											</div>
											{getActionBadge(item.message)}
										</div>
									))}
								</div>
							)}
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
}
