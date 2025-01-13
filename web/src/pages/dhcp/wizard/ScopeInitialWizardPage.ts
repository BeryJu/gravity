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
            subnetCidr: data.subnet as string,
            ttl: 86400,
            _default: false,
            options: [],
            hook: "",
        };
        if (data.router !== "") {
            scope.options?.push({
                tagName: "router",
                value: data.router as string,
            });
        }
        this.host.addActionBefore("Create scope", "create-scope", async (): Promise<boolean> => {
            this.host.state["scope"] = await new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutScopes({
                scope: name,
                dhcpAPIScopesPutInput: scope,
            });
            this.host.state["name"] = name;
            return true;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Name" required name="name">
                <input type="text" value="" class="pf-c-form-control" required />
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Subnet" name="subnet" required>
                <input type="text" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">The IP subnet the DHCP scope manages.</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Router" name="router">
                <input type="text" value="" class="pf-c-form-control" />
                <p class="pf-c-form__helper-text">The router for the specified subnet.</p>
            </ak-form-element-horizontal> `;
    }
}
