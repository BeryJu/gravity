import { ClusterApi, InstanceInstanceInfo } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import PFCard from "@patternfly/patternfly/components/Card/card.css";
import PFGrid from "@patternfly/patternfly/layouts/Grid/grid.css";

import { DEFAULT_CONFIG } from "../../api/Config";
import { AKElement } from "../../elements/Base";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import "../../elements/forms/ModalForm";
import { ModalForm } from "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./roles/RoleAPIConfigForm";
import "./roles/RoleBackupConfigForm";
import "./roles/RoleDHCPConfigForm";
import "./roles/RoleDNSConfigForm";
import "./roles/RoleDiscoveryConfigForm";
import "./roles/RoleMonitoringConfigForm";
import "./roles/RoleTFTPConfigForm";
import "./roles/RoleTSDBConfigForm";

export interface Role {
    id: string;
    name: string;
    settingsAvailable: boolean;
}

export const Roles: Role[] = [
    { id: "api", name: "API", settingsAvailable: true },
    { id: "backup", name: "Backup", settingsAvailable: true },
    { id: "dhcp", name: "DHCP", settingsAvailable: true },
    { id: "discovery", name: "Discovery", settingsAvailable: true },
    { id: "dns", name: "DNS", settingsAvailable: true },
    { id: "etcd", name: "etcd", settingsAvailable: false },
    { id: "monitoring", name: "Monitoring", settingsAvailable: true },
    { id: "tftp", name: "TFTP", settingsAvailable: true },
    { id: "tsdb", name: "TSDB", settingsAvailable: true },
];

@customElement("gravity-cluster-roles")
export class RolesPage extends TablePage<Role> {
    @state()
    instances: InstanceInstanceInfo[] = [];

    pageTitle(): string {
        return "Cluster Role configurations";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }

    static get styles(): CSSResult[] {
        return super.styles.concat(
            PFGrid,
            PFCard,
            css`
                .pf-c-sidebar__content {
                    background-color: transparent;
                }
            `,
            AKElement.GlobalStyle,
        );
    }

    async apiEndpoint(): Promise<PaginatedResponse<Role>> {
        const inst = await new ClusterApi(DEFAULT_CONFIG).clusterGetClusterInfo();
        this.instances = inst.instances || [];
        return PaginationWrapper(Roles);
    }

    columns(): TableColumn[] {
        return [];
    }

    renderRoleConfigForm(role: Role): TemplateResult {
        switch (role.id) {
            case "dns":
                return html`<gravity-cluster-role-dns-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-dns-config>`;
            case "dhcp":
                return html`<gravity-cluster-role-dhcp-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-dhcp-config>`;
            case "api":
                return html`<gravity-cluster-role-api-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-api-config>`;
            case "discovery":
                return html`<gravity-cluster-role-discovery-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-discovery-config>`;
            case "backup":
                return html`<gravity-cluster-role-backup-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-backup-config>`;
            case "monitoring":
                return html`<gravity-cluster-role-monitoring-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-monitoring-config>`;
            case "tsdb":
                return html`<gravity-cluster-role-tsdb-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-tsdb-config>`;
            case "tftp":
                return html`<gravity-cluster-role-tftp-config
                    slot="form"
                    .instancePk=${role.id}
                ></gravity-cluster-role-tftp-config>`;
            default:
                return html`Not yet`;
        }
    }

    row(): TemplateResult[] {
        return [];
    }

    render(): TemplateResult {
        return html`<ak-page-header icon=${this.pageIcon()} header=${this.pageTitle()}>
            </ak-page-header>
            <section class="pf-c-page__main-section pf-m-no-padding-mobile">
                <div class="pf-c-sidebar pf-m-gutter">
                    <div class="pf-c-sidebar__main">
                        ${this.renderSidebarBefore()}
                        <div class="pf-c-sidebar__content pf-l-grid pf-m-gutter">
                            ${this.data?.results.map((role) => {
                                const card = html` <div
                                    class="pf-c-card ${role.settingsAvailable
                                        ? "pf-m-hoverable-raised"
                                        : ""} pf-l-grid__item pf-m-3-col"
                                    @click=${() => {
                                        if (!role.settingsAvailable) {
                                            return;
                                        }
                                        const form = this.shadowRoot?.querySelector<ModalForm>(
                                            `#${role.id}`,
                                        );
                                        if (!form) {
                                            return;
                                        }
                                        form.onClick();
                                    }}
                                    slot="trigger"
                                >
                                    <div class="pf-c-card__title">${role.name}</div>
                                    <div class="pf-c-card__body">
                                        <ak-chip-group
                                            >${this.instances
                                                .filter((inst) => inst.roles?.includes(role.id))
                                                .map((inst) => {
                                                    return html`<ak-chip
                                                        >${inst.identifier}</ak-chip
                                                    >`;
                                                })}</ak-chip-group
                                        >
                                    </div>
                                </div>`;
                                return card;
                            })}
                        </div>
                        ${this.renderSidebarAfter()}
                    </div>
                    ${this.data?.results.map((role) => {
                        return html`<ak-forms-modal id="${role.id}">
                            <span slot="submit"> ${"Update"} </span>
                            <span slot="header"> ${`Update ${role.name} Role config`} </span>
                            ${this.renderRoleConfigForm(role)}
                        </ak-forms-modal>`;
                    })}
                </div>
            </section>`;
    }
}
