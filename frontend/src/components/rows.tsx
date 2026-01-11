import React from 'react';
import { Link } from 'react-router-dom';
import { MoreHorizontal } from 'lucide-react';
import { Checkbox } from './ui/checkbox';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from './ui/dropdown-menu';
import { Button } from './ui/button';
import { Skeleton } from './ui/skeleton';
import { cn } from '@/lib/utils';
import type { CellValue, TableData } from '@/lib/types';

interface RowsProps {
	data: TableData;
	selectedRows: Record<number, boolean>;
	isAllSelected: boolean;
	isSomeSelected: boolean;
	toggleAllSelection: () => void;
	toggleRowSelection: (index: number) => void;
	deleteRow: (hash: string) => void;
}

export const Rows = ({
	data,
	selectedRows,
	isAllSelected,
	isSomeSelected,
	toggleAllSelection,
	toggleRowSelection,
	deleteRow,
}: RowsProps) => {
	return (
		<div className="rounded-md border overflow-auto  relative">
			<div
				className="grid min-w-full"
				style={{
					gridTemplateColumns: `40px ${data.cols.map(() => 'minmax(150px, 1fr)').join(' ')} 80px`,
				}}
			>
				{/* Checkbox Header */}
				<div className="h-10 px-2 flex items-center justify-center border-b bg-muted/50 sticky top-0 z-20">
					<Checkbox
						checked={
							isAllSelected ? true : isSomeSelected ? 'indeterminate' : false
						}
						onCheckedChange={toggleAllSelection}
						aria-label="Select all"
					/>
				</div>

				{/* Data Columns Headers */}
				{data.cols.map((col) => (
					<div
						key={col.columnName}
						className="h-10 px-2 text-left align-middle font-medium text-muted-foreground flex items-center border-b bg-muted/50 sticky top-0 z-20"
					>
						{col.columnName}
					</div>
				))}

				{/* Actions Header */}
				<div className="h-10 px-2 text-right align-middle font-medium text-muted-foreground flex items-center justify-end border-b bg-muted/50 sticky top-0 right-0 z-30 shadow-[-1px_0_0_0_var(--border)]">
					Actions
				</div>

				{/* Body Rows */}
				{data.rows.length > 0 ? (
					data.rows.map((row, rowIndex) => {
						const hash = String(row[0] ?? '');
						const isSelected = !!selectedRows[rowIndex];

						return (
							<React.Fragment key={rowIndex}>
								{/* Checkbox Cell */}
								<div
									className={cn(
										'p-2 flex items-center justify-center border-b bg-card group-hover:bg-muted/50 transition-colors',
										isSelected && 'bg-muted',
									)}
								>
									<Checkbox
										checked={isSelected}
										onCheckedChange={() => toggleRowSelection(rowIndex)}
										aria-label="Select row"
									/>
								</div>

								{row.slice(1).map((cell: CellValue, cellIndex: number) => (
									<div
										key={cellIndex}
										className={cn(
											'p-2 align-middle flex items-center border-b bg-card group-hover:bg-muted/50 transition-colors lowercase',
											isSelected && 'bg-muted',
										)}
									>
										{cell === null ? 'NULL' : String(cell)}
									</div>
								))}

								{/* Actions Cell */}
								<div
									className={cn(
										'p-2 align-middle flex items-center justify-end border-b bg-card group-hover:bg-muted/50 transition-colors sticky right-0 z-10 shadow-[-1px_0_0_0_var(--border)]',
										isSelected && 'bg-muted',
									)}
								>
									<DropdownMenu>
										<DropdownMenuTrigger asChild>
											<Button variant="ghost" className="h-8 w-8 p-0">
												<span className="sr-only">Open menu</span>
												<MoreHorizontal className="h-4 w-4" />
											</Button>
										</DropdownMenuTrigger>
										<DropdownMenuContent align="end" className="px-4 py-2">
											<DropdownMenuLabel>Actions</DropdownMenuLabel>
											<DropdownMenuItem asChild>
												<Link
													to={`/tables/${data.activeTable}/rows?hash=${hash}&page=${data.page}`}
												>
													Edit Row
												</Link>
											</DropdownMenuItem>
											<DropdownMenuSeparator />
											<Button
												variant="destructive"
												className="w-full"
												onClick={() => deleteRow(hash)}
											>
												Delete Row
											</Button>
										</DropdownMenuContent>
									</DropdownMenu>
								</div>
							</React.Fragment>
						);
					})
				) : (
					<div className="col-span-full h-24 flex items-center justify-center">
						No results.
					</div>
				)}
			</div>
		</div>
	);
};

interface RowsSkeletonProps {
	columns?: number;
	rows?: number;
}

export const RowsSkeleton = ({ columns = 5, rows = 10 }: RowsSkeletonProps) => {
	return (
		<div className="rounded-md border overflow-auto relative">
			<div
				className="grid min-w-full"
				style={{
					gridTemplateColumns: `40px ${Array(columns).fill('minmax(150px, 1fr)').join(' ')} 80px`,
				}}
			>
				{/* Checkbox Header Skeleton */}
				<div className="h-10 px-2 flex items-center justify-center border-b bg-muted/50 sticky top-0 z-20">
					<Skeleton className="h-4 w-4 rounded" />
				</div>

				{/* Column Headers Skeleton */}
				{Array.from({ length: columns }).map((_, index) => (
					<div
						key={index}
						className="h-10 px-2 flex items-center border-b bg-muted/50 sticky top-0 z-20"
					>
						<Skeleton className="h-4 w-24" />
					</div>
				))}

				{/* Actions Header Skeleton */}
				<div className="h-10 px-2 flex items-center justify-end border-b bg-muted/50 sticky top-0 right-0 z-30 shadow-[-1px_0_0_0_var(--border)]">
					<Skeleton className="h-4 w-12" />
				</div>

				{/* Body Rows Skeleton */}
				{Array.from({ length: rows }).map((_, rowIndex) => (
					<React.Fragment key={rowIndex}>
						{/* Checkbox Cell Skeleton */}
						<div className="p-2 flex items-center justify-center border-b bg-card">
							<Skeleton className="h-4 w-4 rounded" />
						</div>

						{/* Data Cells Skeleton */}
						{Array.from({ length: columns }).map((_, cellIndex) => (
							<div
								key={cellIndex}
								className="p-2 flex items-center border-b bg-card"
							>
								<Skeleton
									className="h-4"
									style={{
										width: `${Math.floor(Math.random() * 40) + 50}%`,
									}}
								/>
							</div>
						))}

						{/* Actions Cell Skeleton */}
						<div className="p-2 flex items-center justify-end border-b bg-card sticky right-0 z-10 shadow-[-1px_0_0_0_var(--border)]">
							<Skeleton className="h-8 w-8 rounded" />
						</div>
					</React.Fragment>
				))}
			</div>
		</div>
	);
};
