import { ChartData } from "chart.js";
import { RolesTsdbApi, TypesAPIMetricsGetOutput, TypesAPIMetricsRole } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { groupBy } from "../../../common/utils";
import { getColorFromString } from "../../../elements/charts/Chart";
import { AKChart } from "../../../elements/charts/Chart";

@customElement("gravity-overview-charts-cpu-usage")
export class CPUUsageChart extends AKChart<TypesAPIMetricsGetOutput> {
    apiRequest(): Promise<TypesAPIMetricsGetOutput> {
        return new RolesTsdbApi(DEFAULT_CONFIG).tsdbGetMetrics({
            role: TypesAPIMetricsRole.System,
            category: "cpu",
        });
    }

    getChartType(): string {
        return "line";
    }

    getChartData(data: TypesAPIMetricsGetOutput): ChartData {
        const chartData: ChartData = {
            datasets: [],
        };
        groupBy(data?.records || [], (record) => record.node).forEach(([node, records]) => {
            const background = getColorFromString(node);
            background.a = 0.3;
            chartData.datasets.push({
                label: node,
                borderColor: getColorFromString(node).toString(),
                backgroundColor: background.toString(),
                spanGaps: true,
                fill: "origin",
                cubicInterpolationMode: "monotone",
                tension: 0.4,
                data: records.map((record) => {
                    return {
                        x: record.time.getTime(),
                        y: Math.round(record.value),
                    };
                }),
            });
        });
        return chartData;
    }
}
