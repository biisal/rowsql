import * as React from 'react';
import {
	Pagination,
	PaginationContent,
	PaginationItem,
	PaginationLink,
	PaginationPrevious,
	PaginationNext,
	PaginationEllipsis,
} from '@/components/ui/pagination';
import { cn } from '@/lib/utils';
import { Input } from './ui/input';
import { Button } from './ui/button';

interface AppPaginationProps {
	currentPage: number;
	onPageChange: (page: number) => void;
	hasNextPage: boolean;
	hasPreviousPage?: boolean;
	totalPages?: number;
	showPageNumbers?: boolean;
	className?: string;
}

export function AppPagination({
	currentPage,
	onPageChange,
	hasNextPage,
	hasPreviousPage = currentPage > 1,
	totalPages,
	showPageNumbers = true,
	className,
}: AppPaginationProps) {
	const handlePageChange = (page: number, e: React.MouseEvent) => {
		e.preventDefault();
		if (page > 0 && page <= (totalPages || Infinity)) {
			onPageChange(page);
		}
	};

	const renderPageNumbers = () => {
		if (!showPageNumbers) return null;

		const pages: React.ReactNode[] = [];

		if (currentPage > 2) {
			pages.push(
				<PaginationItem key="page-1">
					<PaginationLink href="#" onClick={(e) => handlePageChange(1, e)}>
						1
					</PaginationLink>
				</PaginationItem>,
			);

			if (currentPage > 3) {
				pages.push(
					<PaginationItem key="ellipsis-start">
						<PaginationEllipsis />
					</PaginationItem>,
				);
			}
		}

		if (currentPage > 1) {
			pages.push(
				<PaginationItem key={`page-${currentPage - 1}`}>
					<PaginationLink
						href="#"
						onClick={(e) => handlePageChange(currentPage - 1, e)}
					>
						{currentPage - 1}
					</PaginationLink>
				</PaginationItem>,
			);
		}

		pages.push(
			<PaginationItem key={`page-${currentPage}`}>
				<PaginationLink href="#" isActive onClick={(e) => e.preventDefault()}>
					{currentPage}
				</PaginationLink>
			</PaginationItem>,
		);

		if (hasNextPage) {
			pages.push(
				<PaginationItem key={`page-${currentPage + 1}`}>
					<PaginationLink
						href="#"
						onClick={(e) => handlePageChange(currentPage + 1, e)}
					>
						{currentPage + 1}
					</PaginationLink>
				</PaginationItem>,
			);

			if (!totalPages || currentPage + 1 < totalPages) {
				pages.push(
					<PaginationItem key="ellipsis-end">
						<PaginationEllipsis />
					</PaginationItem>,
				);
			}
		}

		if (totalPages && currentPage < totalPages - 1 && hasNextPage) {
			pages.push(
				<PaginationItem key={`page-${totalPages}`}>
					<PaginationLink
						href="#"
						onClick={(e) => handlePageChange(totalPages, e)}
					>
						{totalPages}
					</PaginationLink>
				</PaginationItem>,
			);
		}

		return pages;
	};

	const [jumpPage, setJumpPage] = React.useState('');

	return (
		<Pagination className={className}>
			<PaginationContent className="flex-wrap items-center justify-center">
				<PaginationItem>
					<PaginationPrevious
						href="#"
						onClick={(e) => {
							e.preventDefault();
							if (hasPreviousPage) {
								onPageChange(currentPage - 1);
							}
						}}
						className={cn(!hasPreviousPage && 'pointer-events-none opacity-50')}
					/>
				</PaginationItem>

				{renderPageNumbers()}

				<PaginationItem>
					<PaginationNext
						href="#"
						onClick={(e) => {
							e.preventDefault();
							if (hasNextPage) {
								onPageChange(currentPage + 1);
							}
						}}
						className={cn(!hasNextPage && 'pointer-events-none opacity-50')}
					/>
				</PaginationItem>
				<PaginationItem>
					<div className="flex items-center space-x-2">
						<Input
							type="number"
							placeholder="Page"
							value={jumpPage}
							onChange={(e) => setJumpPage(e.target.value)}
							onKeyDown={(e) => {
								if (e.key === 'Enter') {
									const page = parseInt(jumpPage);
									if (page > 0 && page <= (totalPages || Infinity)) {
										onPageChange(page);
										setJumpPage('');
									}
								}
							}}
							className="w-24"
						/>
						<Button
							onClick={() => {
								const page = parseInt(jumpPage);
								if (page > 0 && page <= (totalPages || Infinity)) {
									onPageChange(page);
									setJumpPage('');
								}
							}}
							size="sm"
						>
							Go
						</Button>
					</div>
				</PaginationItem>
			</PaginationContent>
		</Pagination>
	);
}
