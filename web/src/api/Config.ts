import {
    Configuration,
    Middleware,
    ResponseContext,
} from "gravity-api";

export class LoggingMiddleware implements Middleware {
    post(context: ResponseContext): Promise<Response | void> {
        let msg = `gravity/api: `;
        msg += `${context.response.status} ${context.init.method} ${context.url}`;
        console.debug(msg);
        return Promise.resolve(context.response);
    }
}

export const DEFAULT_CONFIG = new Configuration({
    basePath: "/api/v1",
    middleware: [
        new LoggingMiddleware(),
    ],
});
