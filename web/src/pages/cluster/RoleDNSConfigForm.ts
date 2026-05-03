import { DnsRoleConfig, RolesDnsApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-dns-config")
export class RoleDNSConfigForm extends ModelForm<DnsRoleConfig, string> {
    async loadInstance(): Promise<DnsRoleConfig> {
        const config = await new RolesDnsApi(DEFAULT_CONFIG).dnsGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: DnsRoleConfig): Promise<unknown> => {
        return new RolesDnsApi(DEFAULT_CONFIG).dnsPutRoleConfig({
            dnsAPIRoleConfigInput: { config: data },
        });
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label="Port" required name="port">
            <input
                type="number"
                value=${first(this.instance?.port, 53)}
                class="pf-c-form-control"
                required
            />
        </ak-form-element-horizontal>`;
    }
}
