import { ClusterApi, InstanceAPIClusterInfoOutput } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { AdminStatus, AdminStatusCard } from "./AdminStatusCard";

@customElement("gravity-overview-card-version")
export class VersionCard extends AdminStatusCard<InstanceAPIClusterInfoOutput> {
    header = "Version";
    headerLink = "#/cluster/nodes";

    async getPrimaryValue(): Promise<InstanceAPIClusterInfoOutput> {
        return await new ClusterApi(DEFAULT_CONFIG).clusterGetClusterInfo();
    }

    getStatus(value: InstanceAPIClusterInfoOutput): Promise<AdminStatus> {
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
            href="https://github.com/BeryJu/gravity/releases/tag/v${this.value
                ?.clusterVersionShort}"
            target="_blank"
        >
            ${this.value?.clusterVersion}
        </a>`;
    }
}
