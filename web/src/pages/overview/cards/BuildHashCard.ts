import { InstanceAPIInstanceInfo, InstancesApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-build-hash")
export class BuildHashCard extends AdminStatusCard<InstanceAPIInstanceInfo> {
    header = "Build Hash";

    getPrimaryValue(): Promise<InstanceAPIInstanceInfo> {
        return new InstancesApi(DEFAULT_CONFIG).rootGetInfo();
    }

    renderValue(): TemplateResult {
        return html`<a
            href="https://github.com/BeryJu/gravity/commit/${this.value?.buildHash}"
            target="_blank"
        >
            ${this.value?.buildHash.substring(0, 7)}
        </a>`;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    getStatus(value: InstanceAPIInstanceInfo): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
        });
    }
}
