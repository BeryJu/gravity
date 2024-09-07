import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFFormControl from "@patternfly/patternfly/components/FormControl/form-control.css";
import PFInputGroup from "@patternfly/patternfly/components/InputGroup/input-group.css";
import PFTextInputGroup from "@patternfly/patternfly/components/TextInputGroup/text-input-group.css";
import PFToolbar from "@patternfly/patternfly/components/Toolbar/toolbar.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-table-search")
export class TableSearch extends AKElement {
    @property()
    value: string | undefined;

    @property()
    onSearch: ((value: string) => void) | undefined;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFButton,
            PFToolbar,
            PFInputGroup,
            PFFormControl,
            PFTextInputGroup,
            AKElement.GlobalStyle,
            css`
                ::-webkit-search-cancel-button {
                    display: none;
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-input-group__item pf-m-fill">
            <form
                class="pf-v6-c-text-input-group"
                method="get"
                @submit=${(e: Event) => {
                    e.preventDefault();
                    if (!this.onSearch) return;
                    const el =
                        this.shadowRoot?.querySelector<HTMLInputElement>("input[type=search]");
                    if (!el) return;
                    if (el.value === "") return;
                    this.onSearch(el?.value);
                }}
            >
                <div class="pf-v6-c-text-input-group__main pf-m-icon">
                    <span class="pf-v6-c-text-input-group__text">
                        <span class="pf-v6-c-text-input-group__icon">
                            <i class="fas fa-fw fa-search"></i>
                        </span>
                        <input
                            class="pf-v6-c-text-input-group__text-input"
                            name="search"
                            type="search"
                            placeholder=${"Search..."}
                            value="${ifDefined(this.value)}"
                            @search=${(ev: Event) => {
                                if (!this.onSearch) return;
                                this.onSearch((ev.target as HTMLInputElement).value);
                            }}
                        />
                    </span>
                </div>
                <div class="pf-v6-c-text-input-group__utilities">
                    <button
                        class="pf-v6-c-button pf-m-plain"
                        type="reset"
                        @click=${() => {
                            if (!this.onSearch) return;
                            this.onSearch("");
                        }}
                    >
                        <span class="pf-v6-c-button__icon">
                            <i class="fas fa-times fa-fw" aria-hidden="true"></i>
                        </span>
                    </button>
                    <button class="pf-v6-c-button pf-m-plain" type="submit">
                        <i class="fas fa-search" aria-hidden="true"></i>
                    </button>
                </div>
            </form>
        </div>`;
    }
}
