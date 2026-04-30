import { listRowsOptions } from "@/client/@tanstack/react-query.gen";
import { useQuery } from "@tanstack/react-query";

interface UseRowsProps {
  tableName: string;
  page: number;
  limit: number;
  column?: string;
  order?: "ASC" | "DESC";
}

export const useRows = ({ tableName, page, column, order }: UseRowsProps) => {
  return useQuery(
    listRowsOptions({
      path: { tableName },
      query: { column, page, order },
    }),
  );
};
