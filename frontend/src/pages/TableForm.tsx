import { TableFormInput } from '@/feature/tables';
import { Button } from '@/components/ui/button';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { useMemo } from 'react';
import { Controller, useForm, type Resolver, useFieldArray } from 'react-hook-form';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from '@/components/ui/field';
import { useMutation, useQuery } from '@tanstack/react-query';
import { createNewTableMutation, newTableFormFiledsOptions } from '@/client/@tanstack/react-query.gen';
import { Loader2 } from 'lucide-react';
import { zColValue } from '@/client/zod.gen';

const zColValueFixed = zColValue.extend({
	size: z.coerce.number().optional(),
});

const formSchema = z.object({
	tableName: z.string().min(1, 'Table name is required'),
	inputs: z.array(zColValueFixed)

})
export type FormValues = z.infer<typeof formSchema>;
type zcolFixedType = z.infer<typeof zColValueFixed>;

function getDefalutInputs(): zcolFixedType {
	return {
		columnName: "",
		columnType: {
			dataType: "INT",
			hasSize: false,
			hasAutoIncrement: false,
			hasDefault: false,
			inputType: "text",
			isNull: false,
			isPk: false,
			isUnique: false,
		},
	}
}
export const TableForm = () => {
	const navigate = useNavigate();

	const { data: serverDataTypes, isLoading: loadingDataTypes } = useQuery(newTableFormFiledsOptions());

	const createTableMutation = useMutation({
		...createNewTableMutation(),
		onSuccess: (_, variable) => {
			const tableName = variable.body.tableName;
			toast.success('Table created successfully');
			navigate(`/tables/${tableName}`);
		},
	});

	const dataTypes = useMemo(() => {
		if (!serverDataTypes) return [];
		return [...(serverDataTypes.numericType || []), ...(serverDataTypes.stringType || [])];
	}, [serverDataTypes]);

	const initialValues = useMemo(() => {
		if (dataTypes.length === 0) return undefined;
		const initialInput = getDefalutInputs();
		return {
			tableName: '',
			inputs: [initialInput],
		};
	}, [dataTypes]);

	const form = useForm<FormValues>({
		resolver: zodResolver(formSchema) as Resolver<FormValues>,
		values: initialValues,
		resetOptions: {
			keepDirtyValues: true,
		},
		mode: 'onChange',
		reValidateMode: 'onChange',
	});

	const { fields, append, remove } = useFieldArray({
		control: form.control,
		name: 'inputs',
	});

	function addCol() {
		const newInput = getDefalutInputs();
		if (dataTypes.length > 0) {
			newInput.columnType.dataType = dataTypes[0].type;
			newInput.columnType.hasSize = dataTypes[0].hasSize;
			newInput.columnType.hasAutoIncrement = dataTypes[0].hasAutoIncrement;
			newInput.columnType.hasDigit = dataTypes[0].hasDigit;
			newInput.columnType.hasValues = dataTypes[0].hasValues;
		}
		append(newInput);
	}

	function removeCol(index: number) {
		if (fields.length <= 1) {
			toast.error('At least one column is required');
			return;
		}
		remove(index);
	}

	async function onSubmit(values: FormValues) {
		const emptyColumns = values.inputs.filter((input) => !input.columnName.trim());
		if (emptyColumns.length > 0) {
			toast.error('Please fill in all column names');
			return;
		}
		await createTableMutation.mutateAsync({
			body: {
				tableName: values.tableName,
				inputs: values.inputs,
			},
		});
	}

	if (!initialValues || loadingDataTypes) {
		return <div className="p-8 flex justify-center"><Loader2 className="w-8 h-8 animate-spin text-primary" /></div>;
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
								{fields.map((field, index) => (
									<div key={field.id} className="relative">
										<TableFormInput
											dataTypes={dataTypes}
											index={index}
											control={form.control}
											setValue={form.setValue}
											getValues={form.getValues}
										/>
										{fields.length > 1 && (
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
											console.log('Fields:', fields);
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
