import { CSSResult, nothing } from "lit";
import { TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly/components/FormControl/form-control.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { convertToSlug } from "../../common/utils";
import { AKElement } from "../Base";
import { FormGroup } from "../forms/FormGroup";

@customElement("ak-form-element-horizontal")
export class HorizontalFormElement extends AKElement {
    static get styles(): CSSResult[] {
        return [PFBase, PFForm, PFFormControl, AKElement.GlobalStyle];
    }

    @property()
    label = "";

    @property({ type: Boolean })
    required = false;

    @property({ type: Boolean })
    writeOnly = false;

    @property({ type: Boolean })
    writeOnlyActivated = false;

    @property({ attribute: false })
    errorMessages: string[] = [];

    @property({ type: Boolean })
    slugMode = false;

    @property()
    helperText?: string;

    @property({ type: Boolean })
    checkbox = false;

    _invalid = false;

    @property({ type: Boolean })
    set invalid(v: boolean) {
        this._invalid = v;
        // check if we're in a form group, and expand that form group
        const parent = this.parentElement?.parentElement;
        if (parent && "expanded" in parent) {
            (parent as FormGroup).expanded = true;
        }
    }
    get invalid(): boolean {
        return this._invalid;
    }

    @property()
    name = "";

    updated(): void {
        this.querySelectorAll<HTMLInputElement>("input[autofocus]").forEach((input) => {
            input.focus();
        });
        if (this.name === "slug" || this.slugMode) {
            this.querySelectorAll<HTMLInputElement>("input[type='text']").forEach((input) => {
                input.addEventListener("keyup", () => {
                    input.value = convertToSlug(input.value);
                });
            });
        }
        this.querySelectorAll("*").forEach((input) => {
            if (this.name && this.name !== "") {
                switch (input.tagName.toLowerCase()) {
                    case "input":
                    case "textarea":
                    case "select":
                    case "ak-codemirror":
                    case "ak-chip-group":
                    case "ak-search-select":
                        (input as HTMLInputElement).name = this.name;
                        break;
                    default:
                        return;
                }
            }
            if (this.writeOnly && !this.writeOnlyActivated) {
                const i = input as HTMLInputElement;
                i.setAttribute("hidden", "true");
                const handler = () => {
                    i.removeAttribute("hidden");
                    this.writeOnlyActivated = true;
                    i.parentElement?.removeEventListener("click", handler);
                };
                i.parentElement?.addEventListener("click", handler);
            }
        });
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-form pf-m-horizontal">
            <div class="pf-v6-c-form__group">
                <div class="pf-v6-c-form__group-label">
                    <label class="pf-v6-c-form__label">
                        <span class="pf-v6-c-form__label-text">${this.label}</span>
                        ${this.required
                            ? html`&nbsp;<span
                                      class="pf-v6-c-form__label-required"
                                      aria-hidden="true"
                                      >*</span
                                  >`
                            : nothing}
                    </label>
                </div>
                <div class="pf-v6-c-form__group-control">
                    ${this.checkbox
                        ? html`<slot></slot>`
                        : html`
                              <span
                                  class="pf-v6-c-form-control ${this.required
                                      ? "pf-m-required"
                                      : ""}"
                              >
                                  <slot></slot>
                              </span>
                          `}
                    <div class="pf-v6-c-form__helper-text" aria-live="polite">
                        <div class="pf-v6-c-helper-text">
                            ${this.helperText
                                ? html`
                                      <div class="pf-v6-c-form__helper-text">
                                          <div class="pf-v6-c-helper-text">
                                              <div class="pf-v6-c-helper-text__item">
                                                  <span class="pf-v6-c-helper-text__item-text"
                                                      >${this.helperText}</span
                                                  >
                                              </div>
                                          </div>
                                      </div>
                                  `
                                : nothing}
                            ${this.errorMessages.map((message) => {
                                return html`<div class="pf-v6-c-helper-text__item pf-m-error">
                                    <span class="pf-v6-c-helper-text__item-icon">
                                        <i
                                            class="fas fa-fw fa-exclamation-circle"
                                            aria-hidden="true"
                                        ></i>
                                    </span>
                                    <span class="pf-v6-c-helper-text__item-text">${message}</span>
                                </div>`;
                            })}
                        </div>
                    </div>
                </div>
            </div>
        </div>`;
    }
}
