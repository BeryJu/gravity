import { LitElement, html, css } from "lit";

import "@spectrum-web-components/split-view/sp-split-view.js";
import "@spectrum-web-components/theme/sp-theme.js";
import "@spectrum-web-components/sidenav/sp-sidenav.js";
import "@spectrum-web-components/sidenav/sp-sidenav-heading.js";
import "@spectrum-web-components/sidenav/sp-sidenav-item.js";

class App extends LitElement {

    static get styles() {
        return css`
            :host {
                display: block;
            }
        `;
    }

    render() {
        return html`
            <sp-theme color="light" scale="medium">
                <sp-split-view resizable collapsible primary-size="auto" class="fill">
                    <sp-sidenav defaultValue="Docs">
                        <sp-sidenav-item value="Docs" href="/components/SideNav">
                            Docs
                        </sp-sidenav-item>
                        <sp-sidenav-item value="Guides" href="/guides/getting_started">
                            Guides
                        </sp-sidenav-item>
                        <sp-sidenav-item value="Community" href="/community">
                            Community
                        </sp-sidenav-item>
                        <sp-sidenav-item value="Storybook" href="/storybook" target="_blank">
                            Storybook
                        </sp-sidenav-item>
                        <sp-sidenav-item
                            value="Releases"
                            href="http://git.corp.adobe.com/React/react-spectrum/releases"
                            target="_blank"
                            disabled
                        >
                            Releases
                        </sp-sidenav-item>
                        <sp-sidenav-item
                            value="GitHub"
                            href="http://git.corp.adobe.com/React/react-spectrum"
                            target="_blank"
                        >
                            GitHub
                        </sp-sidenav-item>
                    </sp-sidenav>
                </sp-split-view>
            </sp-theme>
        `;
    }
}
customElements.define("ddet-app", App);
