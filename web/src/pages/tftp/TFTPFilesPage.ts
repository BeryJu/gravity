import { RolesTftpApi, TftpAPIFile } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";

@customElement("gravity-tftp-files")
export class TFTPFilesPage extends TablePage<TftpAPIFile> {
    pageTitle(): string {
        return "TFTP Files";
    }
    pageDescription(): string | undefined {
        return "";
    }
    pageIcon(): string {
        return "";
    }
    checkbox = true;

    searchEnabled(): boolean {
        return true;
    }

    async apiEndpoint(): Promise<PaginatedResponse<TftpAPIFile>> {
        const files = await new RolesTftpApi(DEFAULT_CONFIG).tftpGetFiles();
        const data = files.files || [];
        data.sort((a, b) => {
            if (a.name > b.name) return 1;
            if (a.name < b.name) return -1;
            return 0;
        });
        return PaginationWrapper(data);
    }

    columns(): TableColumn[] {
        return [
            new TableColumn("Filename"),
            new TableColumn("Host"),
            new TableColumn("Size"),
            new TableColumn("Actions"),
        ];
    }

    row(item: TftpAPIFile): TemplateResult[] {
        return [
            html`${item.name}`,
            html`<pre>${item.host}</pre>`,
            html`${item.sizeBytes} Bytes`,
            html``,
        ];
    }

    renderToolbarSelected(): TemplateResult {
        const disabled = this.selectedElements.length < 1;
        return html`<ak-forms-delete-bulk
            objectLabel=${"TFTP File(s)"}
            .objects=${this.selectedElements}
            .metadata=${(item: TftpAPIFile) => {
                return [
                    { key: "Name", value: item.name },
                    { key: "Host", value: item.host },
                ];
            }}
            .delete=${(item: TftpAPIFile) => {
                return new RolesTftpApi(DEFAULT_CONFIG).tftpDeleteFiles({
                    name: item.name,
                    host: item.host,
                });
            }}
        >
            <button ?disabled=${disabled} slot="trigger" class="pf-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }
}
