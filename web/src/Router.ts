import { CSSResult, LitElement, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { until } from "lit/directives/until.js";

export class Route {
    name: string;
    handler: () => Promise<TemplateResult> = async () => html``;

    constructor(name: string, handler: () => Promise<TemplateResult>) {
        this.name = name;
        this.handler = handler;
    }

    activeCallback(): void {}

    render(): TemplateResult {
        return html`${until(this.handler(), html`Loading...`)}`;
    }
}

@customElement("gravity-router")
export class Router extends LitElement {
    @property({ attribute: false })
    routes: Route[] = [];

    @property({ attribute: false })
    activeRoute?: Route;

    static get styles(): CSSResult {
        return css`
            .wrapper {
                height: 100vh;
                margin: 3rem 13rem;
            }
        `;
    }

    constructor() {
        super();
        window.addEventListener("hashchange", () => {
            this.navigate();
        });
    }

    firstUpdated(): void {
        this.navigate();
    }

    navigate(): void {
        const name = window.location.hash.substring(1, Infinity);
        const route = this.routes.filter((route) => route.name == name)[0];
        this.activeRoute = route;
        if (!this.activeRoute) {
            window.location.hash = `#${this.routes[0].name}`;
            return;
        }
        this.requestUpdate();
    }

    render(): TemplateResult {
        return html`<div class="wrapper">${this.activeRoute?.render()}</div>`;
    }
}
