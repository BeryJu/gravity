import { html, LitElement, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";
import { get } from "src/services/api";

@customElement("ddet-cluster-nodes")
export class ClusterNodePage extends LitElement {
    render(): TemplateResult {
        return html`
            ${until(
                get("/api/v0/etcd/members").then((res) => {
                    return res.map((member: any) => {
                        return html`${member.ID}: ${member.name}`;
                    });
                }),
            )}
        `;
    }
}
