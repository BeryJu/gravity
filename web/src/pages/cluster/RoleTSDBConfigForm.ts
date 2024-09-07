import { RolesTsdbApi, TsdbRoleConfig } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-cluster-role-tsdb-config")
export class RoleTSDBConfigForm extends ModelForm<TsdbRoleConfig, string> {
    async loadInstance(): Promise<TsdbRoleConfig> {
        const config = await new RolesTsdbApi(DEFAULT_CONFIG).tsdbGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: TsdbRoleConfig): Promise<unknown> => {
        return new RolesTsdbApi(DEFAULT_CONFIG).tsdbPutRoleConfig({
            tsdbAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal name="enabled">
                <div class="pf-v6-c-check">
                    <input
                        type="checkbox"
                        class="pf-v6-c-check__input"
                        ?checked=${first(this.instance?.enabled, true)}
                    />
                    <label class="pf-v6-c-check__label"> ${"Enabled"} </label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Metrics expiry"
                ?required=${true}
                name="expire"
                helperText="Time in seconds before oldest metrics are deleted. Defaults to 30 minutes."
            >
                <input type="number" value="${first(this.instance?.expire, 60 * 30)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Scrape interval"
                ?required=${true}
                name="scrape"
                helperText="Interval in seconds of how often metrics are collected and written to the database."
            >
                <input type="number" value="${first(this.instance?.scrape, 30)}" required />
            </ak-form-element-horizontal>`;
    }
}
