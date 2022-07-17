import "@spectrum-web-components/sidenav/sp-sidenav-heading.js";
import "@spectrum-web-components/sidenav/sp-sidenav-item.js";
import "@spectrum-web-components/sidenav/sp-sidenav.js";
import "@spectrum-web-components/split-view/sp-split-view.js";
import "@spectrum-web-components/theme/scale-medium.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/theme/theme-darkest.js";
import "@spectrum-web-components/theme/theme-lightest.js";
import { Route } from "src/elements/router/Route";
import "src/elements/router/RouterOutlet";
import "src/pages/OverviewPage";

import { LitElement, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

export const ROUTES = [
    new Route(new RegExp("^/$")).redirect("/overview"),
    new Route(new RegExp("^/overview$"), async () => {
        await import("src/pages/OverviewPage");
        return html`<gravity-overview></gravity-overview>`;
    }),
    new Route(new RegExp("^/cluster/nodes$"), async () => {
        await import("src/pages/ClusterNodesPage");
        return html`<gravity-cluster-nodes></gravity-cluster-nodes>`;
    }),
    new Route(new RegExp("^/dns/zones$"), async () => {
        await import("src/pages/dns/DNSZonesPage");
        return html`<gravity-dns-zones></gravity-dns-zones>`;
    }),
    new Route(new RegExp("^/dns/zones/(?<zone>.*)$"), async (args) => {
        await import("src/pages/dns/DNSRecordsPage");
        return html`<gravity-dns-records zone=${args.zone}></gravity-dns-records>`;
    }),
    new Route(new RegExp("^/dhcp/scopes$"), async () => {
        await import("src/pages/dhcp/DHCPScopesPage");
        return html`<gravity-dhcp-scopes></gravity-dhcp-scopes>`;
    }),
    new Route(new RegExp("^/dhcp/scopes/(?<scope>.*)$"), async (args) => {
        await import("src/pages/dhcp/DHCPLeasesPage");
        return html`<gravity-dhcp-leases scope=${args.scope}></gravity-dhcp-leases>`;
    }),
];

@customElement("gravity-app")
export class App extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
            }
            sp-split-view,
            sp-sidenav {
                height: 100vh;
            }
        `;
    }

    renderSidebar(): TemplateResult {
        return html`<sp-sidenav variant="multilevel" value=${window.location.hash || "#/overview"}>
            <sp-sidenav-item
                value="#/overview"
                label="Overview"
                href="#/overview"
            ></sp-sidenav-item>
            <sp-sidenav-heading label="DNS">
                <sp-sidenav-item
                    value="#/dns/zones"
                    label="Zones"
                    href="#/dns/zones"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
            <sp-sidenav-heading label="DHCP">
                <sp-sidenav-item
                    value="#/dhcp/scopes"
                    label="Subnets"
                    href="#/dhcp/scopes"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
            <sp-sidenav-heading label="Discovery">
                <sp-sidenav-item
                    value="#/discovery/devices"
                    label="Devices"
                    href="#/discovery/devices"
                ></sp-sidenav-item>
                <sp-sidenav-item
                    value="#/discovery/subnets"
                    label="Subnets"
                    href="#/discovery/subnets"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
            <sp-sidenav-heading label="Backup">
                <sp-sidenav-item
                    value="#/backup/status"
                    label="Status"
                    href="#/backup/status"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
            <sp-sidenav-heading label="Cluster">
                <sp-sidenav-item
                    value="#/cluster/roles"
                    label="Instance Roles"
                    href="#/cluster/roles"
                ></sp-sidenav-item>
                <sp-sidenav-item
                    value="#/cluster/nodes"
                    label="Nodes"
                    href="#/cluster/nodes"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
        </sp-sidenav>`;
    }

    render() {
        let theme = "lightest";
        if (window.location.search.includes("dark")) {
            theme = "darkest";
        }
        return html`
            <sp-theme theme="classic" scale="medium" color=${theme}>
                <sp-split-view primary-min="50" secondary-min="240" primary-size="240">
                    ${this.renderSidebar()}
                    <gravity-router-outlet .routes=${ROUTES}> </gravity-router-outlet>
                </sp-split-view>
            </sp-theme>
        `;
    }
}
