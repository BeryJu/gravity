import { customElement } from "@lit/reactive-element/decorators/custom-element.js";
import { TemplateResult, html } from "lit";

import { KeyUnknown } from "../../../elements/forms/Form";
import "../../../elements/forms/FormGroup";
import "../../../elements/forms/HorizontalFormElement";
import { WizardFormPage } from "../../../elements/wizard/WizardFormPage";

@customElement("gravity-dns-wizard-cache")
export class ZoneCacheWizardPage extends WizardFormPage {
    sidebarLabel = () => "Cache configuration";

    nextDataCallback = async (data: KeyUnknown): Promise<boolean> => {
        if (!(data["enableCache"] as boolean)) {
            return true;
        }
        const config = this.host.state["handlerConfigs"] as KeyUnknown[];
        const forwarderConfig = config.filter((config) =>
            (config.type as string).startsWith("forward_"),
        );
        forwarderConfig.forEach((conf) => {
            conf["cache_ttl"] = (data["cacheTTL"] as number).toString();
        });
        forwarderConfig.splice(0, 1, {
            type: "memory",
        });
        forwarderConfig.splice(1, 1, {
            type: "etcd",
        });
        return true;
    };

    renderForm(): TemplateResult {
        return html`<ak-form-element-horizontal name="enableCache" checkbox>
                <div class="pf-v6-c-check">
                    <input type="checkbox" class="pf-v6-c-check__input" ?checked=${true} />
                    <label class="pf-v6-c-check__label"> ${"Enable cache"} </label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal
                label=${"Cache TTL"}
                required
                name="cacheTTL"
                helperText="Duration in seconds the records will be cached for."
            >
                <input type="number" required value=${3600} />
            </ak-form-element-horizontal>`;
    }
}
