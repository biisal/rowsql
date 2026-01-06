import * as React from 'react';
import { useEffect, useState } from 'react';
import { useParams, Link, useSearchParams } from 'react-router-dom';
import { Plus } from 'lucide-react';
import axios from 'axios';
import api from '@/lib/axios';

import { Button } from '@/components/ui/button';
import { AppPagination } from '@/components/app-pagination';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { toast } from 'sonner';
import { DeletAlert } from '@/components/delete-alert';
import { RowOrderForm } from '@/components/row-order-form';
import { Rows } from '@/components/rows';
import type { TableData } from '@/lib/types';

export function TableRows() {
	const { tableName } = useParams<{ tableName: string }>();
	const [searchParams, setSearchParams] = useSearchParams();
	const [refresh, setRefesh] = useState(0);
	const [data, setData] = useState<TableData | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [selectedRows, setSelectedRows] = useState<Record<number, boolean>>({});
	const page = parseInt(searchParams.get('page') || '1');
	const col = searchParams.get('col');
	const order = searchParams.get('order') || 'asc';

	const fetchData = React.useCallback(async () => {
		if (!tableName) return;
		setError(null);

		try {
			const response = await api.get(
				`/table/${tableName}?page=${page}${col ? `&column=${col}&order=${order}` : ''}`,
			);
			if (response.data.success) {
				setData(response.data.data);
				setSelectedRows({});
			}
		} catch (err) {
			if (axios.isAxiosError(err)) {
				const { response } = err;
				const error =
					response?.data?.error || err.message || 'Something went wrong';
				toast.error(error);
				searchParams.delete('col');
				searchParams.delete('order');
				setSearchParams(searchParams);
				//  TODO : better error handling
			} else {
				setError(
					err instanceof Error ? err.message : 'An unknown error occurred',
				);
			}
		}
	}, [tableName, searchParams, setSearchParams, col, page, order, refresh]); //eslint-disable-line

	useEffect(() => {
		(async function () {
			fetchData();
		})();
	}, [fetchData]);

	const deleteRow = async (hash: string) => {
		try {
			const res = await api.delete(`/table/${tableName}/row/${hash}`);
			if (res.data.success) {
				toast.success('Row deleted successfully');
				setRefesh((r) => r + 1);
				return;
			}
			toast.error('Failed to delete row');
		} catch (err) {
			if (axios.isAxiosError(err)) {
				setError(
					err.response?.data?.error || err.message || 'Something went wrong',
				);
			} else {
				setError(
					err instanceof Error ? err.message : 'An unknown error occurred',
				);
			}
		}
	};

	// Reset search params when table name changes
	useEffect(() => {
		setSearchParams({}, { replace: true });
	}, [tableName]); // eslint-disable-line react-hooks/exhaustive-deps

	const toggleRowSelection = (index: number) => {
		setSelectedRows((prev) => ({
			...prev,
			[index]: !prev[index],
		}));
	};

	const toggleAllSelection = () => {
		if (!data) return;
		const allSelected =
			data.rows.length > 0 && data.rows.every((_, idx) => selectedRows[idx]);
		if (allSelected) {
			setSelectedRows({});
		} else {
			const newSelection: Record<number, boolean> = {};
			data.rows.forEach((_, idx) => {
				newSelection[idx] = true;
			});
			setSelectedRows(newSelection);
		}
	};

	if (error) {
		return (
			<div className="flex items-center justify-center h-full">
				<div className="text-destructive font-medium">Error: {error}</div>
			</div>
		);
	}

	if (!data) {
		return <div className="p-8">No data found.</div>;
	}

	const isAllSelected =
		data.rows.length > 0 && data.rows.every((_, idx) => selectedRows[idx]);
	const isSomeSelected =
		data.rows.some((_, idx) => selectedRows[idx]) && !isAllSelected;

	return (
		<div className="flex-1 p-4 md:p-8 overflow-y-auto w-full max-w-full">
			<div className="space-y-6">
				<div className="flex items-center justify-between">
					<h1 className="text-3xl font-bold tracking-tight">
						{data.activeTable}
					</h1>
					<div className="flex items-center justify-center gap-1">
						<DeletAlert tableName={data.activeTable} />
						<Link to={`/table/${data.activeTable}/form`}>
							<Button className="shadow-lg shadow-primary/20">
								<Plus className="mr-2 h-4 w-4" /> Insert Record
							</Button>
						</Link>
					</div>
				</div>

				<Card className="border-border/50 bg-card/50 backdrop-blur-sm">
					<CardHeader className="px-6 py-4 flex justify-between items-center flex-wrap border-b border-border/50">
						<CardTitle className="text-lg font-medium">Table Data</CardTitle>
						<RowOrderForm
							cols={data.cols}
							initialValue={{ col, order }}
							setUrlParams={setSearchParams}
						/>
					</CardHeader>
					<CardContent className="p-4">
						(
						<Rows
							data={data}
							selectedRows={selectedRows}
							isAllSelected={isAllSelected}
							isSomeSelected={isSomeSelected}
							toggleAllSelection={toggleAllSelection}
							toggleRowSelection={toggleRowSelection}
							deleteRow={deleteRow}
						/>
						)
						<div className="flex items-center flex-col sm:flex-row justify-center md:justify-between  py-4">
							<div className="text-muted-foreground flex-1  text-sm">
								{Object.values(selectedRows).filter(Boolean).length} of{' '}
								{data.rows.length} row(s) selected.
							</div>
							<AppPagination
								currentPage={page}
								onPageChange={(newPage) =>
									setSearchParams({ page: String(newPage) })
								}
								hasNextPage={data.hasNextPage}
								hasPreviousPage={page > 1}
							/>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
