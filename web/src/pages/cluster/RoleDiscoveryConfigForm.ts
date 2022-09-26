import { DiscoveryRoleConfig, RolesDiscoveryApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-discovery-config")
export class RoleDiscoveryConfigForm extends ModelForm<DiscoveryRoleConfig, string> {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    loadInstance(pk: string): Promise<DiscoveryRoleConfig> {
        return new RolesDiscoveryApi(DEFAULT_CONFIG)
            .discoveryGetRoleConfig()
            .then((config) => config.config);
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
            discoveryRoleDiscoveryConfigInput: { config: data },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal ?required=${true} name="enabled">
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${first(this.instance?.enabled, true)}
                    />
                    <label class="pf-c-check__label"> ${"Enabled"} </label>
                </div>
            </ak-form-element-horizontal>
        </form>`;
    }
}
