import { InstancesApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-current-instance")
export class CurrentInstanceCard extends AdminStatusCard<string> {
    header = "Current instance";

    getPrimaryValue(): Promise<string> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInfo().then((info) => {
            return info.currentInstanceIdentifier;
        });
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: string): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
