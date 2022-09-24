import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";



import { DhcpScope, RolesDhcpApi } from "gravity-api";



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
    apiEndpoint(page: number): Promise<PaginatedResponse<DhcpScope>> {
        return new RolesDhcpApi(DEFAULT_CONFIG)
            .dhcpGetScopes()
            .then((scopes) => PaginationWrapper(scopes.scopes || []));
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
