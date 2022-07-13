import "@spectrum-web-components/status-light/sp-status-light.js";
import { DEFAULT_CONFIG } from "src/api/Config";

import { LitElement, TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

import { RolesEtcdApi } from "gravity-api";

@customElement("gravity-cluster-nodes")
export class ClusterNodePage extends LitElement {
    render(): TemplateResult {
        return html`
            ${until(
                new RolesEtcdApi(DEFAULT_CONFIG)
                    .etcdGetMembers()
                    .then((members) => {
                        return members.members?.map((member: any) => {
                            return html`<sp-status-light size="m" variant="positive"
                                >${member.ID}: ${member.name}</sp-status-light
                            > `;
                        });
                    }),
            )}
        `;
    }
}
