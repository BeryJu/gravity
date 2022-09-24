import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DhcpLease, RolesDhcpApi } from "gravity-api";

import { DEFAULT_CONFIG } from "../../api/Config";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-dhcp-leases")
export class DHCPLeasesPage extends TablePage<DhcpLease> {
    @property()
    scope: string = "";

    pageTitle(): string {
        return `DHCP Leases for ${this.scope}`;
    }
    pageDescription(): string | undefined {
        return "";
    }
    pageIcon(): string {
        return "";
    }
    apiEndpoint(page: number): Promise<PaginatedResponse<DhcpLease>> {
        return new RolesDhcpApi(DEFAULT_CONFIG)
            .dhcpGetLeases({
                scope: this.scope,
            })
            .then((leases) => PaginationWrapper(leases.leases || []));
    }
    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Address"),
            new TableColumn("Identifier"),
            new TableColumn("Actions"),
        ];
    }
    row(item: DhcpLease): TemplateResult[] {
        return [
            html`${item.hostname}`,
            html`${item.address}`,
            html`${item.identifier}`,
            html`<sp-button
                size="m"
                @click=${() => {
                    new RolesDhcpApi(DEFAULT_CONFIG)
                        .dhcpWolLeases({
                            identifier: item.identifier || "",
                            scope: this.scope,
                        })
                        .then(() => {
                            alert("Successfully sent WOL");
                        })
                        .catch(() => {
                            alert("failed to send WOL");
                        });
                }}
                >WOL</sp-button
            >`,
        ];
    }
}
