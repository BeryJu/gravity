import { ClusterApi, InstanceInstanceInfo, TypesAPIMetricsRole } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import "../../elements/table/TableChart";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./wizard/ClusterJoinWizard";

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
        const inst = await new ClusterApi(DEFAULT_CONFIG).clusterGetClusterInfo();
        return PaginationWrapper(inst.instances || []);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Identifier"),
            new TableColumn("Roles"),
            new TableColumn("IP"),
            new TableColumn("Version"),
            new TableColumn(""),
        ];
    }

    row(item: InstanceInstanceInfo): TemplateResult[] {
        return [
            html`${item.identifier}`,
            html`<ak-chip-group
                >${item.roles?.map((role) => {
                    return html`<ak-chip>${role}</ak-chip>`;
                })}</ak-chip-group
            >`,
            html`${item.ip}`,
            html`${item.version}`,
            html`<gravity-table-chart
                role=${TypesAPIMetricsRole.System}
                category="cpu"
                .extraKeys=${[item.identifier]}
                legend="${item.identifier} CPU"
            ></gravity-table-chart>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html` <gravity-cluster-join-wizard></gravity-cluster-join-wizard> `;
    }
}
