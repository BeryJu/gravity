import "@spectrum-web-components/status-light/sp-status-light.js";
import { DEFAULT_CONFIG } from "src/api/Config";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { EtcdMember, RolesEtcdApi } from "gravity-api";

import "../elements/Table";

@customElement("gravity-cluster-nodes")
export class ClusterNodePage extends LitElement {
    render(): TemplateResult {
        return html`
            <gravity-header>Cluster nodes</gravity-header>
            <sp-divider size="m"></sp-divider>
            <gravity-table
                .columns=${["Status", "ID", "Name"]}
                .data=${() => {
                    return new RolesEtcdApi(DEFAULT_CONFIG)
                        .etcdGetMembers()
                        .then((members) => members.members || []);
                }}
                .rowRender=${(item: EtcdMember) => {
                    return [
                        html`<sp-status-light size="m" variant="positive"></sp-status-light>`,
                        html`${item.id}`,
                        html`${item.name}`,
                    ];
                }}
            >
            </gravity-table>
        `;
    }
}
