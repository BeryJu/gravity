import Convert from "ansi-to-html";
import { ApiAPILogMessage, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement, state } from "lit/decorators.js";
import { unsafeHTML } from "lit/directives/unsafe-html.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/PageHeader";
import "../../elements/Spinner";
import "../../elements/buttons/SpinnerButton";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-cluster-node-logs")
export class ClusterNodeLogsPage extends TablePage<ApiAPILogMessage> {
    converter = new Convert();

    @state()
    isStructured = false;

    pageTitle(): string {
        return "Node logs";
    }
    pageDescription(): string | undefined {
        return undefined;
    }
    pageIcon(): string {
        return "";
    }
    async apiEndpoint(): Promise<PaginatedResponse<ApiAPILogMessage>> {
        const logs = await new RolesApiApi(DEFAULT_CONFIG).apiGetLogMessages();
        this.isStructured = logs.isJSON || true;
        if (!logs.messages) {
            logs.messages = [];
        }
        logs.messages.reverse();
        return PaginationWrapper(logs.messages);
    }
    columns(): TableColumn[] {
        if (this.isStructured) {
            return [
                new TableColumn("Level"),
                new TableColumn("Timestamp"),
                new TableColumn("Logger"),
                new TableColumn("Message"),
                new TableColumn("Fields"),
            ];
        }
        return [new TableColumn("Message")];
    }
    row(item: ApiAPILogMessage): TemplateResult[] {
        if (this.isStructured) {
            try {
                const payload: LogMessage = JSON.parse(item.message || "");
                const otherFields: { [key: string]: unknown } = { ...payload };
                delete otherFields.level;
                delete otherFields.logger;
                delete otherFields.msg;
                delete otherFields.ts;
                return [
                    html`${payload.level}`,
                    html`${new Date(payload.ts * 1000).toLocaleTimeString()}`,
                    html`${payload.logger}`,
                    html`${payload.msg}`,
                    html`<pre>${JSON.stringify(otherFields)}</pre>`,
                ];
            } catch (error) {
                console.log(error);
                this.isStructured = false;
            }
        }
        return [html`${unsafeHTML(this.converter.toHtml(item.message || ""))}`];
    }
}

interface LogMessage {
    level: string;
    ts: number;
    logger: string;
    msg: string;
    [key: string]: unknown;
}
