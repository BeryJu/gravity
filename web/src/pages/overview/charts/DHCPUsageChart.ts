import { ChartData, ChartOptions } from "chart.js";
import { DhcpAPIScopesGetOutput, RolesDhcpApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AKChart, getColorFromString } from "../../../elements/charts/Chart";

@customElement("gravity-overview-charts-dhcp-usage")
export class DHCPUsageChart extends AKChart<DhcpAPIScopesGetOutput> {
    apiRequest(): Promise<DhcpAPIScopesGetOutput> {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes();
    }

    getChartType(): string {
        return "doughnut";
    }

    getOptions(): ChartOptions {
        return {
            maintainAspectRatio: false,
        };
    }

    getChartData(data: DhcpAPIScopesGetOutput): ChartData {
        return {
            labels: (data.scopes || []).map((scope) => scope.scope),
            datasets: (data.scopes || []).map((d) => {
                const usage = Math.round(d.statistics.used / (100 / d.statistics.usable));
                return {
                    backgroundColor: [getColorFromString(d.scope).toString(), "#ffffff"],
                    borderColor: "rgba(0,0,0,0)",
                    spanGaps: true,
                    data: [usage, 100 - usage],
                    label: d.scope,
                };
            }),
        };
    }
}
