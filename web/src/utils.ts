import { PaginatedResponse } from "./elements/table/Table";

export interface KV {
    [key: string]: string;
}

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
