import { AuthAPIToken, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { MessageLevel } from "../../common/messages";
import { Form } from "../../elements/forms/Form";
import "../../elements/forms/HorizontalFormElement";
import { showMessage } from "../../elements/messages/MessageContainer";

@customElement("gravity-auth-token-form")
export class AuthTokenForm extends Form<AuthAPIToken> {
    getSuccessMessage(): string {
        return "Successfully created token.";
    }

    send = (data: AuthAPIToken): Promise<void> => {
        return new RolesApiApi(DEFAULT_CONFIG)
            .apiPutTokens({
                username: data.username,
            })
            .then((out) => {
                showMessage({
                    level: MessageLevel.success,
                    message: out.key,
                });
            });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="User" ?required=${true} name="username">
                <select class="pf-c-form-control">
                    ${until(
                        new RolesApiApi(DEFAULT_CONFIG).apiGetUsers().then((users) => {
                            return users.users?.map((user) => {
                                return html`<option value="${user.username}">
                                    ${user.username}
                                </option>`;
                            });
                        }),
                    )}
                </select>
            </ak-form-element-horizontal>
        </form>`;
    }
}
