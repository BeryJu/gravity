import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFList from "@patternfly/patternfly/components/List/list.css";

import { EVENT_REFRESH } from "../../common/constants";
import { MessageLevel } from "../../common/messages";
import { ModalButton } from "../buttons/ModalButton";
import "../buttons/SpinnerButton";
import { showMessage } from "../messages/MessageContainer";

@customElement("ak-forms-delete")
export class DeleteForm extends ModalButton {
    static get styles(): CSSResult[] {
        return super.styles.concat(PFList);
    }

    @property({ attribute: false })
    obj?: Record<string, unknown>;

    @property()
    objectLabel?: string;

    @property({ attribute: false })
    delete!: () => Promise<unknown>;

    confirm(): Promise<void> {
        return this.delete()
            .then(() => {
                this.onSuccess();
                this.open = false;
                this.dispatchEvent(
                    new CustomEvent(EVENT_REFRESH, {
                        bubbles: true,
                        composed: true,
                    }),
                );
            })
            .catch((e) => {
                this.onError(e);
                throw e;
            });
    }

    onSuccess(): void {
        showMessage({
            message: `Successfully deleted ${this.objectLabel} ${this.obj?.name}`,
            level: MessageLevel.success,
        });
    }

    onError(e: Error): void {
        showMessage({
            message: `Failed to delete ${this.objectLabel}: ${e.toString()}`,
            level: MessageLevel.error,
        });
    }

    renderModalInner(): TemplateResult {
        let objName = this.obj?.name;
        if (objName) {
            objName = ` "${objName}"`;
        } else {
            objName = "";
        }
        return html`<section class="pf-c-modal-box__header pf-c-page__main-section pf-m-light">
                <div class="pf-c-content">
                    <h1 class="pf-c-title pf-m-2xl">${`Delete ${this.objectLabel}`}</h1>
                </div>
            </section>
            <section class="pf-c-modal-box__body pf-c-page__main-section pf-m-light">
                <form class="pf-c-form pf-m-horizontal">
                    <p>${`Are you sure you want to delete ${this.objectLabel} ${objName} ?`}</p>
                </form>
            </section>
            <footer class="pf-c-modal-box__footer">
                <ak-spinner-button
                    .callAction=${() => {
                        return this.confirm();
                    }}
                    class="pf-m-danger"
                >
                    ${"Delete"} </ak-spinner-button
                >&nbsp;
                <ak-spinner-button
                    .callAction=${async () => {
                        this.open = false;
                    }}
                    class="pf-m-secondary"
                >
                    ${"Cancel"}
                </ak-spinner-button>
            </footer>`;
    }
}
