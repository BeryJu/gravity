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
                .columns=${["Scope", "Subnet"]}
                .data=${() => {
                    return new RolesDhcpApi(DEFAULT_CONFIG)
                        .dhcpGetScopes()
                        .then((scopes) => scopes.scopes || []);
                }}
                .rowRender=${(item: DhcpScope) => {
                    return [
                        html`<a href=${`#/dhcp/scopes/${item.scope}`}>${item.scope}</a>`,
                        html`${item.subnetCidr}`,
                        html`<button @click=${() => {
                            if (confirm(`Delete scope ${item.scope}?`)) {
                                new RolesDhcpApi(DEFAULT_CONFIG).dhcpDeleteScopes({
                                    scope: item.scope,
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
