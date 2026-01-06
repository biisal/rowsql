export interface Column {
	columnName: string;
	dataType: string;
}
export interface Table {
	tableName: string;
	tableSchema: string;
}

export interface FormInputType {
	type: string;
	hasSize: boolean;
	size?: number;
	hasBool?: boolean;
	hasAutoIncrement?: boolean;
	hasDefault?: boolean;
}

export interface Input {
	dataType: FormInputType;
	colName: string;
	isNull: boolean;
	isPk: boolean;
	isUnique: boolean;
}

export interface Form {
	selectedDataType: FormInputType;
	dataTypes: FormInputType[];
	tableName: string;
	inputs: Input[];
}

export interface DbDataTypes {
	numericType: FormInputType[];
	stringType: FormInputType[];
}

export interface ErrorResponse {
	error: string;
	success: boolean;
	status: number;
}

export type CellValue = string | number | boolean | null;
type RowData = CellValue[];

export interface TableData {
	page: number;
	tables: Table[];
	cols: Column[];
	activeTable: string;
	rows: RowData[];
}
