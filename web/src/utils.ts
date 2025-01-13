import { AbstractIPNum, IPv4, IPv6 } from "ip-num";

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

export function sortByIP<T>(getter: (item: T) => string): (a: T, b: T) => number {
    const getIP = (input: string) => {
        let ip: AbstractIPNum;
        try {
            ip = IPv6.fromString(input);
        } catch {
            ip = IPv4.fromString(input);
        }
        return ip;
    };
    return (a: T, b: T) => {
        const aIP = getIP(getter(a));
        const bIP = getIP(getter(b));
        if (aIP > bIP) return 1;
        if (aIP < bIP) return -1;
        return 0;
    };
}
