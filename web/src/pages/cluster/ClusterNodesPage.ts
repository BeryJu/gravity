import {
    ClusterApi,
    EtcdAPIMembersOutput,
    InstanceInstanceInfo,
    RolesEtcdApi,
    TypesAPIMetricsRole,
} from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/chips/Chip";
import "../../elements/chips/ChipGroup";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import "../../elements/table/TableChart";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./ClusterNodeForm";
import "./wizard/ClusterJoinWizard";

@customElement("gravity-cluster-nodes")
export class ClusterNodePage extends TablePage<InstanceInstanceInfo> {
    @state()
    etcdNodes?: EtcdAPIMembersOutput;

    checkbox = true;

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
        this.etcdNodes = await new RolesEtcdApi(DEFAULT_CONFIG).etcdGetMembers();
        return PaginationWrapper(inst.instances || []);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Identifier"),
            new TableColumn("Roles"),
            new TableColumn("IP"),
            new TableColumn("Version"),
            new TableColumn("Actions"),
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
            html`<ak-forms-modal>
                <span slot="submit"> ${"Update"} </span>
                <span slot="header"> ${"Update Lease"} </span>
                <gravity-cluster-node-form slot="form" .instancePk=${item.identifier}>
                </gravity-cluster-node-form>
                <button slot="trigger" class="pf-c-button pf-m-plain">
                    <i class="fas fa-edit"></i>
                </button>
            </ak-forms-modal>`,
            html`<gravity-table-chart
                role=${TypesAPIMetricsRole.System}
                category="cpu"
                .extraKeys=${[item.identifier]}
                legend="${item.identifier} CPU"
            ></gravity-table-chart>`,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"Cluster Node(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: InstanceInstanceInfo) => {
                return [
                    { key: "Identifier", value: item.identifier },
                    { key: "IP", value: item.ip },
                ];
            }}
            .delete=${(item: InstanceInstanceInfo) => {
                const peerId =
                    this.etcdNodes?.members?.filter((member) => member.name === item.identifier) ||
                    [];
                if (peerId?.length < 1) {
                    return;
                }
                return new RolesEtcdApi(DEFAULT_CONFIG).etcdRemoveMember({
                    peerID: peerId[0].id!,
                });
            }}
        >
            <span slot="notice"> </span>
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return html`<gravity-cluster-join-wizard></gravity-cluster-join-wizard>`;
    }
}
