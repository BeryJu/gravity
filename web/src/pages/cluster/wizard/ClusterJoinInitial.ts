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
        this.host.state["roles"] = roles.join(",");

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
        return html`
            <form class="pf-c-form pf-m-horizontal">
                <ak-form-element-horizontal label=${"Name"} ?required=${true} name="name">
                    <input type="text" value="" class="pf-c-form-control" required />
                    <p class="pf-c-form__helper-text">
                        The unique identifier of the node being added to the cluster.
                    </p>
                </ak-form-element-horizontal>
                <ak-form-element-horizontal label=${"Roles"} ?required=${true}>
                    ${this.roles.map((role) => {
                        return html`<div class="pf-c-check">
                            <input
                                type="checkbox"
                                class="pf-c-check__input"
                                ?checked=${true}
                                name=${`role_${role.toLowerCase()}`}
                            />
                            <label class="pf-c-check__label"> ${role} </label>
                        </div>`;
                    })}
                    <p class="pf-c-form__helper-text">
                        Select which roles the new node should provide
                    </p>
                </ak-form-element-horizontal>
            </form>
        `;
    }
}
