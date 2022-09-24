import { DhcpLease, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { MessageLevel } from "../../common/messages";
import "../../elements/buttons/SpinnerButton";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { showMessage } from "../../elements/messages/MessageContainer";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DHCPLeaseForm";

@customElement("gravity-dhcp-leases")
export class DHCPLeasesPage extends TablePage<DhcpLease> {
    @property()
    scope = "";

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

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
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
            objectLabel=${"DHCP Lease(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: DhcpLease) => {
                return [
                    { key: "Scope", value: item.scopeKey },
                    { key: "Name", value: item.hostname },
                    { key: "Address", value: item.address },
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
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    row(item: DhcpLease): TemplateResult[] {
        return [
            html`${item.hostname}`,
            html`${item.address}`,
            html`${item.identifier}`,
            html`<ak-forms-modal>
                    <span slot="submit"> ${"Update"} </span>
                    <span slot="header"> ${"Update Zone"} </span>
                    <gravity-dhcp-lease-form
                        slot="form"
                        scope=${this.scope}
                        .instancePk=${item.identifier}
                    >
                    </gravity-dhcp-lease-form>
                    <button slot="trigger" class="pf-c-button pf-m-plain">
                        <i class="fas fa-edit"></i>
                    </button> </ak-forms-modal
                ><ak-spinner-button
                    .callAction=${() => {
                        return new RolesDhcpApi(DEFAULT_CONFIG)
                            .dhcpWolLeases({
                                identifier: item.identifier || "",
                                scope: this.scope,
                            })
                            .then(() => {
                                showMessage({
                                    message: "Successfully sent WOL.",
                                    level: MessageLevel.success,
                                });
                            })
                            .catch((exc) => {
                                showMessage({
                                    message: exc.toString(),
                                    level: MessageLevel.error,
                                });
                            });
                    }}
                    class="pf-m-primary"
                    >WOL
                </ak-spinner-button>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Create"} </span>
                <span slot="header"> ${"Create lease"} </span>
                <gravity-dhcp-lease-form slot="form" scope=${this.scope}> </gravity-dhcp-lease-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
            </ak-forms-modal>
        `;
    }
}
