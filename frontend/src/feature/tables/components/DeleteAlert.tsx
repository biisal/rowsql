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
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteTableMutation, listTablesQueryKey } from '@/client/@tanstack/react-query.gen';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input } from '@/components/ui/input';
import { toast } from 'sonner';

export function DeleteAlert({ tableName }: { tableName: string }) {
	const queryClient = useQueryClient();
	const navigate = useNavigate();
	const [verificationQuery, setVerificationQuery] = useState('');
	const [open, setOpen] = useState(false);

	const deleteMutation = useMutation({
		...deleteTableMutation(),
		onSuccess: () => {
			toast.success('Table deleted successfully');
			queryClient.invalidateQueries({ queryKey: listTablesQueryKey() });
			setOpen(false);
			navigate('/');
		},
		onError: (err) => {
			toast.error(err?.detail || 'Failed to delete table');
		}
	});

	async function handleDelete(e: React.MouseEvent) {
		e.preventDefault();
		await deleteMutation.mutateAsync({
			body: {
				tableName,
				verificationQuery
			}
		});
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
					<AlertDialogCancel disabled={deleteMutation.isPending}>Cancel</AlertDialogCancel>
					<AlertDialogAction
						onClick={handleDelete}
						disabled={deleteMutation.isPending}
						className={buttonVariants({ variant: 'danger' })}
					>
						Delete
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
}
