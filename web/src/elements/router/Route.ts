import { TemplateResult, html } from "lit";
import { until } from "lit/directives/until.js";

export const SLUG_REGEX = "[-a-zA-Z0-9_]+";
export const ID_REGEX = "\\d+";
export const UUID_REGEX = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}";

export interface RouteArgs {
    [key: string]: string;
}

export class Route {
    url: RegExp;

    private callback?: (args: RouteArgs) => Promise<TemplateResult>;

    constructor(url: RegExp, callback?: (args: RouteArgs) => Promise<TemplateResult>) {
        this.url = url;
        this.callback = callback;
    }

    redirect(to: string, raw = false): Route {
        this.callback = async () => {
            console.debug(`gravity/router: redirecting ${to}`);
            if (!raw) {
                window.location.hash = `#${to}`;
            } else {
                window.location.hash = to;
            }
            return html``;
        };
        return this;
    }

    render(args: RouteArgs): TemplateResult {
        if (this.callback) {
            return html`${until(this.callback(args))}`;
        }
        throw new Error("Route does not have callback or element");
    }

    toString(): string {
        return `<Route url=${this.url} callback=${this.callback ? "true" : "false"}>`;
    }
}
