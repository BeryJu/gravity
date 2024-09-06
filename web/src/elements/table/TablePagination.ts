import { CSSResult, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import PFButton from "@patternfly/patternfly-v6/components/Button/button.css";
import PFPagination from "@patternfly/patternfly-v6/components/Pagination/pagination.css";
import PFBase from "@patternfly/patternfly-v6/patternfly-base.css";

import { AKElement } from "../Base";

export interface Pagination {
    next?: number;
    previous?: number;

    count: number;
    current: number;
    totalPages: number;

    startIndex: number;
    endIndex: number;
}

@customElement("ak-table-pagination")
export class TablePagination extends AKElement {
    @property({ attribute: false })
    pages: Pagination | undefined;

    @property({ attribute: false })
    // eslint-disable-next-line
    pageChangeHandler: (page: number) => void = (page: number) => {};

    static get styles(): CSSResult[] {
        return [PFBase, PFButton, PFPagination, AKElement.GlobalStyle];
    }

    render(): TemplateResult {
        return html`<div class="pf-v6-c-pagination pf-m-compact pf-m-hidden pf-m-visible-on-md">
            <div
                class="pf-v6-c-pagination pf-m-compact pf-m-compact pf-m-hidden pf-m-visible-on-md"
            >
                <div class="pf-v6-c-options-menu">
                    <div class="pf-v6-c-options-menu__toggle pf-m-text pf-m-plain">
                        <span class="pf-v6-c-options-menu__toggle-text">
                            ${`${this.pages?.startIndex} - ${this.pages?.endIndex} of ${this.pages?.count}`}
                        </span>
                    </div>
                </div>
                <nav class="pf-v6-c-pagination__nav" aria-label="Pagination">
                    <div class="pf-v6-c-pagination__nav-control pf-m-prev">
                        <button
                            class="pf-v6-c-button pf-m-plain"
                            @click=${() => {
                                this.pageChangeHandler(this.pages?.previous || 0);
                            }}
                            ?disabled="${(this.pages?.previous || 0) < 1}"
                            aria-label="${"Go to previous page"}"
                        >
                            <i class="fas fa-angle-left" aria-hidden="true"></i>
                        </button>
                    </div>
                    <div class="pf-v6-c-pagination__nav-control pf-m-next">
                        <button
                            class="pf-v6-c-button pf-m-plain"
                            @click=${() => {
                                this.pageChangeHandler(this.pages?.next || 0);
                            }}
                            ?disabled="${(this.pages?.next || 0) <= 0}"
                            aria-label="${"Go to next page"}"
                        >
                            <i class="fas fa-angle-right" aria-hidden="true"></i>
                        </button>
                    </div>
                </nav>
            </div>
        </div>`;
    }
}
