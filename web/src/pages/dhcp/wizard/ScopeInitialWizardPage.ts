import { DhcpAPIScopesPutInput, RolesDhcpApi } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dhcp-wizard-initial")
export class ScopeInitialWizardPage extends WizardFormPage {
    sidebarLabel = () => "Scope details";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const name = data.name as string;
        const scope: DhcpAPIScopesPutInput = {
            // placeholder, this is overwritten later
            subnetCidr: "10.0.0.0/8",
            ttl: 86400,
            _default: false,
            options: [],
            hook: "",
        };
        this.host.state["scopeReq"] = scope;
        this.host.addActionBefore("Create scope", "create-scope", async (): Promise<boolean> => {
            this.host.state["scope"] = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutScopes({
                scope: name,
                dhcpAPIScopesPutInput: this.host.state["scopeReq"] as DhcpAPIScopesPutInput,
            });
            this.host.state["name"] = name;
            return true;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Name" required name="name">
            <input type="text" value="" class="pf-c-form-control" required />
            <p class="pf-c-form__helper-text">
                Name of the scope. When importing leases, this name must match the name of a scope
                in the the import file.
            </p>
        </ak-form-element-horizontal> `;
    }
}
