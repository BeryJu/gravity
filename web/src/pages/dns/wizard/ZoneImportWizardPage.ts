import { DnsAPIZonesImporterType, RolesDnsApi } from "gravity-api";

import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { DEFAULT_CONFIG } from "../../../api/Config";
import { convertToTitle } from "../../../common/utils";
import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dns-wizard-import")
export class ZoneImportWizardPage extends WizardFormPage {
    sidebarLabel = () => "Import configuration";

    async FileToString(file: File): Promise<string | null> {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsText(file);
            reader.onload = () => {
                resolve(reader.result as string);
            };
            reader.onerror = reject;
        });
    }

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        const file = this.form?.getFormFiles()["file"];
        if (!file) {
            return false;
        }
        const contents = await this.FileToString(file);
        this.host.addActionAfter("Importing records", "import", async () => {
            const name = this.host.state["name"] as string;
            const result = await new RolesDnsApi(DEFAULT_CONFIG).dnsImportZones({
                zone: name,
                dnsAPIZonesImportInput: {
                    payload: contents || "",
                    type: data.Type as DnsAPIZonesImporterType,
                },
            });
            return result.successful || false;
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal label="Type" required name="type">
                ${Object.keys(DnsAPIZonesImporterType).map((type) => {
                    return html`<div class="pf-c-radio">
                        <input
                            class="pf-c-radio__input"
                            type="radio"
                            name="type"
                            id=${type}
                            value=${type}
                        />
                        <label class="pf-c-radio__label" for=${type}>${convertToTitle(type)}</label>
                    </div>`;
                })}
                <p class="pf-c-form__helper-text">Format of the data to import</p>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="File" name="file" required>
                <input type="file" value="" class="pf-c-form-control" required />
                <p class="pf-c-form__helper-text">File to import</p>
            </ak-form-element-horizontal>`;
    }
}
