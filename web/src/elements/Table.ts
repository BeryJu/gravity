import SPTableIndexVars from "@spectrum-css/table/dist/index-vars.css";
import SPTableVars from "@spectrum-css/table/dist/vars.css";
import "@spectrum-web-components/icons/sp-icons-medium.js";

import { CSSResult, LitElement, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

@customElement("gravity-table")
export class Table<T> extends LitElement {
    @property({ attribute: false })
    columns: string[] = [];

    @property({ attribute: false })
    rowRender: (item: T) => TemplateResult[] = (item: T): TemplateResult[] => {
        return [html`${item}`];
    };

    @property({ attribute: false })
    rowLink: (item: T) => string | undefined = (item: T): string | undefined => {
        return undefined;
    };

    @property({ attribute: false })
    data: () => Promise<T[]> = (): Promise<T[]> => {
        return Promise.resolve([] as T[]);
    };

    static get styles(): CSSResult[] {
        return [
            SPTableIndexVars,
            SPTableVars,
            css`
                :host,
                table {
                    width: 100%;
                }
                .spectrum-Table-headCell {
                    text-align: left;
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`
            <table class="spectrum-Table spectrum-Table--sizeM">
                <thead class="spectrum-Table-head">
                    <tr>
                        ${this.columns.map((col) => {
                            return html`<th class="spectrum-Table-headCell">${col}</th>`;
                        })}
                    </tr>
                </thead>
                <tbody class="spectrum-Table-body">
                    ${until(
                        this.data().then((data) => {
                            return data.map((item) => {
                                return html`<tr
                                    class="spectrum-Table-row"
                                    @click=${() => {
                                        const link = this.rowLink(item);
                                        if (link) {
                                            window.location.assign(link);
                                        }
                                    }}
                                >
                                    ${this.rowRender(item).map((col) => {
                                        return html`<td class="spectrum-Table-cell">${col}</td>`;
                                    })}
                                </tr>`;
                            });
                        }),
                        html`Loading...`,
                    )}
                </tbody>
            </table>
        `;
    }
}
