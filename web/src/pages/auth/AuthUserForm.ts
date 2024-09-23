import { AuthAPIUser, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { KV, first } from "../../utils";

@customElement("gravity-auth-user-form")
export class AuthUserForm extends ModelForm<AuthAPIUser, string> {
    async loadInstance(pk: string): Promise<AuthAPIUser> {
        const users = await new RolesApiApi(DEFAULT_CONFIG).apiGetUsers({
            username: pk,
        });
        const user = first(users.users);
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
            },
        });
    };

    renderForm(): TemplateResult {
        return html` ${this.instance
                ? html``
                : html`<ak-form-element-horizontal label="Username" required name="username">
                      <input type="text" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal label="Password" required name="password">
                <input type="password" required />
            </ak-form-element-horizontal>`;
    }
}
