import { RolesTftpApi, TftpAPIFilesPutInput } from "gravity-api";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-tftp-file-form")
export class TFTPFileForm extends ModelForm<TftpAPIFilesPutInput, string> {
    loadInstance(): Promise<TftpAPIFilesPutInput> {
        throw new Error("Method not implemented.");
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return "Successfully updated file.";
        } else {
            return "Successfully created file.";
        }
    }

    async toBase64(file: File): Promise<string | null> {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = () => {
                let data = reader.result as string;
                data = data.split(";base64,")[1];
                resolve(data);
            };
            reader.onerror = reject;
        });
    }

    send = async (data: TftpAPIFilesPutInput): Promise<void> => {
        const file = this.getFormFiles()["file"];
        const b64 = await this.toBase64(file);
        if (!b64) {
            throw new Error("Empty file");
        }
        data.data = b64;
        return new RolesTftpApi(DEFAULT_CONFIG).tftpPutFiles({
            tftpAPIFilesPutInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Name" required name="name">
                <input
                    type="text"
                    value=${this.instance?.name || ""}
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Filename</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="Host" required name="host">
                <input
                    type="text"
                    value=${this.instance?.host || ""}
                    class="pf-c-form-control"
                    required
                />
                <p class="pf-c-form__helper-text">Host namespace</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label=${"File"} name="file">
                <input type="file" value="" class="pf-c-form-control" />
            </ak-form-element-horizontal>`;
    }
}
