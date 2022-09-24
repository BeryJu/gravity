import { RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "./api/Config";
import { AKElement } from "./elements/Base";

@customElement("gravity-login")
export class LoginPage extends AKElement {
    render(): TemplateResult {
        return html`
            <from>
                <button @click=${() => {
                    new RolesApiApi(DEFAULT_CONFIG)
                        .apiUsersLogin({
                            authUserLoginInput: {
                                username: "jens",
                                password: "foo",
                            },
                        })
                        .then((a) => {
                            if (a.successful) {
                                window.location.href = "#/";
                            }
                        });
                }}>Ye</button>
            </form>
        `;
    }
}
