import { ClusterInstancesApi, InstanceInstanceInfo } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./RoleAPIConfigForm";
import "./RoleBackupConfigForm";
import "./RoleDHCPConfigForm";
import "./RoleDNSConfigForm";
import "./RoleDiscoveryConfigForm";
import "./RoleMonitoringConfigForm";
import "./RoleTSDBConfigForm";

export interface Role {
    id: string;
    name: string;
}

export const roles: Role[] = [
    { id: "api", name: "API" },
    { id: "backup", name: "Backup" },
    { id: "dhcp", name: "DHCP" },
    { id: "discovery", name: "Discovery" },
    { id: "dns", name: "DNS" },
    { id: "etcd", name: "etcd" },
    { id: "monitoring", name: "Monitoring" },
    { id: "tsdb", name: "TSDB" },
];

@customElement("gravity-cluster-roles")
export class RolesPage extends TablePage<Role> {
    @state()
    accessor instances: InstanceInstanceInfo[] = [];

    pageTitle(): string {
        return "Cluster Role configurations";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }

    async apiEndpoint(): Promise<PaginatedResponse<Role>> {
        const inst = await new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInstances();
        this.instances = inst.instances || [];
        return PaginationWrapper(roles);
    }

    columns(): TableColumn[] {
        return [new TableColumn("Name"), new TableColumn("Nodes"), new TableColumn("Actions")];
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
            default:
                return html`Not yet`;
        }
    }

    row(item: Role): TemplateResult[] {
        return [
            html`${item.name}`,
            html`<ak-chip-group
                >${this.instances
                    .filter((inst) => inst.roles?.includes(item.id))
                    .map((inst) => {
                        return html`<ak-chip>${inst.identifier}</ak-chip>`;
                    })}</ak-chip-group
            >`,
            html`${item.id === "etcd"
                ? html``
                : html`<ak-forms-modal>
                      <span slot="submit"> ${"Update"} </span>
                      <span slot="header"> ${"Update Role config"} </span>
                      ${this.renderRoleConfigForm(item)}
                      <button slot="trigger" class="pf-c-button pf-m-plain">
                          <i class="fas fa-edit"></i>
                      </button>
                  </ak-forms-modal>`}`,
        ];
    }
}
