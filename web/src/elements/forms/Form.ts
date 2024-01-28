import { ResponseError, RestErrResponse } from "gravity-api";

import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFAlert from "@patternfly/patternfly/components/Alert/alert.css";
import PFButton from "@patternfly/patternfly/components/Button/button.css";
import PFCard from "@patternfly/patternfly/components/Card/card.css";
import PFForm from "@patternfly/patternfly/components/Form/form.css";
import PFFormControl from "@patternfly/patternfly/components/FormControl/form-control.css";
import PFInputGroup from "@patternfly/patternfly/components/InputGroup/input-group.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { EVENT_REFRESH } from "../../common/constants";
import { MessageLevel } from "../../common/messages";
import { convertToSlug } from "../../common/utils";
import { AKElement } from "../Base";
import { HorizontalFormElement } from "../forms/HorizontalFormElement";
import { showMessage } from "../messages/MessageContainer";

export interface KeyUnknown {
    [key: string]: unknown;
}

/**
 * Recursively assign `value` into `json` while interpreting the dot-path of `element.name`
 */
function assignValue(element: HTMLInputElement, value: unknown, json: KeyUnknown): void {
    let parent = json;
    if (!element.name?.includes(".")) {
        parent[element.name] = value;
        return;
    }
    const nameElements = element.name.split(".");
    for (let index = 0; index < nameElements.length - 1; index++) {
        const nameEl = nameElements[index];
        // Ensure all nested structures exist
        if (!(nameEl in parent)) parent[nameEl] = {};
        parent = parent[nameEl] as { [key: string]: unknown };
    }
    parent[nameElements[nameElements.length - 1]] = value;
}

/**
 * Convert the elements of the form to JSON.[4]
 *
 */
export function serializeForm<T extends KeyUnknown>(
    elements: NodeListOf<HorizontalFormElement>,
): T | undefined {
    const json: { [key: string]: unknown } = {};
    elements.forEach((element) => {
        element.requestUpdate();
        const inputElement = element.querySelector<HTMLInputElement>("[name]");
        if (element.hidden || !inputElement) {
            return;
        }
        // Skip elements that are writeOnly where the user hasn't clicked on the value
        if (element.writeOnly && !element.writeOnlyActivated) {
            return;
        }
        if (
            inputElement.tagName.toLowerCase() === "select" &&
            "multiple" in inputElement.attributes
        ) {
            const selectElement = inputElement as unknown as HTMLSelectElement;
            assignValue(
                inputElement,
                Array.from(selectElement.selectedOptions).map((v) => v.value),
                json,
            );
        } else if (inputElement.tagName.toLowerCase() === "input" && inputElement.type === "date") {
            assignValue(inputElement, inputElement.valueAsDate, json);
        } else if (
            inputElement.tagName.toLowerCase() === "input" &&
            inputElement.type === "datetime-local"
        ) {
            assignValue(inputElement, new Date(inputElement.valueAsNumber), json);
        } else if (
            inputElement.tagName.toLowerCase() === "input" &&
            "type" in inputElement.dataset &&
            inputElement.dataset["type"] === "datetime-local"
        ) {
            // Workaround for Firefox <93, since 92 and older don't support
            // datetime-local fields
            assignValue(inputElement, new Date(inputElement.value), json);
        } else if (
            inputElement.tagName.toLowerCase() === "input" &&
            inputElement.type === "checkbox"
        ) {
            assignValue(inputElement, inputElement.checked, json);
        } else if (
            inputElement.tagName.toLowerCase() === "input" &&
            inputElement.type === "number"
        ) {
            assignValue(inputElement, parseInt(inputElement.value, 10), json);
        } else if ("selectedFlow" in inputElement) {
            assignValue(inputElement, inputElement.value, json);
        } else {
            assignValue(inputElement, inputElement.value, json);
        }
    });
    return json as unknown as T;
}

@customElement("ak-form")
export class Form<T> extends AKElement {
    viewportCheck = true;

    @property()
    accessor successMessage = "";

    @property()
    accessor send!: (data: T) => Promise<unknown>;

    @property({ attribute: false })
    accessor nonFieldErrors: string[] | undefined;

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFCard,
            PFButton,
            PFForm,
            PFAlert,
            PFInputGroup,
            PFFormControl,
            AKElement.GlobalStyle,
            css`
                select[multiple] {
                    height: 15em;
                }
            `,
        ];
    }

    get isInViewport(): boolean {
        const rect = this.getBoundingClientRect();
        return !(rect.x + rect.y + rect.width + rect.height === 0);
    }

    getSuccessMessage(): string {
        return this.successMessage;
    }

    updated(): void {
        this.shadowRoot
            ?.querySelectorAll("ak-form-element-horizontal[name=name]")
            .forEach((nameInput) => {
                const input = nameInput.firstElementChild as HTMLInputElement;
                const form = nameInput.closest("form");
                if (form === null) {
                    return;
                }
                const slugFieldWrapper = form.querySelector(
                    "ak-form-element-horizontal[name=slug]",
                );
                if (!slugFieldWrapper) {
                    return;
                }
                const slugField = slugFieldWrapper.firstElementChild as HTMLInputElement;
                // Only attach handler if the slug is already equal to the name
                // if not, they are probably completely different and shouldn't update
                // each other
                if (convertToSlug(input.value) !== slugField.value) {
                    return;
                }
                nameInput.addEventListener("input", () => {
                    slugField.value = convertToSlug(input.value);
                });
            });
    }

    resetForm(): void {
        const form = this.shadowRoot?.querySelector<HTMLFormElement>("form");
        form?.reset();
    }

    serializeForm(): T | undefined {
        const elements = this.shadowRoot?.querySelectorAll<HorizontalFormElement>(
            "ak-form-element-horizontal",
        );
        if (!elements) {
            return {} as T;
        }
        return serializeForm(elements) as T;
    }

    submit(ev: Event): Promise<unknown> | undefined {
        ev.preventDefault();
        const data = this.serializeForm();
        if (!data) {
            return;
        }
        return this.send(data)
            .then((r) => {
                const message = this.getSuccessMessage();
                if (message) {
                    showMessage({
                        level: MessageLevel.success,
                        message: this.getSuccessMessage(),
                    });
                }
                this.dispatchEvent(
                    new CustomEvent(EVENT_REFRESH, {
                        bubbles: true,
                        composed: true,
                    }),
                );
                return r;
            })
            .catch(async (ex: Error | ResponseError) => {
                console.warn(ex);
                if (!(ex instanceof ResponseError)) {
                    throw ex;
                }
                let msg = ex.response.statusText;
                if (ex.response.status > 399 && ex.response.status < 500) {
                    const errorMessage: RestErrResponse = await ex.response.json();
                    if (!errorMessage) return errorMessage;
                    if (errorMessage instanceof Error) {
                        throw errorMessage;
                    }
                    if (errorMessage.error) {
                        this.nonFieldErrors = [errorMessage.error];
                        msg = errorMessage.error;
                    }
                }
                // error is local or not from rest_framework
                showMessage({
                    message: msg,
                    level: MessageLevel.error,
                });
                // rethrow the error so the form doesn't close
                throw ex;
            });
    }

    renderNonFieldErrors(): TemplateResult {
        if (!this.nonFieldErrors) {
            return html``;
        }
        return html`<div class="pf-c-form__alert">
            ${this.nonFieldErrors.map((err) => {
                return html`<div class="pf-c-alert pf-m-inline pf-m-danger">
                    <div class="pf-c-alert__icon">
                        <i class="fas fa-exclamation-circle"></i>
                    </div>
                    <h4 class="pf-c-alert__title">${err}</h4>
                </div>`;
            })}
        </div>`;
    }

    renderFormWrapper(): TemplateResult {
        const inline = this.renderForm();
        if (inline) {
            return html`<form
                class="pf-c-form pf-m-horizontal"
                @submit=${(ev: Event) => {
                    ev.preventDefault();
                }}
            >
                ${inline}
            </form>`;
        }
        return html`<slot></slot>`;
    }

    renderForm(): TemplateResult | undefined {
        return undefined;
    }

    renderVisible(): TemplateResult {
        return html` ${this.renderNonFieldErrors()} ${this.renderFormWrapper()}`;
    }

    render(): TemplateResult {
        if (this.viewportCheck && !this.isInViewport) {
            return html``;
        }
        return this.renderVisible();
    }
}
