import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { property } from "@lit/reactive-element/decorators/property.js";
import { CSSResult, TemplateResult, html } from "lit";
import { state } from "lit/decorators.js";

import PFActionList from "@patternfly/patternfly/components/ActionList/action-list.css";
import PFWizard from "@patternfly/patternfly/components/Wizard/wizard.css";

import { ModalButton } from "../buttons/ModalButton";
import "./ActionWizardPage";
import { WizardPage } from "./WizardPage";

export interface WizardAction {
    displayName: string;
    subText?: string;
    uid: string;
    run: () => Promise<boolean>;
}

export const ApplyActionsSlot = "apply-actions";

@customElement("ak-wizard")
export class Wizard extends ModalButton {
    @property({ type: Boolean })
    canCancel = true;

    @property({ type: Boolean })
    canBack = true;

    @property()
    header: string | undefined;

    @property()
    description: string | undefined;

    @property({ type: Boolean })
    isValid = false;

    static get styles(): CSSResult[] {
        return super.styles.concat(PFWizard, PFActionList);
    }

    @state()
    _steps: string[] = [];

    get steps(): string[] {
        return this._steps;
    }

    set steps(steps: string[]) {
        const addApplyActionsSlot = this.steps.includes(ApplyActionsSlot);
        this._steps = steps;
        if (addApplyActionsSlot) {
            this.steps.push(ApplyActionsSlot);
        }
        this.steps.forEach((step) => {
            const exists = this.querySelector(`[slot=${step}]`) !== null;
            if (!exists) {
                const el = document.createElement(step);
                el.slot = step;
                el.dataset["wizardmanaged"] = "true";
                this.appendChild(el);
            }
        });
        this.requestUpdate();
    }

    _initialSteps: string[] = [];

    @property({ attribute: false })
    actions: WizardAction[] = [];

    @state()
    _currentStep: WizardPage | undefined;

    set currentStep(value: WizardPage | undefined) {
        this._currentStep = value;
        if (this._currentStep) {
            this._currentStep.activeCallback();
            this._currentStep.requestUpdate();
        }
    }

    get currentStep(): WizardPage | undefined {
        return this._currentStep;
    }

    @property({ attribute: false })
    finalHandler: () => Promise<void> = () => {
        return Promise.resolve();
    };

    @property({ attribute: false })
    state: { [key: string]: unknown } = {};

    firstUpdated(): void {
        this._initialSteps = this._steps;
    }

    /**
     * Add action to the beginning of the list
     */
    addActionBefore(displayName: string, uid: string, run: () => Promise<boolean>): void {
        this.actions = this.actions.filter((action) => action.uid !== uid);
        this.actions.unshift({
            displayName,
            run,
            uid,
        });
    }

    /**
     * Add action at the end of the list
     */
    addActionAfter(displayName: string, uid: string, run: () => Promise<boolean>): void {
        this.actions = this.actions.filter((action) => action.uid !== uid);
        this.actions.push({
            displayName,
            run,
            uid,
        });
    }

    renderClose() {
        return html``;
    }

    renderModalInner(): TemplateResult {
        const firstPage = this.querySelector<WizardPage>(`[slot=${this.steps[0]}]`);
        if (!this.currentStep && firstPage) {
            this.currentStep = firstPage;
        }
        const currentIndex = this.currentStep ? this.steps.indexOf(this.currentStep.slot) : 0;
        let lastPage = currentIndex === this.steps.length - 1;
        if (lastPage && !this.steps.includes("ak-wizard-page-action") && this.actions.length > 0) {
            this.steps = this.steps.concat("ak-wizard-page-action");
            lastPage = currentIndex === this.steps.length - 1;
        }
        return html`<div class="pf-v6-c-wizard">
            <div class="pf-v6-c-wizard__header">
                ${this.canCancel
                    ? html`<div class="pf-v6-c-wizard__close">
                          <button
                              class="pf-v6-c-button pf-m-plain"
                              type="button"
                              aria-label="${"Close"}"
                              @click=${() => {
                                  this.reset();
                              }}
                          >
                              <span class="pf-v6-c-button__icon">
                                  <i class="fas fa-times" aria-hidden="true"></i>
                              </span>
                          </button>
                      </div> `
                    : html``}
                <div class="pf-v6-c-wizard__title">
                    <h1 class="pf-v6-c-wizard__title-text">${this.header}</h1>
                </div>
                <div class="pf-v6-c-wizard__description">${this.description}</div>
            </div>
            <div class="pf-v6-c-wizard__outer-wrap">
                <div class="pf-v6-c-wizard__inner-wrap">
                    <nav class="pf-v6-c-wizard__nav">
                        <ol class="pf-v6-c-wizard__nav-list">
                            ${this.steps.map((step, idx) => {
                                const currentIdx = this.currentStep
                                    ? this.steps.indexOf(this.currentStep.slot)
                                    : 0;
                                return html`
                                    <li class="pf-v6-c-wizard__nav-item">
                                        <button
                                            class="pf-v6-c-wizard__nav-link ${idx === currentIdx
                                                ? "pf-m-current"
                                                : ""}"
                                            type="button"
                                            ?disabled=${currentIdx < idx}
                                            @click=${() => {
                                                const stepEl = this.querySelector<WizardPage>(
                                                    `[slot=${step}]`,
                                                );
                                                if (stepEl) {
                                                    this.currentStep = stepEl;
                                                }
                                            }}
                                        >
                                            <span class="pf-v6-c-wizard__nav-link-main">
                                                <span class="pf-v6-c-wizard__nav-link-text">
                                                    ${this.querySelector<WizardPage>(
                                                        `[slot=${step}]`,
                                                    )?.sidebarLabel()}
                                                </span>
                                            </span>
                                        </button>
                                    </li>
                                `;
                            })}
                        </ol>
                    </nav>
                    <main class="pf-v6-c-wizard__main">
                        <div class="pf-v6-c-wizard__main-body">
                            <slot name=${this.currentStep?.slot || this.steps[0]}></slot>
                        </div>
                    </main>
                </div>
                <footer class="pf-v6-c-wizard__footer">
                    <div class="pf-v6-c-action-list">
                        <div class="pf-v6-c-action-list__group">
                            <button
                                class="pf-v6-c-button pf-m-primary"
                                type="submit"
                                ?disabled=${!this.isValid}
                                @click=${async () => {
                                    const cb = await this.currentStep?.nextCallback();
                                    if (!cb) {
                                        return;
                                    }
                                    if (lastPage) {
                                        await this.finalHandler();
                                        this.reset();
                                    } else {
                                        const nextPage = this.querySelector<WizardPage>(
                                            `[slot=${this.steps[currentIndex + 1]}]`,
                                        );
                                        if (nextPage) {
                                            this.currentStep = nextPage;
                                        }
                                    }
                                }}
                            >
                                ${lastPage ? "Finish" : "Next"}
                            </button>
                            ${(this.currentStep ? this.steps.indexOf(this.currentStep.slot) : 0) >
                                0 && this.canBack
                                ? html`
                                      <button
                                          class="pf-v6-c-button pf-m-secondary"
                                          type="button"
                                          @click=${() => {
                                              const prevPage = this.querySelector<WizardPage>(
                                                  `[slot=${this.steps[currentIndex - 1]}]`,
                                              );
                                              if (prevPage) {
                                                  this.currentStep = prevPage;
                                              }
                                          }}
                                      >
                                          ${"Back"}
                                      </button>
                                  `
                                : html``}
                        </div>
                        <div class="pf-v6-c-action-list__group">
                            ${this.canCancel
                                ? html`<div class="pf-v6-c-wizard__footer-cancel">
                                      <button
                                          class="pf-v6-c-button pf-m-link"
                                          type="button"
                                          @click=${() => {
                                              this.reset();
                                          }}
                                      >
                                          ${"Cancel"}
                                      </button>
                                  </div>`
                                : html``}
                        </div>
                    </div>
                </footer>
            </div>
        </div>`;
    }

    reset(): void {
        this.open = false;
        this.querySelectorAll("[data-wizardmanaged=true]").forEach((el) => {
            el.remove();
        });
        this.steps = this._initialSteps;
        this.actions = [];
        this.state = {};
        this.currentStep = undefined;
        this.canBack = true;
        this.canCancel = true;
    }
}
