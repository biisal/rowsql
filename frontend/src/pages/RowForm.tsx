import { useEffect, useState } from 'react';
import {
	useParams,
	useSearchParams,
	useNavigate,
	Link,
} from 'react-router-dom';
import { Controller, useForm } from 'react-hook-form';
import { ArrowLeft, Save, Loader2 } from 'lucide-react';
import api from '@/lib/axios';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card';
import { toast } from 'sonner';
import { Checkbox } from '@/components/ui/checkbox';
import { Textarea } from '@/components/ui/textarea';
import {
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from '@/components/ui/field';

interface Column {
	columnName: string;
	dataType: string;
	inputType: 'text' | 'number' | 'checkbox' | 'textarea' | 'json' | 'select';
	value: string | number | boolean | null;
	isUnique: boolean;
	hasAutoIncrement: boolean;
	hasDefault: boolean;
}

interface FormData {
	Action: string;
	Tables: unknown[];
	Cols: Column[];
	ActiveTable: string;
}

interface CheckedState {
	checked: boolean;
	oldVal: undefined;
}

export function RowForm() {
	const [autoEnabled, setAutoEnabled] = useState<Record<string, CheckedState>>(
		{},
	);

	const [hasDefaults, setHasDefaluts] = useState<Record<string, CheckedState>>(
		{},
	);
	const { tableName } = useParams<{ tableName: string }>();
	const [searchParams] = useSearchParams();
	const navigate = useNavigate();
	const [loading, setLoading] = useState(true);
	const [formData, setFormData] = useState<FormData | null>(null);

	const hash = searchParams.get('hash');
	const page = searchParams.get('page') || '1';

	const form = useForm({
		defaultValues: {},
		mode: 'onChange',
		reValidateMode: 'onChange',
	});

	useEffect(() => {
		const fetchFormData = async () => {
			if (!tableName) return;

			setLoading(true);
			try {
				const url = hash
					? `/tables/${tableName}/form?hash=${hash}&page=${page}`
					: `/tables/${tableName}/form`;
				const response = await api.get(url);

				if (response.data.success) {
					const data = response.data.data;
					setFormData(data);

					// Reset form with default values
					if (data.Cols && data.Cols.length > 0) {
						const defaultValues: Record<string, string | number | boolean> = {};
						data.Cols.forEach((col: Column) => {
							if (col.value !== undefined && col.value !== null) {
								if (col.inputType === 'checkbox') {
									defaultValues[col.columnName] = Boolean(col.value);
								} else if (col.inputType === 'number') {
									defaultValues[col.columnName] = col.value;
								} else {
									defaultValues[col.columnName] = String(col.value);
								}
							} else {
								defaultValues[col.columnName] =
									col.inputType === 'checkbox' ? false : '';
							}
						});

						form.reset(defaultValues);
					}
				}
			} catch (err) {
				const errorMessage =
					(
						err as {
							response?: { data?: { error?: string } };
							message?: string;
						}
					).response?.data?.error ||
					(err as { message?: string }).message ||
					'Failed to load form';
				toast.error(errorMessage);
			} finally {
				setLoading(false);
			}
		};

		fetchFormData();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [tableName, hash, page]);

	const onSubmit = async (data: Record<string, string | number | boolean>) => {
		if (!formData || !tableName) return;

		console.log('=== FORM SUBMISSION STARTED ===');
		console.log('Form submitted with values:', data);

		try {
			const payload: Record<string, { value: string; type: string }> = {};

			formData.Cols.forEach((col) => {
				const value = data[col.columnName];
				payload[col.columnName] = {
					value: value !== undefined && value !== null ? String(value) : '',
					type: col.inputType,
				};
			});

			console.log('Payload:', payload);

			const url = hash
				? `/tables/${tableName}/form?hash=${hash}&page=${page}`
				: `/tables/${tableName}/form`;

			await api.post(url, payload);

			toast.success(`Row ${hash ? 'updated' : 'created'} successfully`);
			navigate(`/tables/${tableName}?page=${page}`);
		} catch (err) {
			console.error('Error saving row:', err);
			const errorMessage =
				(err as { response?: { data?: { error?: string } }; message?: string })
					.response?.data?.error ||
				(err as { message?: string }).message ||
				'Failed to save row';
			toast.error(errorMessage);
		}
	};

	if (loading) {
		return (
			<div className="flex items-center justify-center h-full">
				<Loader2 className="w-8 h-8 animate-spin text-primary" />
			</div>
		);
	}

	if (!formData) {
		return (
			<div className="flex items-center justify-center h-full">
				<p className="text-muted-foreground">No data available</p>
			</div>
		);
	}

	return (
		<div className="bg-background p-8">
			<div className="max-w-3xl mx-auto">
				<div className="flex items-center gap-4 mb-6">
					<Link to={`/tables/${tableName}?page=${page}`}>
						<Button variant="ghost" size="icon">
							<ArrowLeft className="w-5 h-5" />
						</Button>
					</Link>
					<h1 className="text-3xl font-bold tracking-tight">
						{formData.Action} Row - {tableName}
					</h1>
				</div>

				<Card>
					<form onSubmit={form.handleSubmit(onSubmit)}>
						<FieldGroup>
							<CardHeader>
								<CardTitle>
									{formData.Action === 'Insert' ? 'Add New Row' : 'Edit Row'}
								</CardTitle>
							</CardHeader>

							<CardContent className="space-y-6">
								{formData.Cols && Array.isArray(formData.Cols) ? (
									formData.Cols.map((col, index) => (
										<Controller
											key={index}
											control={form.control}
											name={col.columnName as never}
											render={({ field, fieldState }) => (
												<Field>
													<FieldLabel
														htmlFor={`field-${col.columnName}`}
														className="capitalize"
													>
														{col.columnName.replace(/_/g, ' ')}
														<span className="text-muted-foreground ml-2 text-xs font-normal">
															({col.dataType})
														</span>
														{col.isUnique && (
															<span className="ml-2 text-xs font-normal text-primary">
																• Unique
															</span>
														)}
													</FieldLabel>

													{col.hasDefault && (
														<label className="flex items-center mb-2 text-sm font-medium cursor-pointer">
															<Checkbox
																className="mr-2"
																checked={!!hasDefaults[col.columnName]?.checked}
																onCheckedChange={(checked) => {
																	const isChecked = checked === true;

																	setHasDefaluts((prev) => ({
																		...prev,
																		[col.columnName]: {
																			checked: isChecked,
																			oldVal: field.value,
																		},
																	}));

																	if (isChecked) {
																		field.onChange('');
																	} else {
																		field.onChange(
																			hasDefaults[col.columnName]?.oldVal || '',
																		);
																	}
																}}
															/>
															Use Default Value
														</label>
													)}

													{col.inputType === 'checkbox' ? (
														<div className="flex items-center space-x-2 cursor-pointer bg-foreground/5 p-3 rounded">
															<Checkbox
																id={`field-${col.columnName}`}
																disabled={hasDefaults[col.columnName]?.checked}
																checked={field.value || false}
																onCheckedChange={(checked) =>
																	field.onChange(checked === true)
																}
															/>
															<FieldLabel
																htmlFor={`field-${col.columnName}`}
																className="cursor-pointer mt-0! font-normal"
															>
																Enable {col.columnName.replace(/_/g, ' ')}
															</FieldLabel>
														</div>
													) : col.inputType === 'textarea' ||
													  col.inputType === 'json' ? (
														<Textarea
															disabled={hasDefaults[col.columnName]?.checked}
															{...field}
															id={`field-${col.columnName}`}
															placeholder={`Enter ${col.columnName.replace(/_/g, ' ')}`}
															rows={5}
															value={field.value || ''}
														/>
													) : (
														<>
															{col.hasAutoIncrement && (
																<label className="flex items-center mb-2 text-sm font-medium cursor-pointer">
																	<Checkbox
																		className="mr-2"
																		checked={
																			!!autoEnabled[col.columnName]?.checked
																		}
																		onCheckedChange={(checked) => {
																			const isChecked = checked === true;

																			setAutoEnabled((prev) => ({
																				...prev,
																				[col.columnName]: {
																					checked: isChecked,
																					oldVal: field.value,
																				},
																			}));

																			if (isChecked) {
																				field.onChange('');
																			} else {
																				field.onChange(
																					autoEnabled[col.columnName]?.oldVal ||
																						'',
																				);
																			}
																		}}
																	/>
																	Auto Increment
																</label>
															)}
															<Input
																{...field}
																disabled={
																	autoEnabled[col.columnName]?.checked ||
																	hasDefaults[col.columnName]?.checked
																}
																id={`field-${col.columnName}`}
																type={
																	col.inputType === 'number' ? 'number' : 'text'
																}
																placeholder={`Enter ${col.columnName.replace(/_/g, ' ')}`}
																value={field.value || ''}
															/>
														</>
													)}

													{fieldState.invalid && (
														<FieldError errors={[fieldState.error]} />
													)}
												</Field>
											)}
										/>
									))
								) : (
									<div className="text-center text-muted-foreground py-8">
										No columns found for this table.
									</div>
								)}
							</CardContent>

							<CardFooter className="flex justify-end gap-4">
								<Link to={`/tables/${tableName}?page=${page}`}>
									<Button variant="outline" type="button">
										Cancel
									</Button>
								</Link>
								<Button
									type="submit"
									disabled={form.formState.isSubmitting}
									onClick={() => {
										console.log('=== SUBMIT BUTTON CLICKED ===');
										console.log('Form values:', form.getValues());
										console.log('Form errors:', form.formState.errors);
										console.log('Is valid:', form.formState.isValid);
									}}
								>
									{form.formState.isSubmitting ? (
										<>
											<Loader2 className="mr-2 h-4 w-4 animate-spin" />
											Saving...
										</>
									) : (
										<>
											<Save className="mr-2 h-4 w-4" />
											Save
										</>
									)}
								</Button>
							</CardFooter>
						</FieldGroup>
					</form>
				</Card>
			</div>
		</div>
	);
}
