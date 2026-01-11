import TableFromInput from '@/components/table-form-input';
import { Button } from '@/components/ui/button';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import api from '@/lib/axios';
import type {
	DbDataTypes,
	ErrorResponse,
	Form as FormType,
	Input as InputType,
} from '@/lib/types';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';
import useTableStore, { createTable } from '@/lib/store/use-table';
import { useNavigate } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from '@/components/ui/field';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
const formSchema = z.object({
	tableName: z.string().min(1, 'Table name is required'),
	inputs: z
		.array(
			z.object({
				colName: z.string().min(1, 'Column name is required'),
				isNull: z.boolean(),
				isPk: z.boolean(),
				isUnique: z.boolean(),
				dataType: z.object({
					type: z.string().min(1, 'Data type is required'),
					size: z.number().optional(),
					hasSize: z.boolean(),
					hasBool: z.boolean().optional(),
					autoIncrement: z.boolean().optional(),
				}),
			}),
		)
		.min(1, 'At least one column is required'),
});

type FormValues = z.infer<typeof formSchema>;

export const TableForm = () => {
	const queryClient = useQueryClient();
	const navigate = useNavigate();
	const { refreshTables } = useTableStore();

	const createTableMutation = useMutation({
		mutationFn: createTable,
		onSuccess: (_, variable) => {
			const tableName = variable.tableName;
			toast.success('Table created successfully');
			queryClient.invalidateQueries({ queryKey: ['tables'] });
			navigate(`/tables/${tableName}`);
			refreshTables(true);
		},
		onError: (err: AxiosError<ErrorResponse>) => {
			toast.error(
				err.response?.data?.error || err.message || 'Something went wrong',
			);
		},
	});

	// const { createTable } = useTableStore();
	const [formData, setFormData] = useState<FormType>({
		inputs: [],
		dataTypes: [],
		selectedDataType: {
			type: 'VARCHAR',
			size: 255,
			hasSize: true,
		},
		tableName: '',
	});
	const [mount, setMount] = useState(false);

	const form = useForm<FormValues>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			tableName: '',
			inputs: [],
		},
		mode: 'onChange',
		reValidateMode: 'onChange',
	});

	function addCol() {
		const currentInputs = form.getValues('inputs');
		const newInput: InputType = {
			colName: '',
			isNull: false,
			isPk: false,
			isUnique: false,
			dataType: formData.selectedDataType,
		};

		const newInputs = [...currentInputs, newInput];
		form.setValue('inputs', newInputs, {
			shouldValidate: true,
			shouldDirty: true,
			shouldTouch: true,
		});

		const updatedFormData = {
			...formData,
			inputs: [...formData.inputs, newInput],
		};
		setFormData(updatedFormData);
	}

	function removeCol(index: number) {
		const currentInputs = form.getValues('inputs');
		if (currentInputs.length <= 1) {
			toast.error('At least one column is required');
			return;
		}

		const newInputs = currentInputs.filter((_, idx) => idx !== index);
		form.setValue('inputs', newInputs, { shouldValidate: true });

		const updatedFormData = { ...formData };
		updatedFormData.inputs = updatedFormData.inputs.filter(
			(_, idx) => idx !== index,
		);
		setFormData(updatedFormData);
	}

	useEffect(() => {
		async function fetchFormDataTypes() {
			try {
				const res = await api.get('/tables/form/new');
				if (res.status === 200) {
					const data: { data: DbDataTypes } = res.data;
					const types = [...data.data.numericType, ...data.data.stringType];
					if (types.length === 0) {
						toast.error('No data types found');
						return;
					}
					const newForm: FormType = {
						tableName: '',
						dataTypes: types,
						inputs: [
							{
								colName: '',
								isNull: false,
								isPk: false,
								isUnique: false,
								dataType: types[0],
							},
						],
						selectedDataType: types[0],
					};
					setFormData(newForm);

					form.reset(
						{
							tableName: newForm.tableName,
							inputs: newForm.inputs,
						},
						{ keepDefaultValues: false },
					);
				}
			} catch (error) {
				console.error(error);
				toast.error('Failed to fetch data types');
			} finally {
				setMount(true);
			}
		}
		fetchFormDataTypes();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	async function onSubmit(values: FormValues) {
		const emptyColumns = values.inputs.filter((input) => !input.colName.trim());
		if (emptyColumns.length > 0) {
			console.log('Empty columns found:', emptyColumns);
			toast.error('Please fill in all column names');
			return;
		}
		await createTableMutation.mutateAsync({
			tableName: values.tableName,
			inputs: values.inputs,
		});
	}

	if (!mount) {
		return null;
	}

	return (
		<div className="bg-background p-8">
			<div className="mx-auto">
				<Card>
					<form onSubmit={form.handleSubmit(onSubmit)}>
						<FieldGroup>
							<CardHeader>
								<CardTitle>Add Table Column</CardTitle>

								<Controller
									control={form.control}
									name="tableName"
									render={({ field, fieldState }) => (
										<Field>
											<FieldLabel>Table Name</FieldLabel>
											<Input {...field} placeholder="Enter table name" />
											{fieldState.invalid && (
												<FieldError errors={[fieldState.error]} />
											)}
										</Field>
									)}
								/>
							</CardHeader>
							<CardContent className="gap-4 grid lg:grid-cols-2">
								{form.formState.errors.inputs &&
									!Array.isArray(form.formState.errors.inputs) && (
										<div className="col-span-full">
											<FieldError errors={[form.formState.errors.inputs]} />
										</div>
									)}
								{formData?.inputs.map((_, index) => (
									<div key={index} className="relative">
										<TableFromInput
											formData={formData}
											index={index}
											control={form.control}
										/>
										{formData.inputs.length > 1 && (
											<Button
												type="button"
												variant="danger"
												size="sm"
												className="absolute top-4 right-4"
												onClick={() => removeCol(index)}
											>
												Remove
											</Button>
										)}
									</div>
								))}
							</CardContent>
							<CardFooter className="space-x-2">
								<Button type="button" onClick={addCol}>
									Add Column
								</Button>
								<Button
									type="submit"
									disabled={
										createTableMutation.isPending || !form.formState.isValid
									}
									onClick={() => {
										console.log('=== SUBMIT BUTTON CLICKED ===');
										console.log('Form values:', form.getValues());
										console.log('Form errors:', form.formState.errors);
										console.log('Is valid:', form.formState.isValid);
									}}
								>
									{form.formState.isSubmitting ? 'Creating...' : 'Create Table'}
								</Button>
								{import.meta.env.DEV && (
									<Button
										type="button"
										variant="outline"
										onClick={() => {
											console.log('=== DEBUG INFO ===');
											console.log('Form values:', form.getValues());
											console.log('Form errors:', form.formState.errors);
											console.log('Is valid:', form.formState.isValid);
											console.log(
												'Is submitting:',
												form.formState.isSubmitting,
											);
											console.log('Local formData:', formData);
										}}
									>
										Debug
									</Button>
								)}
							</CardFooter>
						</FieldGroup>
					</form>
				</Card>
			</div>
		</div>
	);
};
