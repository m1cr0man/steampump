import { render } from "solid-js/dom"

import { App } from "."

it("renders without crashing", (): void => {
  const div: HTMLElement = document.createElement("div")
  const dispose: () => void = render(App, div)
  div.textContent = ""
  dispose()
})
