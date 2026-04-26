import { type Control, Controller, type UseFormSetValue, type UseFormGetValues, useWatch } from 'react-hook-form';
import { Checkbox } from '@/components/ui/checkbox';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type { FormValues } from '@/pages/TableForm';
import type { VarianDataType } from '@/client';

interface TableInputProps {
  index: number;
  dataTypes: VarianDataType[];
  control: Control<FormValues>;
  setValue: UseFormSetValue<FormValues>;
  getValues: UseFormGetValues<FormValues>;
}

export function TableFormInput({
  index,
  dataTypes,
  control,
  setValue,
  getValues,
}: TableInputProps) {
  const columnType = useWatch({
    control,
    name: `inputs.${index}.columnType`,
  });

  const columnSize = useWatch({
    control,
    name: `inputs.${index}.size`,
  });

  function getDataTypeByType(type: string) {
    return (
      dataTypes.find(({ type: t }) => t === type) ||
      dataTypes[0]
    );
  }

  return (
    <div className="flex flex-col bg-foreground/5 rounded-md p-4 w-full gap-4">
      <Controller
        control={control}
        name={`inputs.${index}.columnName`}
        render={({ field, fieldState }) => (
          <Field>
            <FieldLabel htmlFor={`columnName-${index}`}>Column Name</FieldLabel>
            <Input
              {...field}
              id={`columnName-${index}`}
              placeholder="Enter column name"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <Controller
        control={control}
        name={`inputs.${index}.defaultValue`}
        render={({ field, fieldState }) => (
          <Field>
            <FieldLabel htmlFor={`defaultValue-${index}`}>Default value</FieldLabel>
            <Input
              {...field}
              id={`defaultValue-${index}`}
              placeholder="Enter default value"
              value={(field.value as string) || ''}
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <div className="grid grid-cols-2 gap-4">
        <Controller
          control={control}
          name={`inputs.${index}.columnType.dataType`}
          render={({ field, fieldState }) => (
            <Field>
              <FieldLabel htmlFor={`dataType-${index}`}>Data Type</FieldLabel>
              <Select
                onValueChange={(type) => {
                  const dataType = getDataTypeByType(type);
                  // Update the entire columnType object to ensure all flags are set correctly
                  setValue(`inputs.${index}.columnType`, {
                    ...getValues(`inputs.${index}.columnType`),
                    dataType: dataType.type,
                    hasSize: dataType.hasSize,
                    hasAutoIncrement: dataType.hasAutoIncrement,
                    hasDigit: dataType.hasDigit,
                    hasValues: dataType.hasValues,
                  }, { shouldDirty: true });

                  field.onChange(dataType.type);
                }}
                value={field.value}
              >
                <SelectTrigger id={`dataType-${index}`}>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {dataTypes.map(({ type }, idx) => (
                    <SelectItem key={idx} value={type}>
                      {type}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        {columnType?.hasSize && (
          <Field>
            <FieldLabel htmlFor={`dataSize-${index}`}>
              Size/Length
            </FieldLabel>
            <Input
              id={`dataSize-${index}`}
              type="number"
              value={columnSize?.toString() || ''}
              onChange={(e) => {
                setValue(`inputs.${index}.size`, e.target.value ? Number(e.target.value) : undefined, { shouldDirty: true });
              }}
            />
          </Field>
        )}

        {columnType?.hasAutoIncrement && (
          <Field
            orientation="horizontal"
            className="cursor-pointer p-2 rounded"
          >
            <Checkbox
              className="cursor-pointer"
              id={`dataHasAutoIncrement-${index}`}
              checked={columnType.isPk}
              onCheckedChange={(checked) => {
                setValue(`inputs.${index}.columnType.isPk`, checked === true, { shouldDirty: true });
                if (checked === true) {
                  setValue(`inputs.${index}.columnType.isNull`, false, { shouldDirty: true });
                }
              }}
            />
            <FieldLabel htmlFor={`dataHasAutoIncrement-${index}`}>
              Auto Increment
            </FieldLabel>
          </Field>
        )}
      </div>

      <div className="grid grid-cols-3 gap-2">
        <Controller
          control={control}
          name={`inputs.${index}.columnType.isNull`}
          render={({ field, fieldState }) => (
            <Field
              orientation="horizontal"
              className="cursor-pointer bg-foreground/5 p-2 rounded"
            >
              <Checkbox
                id={`isNull-${index}`}
                checked={field.value}
                onCheckedChange={(checked) => field.onChange(checked === true)}
              />
              <FieldLabel
                htmlFor={`isNull-${index}`}
                className="cursor-pointer font-normal"
              >
                Allow NULL
              </FieldLabel>

              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        <Controller
          control={control}
          name={`inputs.${index}.columnType.isPk`}
          render={({ field, fieldState }) => (
            <Field
              orientation="horizontal"
              className="cursor-pointer bg-foreground/5 p-2 rounded"
            >
              <Checkbox
                id={`isPk-${index}`}
                checked={field.value}
                onCheckedChange={(checked) => field.onChange(checked === true)}
              />
              <FieldLabel
                htmlFor={`isPk-${index}`}
                className="cursor-pointer font-normal"
              >
                Primary Key
              </FieldLabel>

              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />

        <Controller
          control={control}
          name={`inputs.${index}.columnType.isUnique`}
          render={({ field, fieldState }) => (
            <Field
              orientation="horizontal"
              className="cursor-pointer bg-foreground/5 p-2 rounded"
            >
              <Checkbox
                id={`isUnique-${index}`}
                checked={field.value}
                onCheckedChange={(checked) => field.onChange(checked === true)}
              />
              <FieldLabel
                htmlFor={`isUnique-${index}`}
                className="cursor-pointer font-normal"
              >
                Unique
              </FieldLabel>

              {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
            </Field>
          )}
        />
      </div>
    </div>
  );
}


