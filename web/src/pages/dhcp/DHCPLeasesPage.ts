import { DhcpLease, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
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
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    apiEndpoint(page: number): Promise<PaginatedResponse<DhcpLease>> {
        return new RolesDhcpApi(DEFAULT_CONFIG)
            .dhcpGetLeases({
                scope: this.scope,
            })
            .then((leases) => {
                const data = (leases.leases || []).filter(
                    (l) =>
                        l.hostname.toLowerCase().includes(this.search.toLowerCase()) ||
                        l.address.includes(this.search),
                );
                data.sort((a, b) => {
                    if (a.hostname > b.hostname) return 1;
                    if (a.hostname < b.hostname) return -1;
                    return 0;
                });
                return PaginationWrapper(data);
            });
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Address"),
            new TableColumn("Identifier"),
            new TableColumn("Actions"),
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${`DHCP Lease(s)`}
            .objects=${this.selectedElements}
            .metadata=${(item: DhcpLease) => {
                return [
                    { key: `Scope`, value: item.scopeKey },
                    { key: `Name`, value: item.hostname },
                    { key: `Address`, value: item.address },
                ];
            }}
            .delete=${(item: DhcpLease) => {
                return new RolesDhcpApi(DEFAULT_CONFIG).dhcpDeleteLeases({
                    identifier: item.identifier,
                    scope: item.scopeKey,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${`Delete`}
            </button>
        </ak-forms-delete-bulk>`;
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
