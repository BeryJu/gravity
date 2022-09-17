import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFAlertGroup from "@patternfly/patternfly/components/AlertGroup/alert-group.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { EVENT_MESSAGE, EVENT_WS_MESSAGE, WS_MSG_TYPE_MESSAGE } from "../../common/constants";
import { SentryIgnoredError } from "../../common/errors";
import { WSMessage } from "../../common/ws";
import { AKElement } from "../Base";
import "../messages/Message";
import { APIMessage } from "../messages/Message";

export function showMessage(message: APIMessage, unique = false): void {
    const container = document.querySelector<MessageContainer>("ak-message-container");
    if (!container) {
        throw new SentryIgnoredError("failed to find message container");
    }
    container.addMessage(message, unique);
    container.requestUpdate();
}

@customElement("ak-message-container")
export class MessageContainer extends AKElement {
    @property({ attribute: false })
    messages: APIMessage[] = [];

    static get styles(): CSSResult[] {
        return [
            PFBase,
            PFAlertGroup,
            css`
                /* Fix spacing between messages */
                ak-message {
                    display: block;
                }
            `,
        ];
    }

    constructor() {
        super();
        this.addEventListener(EVENT_WS_MESSAGE, ((e: CustomEvent<WSMessage>) => {
            if (e.detail.message_type !== WS_MSG_TYPE_MESSAGE) return;
            this.addMessage(e.detail as unknown as APIMessage);
        }) as EventListener);
        this.addEventListener(EVENT_MESSAGE, ((e: CustomEvent<APIMessage>) => {
            this.addMessage(e.detail);
        }) as EventListener);
    }

    addMessage(message: APIMessage, unique = false): void {
        if (unique) {
            const matchingMessages = this.messages.filter((m) => m.message == message.message);
            if (matchingMessages.length > 0) {
                return;
            }
        }
        this.messages.push(message);
    }

    render(): TemplateResult {
        return html`<ul class="pf-c-alert-group pf-m-toast">
            ${this.messages.map((m) => {
                return html`<ak-message
                    .message=${m}
                    .onRemove=${(m: APIMessage) => {
                        this.messages = this.messages.filter((v) => v !== m);
                        this.requestUpdate();
                    }}
                >
                </ak-message>`;
            })}
        </ul>`;
    }
}
