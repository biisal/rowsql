import type { ColMetaData  } from '@/client/types.gen';

export interface Column extends ColMetaData {
	columnName: string; // Keep for backward compatibility if needed, or remove?
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
	default: unknown,
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

export interface TableData {
	page: number;
	cols: ColMetaData[];
	activeTable: string;
	rows: CellValue[][];
	rowCount: number;
	hasNextPage: boolean;
	totalPages: number;
}
