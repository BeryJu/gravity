import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { InstanceInstanceInfo, InstancesApi } from "gravity-api";

import { DEFAULT_CONFIG } from "../api/Config";
import "../elements/chips/Chip";
import "../elements/chips/ChipGroup";
import { PaginatedResponse, TableColumn } from "../elements/table/Table";
import { TablePage } from "../elements/table/TablePage";
import { PaginationWrapper } from "../utils";

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
    apiEndpoint(page: number): Promise<PaginatedResponse<InstanceInstanceInfo>> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInstances().then((inst) => {
            return PaginationWrapper(inst.instances || []);
        });
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
