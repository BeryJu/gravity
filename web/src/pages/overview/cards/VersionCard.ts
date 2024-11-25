import {
    ClusterInstancesApi,
    InstanceAPIInstanceInfo,
    InstanceAPIInstancesOutput,
} from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-version")
export class VersionCard extends AdminStatusCard<InstanceAPIInstancesOutput> {
    header = "Version";
    headerLink = "#/cluster/nodes";

    clusterInfo?: InstanceAPIInstanceInfo;

    async getPrimaryValue(): Promise<InstanceAPIInstancesOutput> {
        this.clusterInfo = await new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInfo();
        return await new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInstances();
    }

    getStatus(value: InstanceAPIInstancesOutput): Promise<AdminStatus> {
        const matching =
            value.instances?.filter((inst) => {
                return inst.version === value.clusterVersion;
            }).length === value.instances?.length;
        if (!matching) {
            return Promise.resolve<AdminStatus>({
                icon: "fa fa-exclamation-triangle pf-m-warning",
                message: html`Mismatched version in cluster!`,
            });
        }
        return Promise.resolve<AdminStatus>({
            icon: "fa fa-check-circle pf-m-success",
            message: html`Matching versions across nodes.`,
        });
    }

    renderValue(): TemplateResult {
        return html`<a
            href="https://github.com/BeryJu/gravity/commit/${this.clusterInfo?.buildHash}"
            target="_blank"
        >
            ${this.value?.clusterVersion}
        </a>`;
    }
}
