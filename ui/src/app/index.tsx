import { useContext } from "solid-js"
import { Match, Switch } from "solid-js/dom"

import { steamPumpContext, SteamPumpProvider } from "../context"
import { Games } from "../games"

import "./style.css"

// Needs to be separate because the context provider can't be
// In the same element as you use context
function Pages(): JSX.Element {
  const [_, muts] = useContext(steamPumpContext)

  return (
    <Switch>
      <Match when={muts.isPage("games")}>
        <Games />
      </Match>
    </Switch>
  )
}

export function App(): JSX.Element {
  return (
    <SteamPumpProvider>
      <Pages />
    </SteamPumpProvider>
  )
}
