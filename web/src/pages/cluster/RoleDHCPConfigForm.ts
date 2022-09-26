import { DhcpRoleConfig, RolesDhcpApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-dhcp-config")
export class RoleDHCPConfigForm extends ModelForm<DhcpRoleConfig, string> {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    loadInstance(pk: string): Promise<DhcpRoleConfig> {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetRoleConfig().then((config) => config.config);
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: DhcpRoleConfig): Promise<unknown> => {
        console.log(data);
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutRoleConfig({
            dhcpRoleDHCPConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Port" ?required=${true} name="port">
                <input
                    type="number"
                    value="${first(this.instance?.port, 67)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
