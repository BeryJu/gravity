import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { CSSResult, TemplateResult, html } from "lit";
import { property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFRadio from "@patternfly/patternfly/components/Radio/radio.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../../../elements/Base";
import "../../../elements/wizard/Wizard";
import "./ZoneCacheWizardPage";
import "./ZoneForwarderWizardPage";
import "./ZoneInitialWizardPage";
import "./ZonePresetWizardPage";

@customElement("gravity-dns-zone-wizard")
export class DNSZoneWizard extends AKElement {
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
                .steps=${["gravity-dns-wizard-initial", "gravity-dns-wizard-type"]}
                header=${"New zone"}
                description=${"Create a new DNS zone."}
            >
                <button slot="trigger" class="pf-c-button pf-m-primary">${this.createText}</button>
            </ak-wizard>
        `;
    }
}
