import { DiscoveryRoleConfig, RolesDiscoveryApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { first } from "../../../common/utils";
import { ModelForm } from "../../../elements/forms/ModelForm";
import "../../elements/forms/HorizontalFormElement";

@customElement("gravity-cluster-role-discovery-config")
export class RoleDiscoveryConfigForm extends ModelForm<DiscoveryRoleConfig, string> {
    async loadInstance(): Promise<DiscoveryRoleConfig> {
        const config = await new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: DiscoveryRoleConfig): Promise<unknown> => {
        return new RolesDiscoveryApi(DEFAULT_CONFIG).discoveryPutRoleConfig({
            discoveryAPIRoleConfigInput: { config: data },
        });
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal name="enabled">
            <div class="pf-c-check">
                <input
                    type="checkbox"
                    class="pf-c-check__input"
                    ?checked=${first(this.instance?.enabled, true)}
                />
                <label class="pf-c-check__label"> ${"Enabled"} </label>
            </div>
        </ak-form-element-horizontal>`;
    }
}
