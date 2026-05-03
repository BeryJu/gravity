import { DnsAPIRecord, DnsAPIZone, RolesDnsApi } from "gravity-api";

import { TemplateResult, html, nothing } from "lit";
import { customElement, property, state } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DNSRecordForm";

@customElement("gravity-dns-records")
export class DNSRecordsPage extends TablePage<DnsAPIRecord> {
    @property()
    zone: string | undefined;

    @state()
    _zone?: DnsAPIZone;

    @state()
    zoneCanStoreRecords = true;

    @state()
    isReverseZone = false;

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
        return this.zoneCanStoreRecords;
    }

    async apiEndpoint(): Promise<PaginatedResponse<DnsAPIRecord>> {
        if ((this.zone || "").endsWith(".in-addr.arpa.")) {
            this.isReverseZone = true;
        }
        const zone = await new RolesDnsApi(DEFAULT_CONFIG).dnsGetZones({
            name: this.zone,
        });
        this._zone = zone.zones![0];
        this.zoneCanStoreRecords =
            (this._zone.handlerConfigs || []).filter((h) => h.type.toLowerCase() === "etcd")
                .length > 0;
        if (!this.zoneCanStoreRecords) {
            return PaginationWrapper([]);
        }
        const records = await new RolesDnsApi(DEFAULT_CONFIG).dnsGetRecords({
            zone: this.zone || ".",
        });
        const data = (records.records || []).filter(
            (l) =>
                l.fqdn.toLowerCase().includes(this.search.toLowerCase()) ||
                l.type.toLowerCase().includes(this.search.toLowerCase()) ||
                l.data.includes(this.search),
        );
        data.sort((a, b) => {
            if (a.fqdn > b.fqdn) return 1;
            if (a.fqdn < b.fqdn) return -1;
            return parseInt(a.uid) - parseInt(b.uid);
        });
        return PaginationWrapper(data);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Hostname"),
            new TableColumn("Record Type"),
            new TableColumn("Data"),
            new TableColumn("Actions"),
        ];
    }

    row(item: DnsAPIRecord): TemplateResult[] {
        return [
            html`${item.hostname}${item.uid === "" ? html`` : html` (${item.uid})`}`,
            html`${item.type}`,
            html`<pre>${item.data}</pre>`,
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
            .metadata=${(item: DnsAPIRecord) => {
                return [
                    { key: "Hostname", value: item.hostname },
                    { key: "Type", value: item.type },
                    { key: "Data", value: item.data },
                ];
            }}
            .delete=${(item: DnsAPIRecord) => {
                return new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteRecords({
                    zone: this.zone || "",
                    hostname: item.hostname,
                    uid: item.uid,
                    type: item.type,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return this.zoneCanStoreRecords
            ? html`
                  <ak-forms-modal submitKeepOpen="submit-keep-open">
                      <span slot="submit"> ${"Create"} </span>
                      <span slot="submit-keep-open"> ${"Create & stay open"} </span>
                      <span slot="header"> ${"Create Record"} </span>
                      <gravity-dns-record-form
                          zone=${ifDefined(this.zone)}
                          slot="form"
                          recordType=${this.isReverseZone ? "PTR" : "A"}
                      >
                      </gravity-dns-record-form>
                      <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
                  </ak-forms-modal>
              `
            : html``;
    }

    renderEmpty(inner?: TemplateResult): TemplateResult {
        return super.renderEmpty(html`
            ${inner
                ? inner
                : html`<ak-empty-state
                      icon=${this.zoneCanStoreRecords ? this.pageIcon() : "fa fa-times"}
                      header=${this.zoneCanStoreRecords
                          ? "No objects found."
                          : "Zone cannot store records."}
                  >
                      <div slot="body">
                          ${this.zoneCanStoreRecords
                              ? nothing
                              : html`<span
                                    >Zone is not configured with an <code>etcd</code> handler and
                                    cannot store records.</span
                                >`}
                          ${this.searchEnabled() ? this.renderEmptyClearSearch() : nothing}
                      </div>
                      <div slot="primary">
                          ${this.zoneCanStoreRecords ? this.renderObjectCreate() : nothing}
                      </div>
                  </ak-empty-state>`}
        `);
    }
}
