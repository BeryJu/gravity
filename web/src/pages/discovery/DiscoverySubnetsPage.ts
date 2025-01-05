import { DiscoveryAPISubnet, RolesDiscoveryApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { MessageLevel } from "../../common/messages";
import "../../elements/buttons/ActionButton";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { showMessage } from "../../elements/messages/MessageContainer";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DiscoverySubnetForm";

@customElement("gravity-discovery-subnets")
export class DiscoverySubnetsPage extends TablePage<DiscoveryAPISubnet> {
    pageTitle(): string {
        return "Discovery subnets";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    async apiEndpoint(): Promise<PaginatedResponse<DiscoveryAPISubnet>> {
        const subnets = await new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryGetSubnets();
        const data = (subnets.subnets || []).filter((l) =>
            l.name.toLowerCase().includes(this.search.toLowerCase()),
        );
        data.sort((a, b) => {
            if (a.name > b.name) return 1;
            if (a.name < b.name) return -1;
            return 0;
        });
        return PaginationWrapper(data);
    }

    columns(): TableColumn[] {
        return [new TableColumn("Name"), new TableColumn("CIDR"), new TableColumn("Actions")];
    }

    row(item: DiscoveryAPISubnet): TemplateResult[] {
        return [
            html`${item.name}`,
            html`${item.subnetCidr}`,
            html`<ak-forms-modal>
                    <span slot="submit"> ${"Update"} </span>
                    <span slot="header"> ${"Update Subnet"} </span>
                    <gravity-discovery-subnet-form slot="form" .instancePk=${item.name}>
                    </gravity-discovery-subnet-form>
                    <button slot="trigger" class="pf-c-button pf-m-plain">
                        <i class="fas fa-edit"></i>
                    </button>
                </ak-forms-modal>
                <ak-action-button
                    class="pf-m-plain"
                    .apiRequest=${async () => {
                        await new RolesDiscoveryApi(DEFAULT_CONFIG).discoverySubnetStart({
                            identifier: item.name,
                        });
                        showMessage({
                            level: MessageLevel.info,
                            message: "Started discovery",
                        });
                    }}
                >
                    <i class="fas fa-play" aria-hidden="true"></i>
                </ak-action-button>`,
        ];
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Create"} </span>
                <span slot="header"> ${"Create Subnet"} </span>
                <gravity-discovery-subnet-form slot="form"> </gravity-discovery-subnet-form>
                <button slot="trigger" class="pf-c-button pf-m-primary">${"Create"}</button>
            </ak-forms-modal>
        `;
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"Discovery subnets(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: DiscoveryAPISubnet) => {
                return [{ key: "Name", value: item.name }];
            }}
            .delete=${(item: DiscoveryAPISubnet) => {
                return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryDeleteSubnets({
                    identifier: item.name,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }
}
