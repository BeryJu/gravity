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
        data.port = parseInt(data.port as unknown as string);
        return new RolesApiApi(DEFAULT_CONFIG).apiPutRoleConfig({
            apiAPIRoleConfigInput: {
                config: data,
            },
        });
    };

    renderForm(): TemplateResult {
        return html` <ak-form-element-horizontal label="Port" required name="port">
                <input type="number" value="${first(this.instance?.port, 8008)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Cookie Secret"
                required
                name="cookieSecret"
                helperText="Secret used to sign cookies."
            >
                <input type="text" value="${ifDefined(this.instance?.cookieSecret)}" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label="Session Duration"
                ?required=${true}
                name="sessionDuration"
                helperText="Duration for which a session is valid for."
            >
                <input
                    type="text"
                    value="${first(this.instance?.sessionDuration, "24h")}"
                    required
                />
            </ak-form-element-horizontal>
            <ak-form-group ?expanded=${true}>
                <span slot="header">OIDC</span>
                <div slot="body" class="pf-c-form">
                    <ak-form-element-horizontal
                        label="Issuer"
                        required
                        name="oidc.issuer"
                        helperText='OpenID Issuer, sometimes also called "Configuration URL". Ensure "".well-known/openid-configuration" suffix is removed."'
                    >
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.issuer)}"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal label="Client ID" required name="oidc.clientID">
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.clientID)}"
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
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Redirect URL"
                        required
                        name="oidc.redirectURL"
                        helperText="Redirect URL Gravity is reachable under, should end in
                            '/auth/oidc/callback'. The placeholder '$INSTANCE_IDENTIFIER' will be replaced by the
                            instance's name and '$INSTANCE_IP' will be replaced by the instances IP."
                    >
                        <input
                            type="text"
                            value="${ifDefined(this.instance?.oidc?.redirectURL)}"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Scopes"
                        required
                        name="oidc.scopesList"
                        helperText="Space-separated list of OpenID scopes to request"
                    >
                        <input
                            type="text"
                            value="${first(
                                this.instance?.oidc?.scopes?.join(" "),
                                "openid profile email",
                            )}"
                            required
                        />
                    </ak-form-element-horizontal>
                    <ak-form-element-horizontal
                        label="Token username field"
                        ?required=${true}
                        name="oidc.tokenUsernameField"
                        helperText="Field in JWT tokens used to lookup user when using Token API authentication."
                    >
                        <input
                            type="text"
                            value="${first(this.instance?.oidc?.tokenUsernameField, "email")}"
                            required
                        />
                    </ak-form-element-horizontal>
                </div>
            </ak-form-group>`;
    }
}
