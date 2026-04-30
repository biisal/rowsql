import { useEffect } from "react";
import {
  useParams,
  useSearchParams,
  useNavigate,
  Link,
} from "react-router-dom";
import { Controller, useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowLeft, Save, Loader2 } from "lucide-react";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { toast } from "sonner";
import { Checkbox } from "@/components/ui/checkbox";
import { Textarea } from "@/components/ui/textarea";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  insertOrUpdateRowMutation,
  rowInsertOrUpdateFormOptions,
} from "@/client/@tanstack/react-query.gen";
import type { ColValue, ErrorModel } from "@/client";
import { zColValue } from "@/client/zod.gen";

const zColField = zColValue.extend({
  size: z.coerce.number().optional(),
  value: z.union([z.string(), z.boolean(), z.number(), z.null()]).default(null),
  useDefault: z.boolean().default(false),
  useAutoIncrement: z.boolean().default(false),
});

const formSchema = z.object({
  cols: z.array(zColField).min(1, "At least one column required"),
});

type FormInput = z.input<typeof formSchema>;
type FormSchema = z.output<typeof formSchema>;
type ColField = z.infer<typeof zColField>;

const colFieldValueSchema = z
  .union([z.string(), z.boolean(), z.number(), z.null()])
  .catch(null);

function buildDefaultCols(cols: ColValue[]): ColField[] {
  return cols.map((col) => ({
    columnName: col.columnName,
    columnType: col.columnType,
    defaultValue: col.defaultValue,
    size: col.size !== undefined ? Number(col.size) : undefined,
    value: colFieldValueSchema.parse(col.value),
    useDefault: false,
    useAutoIncrement: false,
  }));
}

export function RowForm() {
  const { tableName } = useParams<{ tableName: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const hash = searchParams.get("hash");
  const page = Math.max(1, Number(searchParams.get("page")) || 1);
  const isEdit = !!hash;

  const { data, isPending } = useQuery(
    rowInsertOrUpdateFormOptions({
      path: { tableName: tableName || "" },
      query: { hash: hash || undefined, page },
    }),
  );

  const { mutateAsync } = useMutation(insertOrUpdateRowMutation());

  const cols: ColValue[] = data?.cols ?? [];

  const form = useForm<FormInput, unknown, FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      cols: [],
    },
    mode: "onChange",
  });

  const { fields, replace } = useFieldArray({
    control: form.control,
    name: "cols",
  });

  useEffect(() => {
    if (cols.length) {
      replace(buildDefaultCols(cols));
    }
  }, [data]);

  const onSubmit = async (formValues: FormSchema) => {
    if (!tableName) return;
    await mutateAsync(
      {
        path: {
          tableName: tableName,
        },
        body: formValues.cols,
        query: {
          hash: hash || undefined,
          page: page,
        },
      },
      {
        onError: (error: ErrorModel) => {
          console.log(error.errors?.[0]?.message);
        },
        onSuccess: () => {
          toast.success("Row inserted/updated successfully");
          navigate(`/tables/${tableName}?page=${page}`);
        },
      },
    );
  };

  if (isPending) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 animate-spin text-primary" />
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
            {isEdit ? "Edit" : "Insert"} Row — {tableName}
          </h1>
        </div>

        <Card>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <FieldGroup>
              <CardHeader>
                <CardTitle>{isEdit ? "Edit Row" : "Add New Row"}</CardTitle>
              </CardHeader>

              <CardContent className="space-y-6">
                {fields.length === 0 ? (
                  <div className="text-center text-muted-foreground py-8">
                    No columns found for this table.
                  </div>
                ) : (
                  fields.map((field, index) => {
                    const { columnName, columnType } = field;

                    return (
                      <Controller
                        key={field.id}
                        control={form.control}
                        name={`cols.${index}.value`}
                        render={({ field: inputField, fieldState }) => {
                          const useDefault = form.watch(
                            `cols.${index}.useDefault`,
                          );
                          const useAutoIncrement = form.watch(
                            `cols.${index}.useAutoIncrement`,
                          );
                          const isDisabled = useDefault || useAutoIncrement;

                          return (
                            <Field>
                              <FieldLabel
                                htmlFor={`field-${columnName}`}
                                className="capitalize"
                              >
                                {columnName.replace(/_/g, " ")}
                                <span className="text-muted-foreground ml-2 text-xs font-normal">
                                  ({columnType.dataType})
                                </span>
                                {columnType.isUnique && (
                                  <span className="ml-2 text-xs font-normal text-primary">
                                    • Unique
                                  </span>
                                )}
                              </FieldLabel>

                              {columnType.hasDefault && (
                                <Controller
                                  control={form.control}
                                  name={`cols.${index}.useDefault`}
                                  render={({ field: f }) => (
                                    <label className="flex items-center mb-2 text-sm font-medium cursor-pointer">
                                      <Checkbox
                                        className="mr-2"
                                        checked={f.value}
                                        onCheckedChange={(checked) => {
                                          f.onChange(checked === true);
                                          if (checked) inputField.onChange("");
                                        }}
                                      />
                                      Use Default Value
                                    </label>
                                  )}
                                />
                              )}

                              {columnType.hasAutoIncrement && (
                                <Controller
                                  control={form.control}
                                  name={`cols.${index}.useAutoIncrement`}
                                  render={({ field: f }) => (
                                    <label className="flex items-center mb-2 text-sm font-medium cursor-pointer">
                                      <Checkbox
                                        className="mr-2"
                                        checked={f.value}
                                        onCheckedChange={(checked) => {
                                          f.onChange(checked === true);
                                          if (checked)
                                            inputField.onChange(undefined);
                                        }}
                                      />
                                      Auto Increment
                                    </label>
                                  )}
                                />
                              )}

                              {columnType.inputType === "checkbox" ? (
                                <div className="flex items-center space-x-2 bg-foreground/5 p-3 rounded">
                                  <Checkbox
                                    id={`field-${columnName}`}
                                    disabled={isDisabled}
                                    checked={!!inputField.value}
                                    onCheckedChange={(c) =>
                                      inputField.onChange(c === true)
                                    }
                                  />
                                  <FieldLabel
                                    htmlFor={`field-${columnName}`}
                                    className="cursor-pointer mt-0! font-normal"
                                  >
                                    Enable {columnName.replace(/_/g, " ")}
                                  </FieldLabel>
                                </div>
                              ) : columnType.inputType === "textarea" ||
                                columnType.inputType === "json" ? (
                                <Textarea
                                  {...inputField}
                                  id={`field-${columnName}`}
                                  disabled={isDisabled}
                                  placeholder={`Enter ${columnName.replace(/_/g, " ")}`}
                                  rows={5}
                                  value={(inputField.value as string) || ""}
                                />
                              ) : (
                                <Input
                                  {...inputField}
                                  id={`field-${columnName}`}
                                  disabled={isDisabled}
                                  type={
                                    columnType.inputType === "number"
                                      ? "number"
                                      : "text"
                                  }
                                  placeholder={`Enter ${columnName.replace(/_/g, " ")}`}
                                  value={(inputField.value as string) || ""}
                                  onChange={(e) =>
                                    inputField.onChange(
                                      columnType.inputType === "number"
                                        ? e.target.valueAsNumber
                                        : e.target.value,
                                    )
                                  }
                                />
                              )}

                              {fieldState.invalid && (
                                <FieldError errors={[fieldState.error]} />
                              )}
                            </Field>
                          );
                        }}
                      />
                    );
                  })
                )}
              </CardContent>

              <CardFooter className="flex justify-end gap-4">
                <Link to={`/tables/${tableName}?page=${page}`}>
                  <Button variant="outline" type="button">
                    Cancel
                  </Button>
                </Link>
                <Button type="submit" disabled={form.formState.isSubmitting}>
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
