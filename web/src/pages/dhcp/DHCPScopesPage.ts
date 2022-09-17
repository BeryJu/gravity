import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DhcpScope, RolesDhcpApi } from "@beryju/gravity-api";

import { DEFAULT_CONFIG } from "../../api/Config";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

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
}
