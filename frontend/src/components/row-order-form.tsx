import { useEffect, useState } from 'react';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
} from './ui/select';
import type { Column } from '@/lib/types';
import { type SetURLSearchParams } from 'react-router-dom';

interface InitialVal {
	col: string | null;
	order: string | null;
}

interface RowOrderFormProps {
	cols: Column[];
	setUrlParams: SetURLSearchParams;
	initialValue?: InitialVal;
}

export const RowOrderForm = ({
	cols,
	setUrlParams,
	initialValue,
}: RowOrderFormProps) => {
	const [selectedCol, setSelectedCol] = useState<string | null>(
		(initialValue && initialValue.col) ||
			(cols && cols.length > 0 ? cols[0].columnName : null),
	);
	const [selectedOrder, setSelectedOrder] = useState<string | null>(
		(initialValue && initialValue.order) || 'asc',
	);

	useEffect(() => {
		setUrlParams((prev) => {
			prev.set('col', selectedCol || '');
			prev.set('order', selectedOrder || '');
			return prev;
		});
	}, [selectedCol, selectedOrder, setUrlParams]);

	return (
		<div className="flex flex-col gap-1">
			<p className="text-sm font-medium text-muted-foreground">Order by:</p>
			<div className="flex items-center justify-center gap-4 flex-wrap">
				<Select onValueChange={(value) => setSelectedCol(value)}>
					<SelectTrigger className="w-[180px]">
						<SelectValue placeholder={selectedCol || 'Select a column'} />
						<SelectContent>
							<SelectGroup>
								<SelectLabel>Columns</SelectLabel>
								{cols.map((col, idx) => (
									<SelectItem key={idx} value={col.columnName}>
										{col.columnName}
									</SelectItem>
								))}
							</SelectGroup>
						</SelectContent>
					</SelectTrigger>
				</Select>

				<Select onValueChange={(value) => setSelectedOrder(value)}>
					<SelectTrigger className="w-[180px]">
						<SelectValue placeholder={'Select a Order'} />
						<SelectContent>
							<SelectGroup>
								<SelectLabel>{selectedOrder || 'Select a Order'}</SelectLabel>
								<SelectItem value={'asc'}>Ascending</SelectItem>
								<SelectItem value={'desc'}>Descending</SelectItem>
							</SelectGroup>
						</SelectContent>
					</SelectTrigger>
				</Select>
			</div>
		</div>
	);
};
