import { AuthUser, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";
import { KV } from "../../utils";

@customElement("gravity-auth-user-form")
export class AuthUserForm extends ModelForm<AuthUser, string> {
    loadInstance(pk: string): Promise<AuthUser> {
        return new RolesApiApi(DEFAULT_CONFIG).apiGetUsers().then((users) => {
            const user = users.users?.find((z) => z.username === pk);
            if (!user) throw new Error("No user");
            return user;
        });
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated user.";
        } else {
            return "Successfully created user.";
        }
    }

    send = (data: AuthUser): Promise<void> => {
        return new RolesApiApi(DEFAULT_CONFIG).apiPutUsers({
            username: data.username,
            authAuthUsersPut: {
                password: (data as KV).password,
            },
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            ${this.instance
                ? html``
                : html`<ak-form-element-horizontal
                      label="Username"
                      ?required=${true}
                      name="username"
                  >
                      <input type="text" class="pf-c-form-control" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal label="Password" ?required=${true} name="password">
                <input type="password" class="pf-c-form-control" required />
            </ak-form-element-horizontal>
        </form>`;
    }
}
