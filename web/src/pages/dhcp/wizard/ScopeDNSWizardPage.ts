import { DhcpAPIScopesPutInput } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dhcp-wizard-dns")
export class ScopeDNSWizardPage extends WizardFormPage {
    sidebarLabel = () => "DNS configuration";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const req = this.host.state["scopeReq"] as DhcpAPIScopesPutInput;
        req.dns = {
            zone: data.zone as string,
        };
        this.host.state["scopeReq"] = req;
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="DNS Zone" name="zone" required>
                <input type="text" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">DNS Zone which records in this scope should be registered in. Reverse records are created automatically if a matching reverse-DNS zone exists.</p>
            </ak-form-element-horizontal>`;
    }
}
