import { RolesTftpApi, TftpAPIFile } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/buttons/ActionButton";
import "../../elements/forms/DeleteBulkForm";
import "../../elements/forms/ModalForm";
import { PaginatedResponse, TableColumn } from "../../elements/table/Table";
import { TablePage } from "../../elements/table/TablePage";
import { PaginationWrapper } from "../../utils";
import "./TFTPFileForm";

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

    b64toBlob(b64Data: string, contentType = "", sliceSize = 512) {
        const byteCharacters = atob(b64Data);
        const byteArrays = [];

        for (let offset = 0; offset < byteCharacters.length; offset += sliceSize) {
            const slice = byteCharacters.slice(offset, offset + sliceSize);

            const byteNumbers = new Array(slice.length);
            for (let i = 0; i < slice.length; i++) {
                byteNumbers[i] = slice.charCodeAt(i);
            }

            const byteArray = new Uint8Array(byteNumbers);
            byteArrays.push(byteArray);
        }

        const blob = new Blob(byteArrays, { type: contentType });
        return blob;
    }

    download(content: string, fileName: string) {
        const blob = this.b64toBlob(content);
        const urlForDownload = window.URL.createObjectURL(blob);
        const linkElement = document.createElement("a");
        linkElement.href = urlForDownload;
        linkElement.download = fileName;
        linkElement.click();
        URL.revokeObjectURL(urlForDownload);
    }

    row(item: TftpAPIFile): TemplateResult[] {
        return [
            html`<pre>${item.name}</pre>`,
            html`<pre>${item.host}</pre>`,
            html`${item.sizeBytes} Bytes`,
            html`<ak-action-button
                class="pf-m-plain"
                .apiRequest=${async () => {
                    const data = await new RolesTftpApi(DEFAULT_CONFIG).tftpDownloadFiles({
                        host: item.host,
                        name: item.name,
                    });
                    this.download(data.data, item.name);
                }}
            >
                <i class="fas fa-download" aria-hidden="true"></i>
            </ak-action-button>`,
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
            <button ?disabled=${disabled} slot="trigger" class="pf-v6-c-button pf-m-danger">
                ${"Delete"}
            </button>
        </ak-forms-delete-bulk>`;
    }

    renderObjectCreate(): TemplateResult {
        return html`
            <ak-forms-modal>
                <span slot="submit"> ${"Upload"} </span>
                <span slot="header"> ${"Upload File"} </span>
                <gravity-tftp-file-form slot="form"> </gravity-tftp-file-form>
                <button slot="trigger" class="pf-v6-c-button pf-m-primary">${"Upload"}</button>
            </ak-forms-modal>
        `;
    }
}
