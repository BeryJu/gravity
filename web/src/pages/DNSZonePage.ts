import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

import { RolesDnsApi } from "gravity-api";

import "../elements/Header";
import "../elements/Table";

@customElement("gravity-dns-zones")
export class DNSZonePage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>DNS Zones</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table></gravity-table>
            ${until(
                new RolesDnsApi(DEFAULT_CONFIG).dnsZones().then((zones) => {
                    return zones.zones?.map((member: any) => {
                        return html`${JSON.stringify(member)}`;
                    });
                }),
            )}
        `;
    }
}
