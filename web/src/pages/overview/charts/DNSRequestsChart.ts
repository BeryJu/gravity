import { ChartData } from "chart.js";
import { RolesTsdbApi, TypesAPIMetricsGetOutput, TypesAPIMetricsRole } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { groupBy } from "../../../common/utils";
import { getColorFromString } from "../../../elements/charts/Chart";
import { AKChart } from "../../../elements/charts/Chart";

@customElement("gravity-overview-charts-dns-requests")
export class DNSRequestsChart extends AKChart<TypesAPIMetricsGetOutput> {
    apiRequest(): Promise<TypesAPIMetricsGetOutput> {
        return new RolesTsdbApi(DEFAULT_CONFIG).tsdbGetMetrics({ role: TypesAPIMetricsRole.Dns });
    }

    getChartType(): string {
        return "line";
    }

    getChartData(data: TypesAPIMetricsGetOutput): ChartData {
        const chartData: ChartData = {
            datasets: [],
        };
        groupBy(data?.records || [], (record) => record.keys![1]).forEach(([handler, records]) => {
            const background = getColorFromString(handler);
            background.a = 0.3;
            chartData.datasets.push({
                label: handler,
                borderColor: getColorFromString(handler).toString(),
                backgroundColor: background.toString(),
                spanGaps: true,
                fill: "origin",
                cubicInterpolationMode: "monotone",
                tension: 0.4,
                data: records.map((record) => {
                    return {
                        x: record.time.getTime(),
                        y: record.value,
                    };
                }),
            });
        });
        return chartData;
    }
}
