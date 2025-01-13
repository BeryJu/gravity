import { DhcpAPIScopesPutInput } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dhcp-wizard-internal")
export class ScopeInternalWizardPage extends WizardFormPage {
    sidebarLabel = () => "Scope configuration";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const req = this.host.state["scopeReq"] as DhcpAPIScopesPutInput;
        req.subnetCidr = data.subnet as string;
        if (data.router !== "") {
            req.options?.push({
                tagName: "router",
                value: data.router as string,
            });
        }
        this.host.state["scopeReq"] = req;
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Subnet" name="subnet" required>
                <input type="text" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">The IP subnet the DHCP scope manages.</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Router" name="router">
                <input type="text" value="" class="pf-c-form-control" />
                <p class="pf-c-form__helper-text">The router for the specified subnet.</p>
            </ak-form-element-horizontal>`;
    }
}
