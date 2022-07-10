import { LitElement, html, css } from "lit";
import { customElement } from "lit/decorators.js";
import "./router";
import "./ddet-login";
import "@spectrum-web-components/theme/theme-light.js";
import "@spectrum-web-components/theme/scale-medium.js";
import "@spectrum-web-components/theme/sp-theme.js";

import "@spectrum-web-components/split-view/sp-split-view.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/sidenav/sp-sidenav.js";
import "@spectrum-web-components/sidenav/sp-sidenav-heading.js";
import "@spectrum-web-components/sidenav/sp-sidenav-item.js";
import { isLoggedIn } from "src/services/api";

@customElement("ddet-app")
export class App extends LitElement {

    static get styles() {
        return css`
            :host {
                display: block;
            }
            sp-split-view {
                height: 100vh;
            }
        `;
    }

    render() {
        return html`
            <sp-theme color="light" scale="medium">
                <sp-sidenav variant="multilevel" value="Overview">
                    <sp-sidenav-item
                        value="Overview"
                        label="Overview"
                        href="#/overview"
                    ></sp-sidenav-item>
                    <sp-sidenav-heading label="DNS">
                        <sp-sidenav-item
                            value="Zones"
                            label="Zones"
                            href="#/dns/zones"
                        ></sp-sidenav-item>
                        <sp-sidenav-item
                            value="Records"
                            label="Records"
                            href="#/dns/records"
                        ></sp-sidenav-item>
                    </sp-sidenav-heading>
                    <sp-sidenav-heading label="DHCP">
                        <sp-sidenav-item
                            value="Subnets"
                            label="Subnets"
                            href="#/dhcp/subnets"
                        ></sp-sidenav-item>
                    </sp-sidenav-heading>
                    <sp-sidenav-heading label="Discovery">
                        <sp-sidenav-item
                            value="Devices"
                            label="Devices"
                            href="#/discovery/devices"
                        ></sp-sidenav-item>
                        <sp-sidenav-item
                            value="Subnets"
                            label="Subnets"
                            href="#/discovery/subnets"
                        ></sp-sidenav-item>
                    </sp-sidenav-heading>
                    <sp-sidenav-heading label="Backup">
                        <sp-sidenav-item
                            value="Status"
                            label="Status"
                            href="#/backup/status"
                        ></sp-sidenav-item>
                    </sp-sidenav-heading>
                    <sp-sidenav-heading label="Cluster">
                        <sp-sidenav-item
                            value="Roles"
                            label="Roles"
                            href="#/cluster/roles"
                        ></sp-sidenav-item>
                        <sp-sidenav-item
                            value="Nodes"
                            label="Nodes"
                            href="#/cluster/nodes"
                        ></sp-sidenav-item>
                    </sp-sidenav-heading>
                </sp-sidenav>
                <ddet-router></ddet-router>
            </sp-theme>
        `;
    }
}
