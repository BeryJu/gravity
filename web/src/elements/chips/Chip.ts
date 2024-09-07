import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFLabelGroup from "@patternfly/patternfly-v6/components/Label/label-group.css";
import PFLabel from "@patternfly/patternfly-v6/components/Label/label.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-chip")
export class Chip extends AKElement {
    @property()
    value: number | string | undefined;

    @property({ type: Boolean })
    removable = false;

    static get styles(): CSSResult[] {
        return [PFBase, PFButton, PFLabel, PFLabelGroup, AKElement.GlobalStyle];
    }

    render(): TemplateResult {
        return html`
            <li class="pf-v6-c-label-group__list-item">
                <span class="pf-v6-c-label">
                    <span class="pf-v6-c-label__content">
                        <span class="pf-v6-c-label__icon">
                            <i class="fas fa-fw fa-info-circle" aria-hidden="true"></i>
                        </span>
                        <span class="pf-v6-c-label__text"><slot></slot></span>
                    </span>
                </span>
            </li>
        `;
    }
}
