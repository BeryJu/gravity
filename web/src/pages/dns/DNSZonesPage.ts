import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DnsZone, RolesDnsApi } from "gravity-api";

import { DEFAULT_CONFIG } from "../../api/Config";
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
    apiEndpoint(page: number): Promise<PaginatedResponse<DnsZone>> {
        return new RolesDnsApi(DEFAULT_CONFIG)
            .dnsGetZones()
            .then((zones) => PaginationWrapper(zones.zones || []));
    }
    columns(): TableColumn[] {
        return [new TableColumn("Zone"), new TableColumn("Authoritative")];
    }
    row(item: DnsZone): TemplateResult[] {
        return [
            html`<a href=${`#/dns/zones/${item.name}`}>${item.name}</a>`,
            html`${item.authoritative}`,
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
}
