import { LitElement, html, css } from "lit";
import { logout, request } from "./services/api.js";

class App extends LitElement {
    static get styles() {
        return css`
            :host {
                display: block;
            }
        `;
    }
    render() {
        return html`
        `;
    }
}
customElements.define("ddet-app", App);
