import { InstancesApi } from "gravity-api";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-build-hash")
export class BuildHashCard extends AdminStatusCard<string> {
    header = "Build Hash";
    headerLink = "#/cluster/nodes";

    getPrimaryValue(): Promise<string> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInfo().then((info) => {
            return info.buildHash;
        });
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: string): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
