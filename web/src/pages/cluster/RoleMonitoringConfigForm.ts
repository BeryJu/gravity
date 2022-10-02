import { MonitoringRoleConfig, RolesMonitoringApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-monitoring-config")
export class RoleMonitoringConfigForm extends ModelForm<MonitoringRoleConfig, string> {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    loadInstance(pk: string): Promise<MonitoringRoleConfig> {
        return new RolesMonitoringApi(DEFAULT_CONFIG)
            .monitoringGetRoleConfig()
            .then((config) => config.config);
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: MonitoringRoleConfig): Promise<unknown> => {
        return new RolesMonitoringApi(DEFAULT_CONFIG).monitoringPutRoleConfig({
            monitoringAPIRoleConfigInput: { config: data },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Port" ?required=${true} name="port">
                <input
                    type="number"
                    value="${first(this.instance?.port, 8009)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
