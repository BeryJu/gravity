import { BackupRoleConfig, RolesBackupApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-backup-config")
export class RoleBackupConfigForm extends ModelForm<BackupRoleConfig, string> {
    async loadInstance(): Promise<BackupRoleConfig> {
        const config = await new RolesBackupApi(DEFAULT_CONFIG).backupGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: BackupRoleConfig): Promise<unknown> => {
        return new RolesBackupApi(DEFAULT_CONFIG).backupPutRoleConfig({
            backupAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Cron Schedule" ?required=${true} name="cronExpr">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.cronExpr)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Cron backup schedule</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Endpoint" ?required=${true} name="endpoint">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.endpoint)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">S3 Endpoint, including schema.</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Bucket" ?required=${true} name="bucket">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.bucket)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Access key" ?required=${true} name="accessKey">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.accessKey)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Secret Key" ?required=${true} name="secretKey">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.secretKey)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Path" ?required=${true} name="path">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.path)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
