import { html, LitElement, TemplateResult } from "lit";
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

@customElement("ddet-router")
export class Router extends LitElement {
    @property({ attribute: false })
    routes: Route[] = [];

    @property({ attribute: false })
    activeRoute?: Route;

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
        console.debug(route);
        this.activeRoute = route;
        this.requestUpdate();
    }

    render(): TemplateResult {
        return html` ${this.activeRoute?.render()} `;
    }
}
