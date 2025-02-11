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

export function firstElement<T>(items: T[] | null | undefined): T | undefined {
    if (items !== undefined && items !== null && items.length > 0) {
        return items[0];
    }
    return undefined;
}

export function first<T>(...args: Array<T | undefined | null>): T {
    for (let index = 0; index < args.length; index++) {
        const element = args[index];
        if (element !== undefined && element !== null) {
            return element;
        }
    }
    throw new Error(`No compatible arg given: ${args}`);
}

export function ip(raw: string): AbstractIPNum & { getValue(): bigint } {
    try {
        return IPv6.fromString(raw);
    } catch {
        try {
            return IPv4.fromString(raw);
        } catch {
            return new IPv4(0);
        }
    }
}

export function sortByIP<T>(getter: (item: T) => string): (a: T, b: T) => number {
    return (a: T, b: T) => {
        const aIP = ip(getter(a));
        const bIP = ip(getter(b));
        if (aIP.getValue() > bIP.getValue()) return 1;
        if (aIP.getValue() < bIP.getValue()) return -1;
        return 0;
    };
}
