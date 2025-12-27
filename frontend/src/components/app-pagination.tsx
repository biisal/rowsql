import * as React from "react";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationPrevious,
  PaginationNext,
  PaginationEllipsis,
} from "@/components/ui/pagination";
import { cn } from "@/lib/utils";

/**
 * Props for the AppPagination component
 */
interface AppPaginationProps {
  /** Current active page number (1-indexed) */
  currentPage: number;
  /** Callback function when page changes */
  onPageChange: (page: number) => void;
  /** Whether there is a next page available */
  hasNextPage: boolean;
  /** Whether there is a previous page available (defaults to currentPage > 1) */
  hasPreviousPage?: boolean;
  /** Total number of pages (optional, for showing last page) */
  totalPages?: number;
  /** Whether to show page numbers between prev/next buttons */
  showPageNumbers?: boolean;
  /** Additional CSS classes */
  className?: string;
}

/**
 * AppPagination - A reusable pagination component
 *
 * @example
 * ```tsx
 * <AppPagination
 *   currentPage={page}
 *   onPageChange={(newPage) => setSearchParams({ page: String(newPage) })}
 *   hasNextPage={data.length > 0}
 *   hasPreviousPage={page > 1}
 * />
 * ```
 */
export function AppPagination({
  currentPage,
  onPageChange,
  hasNextPage,
  hasPreviousPage = currentPage > 1,
  totalPages,
  showPageNumbers = true,
  className,
}: AppPaginationProps) {
  const handlePageChange = (page: number, e: React.MouseEvent) => {
    e.preventDefault();
    if (page > 0 && page <= (totalPages || Infinity)) {
      onPageChange(page);
    }
  };

  const renderPageNumbers = () => {
    if (!showPageNumbers) return null;

    const pages: React.ReactNode[] = [];

    // Show first page and ellipsis if current page > 2
    if (currentPage > 2) {
      pages.push(
        <PaginationItem key="page-1">
          <PaginationLink href="#" onClick={(e) => handlePageChange(1, e)}>
            1
          </PaginationLink>
        </PaginationItem>,
      );

      if (currentPage > 3) {
        pages.push(
          <PaginationItem key="ellipsis-start">
            <PaginationEllipsis />
          </PaginationItem>,
        );
      }
    }

    // Show previous page
    if (currentPage > 1) {
      pages.push(
        <PaginationItem key={`page-${currentPage - 1}`}>
          <PaginationLink
            href="#"
            onClick={(e) => handlePageChange(currentPage - 1, e)}
          >
            {currentPage - 1}
          </PaginationLink>
        </PaginationItem>,
      );
    }

    // Show current page
    pages.push(
      <PaginationItem key={`page-${currentPage}`}>
        <PaginationLink href="#" isActive onClick={(e) => e.preventDefault()}>
          {currentPage}
        </PaginationLink>
      </PaginationItem>,
    );

    // Show next page if available
    if (hasNextPage) {
      pages.push(
        <PaginationItem key={`page-${currentPage + 1}`}>
          <PaginationLink
            href="#"
            onClick={(e) => handlePageChange(currentPage + 1, e)}
          >
            {currentPage + 1}
          </PaginationLink>
        </PaginationItem>,
      );

      // Show ellipsis if we don't know total pages or if there are more pages
      if (!totalPages || currentPage + 1 < totalPages) {
        pages.push(
          <PaginationItem key="ellipsis-end">
            <PaginationEllipsis />
          </PaginationItem>,
        );
      }
    }

    // Show last page if we know total pages and we're not close to it
    if (totalPages && currentPage < totalPages - 1 && hasNextPage) {
      pages.push(
        <PaginationItem key={`page-${totalPages}`}>
          <PaginationLink
            href="#"
            onClick={(e) => handlePageChange(totalPages, e)}
          >
            {totalPages}
          </PaginationLink>
        </PaginationItem>,
      );
    }

    return pages;
  };

  return (
    <Pagination className={className}>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (hasPreviousPage) {
                onPageChange(currentPage - 1);
              }
            }}
            className={cn(!hasPreviousPage && "pointer-events-none opacity-50")}
          />
        </PaginationItem>

        {renderPageNumbers()}

        <PaginationItem>
          <PaginationNext
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (hasNextPage) {
                onPageChange(currentPage + 1);
              }
            }}
            className={cn(!hasNextPage && "pointer-events-none opacity-50")}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
