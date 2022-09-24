import { DnsRecord, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-dns-records")
export class DNSRecordsPage extends TablePage<DnsRecord> {
    @property()
    zone?: string;

    pageTitle(): string {
        return `DNS Records for ${this.zone}`;
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;
    apiEndpoint(page: number): Promise<PaginatedResponse<DnsRecord>> {
        return new RolesDnsApi(DEFAULT_CONFIG)
            .dnsGetRecords({
                zone: this.zone || ".",
            })
            .then((records) => PaginationWrapper(records.records || []));
    }
    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Record Type"),
            new TableColumn("Data"),
        ];
    }
    row(item: DnsRecord): TemplateResult<1 | 2>[] {
        return [
            html`${item.hostname}${item.uid === "" ? html`` : html` (${item.uid})`}`,
            html`${item.type}`,
            html`${item.data}`,
        ];
    }
}
