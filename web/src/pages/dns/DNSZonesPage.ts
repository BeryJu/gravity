import { DnsAPIZone, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DNSZoneForm";
import "./wizard/DNSZoneWizard";

@customElement("gravity-dns-zones")
export class DNSZonesPage extends TablePage<DnsAPIZone> {
    pageTitle(): string {
        return "DNS Zones";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    apiEndpoint(page: number): Promise<PaginatedResponse<DnsAPIZone>> {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones().then((zones) => {
            const data = (zones.zones || []).filter((l) =>
                l.name.toLowerCase().includes(this.search.toLowerCase()),
            );
            data.sort((a, b) => {
                if (a.name > b.name) return 1;
                if (a.name < b.name) return -1;
                return 0;
            });
            return PaginationWrapper(data);
        });
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Zone"),
            new TableColumn("Authoritative"),
            new TableColumn("Actions"),
        ];
    }

    row(item: DnsAPIZone): TemplateResult[] {
        return [
            html`<a href=${`#/dns/zones/${item.name}`}>
                ${item.name === "." ? "Root Zone" : item.name}
            </a>`,
            html`${item.authoritative ? "Yes" : "No"}`,
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update Zone"} </span>
                <gravity-dns-zone-form slot="form" .instancePk=${item.name}>
                </gravity-dns-zone-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html` <gravity-dns-zone-wizard></gravity-dns-zone-wizard> `;
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"DNS Zone(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: DnsAPIZone) => {
                return [{ key: "Name", value: item.name }];
            }}
            .delete=${(item: DnsAPIZone) => {
                return new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteZones({
                    zone: item.name,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }
}
