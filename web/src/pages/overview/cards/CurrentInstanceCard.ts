import { ClusterInstancesApi, InstanceAPIInstanceInfo } from "gravity-api";

import { html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-current-instance")
export class CurrentInstanceCard extends AdminStatusCard<InstanceAPIInstanceInfo> {
    header = "Current instance";

    getPrimaryValue(): Promise<InstanceAPIInstanceInfo> {
        return new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInstanceInfo();
    }

    getStatus(data: InstanceAPIInstanceInfo): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
            message: html`${data.instanceIP}`,
        });
    }
    renderValue() {
        return html`${this.value?.instanceIdentifier}`;
    }
}
