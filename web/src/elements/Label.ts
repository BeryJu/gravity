import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFLabel from "@patternfly/patternfly-v6/components/Label/label.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "./Base";

export enum PFColor {
    Green = "pf-m-green",
    Orange = "pf-m-orange",
    Red = "pf-m-red",
    Grey = "",
}

@customElement("ak-label")
export class Label extends AKElement {
    @property()
    color: PFColor = PFColor.Grey;

    @property()
    icon: string | undefined;

    @property({ type: Boolean })
    compact = false;

    static get styles(): CSSResult[] {
        return [PFBase, PFLabel, AKElement.GlobalStyle];
    }

    getDefaultIcon(): string {
        switch (this.color) {
            case PFColor.Green:
                return "fa-check";
            case PFColor.Orange:
                return "fa-exclamation-triangle";
            case PFColor.Red:
                return "fa-times";
            case PFColor.Grey:
                return "fa-info-circle";
            default:
                return "";
        }
    }

    render(): TemplateResult {
        return html`<span class="pf-v6-c-label ${this.color} ${this.compact ? "pf-m-compact" : ""}">
            <span class="pf-v6-c-label__content">
                <span class="pf-v6-c-label__icon">
                    <i
                        class="fas fa-fw ${this.icon || this.getDefaultIcon()}"
                        aria-hidden="true"
                    ></i>
                </span>
                <slot></slot>
            </span>
        </span>`;
    }
}
