import { ClusterInstancesApi, InstanceInstanceInfo } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-cluster-nodes")
export class ClusterNodePage extends TablePage<InstanceInstanceInfo> {
    pageTitle(): string {
        return "Cluster nodes";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }

    async apiEndpoint(): Promise<PaginatedResponse<InstanceInstanceInfo>> {
        const inst = await new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInstances();
        return PaginationWrapper(inst.instances || []);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Identifier"),
            new TableColumn("Roles"),
            new TableColumn("IP"),
            new TableColumn("Version"),
        ];
    }

    row(item: InstanceInstanceInfo): TemplateResult[] {
        return [
            html`${item.identifier}`,
            html`<ak-chip-group
                >${item.roles?.split(";").map((role) => {
                    return html`<ak-chip>${role}</ak-chip>`;
                })}</ak-chip-group
            >`,
            html`${item.ip}`,
            html`${item.version}`,
        ];
    }
}
