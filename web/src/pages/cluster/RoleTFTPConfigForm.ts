import { RolesTftpApi, TftpRoleConfig } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-tftp-config")
export class RoleTFTPConfigForm extends ModelForm<TftpRoleConfig, string> {
    async loadInstance(): Promise<TftpRoleConfig> {
        const config = await new RolesTftpApi(DEFAULT_CONFIG).tftpGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: TftpRoleConfig): Promise<unknown> => {
        return new RolesTftpApi(DEFAULT_CONFIG).tftpPutRoleConfig({
            tftpAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Port" ?required=${true} name="port">
                <input type="number" value="${first(this.instance?.port, 69)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal name="enableLocal">
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${first(this.instance?.enableLocal, false)}
                    />
                    <label class="pf-c-check__label"> ${"Enable Local access"} </label>
                </div>
                <p class="pf-c-form__helper-text">Enable access to node-local files.</p>
            </ak-form-element-horizontal>`;
    }
}
