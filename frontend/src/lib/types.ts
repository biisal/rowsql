import type { ColValue, VarianDataType } from "@/client/types.gen";

export interface Form {
  tableName: string;
  inputs: ColValue[];
  dataTypes: VarianDataType[];
  selectedDataType: VarianDataType;
}
