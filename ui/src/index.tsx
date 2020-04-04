import { render } from "solid-js/dom"

import { App } from "./app"
import * as serviceWorker from "./serviceWorker"

// SetTimeout(() => render(App, document.getElementById('root') as Node), 3000)
render(App, document.getElementById("root") as Node)

// If you want your app to work offline and load faster, you can change
// Unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister()
