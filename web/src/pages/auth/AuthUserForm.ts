import { AuthAPIUser, AuthPermission, RolesApiApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { KV, firstElement } from "../../utils";

export const DEFAULT_ADMIN_PERMISSIONS: AuthPermission[] = [
    {
        path: "/*",
        methods: ["GET", "POST", "PUT", "HEAD", "DELETE"],
    },
];

@customElement("gravity-auth-user-form")
export class AuthUserForm extends ModelForm<AuthAPIUser, string> {
    async loadInstance(pk: string): Promise<AuthAPIUser> {
        const users = await new RolesApiApi(DEFAULT_CONFIG).apiGetUsers({
            username: pk,
        });
        const user = firstElement(users.users);
        if (!user) throw new Error("No user");
        return user;
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated user.";
        } else {
            return "Successfully created user.";
        }
    }

    send = (data: AuthAPIUser): Promise<void> => {
        return new RolesApiApi(DEFAULT_CONFIG).apiPutUsers({
            username: this.instance?.username || data.username,
            authAPIUsersPutInput: {
                password: (data as unknown as KV).password,
                permissions: data.permissions,
            },
        });
    };

    renderForm(): TemplateResult {
        return html` ${this.instance
                ? html``
                : html`<ak-form-element-horizontal label="Username" required name="username">
                      <input type="text" class="pf-c-form-control" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal
                label="Password"
                ?required=${this.instance === undefined}
                name="password"
            >
                <input type="password" class="pf-c-form-control" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"Permissions"} name="permissions">
                <ak-codemirror
                    mode="yaml"
                    value=${YAML.stringify(this.instance?.permissions || DEFAULT_ADMIN_PERMISSIONS)}
                >
                </ak-codemirror>
            </ak-form-element-horizontal>`;
    }
}
