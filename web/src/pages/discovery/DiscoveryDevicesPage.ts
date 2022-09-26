import { DiscoveryDevice, RolesDiscoveryApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./DiscoveryDeviceApply";

@customElement("gravity-discovery-devices")
export class DiscoveryDevicesPage extends TablePage<DiscoveryDevice> {
    pageTitle(): string {
        return "Discovered Devices";
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

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    apiEndpoint(page: number): Promise<PaginatedResponse<DiscoveryDevice>> {
        return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryGetDevices().then((devices) => {
            const data = (devices.devices || []).filter(
                (l) =>
                    l.hostname.toLowerCase().includes(this.search.toLowerCase()) ||
                    l.mac.toLowerCase().includes(this.search.toLowerCase()) ||
                    l.ip.includes(this.search),
            );
            data.sort((a, b) => {
                if (a.hostname > b.hostname) return 1;
                if (a.hostname < b.hostname) return -1;
                return 0;
            });
            return PaginationWrapper(data);
        });
    }

    columns(): TableColumn[] {
        return [new TableColumn("IP"), new TableColumn("Hostname"), new TableColumn("MAC")];
    }

    row(item: DiscoveryDevice): TemplateResult<1 | 2>[] {
        return [
            html`<pre>${item.ip}</pre>`,
            html`${item.hostname || "-"}`,
            html`<pre>${item.mac || "-"}</pre>`,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<gravity-discovery-apply
                objectLabel=${"Discovered Device(s)"}
                .objects=${this.selectedElements}
            >
                <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-primary">
                    ${"Apply"}
                </button> </gravity-discovery-apply
            >&nbsp;
            <ak-forms-delete-bulk
                objectLabel=${"Discovered Device(s)"}
                .objects=${this.selectedElements}
                .metadata=${(item: DiscoveryDevice) => {
                    return [
                        { key: "IP", value: item.ip },
                        { key: "Hostname", value: item.hostname || "-" },
                        { key: "MAC", value: item.mac || "-" },
                    ];
                }}
                .delete=${(item: DiscoveryDevice) => {
                    return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryDeleteDevices({
                        identifier: item.identifier,
                    });
                }}
            >
                <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                    ${"Delete"}
                </button>
            </ak-forms-delete-bulk> `;
    }
}
