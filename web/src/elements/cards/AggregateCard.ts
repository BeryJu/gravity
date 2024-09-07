import { CSSResult, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import PFCard from "@patternfly/patternfly/components/Card/card.css";
import PFFlex from "@patternfly/patternfly/layouts/Flex/flex.css";
import PFBase from "@patternfly/patternfly/patternfly-base.css";

import { AKElement } from "../Base";

@customElement("ak-aggregate-card")
export class AggregateCard extends AKElement {
    @property()
    icon: string | undefined;

    @property()
    header: string | undefined;

    @property()
    headerLink: string | undefined;

    @property({ type: Boolean })
    isCenter = true;

    static get styles(): CSSResult[] {
        return [PFBase, PFCard, PFFlex, AKElement.GlobalStyle].concat([
            css`
                .pf-v6-c-card.pf-v6-c-card-aggregate {
                    height: 100%;
                }
                .pf-c-card__header {
                    flex-wrap: nowrap;
                }
                .center-value {
                    font-size: var(--pf-t--global--font--size--heading--h1);
                    text-align: center;
                    color: var(--pf-global--Color--100);
                }
                .subtext {
                    font-size: var(--pf-t--global--font--size--heading--h4);
                }
            `,
        ]);
    }

    renderInner(): TemplateResult {
        return html`<slot></slot>`;
    }

    renderHeaderLink(): TemplateResult {
        return html`${this.headerLink
            ? html`<a class="pf-v6-c-menu-toggle pf-m-plain" href="${this.headerLink}">
                  <span class="pf-v6-c-menu-toggle__icon">
                      <i class="fa fa-link"> </i>
                  </span>
              </a>`
            : ""}`;
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-card pf-v6-c-card-aggregate">
            <div class="pf-v6-c-card__header">
                <div class="pf-v6-c-card__actions">${this.renderHeaderLink()}</div>
                <div class="pf-v6-c-card__header-main">
                    <div class="pf-v6-c-card__title">
                        <h2 class="pf-v6-c-card__title-text">
                            <i class="${ifDefined(this.icon)}"></i>&nbsp;${this.header
                                ? this.header
                                : ""}
                        </h2>
                    </div>
                </div>
            </div>
            <div class="pf-v6-c-card__body ${this.isCenter ? "center-value" : ""}">
                ${this.renderInner()}
            </div>
        </div>`;
    }
}
