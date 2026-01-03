import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Field, FieldError, FieldLabel } from '@/components/ui/field';
import type { Form } from '@/lib/types';
import { Controller, type Control } from 'react-hook-form';

interface TableInputProps {
  index: number;
  formData: Form;
  control: Control<{
    tableName: string;
    inputs: Array<{
      colName: string;
      isNull: boolean;
      isPk: boolean;
      isUnique: boolean;
      dataType: {
        type: string;
        size?: number;
        hasSize: boolean;
        hasBool?: boolean;
        autoIncrement?: boolean;
        hasAutoIncrement?: boolean;
      };
    }>;
  }>;
}

export default function TableFromInput({
  index,
  formData,
  control,
}: TableInputProps) {
  function getDataTypeByType(type: string) {
    return (
      formData.dataTypes.find(({ type: t }) => t === type) ||
      formData.dataTypes[0]
    );
  }

  return (
    <div className="flex flex-col bg-foreground/5 rounded-md p-4 w-full gap-4">
      <Controller
        control={control}
        name={`inputs.${index}.colName`}
        render={({ field, fieldState }) => (
          <Field>
            <FieldLabel htmlFor={`colName-${index}`}>Column Name</FieldLabel>
            <Input
              {...field}
              id={`colName-${index}`}
              placeholder="Enter column name"
            />
            {fieldState.invalid && <FieldError errors={[fieldState.error]} />}
          </Field>
        )}
      />

      <div className="grid grid-cols-2 gap-4">
        <Controller
          control={control}
          name={`inputs.${index}.dataType`}
          render={({ field, fieldState }) => (
            <Field>
              <FieldLabel htmlFor={`dataType-${index}`}>Data Type</FieldLabel>
              <Select
                onValueChange={(type) => {
                  const dataType = getDataTypeByType(type);
                  field.onChange(dataType);
                }}
                value={field.value?.type}
              >
                <SelectTrigger id={`dataType-${index}`}>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {formData.dataTypes.map(({ type }, idx) => (
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

        <Controller
          control={control}
          name={`inputs.${index}.dataType`}
          render={({ field, fieldState }) => {
            if (!field.value?.hasSize) return <></>;

            return (
              <Field>
                <FieldLabel htmlFor={`dataSize-${index}`}>
                  Size/Length
                </FieldLabel>
                <Input
                  id={`dataSize-${index}`}
                  type="number"
                  value={field.value.size || ''}
                  onChange={(e) => {
                    field.onChange({
                      ...field.value,
                      size: Number(e.target.value),
                    });
                  }}
                />

                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            );
          }}
        />
        <Controller
          control={control}
          name={`inputs.${index}.dataType`}
          render={({ field, fieldState }) => {
            if (!field.value?.hasAutoIncrement) return <></>;

            return (
              <Field
                orientation="horizontal"
                className="cursor-pointer  p-2 rounded"
              >
                <Checkbox
                  className="cursor-pointer"
                  id={`dataHasAutoIncrement-${index}`}
                  checked={field.value.autoIncrement}
                  onCheckedChange={(checked) =>
                    field.onChange({
                      ...field.value,
                      autoIncrement: checked === true,
                    })
                  }
                />
                <FieldLabel htmlFor={`dataHasAutoIncrement-${index}`}>
                  Auto Increment
                </FieldLabel>
                {fieldState.invalid && (
                  <FieldError errors={[fieldState.error]} />
                )}
              </Field>
            );
          }}
        />
      </div>

      <div className="grid grid-cols-3 gap-2">
        <Controller
          control={control}
          name={`inputs.${index}.isNull`}
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
          name={`inputs.${index}.isPk`}
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
          name={`inputs.${index}.isUnique`}
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
