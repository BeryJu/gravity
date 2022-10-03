import { InstanceAPIInstanceInfo, InstancesApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-version")
export class VersionCard extends AdminStatusCard<InstanceAPIInstanceInfo> {
    header = "Version";
    headerLink = "#/cluster/nodes";

    getPrimaryValue(): Promise<InstanceAPIInstanceInfo> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInfo();
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: InstanceAPIInstanceInfo): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }

    renderValue(): TemplateResult {
        return html`<a
            href="https://github.com/BeryJu/gravity/commit/${this.value?.buildHash}"
            target="_blank"
        >
            ${this.value?.version} (${this.value?.buildHash.substring(0, 7)})
        </a>`;
    }
}
