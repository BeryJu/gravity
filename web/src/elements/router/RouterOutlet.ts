import "src/elements/router/Router404";

import { CSSResult, LitElement, TemplateResult, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

import { Route } from "./Route";
import { RouteMatch } from "./RouteMatch";

// Poliyfill for hashchange.newURL,
// https://developer.mozilla.org/en-US/docs/Web/API/WindowEventHandlers/onhashchange
window.addEventListener("load", () => {
    if (!window.HashChangeEvent)
        (function () {
            let lastURL = document.URL;
            window.addEventListener("hashchange", function (event) {
                Object.defineProperty(event, "oldURL", {
                    enumerable: true,
                    configurable: true,
                    value: lastURL,
                });
                Object.defineProperty(event, "newURL", {
                    enumerable: true,
                    configurable: true,
                    value: document.URL,
                });
                lastURL = document.URL;
            });
        })();
});

@customElement("gravity-router-outlet")
export class RouterOutlet extends LitElement {
    @property({ attribute: false })
    current?: RouteMatch;

    @property()
    defaultUrl?: string;

    @property({ attribute: false })
    routes: Route[] = [];

    constructor() {
        super();
        window.addEventListener("hashchange", (ev: HashChangeEvent) => this.navigate(ev));
    }

    static get styles(): CSSResult {
        return css`
            :host {
                height: 100vh;
                margin: 3rem 13rem;
            }
        `;
    }

    firstUpdated(): void {
        this.navigate();
    }

    navigate(ev?: HashChangeEvent): void {
        let activeUrl = window.location.hash.slice(1, Infinity);
        if (ev) {
            // Check if we've actually changed paths
            const oldPath = new URL(ev.oldURL).hash.slice(1, Infinity);
            if (oldPath === activeUrl) return;
        }
        if (activeUrl === "") {
            activeUrl = this.defaultUrl || "/";
            window.location.hash = `#${activeUrl}`;
            console.debug(`gravity/router: defaulted URL to ${window.location.hash}`);
            return;
        }
        let matchedRoute: RouteMatch | null = null;
        this.routes.some((route) => {
            const match = route.url.exec(activeUrl);
            if (match != null) {
                matchedRoute = new RouteMatch(route);
                matchedRoute.arguments = match.groups || {};
                matchedRoute.fullUrl = activeUrl;
                console.debug("gravity/router: found match ", matchedRoute);
                return true;
            }
        });
        if (!matchedRoute) {
            console.debug(`gravity/router: route "${activeUrl}" not defined`);
            const route = new Route(RegExp(""), async () => {
                return html`<div class="pf-c-page__main">
                    <gravity-router-404 url=${activeUrl}></gravity-router-404>
                </div>`;
            });
            matchedRoute = new RouteMatch(route);
            matchedRoute.arguments = route.url.exec(activeUrl)?.groups || {};
            matchedRoute.fullUrl = activeUrl;
        }
        this.current = matchedRoute;
        this.requestUpdate();
    }

    render(): TemplateResult | undefined {
        return this.current?.render();
    }
}
