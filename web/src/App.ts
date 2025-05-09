import { RolesApiApi } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFDrawer from "@patternfly/patternfly/components/Drawer/drawer.css";
import PFPage from "@patternfly/patternfly/components/Page/page.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { DEFAULT_CONFIG } from "./api/Config";
import { EVENT_SIDEBAR_TOGGLE } from "./common/constants";
import { AKElement } from "./elements/Base";
import { Route } from "./elements/router/Route";
import "./elements/router/RouterOutlet";
import "./elements/sidebar/Sidebar";
import "./elements/sidebar/SidebarItem";
import "./pages/overview/OverviewPage";

export const ROUTES = [
    new Route(new RegExp("^/$")).redirect("/overview"),
    new Route(new RegExp("^/login$"), async () => {
        await import("./pages/LoginPage");
        return html`<gravity-login></gravity-login>`;
    }),
    new Route(new RegExp("^/overview$"), async () => {
        await import("./pages/overview/OverviewPage");
        return html`<gravity-overview></gravity-overview>`;
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
    new Route(new RegExp("^/discovery/devices$"), async () => {
        await import("./pages/discovery/DiscoveryDevicesPage");
        return html`<gravity-discovery-devices></gravity-discovery-devices>`;
    }),
    new Route(new RegExp("^/discovery/subnets$"), async () => {
        await import("./pages/discovery/DiscoverySubnetsPage");
        return html`<gravity-discovery-subnets></gravity-discovery-subnets>`;
    }),
    new Route(new RegExp("^/tftp/files$"), async () => {
        await import("./pages/tftp/TFTPFilesPage");
        return html`<gravity-tftp-files></gravity-tftp-files>`;
    }),
    new Route(new RegExp("^/cluster/roles$"), async () => {
        await import("./pages/cluster/RolesPage");
        return html`<gravity-cluster-roles></gravity-cluster-roles>`;
    }),
    new Route(new RegExp("^/cluster/nodes/logs$"), async () => {
        await import("./pages/cluster/ClusterNodeLogsPage");
        return html`<gravity-cluster-node-logs></gravity-cluster-node-logs>`;
    }),
    new Route(new RegExp("^/cluster/nodes$"), async () => {
        await import("./pages/cluster/ClusterNodesPage");
        return html`<gravity-cluster-nodes></gravity-cluster-nodes>`;
    }),
    new Route(new RegExp("^/auth/users$"), async () => {
        await import("./pages/auth/AuthUsersPage");
        return html`<gravity-auth-users></gravity-auth-users>`;
    }),
    new Route(new RegExp("^/auth/tokens$"), async () => {
        await import("./pages/auth/AuthTokensPage");
        return html`<gravity-auth-tokens></gravity-auth-tokens>`;
    }),
    new Route(new RegExp("^/tools$"), async () => {
        await import("./pages/tools/ToolPage");
        return html`<gravity-tools></gravity-tools>`;
    }),
];

@customElement("gravity-app")
export class AdminInterface extends AKElement {
    @state()
    showSidebar = true;

    @state()
    isAuthenticated = false;

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

    constructor() {
        super();
        this.showSidebar = window.innerWidth >= 1280;
        window.addEventListener("resize", () => {
            if (!this.isAuthenticated) return;
            this.showSidebar = window.innerWidth >= 1280;
        });
        window.addEventListener(EVENT_SIDEBAR_TOGGLE, () => {
            if (!this.isAuthenticated) return;
            this.showSidebar = !this.showSidebar;
        });
    }

    firstUpdated(): void {
        new RolesApiApi(DEFAULT_CONFIG).apiUsersMe().then((me) => {
            this.isAuthenticated = me.authenticated;
            if (!me.authenticated) {
                this.showSidebar = false;
                if (window.location.hash !== "#/login") {
                    window.location.hash = "#/login";
                    window.location.reload();
                }
            }
        });
    }

    render(): TemplateResult {
        if (!this.isAuthenticated) {
            return html`<ak-router-outlet
                role="main"
                class="pf-c-page__main"
                tabindex="-1"
                defaultUrl="/overview"
                .routes=${ROUTES}
            >
            </ak-router-outlet>`;
        }
        return html`<div class="pf-c-page">
            <ak-sidebar
                class="pf-c-page__sidebar ${this.showSidebar ? "pf-m-expanded" : "pf-m-collapsed"}"
            >
                ${this.renderSidebarItems()}
            </ak-sidebar>
            <main class="pf-c-page__main">
                <ak-router-outlet
                    role="main"
                    class="pf-c-page__main"
                    tabindex="-1"
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
            <ak-sidebar-item expanded>
                <span slot="label">DNS</span>
                <ak-sidebar-item path="/dns/zones" .activeWhen=${["^/dns/zones/(?<zone>.*)$"]}>
                    <span slot="label">Zones</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item expanded>
                <span slot="label">DHCP</span>
                <ak-sidebar-item path="/dhcp/scopes" .activeWhen=${["^/dhcp/scopes/(?<scope>.*)$"]}>
                    <span slot="label">Scopes</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${"TFTP"}</span>
                <ak-sidebar-item path="/tftp/files">
                    <span slot="label">Files</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${"Discovery"}</span>
                <ak-sidebar-item path="/discovery/devices">
                    <span slot="label">Devices</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/discovery/subnets">
                    <span slot="label">Subnets</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${"Cluster"}</span>
                <ak-sidebar-item path="/cluster/roles">
                    <span slot="label">Roles</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/cluster/nodes">
                    <span slot="label">Nodes</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/cluster/nodes/logs">
                    <span slot="label">Logs</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item>
                <span slot="label">${"Auth"}</span>
                <ak-sidebar-item path="/auth/users">
                    <span slot="label">Users</span>
                </ak-sidebar-item>
                <ak-sidebar-item path="/auth/tokens">
                    <span slot="label">Tokens</span>
                </ak-sidebar-item>
            </ak-sidebar-item>
            <ak-sidebar-item path="/tools">
                <span slot="label">Tools</span>
            </ak-sidebar-item>
        `;
    }
}
