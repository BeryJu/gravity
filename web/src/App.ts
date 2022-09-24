import { RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDrawer from "@patternfly/patternfly/components/Drawer/drawer.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { DEFAULT_CONFIG } from "./api/Config";
import { AKElement } from "./elements/Base";
import { Route } from "./elements/router/Route";
import "./elements/router/RouterOutlet";
import "./elements/sidebar/Sidebar";
import "./elements/sidebar/SidebarItem";
import "./pages/OverviewPage";

export const ROUTES = [
    new Route(new RegExp("^/$")).redirect("/overview"),
    new Route(new RegExp("^/login$"), async () => {
        await import("./LoginPage");
        return html`<gravity-login></gravity-login>`;
    }),
    new Route(new RegExp("^/overview$"), async () => {
        await import("./pages/OverviewPage");
        return html`<gravity-overview></gravity-overview>`;
    }),
    new Route(new RegExp("^/cluster/nodes$"), async () => {
        await import("./pages/ClusterNodesPage");
        return html`<gravity-cluster-nodes></gravity-cluster-nodes>`;
    }),
    new Route(new RegExp("^/dns/zones$"), async () => {
        await import("./pages/dns/DNSZonesPage");
        return html`<gravity-dns-zones></gravity-dns-zones>`;
    }),
    new Route(new RegExp("^/dns/zones/(?<zone>.*)$"), async (args) => {
        await import("./pages/dns/DNSRecordsPage");
        return html`<gravity-dns-records zone=${args.zone}></gravity-dns-records>`;
    }),
    new Route(new RegExp("^/dhcp/scopes$"), async () => {
        await import("./pages/dhcp/DHCPScopesPage");
        return html`<gravity-dhcp-scopes></gravity-dhcp-scopes>`;
    }),
    new Route(new RegExp("^/dhcp/scopes/(?<scope>.*)$"), async (args) => {
        await import("./pages/dhcp/DHCPLeasesPage");
        return html`<gravity-dhcp-leases scope=${args.scope}></gravity-dhcp-leases>`;
    }),
];

@customElement("gravity-app")
export class AdminInterface extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFPage,
            PFButton,
            PFDrawer,
            AKElement.GlobalStyle,
            css`
                .pf-c-page__main,
                .pf-c-drawer__content,
                .pf-c-page__drawer {
                    z-index: auto !important;
                    background-color: transparent;
                }
                .display-none {
                    display: none;
                }
                .pf-c-page {
                    background-color: var(--pf-c-page--BackgroundColor) !important;
                }
                @media (prefers-color-scheme: dark) {
                    /* Global page background colour */
                    .pf-c-page {
                        --pf-c-page--BackgroundColor: var(--ak-dark-background);
                    }
                }
            `,
        ];
    }

    firstUpdated(): void {
        new RolesApiApi(DEFAULT_CONFIG).apiUsersMe().then((me) => {
            if (!me.authenticated) {
                window.location.href = "#/login";
            }
        });
    }

    render(): TemplateResult {
        return html` <div class="pf-c-page">
            <ak-sidebar class="pf-c-page__sidebar pf-m-expanded">
                ${this.renderSidebarItems()}
            </ak-sidebar>
            <main class="pf-c-page__main">
                <ak-router-outlet
                    role="main"
                    class="pf-c-page__main"
                    tabindex="-1"
                    id="main-content"
                    defaultUrl="/overview"
                    .routes=${ROUTES}
                >
                </ak-router-outlet>
            </main>
        </div>`;
    }

    renderSidebarItems(): TemplateResult {
        return html`
            <ak-sidebar-item path="/overview">
                <span slot="label">Overview</span>
            </ak-sidebar-item>
            <ak-sidebar-item .expanded=${true}>
                <span slot="label">DNS</span>
                <ak-sidebar-item path="/dns/zones" .activeWhen=${[`^/dhcp/zones/(?<zone>.*)$`]}>
                    <span slot="label">Zones</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item .expanded=${true}>
                <span slot="label">DHCP</span>
                <ak-sidebar-item path="/dhcp/scopes" .activeWhen=${[`^/dhcp/scopes/(?<scope>.*)$`]}>
                    <span slot="label">Scopes</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${`Discovery`}</span>
                <ak-sidebar-item path="/discovery/devices">
                    <span slot="label">Devices</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/discovery/subnets">
                    <span slot="label">Subnets</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item .expanded=${true}>
                <span slot="label">Backup</span>
                <ak-sidebar-item path="/backup/status">
                    <span slot="label">Status</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${`Cluster`}</span>
                <ak-sidebar-item path="/cluster/roles">
                    <span slot="label">Instance Roles</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/cluster/nodes">
                    <span slot="label">Nodes</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
        `;
    }
}
