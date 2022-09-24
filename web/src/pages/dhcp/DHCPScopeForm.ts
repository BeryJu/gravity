import { DhcpScope, RolesDhcpApi } from "gravity-api";
import YAML from "yaml";

import { TemplateResult, html } from "lit";
import { customElement } from "lit/decorators.js";
import { ifDefined } from "lit/directives/if-defined.js";

import { DEFAULT_CONFIG } from "../../api/Config";
import "../../elements/CodeMirror";
import "../../elements/forms/HorizontalFormElement";
import { ModelForm } from "../../elements/forms/ModelForm";

@customElement("gravity-dhcp-scope-form")
export class DHCPScopeForm extends ModelForm<DhcpScope, string> {
    loadInstance(pk: string): Promise<DhcpScope> {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpGetScopes().then((scopes) => {
            const zone = scopes.scopes?.find((z) => z.scope === pk);
            if (!zone) throw new Error("No scope");
            return zone;
        });
    }

    getSuccessMessage(): string {
        if (this.instance) {
            return `Successfully updated scope.`;
        } else {
            return `Successfully created scope.`;
        }
    }

    send = (data: DhcpScope): Promise<void> => {
        return new RolesDhcpApi(DEFAULT_CONFIG).dhcpPutScopes({
            scope: data.scope,
            dhcpScopesInput: data,
        });
    };

    renderForm(): TemplateResult {
        return html`<form class="pf-c-form pf-m-horizontal">
            ${this.instance
                ? html``
                : html`<ak-form-element-horizontal label="Name" ?required=${true} name="scope">
                      <input type="text" class="pf-c-form-control" required />
                  </ak-form-element-horizontal>`}
            <ak-form-element-horizontal name="_default">
                <div class="pf-c-check">
                    <input
                        type="checkbox"
                        class="pf-c-check__input"
                        ?checked=${this.instance?._default}
                    />
                    <label class="pf-c-check__label"> ${`Default`} </label>
                </div>
            </ak-form-element-horizontal>
            <ak-form-element-horizontal label="TTL" ?required=${true} name="ttl">
                <input
                    type="number"
                    value="${this.instance?.ttl || 86400}"
                    class="pf-c-form-control"
                    required
                />
            </ak-form-element-horizontal>
        </form>`;
    }
}
