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

/**
 * @file Temporal utilitie for working with dates and times.
 */

/**
 * Duration in milliseconds for time units used by the `Intl.RelativeTimeFormat` API.
 */
export const Duration = {
    /**
     * The number of milliseconds in a year.
     */
    year: 1000 * 60 * 60 * 24 * 365,
    /**
     * The number of milliseconds in a month.
     */
    month: (24 * 60 * 60 * 1000 * 365) / 12,
    /**
     * The number of milliseconds in a day.
     */
    day: 1000 * 60 * 60 * 24,
    /**
     * The number of milliseconds in an hour.
     */
    hour: 1000 * 60 * 60,
    /**
     * The number of milliseconds in a minute.
     */
    minute: 1000 * 60,
    /**
     * The number of milliseconds in a second.
     */
    second: 1000,
} as const satisfies Partial<Record<Intl.RelativeTimeFormatUnit, number>>;

export type DurationUnit = keyof typeof Duration;

/**
 * The order of time units used by the `Intl.RelativeTimeFormat` API.
 */
const DurationGranularity = [
    "year",
    "month",
    "day",
    "hour",
    "minute",
    "second",
] as const satisfies DurationUnit[];

/**
 * Given two dates, return a human-readable string describing the time elapsed between them.
 */
export function formatElapsedTime(d1: Date, d2: Date = new Date()): string {
    const elapsed = d1.getTime() - d2.getTime();
    const rtf = new Intl.RelativeTimeFormat("default", { numeric: "auto" });

    for (const unit of DurationGranularity) {
        const duration = Duration[unit];

        if (Math.abs(elapsed) > duration || unit === "second") {
            let rounded = Math.round(elapsed / duration);

            if (!isFinite(rounded)) {
                rounded = 0;
            }

            return rtf.format(rounded, unit);
        }
    }
    return rtf.format(Math.round(elapsed / 1000), "second");
}
