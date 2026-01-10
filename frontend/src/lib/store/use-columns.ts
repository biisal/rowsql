import api from '../axios';

export const getColumns = async (tableName: string) => {
	const data = await api.get(`/tables/${tableName}/columns`);
	return data.data;
};
