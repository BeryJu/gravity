import { RolesTsdbApi, TsdbRoleConfig } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { first } from "../../../common/utils";
import { ModelForm } from "../../../elements/forms/ModelForm";
import "../../elements/forms/HorizontalFormElement";

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
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${first(this.instance?.enabled, true)}
                    />
                    <label class="pf-c-check__label"> ${"Enabled"} </label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Metrics expiry" required name="expire">
                <input
                    type="number"
                    value="${first(this.instance?.expire, 60 * 30)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    Time in seconds before oldest metrics are deleted. Defaults to 30 minutes.
                </p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Scrape interval" required name="scrape">
                <input
                    type="number"
                    value="${first(this.instance?.scrape, 30)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    Interval in seconds of how often metrics are collected and written to the
                    database.
                </p>
            </ak-form-element-horizontal>`;
    }
}
