import { DnsRecord, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DNSRecordForm";

@customElement("gravity-dns-records")
export class DNSRecordsPage extends TablePage<DnsRecord> {
    @property()
    zone?: string;

    pageTitle(): string {
        return `DNS Records for ${this.zone === "." ? "Root Zone" : this.zone}`;
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
    apiEndpoint(page: number): Promise<PaginatedResponse<DnsRecord>> {
        return new RolesDnsApi(DEFAULT_CONFIG)
            .dnsGetRecords({
                zone: this.zone || ".",
            })
            .then((records) => {
                const data = (records.records || []).filter(
                    (l) =>
                        l.fqdn.toLowerCase().includes(this.search.toLowerCase()) ||
                        l.type.toLowerCase().includes(this.search.toLowerCase()) ||
                        l.data.includes(this.search),
                );
                data.sort((a, b) => {
                    if (a.fqdn > b.fqdn) return 1;
                    if (a.fqdn < b.fqdn) return -1;
                    return 0;
                });
                return PaginationWrapper(data);
            });
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Record Type"),
            new TableColumn("Data"),
            new TableColumn("Actions"),
        ];
    }

    row(item: DnsRecord): TemplateResult<1 | 2>[] {
        return [
            html`${item.hostname}${item.uid === "" ? html`` : html` (${item.uid})`}`,
            html`${item.type}`,
            html`${item.data}`,
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update Zone"} </span>
                <gravity-dns-record-form
                    slot="form"
                    zone=${ifDefined(this.zone)}
                    .instancePk=${item.hostname + item.uid}
                >
                </gravity-dns-record-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"DNS Record(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: DnsRecord) => {
                return [
                    { key: "Hostname", value: item.hostname },
                    { key: "Type", value: item.type },
                    { key: "Data", value: item.data },
                ];
            }}
            .delete=${(item: DnsRecord) => {
                return new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteRecords({
                    zone: this.zone,
                    hostname: item.hostname,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Create"} </span>
                <span slot="header"> ${"Create Record"} </span>
                <gravity-dns-record-form zone=${ifDefined(this.zone)} slot="form">
                </gravity-dns-record-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
            </ak-forms-modal>
        `;
    }
}
