import { DhcpAPIScope, RolesDhcpApi } from "gravity-api";

import { CSSResult, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFProgress from "@patternfly/patternfly/components/Progress/progress.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DHCPScopeForm";
import "./wizard/DHCPScopeWizard";

@customElement("gravity-dhcp-scopes")
export class DHCPScopesPage extends TablePage<DhcpAPIScope> {
    pageTitle(): string {
        return "DHCP Scopes";
    }
    pageDescription(): string | undefined {
        return "";
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    static get styles(): CSSResult[] {
        return super.styles.concat(PFProgress);
    }

    async apiEndpoint(): Promise<PaginatedResponse<DhcpAPIScope>> {
        const scopes = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes();
        const data = (scopes.scopes || []).filter(
            (l) =>
                l.scope.toLowerCase().includes(this.search.toLowerCase()) ||
                l.dns?.zone?.toLowerCase().includes(this.search.toLowerCase()) ||
                l.subnetCidr.includes(this.search),
        );
        data.sort((a, b) => {
            if (a.scope > b.scope) return 1;
            if (a.scope < b.scope) return -1;
            return 0;
        });
        return PaginationWrapper(data);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Scope"),
            new TableColumn("Subnet"),
            new TableColumn("Usage"),
            new TableColumn("Actions"),
        ];
    }

    row(item: DhcpAPIScope): TemplateResult[] {
        const usage = Math.round((item.statistics.used * 100) / item.statistics.usable);
        return [
            html`<a href=${`#/dhcp/scopes/${item.scope}`}>${item.scope}</a>`,
            html`<pre>${item.subnetCidr}</pre>`,
            html`<div class="pf-c-progress pf-m-sm">
                <div class="pf-c-progress__status" aria-hidden="true">
                    <span class="pf-c-progress__measure">${usage}%</span>
                </div>
                <div
                    class="pf-c-progress__bar"
                    role="progressbar"
                    aria-valuemin="0"
                    aria-valuemax="100"
                    aria-valuenow=${usage}
                    aria-labelledby="progress-sm-example-description"
                >
                    <div class="pf-c-progress__indicator" style="width:${usage}%;"></div>
                </div>
            </div>`,
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update Scope"} </span>
                <gravity-dhcp-scope-form slot="form" .instancePk=${item.scope}>
                </gravity-dhcp-scope-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"DHCP Scope(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: DhcpAPIScope) => {
                return [
                    { key: "Scope", value: item.scope },
                    { key: "CIDR", value: item.subnetCidr },
                ];
            }}
            .delete=${(item: DhcpAPIScope) => {
                return new RolesDhcpApi(DEFAULT_CONFIG).dhcpDeleteScopes({
                    scope: item.scope,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return html` <gravity-dhcp-scope-wizard></gravity-dhcp-scope-wizard> `;
    }
}
