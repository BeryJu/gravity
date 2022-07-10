import { LitElement, html, css, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import "./router";
import "./gravity-login";
import "@spectrum-web-components/theme/theme-light.js";
import "@spectrum-web-components/theme/theme-darkest.js";
import "@spectrum-web-components/theme/scale-medium.js";
import "@spectrum-web-components/theme/sp-theme.js";;

import "@spectrum-web-components/split-view/sp-split-view.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/sidenav/sp-sidenav.js";
import "@spectrum-web-components/sidenav/sp-sidenav-heading.js";
import "@spectrum-web-components/sidenav/sp-sidenav-item.js";
import { Route } from "./router";

export const ROUTES = [
    new Route("/overview", async () => {
        await import("./pages/OverviewPage");
        return html`<gravity-overview></gravity-overview>`;
    }),
    new Route("/cluster/nodes", async () => {
        await import("./pages/ClusterNodesPage");
        return html`<gravity-cluster-nodes></gravity-cluster-nodes>`;
    }),
    new Route("/dns/zones", async () => {
        await import("./pages/DNSZonePage");
        return html`<gravity-dns-zones></gravity-dns-zones>`;
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
                <sp-sidenav-item
                    value="#/dns/records"
                    label="Records"
                    href="#/dns/records"
                ></sp-sidenav-item>
            </sp-sidenav-heading>
            <sp-sidenav-heading label="DHCP">
                <sp-sidenav-item
                    value="#/dhcp/subnets"
                    label="Subnets"
                    href="#/dhcp/subnets"
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
        return html`
            <sp-theme theme="classic" scale="medium" color="darkest">
                <sp-split-view primary-min="50" secondary-min="240" primary-size="240">
                    ${this.renderSidebar()}
                    <gravity-router .routes=${ROUTES}> </gravity-router>
                </sp-split-view>
            </sp-theme>
        `;
    }
}
