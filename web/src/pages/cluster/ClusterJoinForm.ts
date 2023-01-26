import { ApiAPIMemberJoinInput, RolesEtcdApi } from "gravity-api";
import { TemplateResult, html } from "lit";

import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import { Form } from "../../elements/forms/Form";
import "../../elements/forms/HorizontalFormElement";

@customElement("gravity-cluster-join")
export class ClusterJoinForm extends Form<ApiAPIMemberJoinInput> {

    send = (data: ApiAPIMemberJoinInput): Promise<unknown> => {
        return new RolesEtcdApi(DEFAULT_CONFIG).etcdJoinMember({
            apiAPIMemberJoinInput: data,
        });
    };

    getSuccessMessage(): string {
        return "Successfully joined node.";
    }

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            <ak-form-element-horizontal label="Peer IP address" ?required=${true} name="peer">
                <input
                    type="text"
                    value=""
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">
                    Peer serer IP address
                </p>
            </ak-form-element-horizontal>
        </form>`;
    }
}
