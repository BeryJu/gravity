import { ClusterInstancesApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-current-instance")
export class CurrentInstanceCard extends AdminStatusCard<string> {
    header = "Current instance";

    getPrimaryValue(): Promise<string> {
        return new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInfo().then((info) => {
            return info.currentInstanceIdentifier;
        });
    }

    getStatus(): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
