import { html, LitElement, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";
import { get } from "src/services/api";
import "../elements/Table";
import "../elements/Header";
import "@spectrum-web-components/divider/sp-divider.js";

@customElement("gravity-dns-zones")
export class DNSZonePage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>DNS Zones</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table></gravity-table>
            ${until(
                get("/api/v0/dns/zones").then((res) => {
                    return res.map((member: any) => {
                        return html`${JSON.stringify(member)}`;
                    });
                }),
            )}
        `;
    }
}
