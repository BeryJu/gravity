import { Route } from "src/elements/router/Route";

import { TemplateResult } from "lit";

export class RouteMatch {
    route: Route;
    arguments: { [key: string]: string };
    fullUrl?: string;

    constructor(route: Route) {
        this.route = route;
        this.arguments = {};
    }

    render(): TemplateResult {
        return this.route.render(this.arguments);
    }

    toString(): string {
        return `<RouteMatch url=${this.fullUrl} route=${this.route} arguments=${JSON.stringify(
            this.arguments,
        )}>`;
    }
}
