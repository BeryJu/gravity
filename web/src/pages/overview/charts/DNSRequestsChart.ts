import { ChartData } from "chart.js";
import { RolesDnsApi, TypesAPIMetricsGetOutput } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { groupBy } from "../../../common/utils";
import { getColorFromString } from "../../../elements/charts/Chart";
import { AKChart } from "../../../elements/charts/Chart";

@customElement("gravity-overview-charts-dns-requests")
export class DNSRequestsChart extends AKChart<TypesAPIMetricsGetOutput> {
    apiRequest(): Promise<TypesAPIMetricsGetOutput> {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsGetMetrics();
    }

    getChartType(): string {
        return "line";
    }

    getChartData(data: TypesAPIMetricsGetOutput): ChartData {
        const chartData: ChartData = {
            datasets: [],
        };
        groupBy(data?.records || [], (record) => record.handler).forEach(([handler, records]) => {
            chartData.datasets.push({
                label: handler,
                backgroundColor: getColorFromString(handler),
                spanGaps: true,
                fill: true,
                data: records.map((record) => {
                    return {
                        x: parseInt(record.time, 10) * 1000,
                        y: record.value,
                    };
                }),
            });
        });
        return chartData;
    }
}
