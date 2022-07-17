import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DnsZone, RolesDnsApi } from "gravity-api";

import "../elements/Header";
import "../elements/Table";

@customElement("gravity-dns-zones")
export class DNSZonePage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>DNS Zones</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Zone", "Authoritative"]}
                .data=${() => {
                    return new RolesDnsApi(DEFAULT_CONFIG)
                        .dnsGetZones()
                        .then((zones) => zones.zones || []);
                }}
                .rowLink=${(item: DnsZone) => {
                    return `#/dns/zones/${item.name}`;
                }}
                .rowRender=${(item: DnsZone) => {
                    return [html`${item.name}`, html`${item.authoritative}`];
                }}
            >
            </gravity-table>
        `;
    }
}
