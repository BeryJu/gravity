import { DhcpAPILease, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html, nothing } from "lit";
import { customElement, property, state } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { MessageLevel } from "../../common/messages";
import "../../elements/buttons/SpinnerButton";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { showMessage } from "../../elements/messages/MessageContainer";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper, firstElement, formatElapsedTime, ip, sortByIP } from "../../utils";
import "./DHCPLeaseForm";

@customElement("gravity-dhcp-leases")
export class DHCPLeasesPage extends TablePage<DhcpAPILease> {
    @property()
    scope = "";

    @state()
    nextAddress?: string;

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

    async apiEndpoint(): Promise<PaginatedResponse<DhcpAPILease>> {
        const leases = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetLeases({
            scope: this.scope,
        });
        const data = (leases.leases || []).filter(
            (l) =>
                l.hostname.toLowerCase().includes(this.search.toLowerCase()) ||
                l.address.includes(this.search),
        );
        data.sort(sortByIP((i) => i.address));
        try {
            this.nextAddress = await this.getNextFreeIP(data);
        } catch {
            /* */
        }
        return PaginationWrapper(data);
    }

    // Figure out highest free address after IPAM start
    async getNextFreeIP(leases: DhcpAPILease[]) {
        // Fetch scope details
        const scopes = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes({
            name: this.scope,
        });
        const scope = firstElement(scopes.scopes);
        if (!scope || !scope?.ipam?.range_start) {
            return;
        }
        // Start of IPAM but previous by one IP
        const ipamStart = ip(scope.ipam.range_start as string).previousIPNumber();
        const afterIPAM = leases.filter((v) => ip(v.address) < ipamStart);
        while (true) {
            const nextIP = ipamStart.nextIPNumber();
            if (afterIPAM.filter((v) => ip(v.address) === nextIP).length > 0) {
                continue;
            }
            return nextIP.toString();
        }
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Address"),
            new TableColumn("Description"),
            new TableColumn("Identifier"),
            new TableColumn("Expiry"),
            new TableColumn("Actions"),
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
                objectLabel=${"DHCP Lease(s)"}
                .objects=${this.selectedElements}
                .metadata=${(item: DhcpAPILease) => {
                    return [
                        { key: "Scope", value: item.scopeKey },
                        { key: "Name", value: item.hostname },
                        { key: "Address", value: item.address },
                    ];
                }}
                .delete=${(item: DhcpAPILease) => {
                    return new RolesDhcpApi(DEFAULT_CONFIG).dhcpDeleteLeases({
                        identifier: item.identifier,
                        scope: item.scopeKey,
                    });
                }}
            >
                <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                    ${"Delete"}
                </button> </ak-forms-delete-bulk
            >&nbsp;
            <ak-spinner-button
                ?disabled=${disabled}
                .callAction=${async () => {
                    return Promise.all(
                        this.selectedElements.map((item) => {
                            item.expiry = -1;
                            return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutLeases({
                                identifier: item.identifier,
                                scope: this.scope,
                                dhcpAPILeasesPutInput: item,
                            });
                        }),
                    )
                        .then(() => {
                            showMessage({
                                message: `Successfully converted ${this.selectedElements.length} lease(s) to reservation(s).`,
                                level: MessageLevel.success,
                            });
                            this.fetch();
                        })
                        .catch((exc: Error) => {
                            showMessage({
                                message: (exc as Error).toString(),
                                level: MessageLevel.error,
                            });
                        });
                }}
                class="pf-m-secondary"
                >Turn to reservation </ak-spinner-button
            >&nbsp;<ak-spinner-button
                .callAction=${() => {
                    return Promise.all(
                        this.selectedElements.map((item) => {
                            return new RolesDhcpApi(DEFAULT_CONFIG).dhcpWolLeases({
                                identifier: item.identifier,
                                scope: this.scope,
                            });
                        }),
                    )
                        .then(() => {
                            showMessage({
                                message: `Successfully sent ${this.selectedElements.length} WOL packet(s).`,
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
                ?disabled=${disabled}
                class="pf-m-secondary"
                >WOL
            </ak-spinner-button>`;
    }

    row(item: DhcpAPILease): TemplateResult[] {
        return [
            html`${item.hostname}`,
            html`<pre>${item.address}</pre>`,
            html`${item.description !== "" ? item.description : "-"}`,
            html`<pre>${item.identifier}</pre>
                ${item.info ? html` (${item.info.vendor})` : nothing}`,
            html`${(item.expiry || 0) <= -1
                ? html`Reservation`
                : html`<div>${new Date((item.expiry || 0) * 1000).toLocaleString()}</div>
                      <small>${formatElapsedTime(new Date((item.expiry || 0) * 1000))}</small>`}`,
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update Lease"} </span>
                <gravity-dhcp-lease-form
                    slot="form"
                    scope=${this.scope}
                    .instancePk=${item.identifier}
                >
                </gravity-dhcp-lease-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Create"} </span>
                <span slot="header"> ${"Create lease"} </span>
                <gravity-dhcp-lease-form
                    .nextAddress=${this.nextAddress}
                    slot="form"
                    scope=${this.scope}
                >
                </gravity-dhcp-lease-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
            </ak-forms-modal>
        `;
    }
}
