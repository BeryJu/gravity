import { ApiRoleConfig, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { first } from "../../common/utils";
import "../../elements/forms/FormGroup";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { KV } from "../../utils";

@customElement("gravity-cluster-role-api-config")
export class RoleAPIConfigForm extends ModelForm<ApiRoleConfig, string> {
    async loadInstance(): Promise<ApiRoleConfig> {
        const config = await new RolesApiApi(DEFAULT_CONFIG).apiGetRoleConfig();
        return config.config;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated role config.";
        } else {
            return "Successfully created role config.";
        }
    }

    send = (data: ApiRoleConfig): Promise<unknown> => {
        if (data.oidc) {
            const kv = data.oidc as unknown as KV;
            data.oidc.scopes = kv.scopesList.split(" ");
        }
        return new RolesApiApi(DEFAULT_CONFIG).apiPutRoleConfig({
            apiAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Port" ?required=${true} name="port">
                <input
                    type="number"
                    value="${first(this.instance?.port, 8008)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Cookie Secret" ?required=${true} name="cookieSecret">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.cookieSecret)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Secret used to sign cookies.</p>
            </ak-form-element-horizontal>
            <ak-form-group ?expanded=${true}>
                <span slot="header">OIDC</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal label="Issuer" ?required=${true} name="issuer">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.issuer)}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">
                            OpenID Issuer, sometimes also called "Configuration URL". Ensure
                            "".well-known/openid-configuration" suffix is removed.
                        </p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="Client ID" ?required=${true} name="clientID">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.clientID)}"
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Client Secret"
                        ?required=${true}
                        name="clientSecret"
                    >
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.clientSecret)}"
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Redirect URL"
                        ?required=${true}
                        name="redirectURL"
                    >
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.redirectURL)}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">
                            Redirect URL Gravity is reachable under, should end in
                            "/auth/oidc/callback".
                        </p>
                        <p class="pf-c-form__helper-text">
                            The placeholder '$INSTANCE_IDENTIFIER' will be replaced by the
                            instance's name and '$INSTANCE_IP' will be replaced by the instances IP.
                        </p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="Scopes" ?required=${true} name="scopesList">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.scopes?.join(" "))}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">
                            Space-separated list of OpenID scopes to request
                        </p>
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>
        </form>`;
    }
}
