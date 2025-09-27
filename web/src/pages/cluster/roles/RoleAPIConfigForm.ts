import { ApiRoleConfig, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { first } from "../../../common/utils";
import { ModelForm } from "../../../elements/forms/ModelForm";
import { KV } from "../../../utils";
import "../../elements/forms/FormGroup";
import "../../elements/forms/HorizontalFormElement";

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
        data.port = parseInt(data.port as unknown as string);
        return new RolesApiApi(DEFAULT_CONFIG).apiPutRoleConfig({
            apiAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label="Port" required name="port">
                <input
                    type="number"
                    value="${first(this.instance?.port, 8008)}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Cookie Secret" required name="cookieSecret">
                <input
                    type="text"
                    value="${ifDefined(this.instance?.cookieSecret)}"
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Secret used to sign cookies.</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Session Duration"
                required
                name="sessionDuration"
                helperText="Duration for which a session is valid for."
            >
                <input
                    type="text"
                    value="${ifDefined(
                        this.instance?.sessionDuration === ""
                            ? "24h"
                            : this.instance?.sessionDuration,
                    )}"
                    required
                    class="pf-c-form-control"
                />
                <p class="pf-c-form__helper-text">
                    Duration sessions will be valid for. See
                    <a target="_blank" href="https://pkg.go.dev/time#ParseDuration">here</a> for a
                    format reference.
                </p>
            </ak-form-element-horizontal>
            <ak-form-group expanded>
                <span slot="header">OIDC</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal label="Issuer" required name="oidc.issuer">
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
                    <ak-form-element-horizontal label="Client ID" required name="oidc.clientID">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.clientID)}"
                            class="pf-c-form-control"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Client Secret"
                        required
                        name="oidc.clientSecret"
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
                        required
                        name="oidc.redirectURL"
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
                    <ak-form-element-horizontal label="Scopes" required name="oidc.scopesList">
                        <input
                            type="text"
                            value="${first(
                                this.instance?.oidc?.scopes?.join(" "),
                                "openid profile email",
                            )}"
                            class="pf-c-form-control"
                            required
                        />
                        <p class="pf-c-form__helper-text">
                            Space-separated list of OpenID scopes to request.
                        </p>
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Token username field"
                        required
                        name="oidc.tokenUsernameField"
                        helperText="Field in JWT tokens used to lookup user when using Token API authentication."
                    >
                        <input
                            type="text"
                            value="${first(this.instance?.oidc?.tokenUsernameField, "email")}"
                            required
                            class="pf-c-form-control"
                        />
                        <p class="pf-c-form__helper-text">
                            When using JWTs to authenticate to the API, the name of the field used
                            to lookup the user in Gravity.
                        </p>
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>`;
    }
}
