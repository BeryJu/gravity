import { DhcpAPIScope, RolesDhcpApi } from "gravity-api";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFProgress from "@patternfly/patternfly/components/Progress/progress.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { PaginatedResponse, Table, TableColumn } from "../../elements/table/Table";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-overview-dhcp-usage-table")
export class DHCPUsageTable extends Table<DhcpAPIScope & { statistics: { usage: number } }> {
    static get styles(): CSSResult[] {
        return super.styles.concat(PFProgress);
    }

    async apiEndpoint(): Promise<
        PaginatedResponse<DhcpAPIScope & { statistics: { usage: number } }>
    > {
        const scopes = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes();
        const data = (scopes.scopes || []).map((sc) => {
            const ssc = {
                ...sc,
                statistics: {
                    ...sc.statistics,
                    usage: Math.round((sc.statistics.used * 100) / sc.statistics.usable),
                },
            };
            return ssc;
        });
        data.sort((a, b) => {
            if (a.scope > b.scope) return 1;
            if (a.scope < b.scope) return -1;
            return 0;
        });
        return PaginationWrapper(data);
    }

    columns(): TableColumn[] {
        return [new TableColumn("Scope"), new TableColumn("Usage")];
    }

    row(item: DhcpAPIScope & { statistics: { usage: number } }): TemplateResult[] {
        return [
            html`${item.scope}`,
            html`<div class="pf-c-progress pf-m-sm pf-m-singleline">
                <div class="pf-c-progress__status" aria-hidden="true">
                    <span class="pf-c-progress__measure">${item.statistics.usage}%</span>
                </div>
                <div
                    class="pf-c-progress__bar"
                    role="progressbar"
                    aria-valuemin="0"
                    aria-valuemax="100"
                    aria-valuenow=${item.statistics.usage}
                    aria-labelledby="progress-sm-example-description"
                >
                    <div
                        class="pf-c-progress__indicator"
                        style="width:${item.statistics.usage}%;"
                    ></div>
                </div>
            </div>`,
        ];
    }

    renderToolbarContainer(): TemplateResult {
        return html``;
    }
}
