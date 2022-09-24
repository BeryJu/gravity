import { DnsZone, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DNSZoneForm";

@customElement("gravity-dns-zones")
export class DNSZonesPage extends TablePage<DnsZone> {
    pageTitle(): string {
        return "DNS Zones";
    }
    pageDescription(): string | undefined {
        return "DNS Zones innit";
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    apiEndpoint(page: number): Promise<PaginatedResponse<DnsZone>> {
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
    row(item: DnsZone): TemplateResult[] {
        return [
            html`<a href=${`#/dns/zones/${item.name}`}>
                ${item.name === "." ? "Root Zone" : item.name}
            </a>`,
            html`${item.authoritative ? "Yes" : "No"}`,
            html`<ak-forms-modal>
                <span slot="submit"> ${`Update`} </span>
                <span slot="header"> ${`Update Zone`} </span>
                <gravity-dns-zone-form slot="form" .instancePk=${item.name}>
                </gravity-dns-zone-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${`Create`} </span>
                <span slot="header"> ${`Create Zone`} </span>
                <gravity-dns-zone-form slot="form"> </gravity-dns-zone-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${`Create`}</button>
            </ak-forms-modal>
        `;
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${`DNS Zone(s)`}
            .objects=${this.selectedElements}
            .metadata=${(item: DnsZone) => {
                return [{ key: `Name`, value: item.name }];
            }}
            .delete=${(item: DnsZone) => {
                return new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteZones({
                    zone: item.name,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${`Delete`}
            </button>
        </ak-forms-delete-bulk>`;
    }
}
