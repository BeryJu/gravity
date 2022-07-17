import { DEFAULT_CONFIG } from "src/api/Config";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { InstanceInstanceInfo, InstancesApi } from "gravity-api";

import "../elements/Table";

@customElement("gravity-cluster-nodes")
export class ClusterNodePage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>Cluster nodes</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Identifier", "Roles", "IP", "Version"]}
                .data=${() => {
                    return new InstancesApi(DEFAULT_CONFIG)
                        .rootGetInstances()
                        .then((instances) => instances.instances || []);
                }}
                .rowRender=${(item: InstanceInstanceInfo) => {
                    return [
                        html`${item.identifier}`,
                        html`${item.roles}`,
                        html`${item.ip}`,
                        html`${item.version}`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
