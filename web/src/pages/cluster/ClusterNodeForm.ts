import { ClusterApi, ClusterInstancesApi, InstanceAPIClusterInfoOutput } from "gravity-api";
import { KeyUnknown } from "src/elements/forms/Form";
import { ModelForm } from "src/elements/forms/ModelForm";
import { Roles } from "src/pages/cluster/RolesPage";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/FormGroup";
import "../../elements/forms/HorizontalFormElement";

@customElement("gravity-cluster-node-form")
export class ClusterNodeForm extends ModelForm<InstanceAPIClusterInfoOutput, string> {
    async loadInstance(): Promise<InstanceAPIClusterInfoOutput> {
        const config = await new ClusterApi(DEFAULT_CONFIG).clusterGetClusterInfo();
        return config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated instance.";
        } else {
            return "Successfully created instance.";
        }
    }

    send = (data: InstanceAPIClusterInfoOutput): Promise<unknown> => {
        const d = data as unknown as KeyUnknown;
        return new ClusterInstancesApi(DEFAULT_CONFIG).clusterPutInstance({
            identifier: this.instancePk,
            instanceAPIInstancesPutInput: {
                roles: d.roles as string[],
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label=${"Roles"} required>
            ${Roles.map((role) => {
                return html`<div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        checked
                        name=${`role_${role.id}`}
                    />
                    <label class="pf-c-check__label"> ${role.name} </label>
                </div>`;
            })}
            <p class="pf-c-form__helper-text">Select which roles the new node should provide</p>
        </ak-form-element-horizontal>`;
    }
}
