import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Controller, useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Input } from "@/components/ui/input";
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
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

import { Button } from "@/components/ui/button";
import { Loader2, Plus, Save } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useRowContext } from "@/hooks/useRows";

interface RowUpsertFormProps {
  children?: React.ReactNode;
  hash?: string;
}

export const RowUpsertForm = ({ hash, children }: RowUpsertFormProps) => {
  const { page, tableName } = useRowContext();
  const navigate = useNavigate();

  const isEdit = !!hash;

  const { data, isPending } = useQuery(
    rowInsertOrUpdateFormOptions({
      path: { tableName },
      query: { hash: hash || undefined, page },
    }),
  );

  const { mutateAsync: insertOrUpdateMutation } = useMutation(
    insertOrUpdateRowMutation(),
  );

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
    if (data?.cols) {
      replace(buildDefaultCols(data.cols));
    }
  }, [data, replace]);

  const onSubmit = async (formValues: FormSchema) => {
    if (!tableName) return;
    await insertOrUpdateMutation(
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
          console.error(
            "Mutation error:",
            error.errors?.[0]?.message || error.detail,
          );
        },
        onSuccess: () => {
          toast.success("Row inserted/updated successfully");
          navigate(`/tables/${tableName}?page=${page}`);
        },
      },
    );
  };
  return (
    <Sheet>
      <SheetTrigger asChild>
        {children || (
          <Button className="shadow-lg shadow-primary/20">
            <Plus className="mr-2 h-4 w-4" /> {isEdit ? "Update" : "Record"}
          </Button>
        )}
      </SheetTrigger>
      <SheetContent
        onPointerDownOutside={(e) => e.preventDefault()}
        className="min-w-[90%] md:min-w-xl flex flex-col p-0"
      >
        {isPending ? (
          <div className="flex items-center justify-center flex-1">
            <Loader2 className="w-8 h-8 animate-spin text-primary" />
          </div>
        ) : (
          <form
            key={hash || "new"}
            onSubmit={form.handleSubmit(onSubmit)}
            className="flex-1 flex flex-col min-h-0"
          >
            <FieldGroup className="flex-1 flex flex-col min-h-0 gap-0">
              <SheetHeader className="p-6 bg-muted/30 border-b">
                <SheetTitle>{isEdit ? "Update" : "Insert"} Row</SheetTitle>
                <SheetDescription>
                  {isEdit
                    ? "Update the details of the existing record."
                    : "Fill in the fields below to add a new record to the table."}
                </SheetDescription>
              </SheetHeader>

              {form.formState.errors.cols?.message && (
                <div className="mx-6 mt-4 p-3 bg-destructive/10 text-destructive text-sm rounded-md border border-destructive/20">
                  {form.formState.errors.cols.message}
                </div>
              )}

              <ScrollArea type="always" className="min-h-0 flex-1">
                <div className="p-6 flex flex-col gap-6">
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
                                            if (checked)
                                              inputField.onChange("");
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
                </div>
              </ScrollArea>
              <SheetFooter>
                {import.meta.env.DEV && (
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => {
                      console.log("=== DEBUG INFO ===");
                      console.log("Form values:", form.getValues());
                      console.log("Form errors:", form.formState.errors);
                      console.log("Is valid:", form.formState.isValid);
                      console.log(
                        "Is submitting:",
                        form.formState.isSubmitting,
                      );
                      console.log("Fields:", fields);
                    }}
                  >
                    Debug
                  </Button>
                )}
                <div className="grid md:grid-cols-2 gap-2">
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
                  <SheetClose asChild>
                    <Button variant="outline">Cancel</Button>
                  </SheetClose>
                </div>
              </SheetFooter>
            </FieldGroup>
          </form>
        )}
      </SheetContent>
    </Sheet>
  );
};
