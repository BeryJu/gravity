import { PaginatedResponse } from "./elements/table/Table";

export function PaginationWrapper<T>(items: T[]): PaginatedResponse<T> {
    return {
        pagination: {
            count: items.length,
            current: 1,
            totalPages: 1,
            startIndex: 1,
            endIndex: items.length,
        },
        results: items,
    };
}
