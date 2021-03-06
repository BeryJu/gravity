import "@spectrum-web-components/divider/sp-divider.js";
import { DEFAULT_CONFIG } from "src/api/Config";
import "src/elements/Header";
import "src/elements/Table";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DhcpScope, RolesDhcpApi } from "gravity-api";

@customElement("gravity-dhcp-scopes")
export class DHCPScopesPage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>DHCP Scopes</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Scope", "Subnet", ""]}
                .data=${() => {
                    return new RolesDhcpApi(DEFAULT_CONFIG)
                        .dhcpGetScopes()
                        .then((scopes) => scopes.scopes || []);
                }}
                .rowLink=${(item: DhcpScope) => {
                    return `#/dhcp/scopes/${item.scope}`;
                }}
                .rowRender=${(item: DhcpScope) => {
                    return [
                        html`${item.scope}`,
                        html`${item.subnetCidr}`,
                        html`<a href="foo">Edit</a>`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
