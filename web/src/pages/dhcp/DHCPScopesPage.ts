import { DhcpScope, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DHCPScopeForm";

@customElement("gravity-dhcp-scopes")
export class DHCPScopesPage extends TablePage<DhcpScope> {
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

    apiEndpoint(page: number): Promise<PaginatedResponse<DhcpScope>> {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes().then((scopes) => {
            const data = (scopes.scopes || []).filter(
                (l) =>
                    l.scope?.toLowerCase().includes(this.search.toLowerCase()) ||
                    l.dns?.zone?.toLowerCase().includes(this.search.toLowerCase()) ||
                    l.subnetCidr?.includes(this.search),
            );
            data.sort((a, b) => {
                if ((a.scope || "") > (b.scope || "")) return 1;
                if ((a.scope || "") < (b.scope || "")) return -1;
                return 0;
            });
            return PaginationWrapper(data);
        });
    }

    columns(): TableColumn[] {
        return [new TableColumn("Scope"), new TableColumn("Subnet")];
    }
    row(item: DhcpScope): TemplateResult[] {
        return [
            html`<a href=${`#/dhcp/scopes/${item.scope}`}>${item.scope}</a>`,
            html`${item.subnetCidr}`,
        ];
    }
    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${`Create`} </span>
                <span slot="header"> ${`Create Scope`} </span>
                <gravity-dhcp-zone-form slot="form"> </gravity-dhcp-zone-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${`Create`}</button>
            </ak-forms-modal>
        `;
    }
}
