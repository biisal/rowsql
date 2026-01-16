import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button, buttonVariants } from '@/components/ui/button';
import useTableStore from '@/lib/store/use-table';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input } from './ui/input';

export function DeletAlert({ tableName }: { tableName: string }) {
	const { deleteTable, tableDeleting } = useTableStore();
	const [verificationQuery, setVerificationQuery] = useState('');
	const [open, setOpen] = useState(false);
	const navigate = useNavigate();

	async function handleDelete(e: React.MouseEvent) {
		e.preventDefault();

		const success = await deleteTable(tableName, verificationQuery);
		if (success) {
			setOpen(false);
			navigate('/');
		}
	}

	return (
		<AlertDialog open={open} onOpenChange={setOpen}>
			<AlertDialogTrigger asChild>
				<Button variant="danger">Delete Table</Button>
			</AlertDialogTrigger>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
					<AlertDialogDescription className="text-lg">
						enter{' '}
						<span className="text-foreground">
							DROP TABLE <strong>{tableName}</strong>
						</span>{' '}
						to confirm
					</AlertDialogDescription>
				</AlertDialogHeader>

				<Input
					placeholder="Enter verification query"
					value={verificationQuery}
					onChange={(e) => setVerificationQuery(e.target.value)}
				/>
				<AlertDialogFooter>
					<AlertDialogCancel disabled={tableDeleting}>Cancel</AlertDialogCancel>
					<AlertDialogAction
						onClick={handleDelete}
						disabled={tableDeleting}
						className={buttonVariants({ variant: 'danger' })}
					>
						Delete
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
}
