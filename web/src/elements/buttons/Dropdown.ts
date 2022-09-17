import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { EVENT_REFRESH } from "../../common/constants";
import { AKElement } from "../Base";

@customElement("ak-dropdown")
export class DropdownButton extends AKElement {
    menu: HTMLElement | null;

    constructor() {
        super();
        this.menu = this.querySelector<HTMLElement>(".pf-c-dropdown__menu");
        this.querySelectorAll("button.pf-c-dropdown__toggle").forEach((btn) => {
            btn.addEventListener("click", () => {
                if (!this.menu) return;
                this.menu.hidden = !this.menu.hidden;
            });
        });
        window.addEventListener(EVENT_REFRESH, this.clickHandler);
    }

    clickHandler = (): void => {
        if (!this.menu) return;
        this.menu.hidden = true;
    };

    disconnectedCallback(): void {
        super.disconnectedCallback();
        window.removeEventListener(EVENT_REFRESH, this.clickHandler);
    }

    render(): TemplateResult {
        return html`<slot></slot>`;
    }
}
