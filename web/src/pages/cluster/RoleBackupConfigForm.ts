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
        return html` <ak-form-element-horizontal
                label="Cron Schedule"
                ?required=${true}
                name="cronExpr"
                helperText="Cron backup schedule"
            >
                <input type="text" value="${ifDefined(this.instance?.cronExpr)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Endpoint"
                ?required=${true}
                name="endpoint"
                helperText="S3 Endpoint, including schema."
            >
                <input type="text" value="${ifDefined(this.instance?.endpoint)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Bucket" ?required=${true} name="bucket">
                <input type="text" value="${ifDefined(this.instance?.bucket)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Access key" ?required=${true} name="accessKey">
                <input type="text" value="${ifDefined(this.instance?.accessKey)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Secret Key" ?required=${true} name="secretKey">
                <input type="text" value="${ifDefined(this.instance?.secretKey)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Path" ?required=${true} name="path">
                <input type="text" value="${ifDefined(this.instance?.path)}" required />
            </ak-form-element-horizontal>`;
    }
}
