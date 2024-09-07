import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dns-wizard-forward")
export class ZoneForwarderWizardPage extends WizardFormPage {
    sidebarLabel = () => "Forwarder configuration";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const config = this.host.state["handlerConfigs"] as KeyUnknown[];
        const forwarderConfig = config.filter((config) =>
            (config.type as string).startsWith("forward_"),
        );
        forwarderConfig.forEach((conf) => {
            conf["to"] = data.to;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label=${"To"} ?required=${true} name="to">
            <input type="text" value="8.8.8.8:53" required />
            <p class="pf-c-form__helper-text">
                DNS Server to forward queries to, optionally specifying the port.
            </p>
            <p class="pf-c-form__helper-text">
                To specify multiple servers, separate their IPs with a semicolon.
            </p>
        </ak-form-element-horizontal>`;
    }
}
