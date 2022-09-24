import { InstanceInstanceInfo, InstancesApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

export interface Role {
    id: string;
    name: string;
}

export const roles: Role[] = [
    { id: "dhcp", name: "DHCP" },
    { id: "dns", name: "DNS" },
    { id: "api", name: "API" },
    { id: "discovery", name: "Discovery" },
    { id: "backup", name: "Backup" },
    { id: "monitoring", name: "Monitoring" },
    { id: "etcd", name: "etcd" },
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

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    apiEndpoint(page: number): Promise<PaginatedResponse<Role>> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInstances().then((inst) => {
            this.instances = inst.instances || [];
            return PaginationWrapper(roles);
        });
    }

    columns(): TableColumn[] {
        return [new TableColumn("Name"), new TableColumn("Nodes"), new TableColumn("Actions")];
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
            html``,
        ];
    }
}
