import React from 'react';
import { Button } from './button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './select';
import { ChevronLeft, ChevronRight } from 'lucide-react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  pageSizeOptions?: number[];
  showPageSizeSelector?: boolean;
  showTotalInfo?: boolean;
}

export function Pagination({
  currentPage,
  totalPages,
  totalItems,
  pageSize,
  onPageChange,
  onPageSizeChange,
  pageSizeOptions = [10, 20, 50, 100],
  showPageSizeSelector = true,
  showTotalInfo = true,
}: PaginationProps) {
  const startItem = totalItems === 0 ? 0 : (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

  // Always show pagination for now (disable hiding check)
  // if (totalPages <= 1 || totalItems <= pageSize) {
  //     return null;
  // }

  return (
    <div className="w-full">
      {/* Mobile-first: Simple horizontal layout */}
      <div className="flex flex-col space-y-4 sm:flex-row md:flex-col md:gap-y-2 sm:items-center sm:justify-between sm:space-y-0">
        {/* Total info - Always visible but compact on mobile */}
        {showTotalInfo && (
          <div className="text-sm text-gray-600">
            <span className="sm:hidden">
              Page {currentPage} of {totalPages} ({totalItems} total)
            </span>
            <span className="hidden sm:inline">
              Showing {startItem} to {endItem} of {totalItems} series
            </span>
          </div>
        )}

        {/* Pagination controls - Mobile optimized */}
        <div className="flex items-center justify-between">
          {/* Page size selector - Only on larger screens */}
          {showPageSizeSelector && (
            <div className="hidden items-center space-x-2 sm:flex">
              <span className="text-sm text-gray-700">Show:</span>
              <Select
                value={pageSize.toString()}
                onValueChange={value => onPageSizeChange(parseInt(value))}
              >
                <SelectTrigger className="w-17 mr-4">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {pageSizeOptions.map(size => (
                    <SelectItem key={size} value={size.toString()}>
                      {size}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}

          {/* Main pagination controls */}
          <div className="flex items-center space-x-1">
            {/* Previous button */}
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="h-9 px-3"
            >
              <ChevronLeft className="h-4 w-4" />
              <span className="ml-1 hidden xs:inline">Prev</span>
            </Button>

            {/* Page numbers - Responsive display */}
            <div className="flex items-center space-x-1">
              {/* Mobile: Show current page and nearby pages */}
              {totalPages <= 5 ? (
                // Show all pages if 5 or fewer
                Array.from({ length: totalPages }, (_, i) => i + 1).map(
                  pageNum => (
                    <Button
                      key={pageNum}
                      variant={currentPage === pageNum ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => onPageChange(pageNum)}
                      className="h-9 w-9"
                    >
                      {pageNum}
                    </Button>
                  )
                )
              ) : (
                // Show smart pagination for more pages
                <>
                  {/* Always show first page */}
                  <Button
                    variant={currentPage === 1 ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => onPageChange(1)}
                    className="h-9 w-9"
                  >
                    1
                  </Button>

                  {/* Show ellipsis if current page is far from start */}
                  {currentPage > 3 && (
                    <span className="px-2 text-sm text-gray-500">...</span>
                  )}

                  {/* Show pages around current page */}
                  {currentPage > 2 && currentPage < totalPages - 1 && (
                    <Button
                      variant="default"
                      size="sm"
                      onClick={() => onPageChange(currentPage)}
                      className="h-9 w-9"
                    >
                      {currentPage}
                    </Button>
                  )}

                  {/* Show ellipsis if current page is far from end */}
                  {currentPage < totalPages - 2 && (
                    <span className="px-2 text-sm text-gray-500">...</span>
                  )}

                  {/* Always show last page */}
                  <Button
                    variant={currentPage === totalPages ? 'default' : 'outline'}
                    size="sm"
                    onClick={() => onPageChange(totalPages)}
                    className="h-9 w-9"
                  >
                    {totalPages}
                  </Button>
                </>
              )}
            </div>

            {/* Next button */}
            <Button
              variant="outline"
              size="sm"
              onClick={() => onPageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="h-9 px-3"
            >
              <span className="mr-1 hidden xs:inline">Next</span>
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      {/* Mobile page size selector - Bottom row on mobile */}
      {showPageSizeSelector && (
        <div className="mt-3 flex items-center justify-center space-x-2 sm:hidden">
          <span className="text-sm text-gray-700">Items per page:</span>
          <Select
            value={pageSize.toString()}
            onValueChange={value => onPageSizeChange(parseInt(value))}
          >
            <SelectTrigger className="w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {pageSizeOptions.map(size => (
                <SelectItem key={size} value={size.toString()}>
                  {size}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      )}
    </div>
  );
}
