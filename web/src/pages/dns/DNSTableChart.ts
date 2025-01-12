import { ChartData, ChartOptions } from "chart.js";
import { RolesTsdbApi, TypesAPIMetricsGetOutput, TypesAPIMetricsRole } from "gravity-api";

import { css } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKChart } from "../../elements/charts/Chart";

@customElement("gravity-dns-zone-chart")
export class DNSTableChart extends AKChart<TypesAPIMetricsGetOutput> {
    @property()
    zone?: string;

    static get styles() {
        return super.styles.concat(css`
            :host {
                display: flex;
                flex-direction: row;
                justify-content: end;
            }
            .container {
                width: 20rem;
                height: 3rem;
            }
        `);
    }

    async apiRequest(): Promise<TypesAPIMetricsGetOutput> {
        return new RolesTsdbApi(DEFAULT_CONFIG).tsdbGetMetrics({
            role: TypesAPIMetricsRole.Dns,
            category: "zones",
            extraKeys: [this.zone || ""],
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
        this.parentElement.style.width = "28rem";
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
            label: this.zone,
            backgroundColor: "rgba(0,0,0,0)",
            borderColor: "#3873e0",
            spanGaps: true,
            fill: "origin",
            cubicInterpolationMode: "monotone",
            tension: 0.4,
            pointStyle: false,
            data:
                data.records?.map((record) => {
                    return {
                        x: record.time.getTime(),
                        y: record.value,
                    };
                }) || [],
        });
        return chartData;
    }
}
