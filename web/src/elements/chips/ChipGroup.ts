import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFLabelGroup from "@patternfly/patternfly/components/Label/label-group.css";
import PFLabel from "@patternfly/patternfly/components/Label/label.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../Base";
import { Chip } from "../chips/Chip";

@customElement("ak-chip-group")
export class ChipGroup extends AKElement {
    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFLabel,
            PFLabelGroup,
            PFButton,
            AKElement.GlobalStyle,
            css`
                ::slotted(*) {
                    margin: 0 2px;
                }
                .pf-c-chip-group {
                    margin-bottom: 8px;
                }
            `,
        ];
    }

    set value(v: (string | number | undefined)[]) {
        return;
    }

    get value(): (string | number | undefined)[] {
        const values: (string | number | undefined)[] = [];
        this.querySelectorAll<Chip>("ak-chip").forEach((chip) => {
            values.push(chip.value);
        });
        return values;
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-label-group">
            <div class="pf-v6-c-label-group__main">
                <ul class="pf-v6-c-label-group__list" role="list">
                    <slot></slot>
                </ul>
            </div>
        </div>`;
    }
}
