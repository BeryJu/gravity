import { css, CSSResult, html, LitElement, TemplateResult } from "lit";
import SPTableIndexVars from "@spectrum-css/table/dist/index-vars.css";
import SPTableVars from "@spectrum-css/table/dist/vars.css";
import { customElement } from "lit/decorators.js";
import "@spectrum-web-components/icons/sp-icons-medium.js";

@customElement("ddet-table")
export class Table extends LitElement {
    static get styles(): CSSResult[] {
        return [
            SPTableIndexVars,
            SPTableVars,
            css`
                :host,
                table {
                    width: 100%;
                }
            `,
        ];
    }

    render(): TemplateResult {
        return html`
            <table class="spectrum-Table spectrum-Table--sizeM">
                <thead class="spectrum-Table-head">
                    <tr>
                        <th
                            class="spectrum-Table-headCell is-sortable is-sorted-desc"
                            aria-sort="descending"
                            tabindex="0"
                        >
                            Column Title
                            <sp-icon name="ui:Arrow100"></sp-icon>
                        </th>
                        <th class="spectrum-Table-headCell is-sortable" aria-sort="none">
                            Column Title
                            <sp-icon name="ui:Arrow100"></sp-icon>
                        </th>
                        <th class="spectrum-Table-headCell">Column Title</th>
                    </tr>
                </thead>
                <tbody class="spectrum-Table-body">
                    <tr class="spectrum-Table-row">
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Alpha</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Alpha</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Alpha</td>
                    </tr>
                    <tr class="spectrum-Table-row">
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Bravo</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Bravo</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Bravo</td>
                    </tr>
                    <tr class="spectrum-Table-row">
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Charlie</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Charlie</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Charlie</td>
                    </tr>
                    <tr class="spectrum-Table-row">
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Delta</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Delta</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Delta</td>
                    </tr>
                    <tr class="spectrum-Table-row">
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Echo</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Echo</td>
                        <td class="spectrum-Table-cell" tabindex="0">Row Item Echo</td>
                    </tr>
                </tbody>
            </table>
        `;
    }
}
