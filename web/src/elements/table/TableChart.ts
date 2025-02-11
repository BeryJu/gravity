import { ChartData, ChartOptions } from "chart.js";
import { RolesTsdbApi, TypesAPIMetricsGetOutput, TypesAPIMetricsRole } from "gravity-api";

import { css } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKChart } from "../charts/Chart";

@customElement("gravity-table-chart")
export class TableChart extends AKChart<TypesAPIMetricsGetOutput> {
    @property({ type: Array })
    extraKeys: string[] = [];

    @property({ attribute: "role" })
    metricRole?: TypesAPIMetricsRole;

    @property()
    category?: string;

    @property()
    label?: string;

    static get styles() {
        return super.styles.concat(css`
            :host {
                display: flex;
                flex-direction: row;
                justify-content: end;
            }
            .container {
                width: 25rem;
                height: 3rem;
            }
        `);
    }

    async apiRequest(): Promise<TypesAPIMetricsGetOutput> {
        return new RolesTsdbApi(DEFAULT_CONFIG).tsdbGetMetrics({
            role: this.role as TypesAPIMetricsRole,
            category: this.category,
            extraKeys: this.extraKeys,
        });
    }

    getChartType(): string {
        return "line";
    }

    firstUpdated(): void {
        super.firstUpdated();
        // This is a bit hacky but required to make the chart sit well in the table
        if (!this.parentElement) {
            return;
        }
        this.parentElement.style.width = "25rem";
        this.parentElement.style.padding = "0";
    }

    getOptions(): ChartOptions {
        return {
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false,
                },
            },
            layout: {
                padding: 0,
            },
            scales: {
                x: {
                    type: "time",
                    display: false,
                },
                y: {
                    type: "linear",
                    display: false,
                },
            },
        } as ChartOptions;
    }

    getChartData(data: TypesAPIMetricsGetOutput): ChartData {
        const chartData: ChartData = {
            datasets: [],
        };
        chartData.datasets.push({
            label: this.label ?? this.extraKeys[0],
            backgroundColor: "rgba(0,0,0,0)",
            borderColor: "#3873e0",
            spanGaps: true,
            fill: "origin",
            cubicInterpolationMode: "monotone",
            tension: 0.4,
            pointStyle: false,
            data:
                // Data is sorted by timestamp here as we might get data for multiple nodes
                // however we only want to show a single dataset so we've got to make sure the data is linear
                // otherwise we get a chart that's jumping around
                data.records?.sort((a, b) => a.time.getTime() - b.time.getTime()).map((record) => {
                    return {
                        x: record.time.getTime(),
                        y: record.value,
                    };
                }) || [],
        });
        return chartData;
    }
}
