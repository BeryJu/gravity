import "@spectrum-web-components/button/sp-button.js";
import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";
import "src/elements/Header";
import "src/elements/Table";

import { LitElement, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { DhcpLease, RolesDhcpApi } from "gravity-api";

@customElement("gravity-dhcp-leases")
export class DHCPLeasesPage extends LitElement {
    @property()
    scope: string = "";

    render(): TemplateResult {
        return html`
            <gravity-header>DHCP Scopes for ${this.scope}</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Hostname", "Address", "Identifier", "Actions"]}
                .data=${() => {
                    return new RolesDhcpApi(DEFAULT_CONFIG)
                        .dhcpGetLeases({
                            scope: this.scope,
                        })
                        .then((leases) => leases.leases || []);
                }}
                .rowRender=${(item: DhcpLease) => {
                    return [
                        html`${item.hostname}`,
                        html`${item.address}`,
                        html`${item.identifier}`,
                        html`<sp-button
                            size="m"
                            @click=${() => {
                                new RolesDhcpApi(DEFAULT_CONFIG)
                                    .dhcpWolLeases({
                                        identifier: item.identifier || "",
                                        scope: this.scope,
                                    })
                                    .then(() => {
                                        alert("Successfully sent WOL");
                                    })
                                    .catch(() => {
                                        alert("failed to send WOL");
                                    });
                            }}
                            >WOL</sp-button
                        >`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
