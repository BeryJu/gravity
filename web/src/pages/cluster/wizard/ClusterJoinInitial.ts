import { ClusterInstancesApi, RolesApiApi } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";
import { state } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-cluster-join-initial")
export class ClusterJoinInitial extends WizardFormPage {
    sidebarLabel = () => "Node details";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        this.host.state["identifier"] = data.name;

        const roles: string[] = [];
        const prefix = "role_";
        Object.keys(data).forEach((key) => {
            if (!key.startsWith(prefix)) {
                return;
            }
            if (!data[key]) {
                return;
            }
            roles.push(key.replaceAll(prefix, ""));
        });
        this.host.state["roles"] = roles.join(";");

        const info = await new ClusterInstancesApi(DEFAULT_CONFIG).clusterGetInfo();
        this.host.state["node_ip"] = info.currentInstanceIP;

        const user = await new RolesApiApi(DEFAULT_CONFIG).apiUsersMe();

        const token = await new RolesApiApi(DEFAULT_CONFIG).apiPutTokens({
            username: user.username,
        });
        this.host.state["join_token"] = token.key;
        return true;
    };

    @state()
    roles: string[] = [
        "DNS",
        "DHCP",
        "API",
        "etcd",
        "Discovery",
        "Backup",
        "Monitoring",
        "Debug",
        "TSDB",
    ];

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal
                label=${"Name"}
                required
                name="name"
                helperText="The unique identifier of the node being added to the cluster."
            >
                <input type="text" value="" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label=${"Roles"}
                required
                helperText="Select which roles the new node should provide."
                checkbox
            >
                ${this.roles.map((role) => {
                    return html`<div class="pf-v6-c-check">
                        <input
                            type="checkbox"
                            class="pf-v6-c-check__input"
                            ?checked=${true}
                            name=${`role_${role.toLowerCase()}`}
                        />
                        <label class="pf-v6-c-check__label"> ${role} </label>
                    </div>`;
                })}
            </ak-form-element-horizontal>`;
    }
}
