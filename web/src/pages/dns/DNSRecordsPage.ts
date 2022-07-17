import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";
import "src/elements/Header";
import "src/elements/Table";

import { LitElement, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DnsRecord, RolesDnsApi } from "gravity-api";

@customElement("gravity-dns-records")
export class DNSRecordsPage extends LitElement {
    @property()
    zone?: string;

    render(): TemplateResult {
        return html`
            <gravity-header>DNS Records for ${this.zone}</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Hostname", "Type", "Data"]}
                .data=${() => {
                    return new RolesDnsApi(DEFAULT_CONFIG)
                        .dnsGetRecords({
                            zone: this.zone || ".",
                        })
                        .then((records) => records.records || []);
                }}
                .rowRender=${(item: DnsRecord) => {
                    return [
                        html`${item.hostname}${item.uid === "" ? html`` : html` (${item.uid})`}`,
                        html`${item.type}`,
                        html`${item.data}`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
