import { parse } from "ipaddr.js";

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

export function first<T>(items: T[] | null | undefined): T | undefined {
    if (items !== undefined && items !== null && items.length > 0) {
        return items[0];
    }
    return undefined;
}

export function ip2int(ip: string): number {
    return parse(ip)
        .toByteArray()
        .reduce((acc, c) => acc + c);
}
