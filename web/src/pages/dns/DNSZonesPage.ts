import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";
import "src/elements/Header";
import "src/elements/Table";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DnsZone, RolesDnsApi } from "gravity-api";

@customElement("gravity-dns-zones")
export class DNSZonesPage extends LitElement {
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
                .rowRender=${(item: DnsZone) => {
                    return [
                        html`<a href=${`#/dns/zones/${item.name}`}>${item.name}</a>`,
                        html`${item.authoritative}`,
                        html`<button @click=${() => {
                            if (confirm(`Delete zone ${item.name}?`)) {
                                new RolesDnsApi(DEFAULT_CONFIG).dnsDeleteZones({
                                    zone: item.name,
                                }).finally(() => {
                                    this.requestUpdate();
                                });
                            }
                        }}>x</button>`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
