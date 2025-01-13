import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";
import { property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../../elements/Base";
import "../../../elements/wizard/Wizard";
// import "./ZoneCacheWizardPage";
// import "./ZoneForwarderWizardPage";
import "./ScopeImportWizardPage";
import "./ScopeInitialWizardPage";
import "./ScopePresetWizardPage";

@customElement("gravity-dhcp-scope-wizard")
export class DHCPScopeWizard extends AKElement {
    static get styles(): CSSResult[] {
        return [PFBase, PFButton, AKElement.GlobalStyle, PFRadio];
    }

    @property({ type: Boolean })
    open = false;

    @property()
    createText = "Create";

    @property({ type: Boolean })
    showButton = true;

    render(): TemplateResult {
        return html`
            <ak-wizard
                .open=${this.open}
                .steps=${["gravity-dhcp-wizard-initial", "gravity-dhcp-wizard-type"]}
                header=${"New scope"}
                description=${"Create a new DHCP scope."}
            >
                <button slot="trigger" class="pf-c-button pf-m-primary">${this.createText}</button>
            </ak-wizard>
        `;
    }
}
