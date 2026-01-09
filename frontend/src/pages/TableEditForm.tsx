import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select';
// import type { Column } from '@/lib/types';
import { useEffect } from 'react';

const TableEditForm = () => {
	// const [columns, setColumns] = useState<Column[]>([]);

	useEffect(() => {}, []);

	const onSelect = (value: string) => {
		console.log('Selected value:', value);
	};

	return (
		<Card className="">
			<CardHeader>
				<CardTitle>Add Table Column</CardTitle>
			</CardHeader>
			<CardContent>
				<Select onValueChange={onSelect}>
					<SelectTrigger>
						<SelectValue placeholder="Choose a column">
							Choose a column
						</SelectValue>
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectItem value="string">String</SelectItem>
							<SelectItem value="number">Number</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			</CardContent>
		</Card>
	);
};

export default TableEditForm;
