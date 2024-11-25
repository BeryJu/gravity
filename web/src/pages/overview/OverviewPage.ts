import { AuthAPIMeOutput, RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFContent from "@patternfly/patternfly/components/Content/content.css";
import PFList from "@patternfly/patternfly/components/List/list.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFGrid from "@patternfly/patternfly/layouts/Grid/grid.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../../elements/Base";
import "../../elements/PageHeader";
import "../../elements/cards/AggregateCard";
import "./cards/BackupCard";
import "./cards/CurrentInstanceCard";
import "./cards/DHCPScopeCard";
import "./cards/DNSZoneCard";
import "./cards/VersionCard";
import "./charts/CPUUsageChart";
import "./charts/DHCPUsageChart";
import "./charts/DNSRequestsChart";
import "./charts/MemoryUsageChart";

@customElement("gravity-overview")
export class OverviewPage extends AKElement {
    @state()
    me: AuthAPIMeOutput | undefined;

    static get styles(): CSSResult[] {
        return [
            PFGrid,
            PFPage,
            PFContent,
            PFList,
            AKElement.GlobalStyle,
            css`
                .big-graph-container {
                    height: 35em;
                }
            `,
        ];
    }

    firstUpdated(): void {
        new RolesApiApi(DEFAULT_CONFIG).apiUsersMe().then((me) => (this.me = me));
    }

    render(): TemplateResult {
        return html` <ak-page-header>
                <span slot="header"> ${this.me ? html`Hello, ${this.me.username}` : html``} </span>
            </ak-page-header>
            <section class="pf-c-page__main-section">
                <div class="pf-l-grid pf-m-gutter">
                    <div class="pf-l-grid__item pf-m-6-col pf-m-2-col-on-2xl">
                        <gravity-overview-card-dhcp-scopes></gravity-overview-card-dhcp-scopes>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-2-col-on-2xl">
                        <gravity-overview-card-dns-zones></gravity-overview-card-dns-zones>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-2-col-on-2xl">
                        <gravity-overview-card-backup></gravity-overview-card-backup>
                    </div>
                    <div class="pf-l-grid__item pf-m-6-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-version></gravity-overview-card-version>
                    </div>
                    <div class="pf-l-grid__item pf-m-12-col pf-m-3-col-on-2xl">
                        <gravity-overview-card-current-instance></gravity-overview-card-current-instance>
                    </div>
                    <div
                        class="pf-l-grid__item pf-m-12-col pf-m-9-col-on-xl pf-m-9-col-on-2xl big-graph-container"
                    >
                        <ak-aggregate-card
                            icon="pf-icon pf-icon-server"
                            header="DNS requests per handler over the last 30 minutes"
                        >
                            <gravity-overview-charts-dns-requests></gravity-overview-charts-dns-requests>
                        </ak-aggregate-card>
                    </div>
                    <div
                        class="pf-l-grid__item pf-m-12-col pf-m-3-col-on-xl pf-m-3-col-on-2xl big-graph-container"
                    >
                        <ak-aggregate-card icon="pf-icon pf-icon-server" header="DHCP Scope usage">
                            <gravity-overview-charts-dhcp-usage></gravity-overview-charts-dhcp-usage>
                        </ak-aggregate-card>
                    </div>
                    <div
                        class="pf-l-grid__item pf-m-12-col pf-m-6-col-on-xl pf-m-6-col-on-2xl big-graph-container"
                    >
                        <ak-aggregate-card
                            icon="pf-icon pf-icon-server"
                            header="Memory usage per node (MB)"
                        >
                            <gravity-overview-charts-memory-usage></gravity-overview-charts-memory-usage>
                        </ak-aggregate-card>
                    </div>
                    <div
                        class="pf-l-grid__item pf-m-12-col pf-m-6-col-on-xl pf-m-6-col-on-2xl big-graph-container"
                    >
                        <ak-aggregate-card
                            icon="pf-icon pf-icon-server"
                            header="CPU usage per node (%)"
                        >
                            <gravity-overview-charts-cpu-usage></gravity-overview-charts-cpu-usage>
                        </ak-aggregate-card>
                    </div>
                </div>
            </section>`;
    }
}
