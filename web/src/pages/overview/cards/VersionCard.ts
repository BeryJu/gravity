import { ClusterInstancesApi, InstanceAPIInstanceInfo } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-version")
export class VersionCard extends AdminStatusCard<InstanceAPIInstanceInfo> {
    header = "Version";
    headerLink = "#/cluster/nodes";

    getPrimaryValue(): Promise<InstanceAPIInstanceInfo> {
        return new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInfo();
    }

    getStatus(value: InstanceAPIInstanceInfo): Promise<AdminStatus> {
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
            message: html`${value?.buildHash.substring(0, 7)}`,
        });
    }

    renderValue(): TemplateResult {
        return html`<a
            href="https://github.com/BeryJu/gravity/commit/${this.value?.buildHash}"
            target="_blank"
        >
            ${this.value?.version}
        </a>`;
    }
}
