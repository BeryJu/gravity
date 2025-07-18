import { ApiAPILogMessage, RolesApiApi } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/PageHeader";
import "../../elements/Spinner";
import "../../elements/buttons/SpinnerButton";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-cluster-node-logs")
export class ClusterNodeLogsPage extends TablePage<ApiAPILogMessage> {
    expandable = true;

    pageTitle(): string {
        return "Cluster logs";
    }

    pageDescription(): string | undefined {
        return undefined;
    }

    pageIcon(): string {
        return "";
    }

    async apiEndpoint(): Promise<PaginatedResponse<ApiAPILogMessage>> {
        const logs = await new RolesApiApi(DEFAULT_CONFIG).apiGetLogMessages();
        if (!logs.messages) {
            logs.messages = [];
        }
        logs.messages.reverse();
        return PaginationWrapper(logs.messages);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Level"),
            new TableColumn("Timestamp"),
            new TableColumn("Logger"),
            new TableColumn("Message"),
        ];
    }

    renderExpanded(item: ApiAPILogMessage): TemplateResult {
        return html`<td role="cell" colspan="4">
            <div class="pf-c-table__expandable-row-content">
                <pre>${JSON.stringify(item.fields, null, 4)}</pre>
            </div>
        </td>`;
    }

    row(item: ApiAPILogMessage): TemplateResult[] {
        return [
            html`${item.level}`,
            html`${item.time?.toLocaleString()}`,
            html`<pre>${item.logger}</pre>`,
            html`<pre>${item.message}</pre>`,
        ];
    }
}
