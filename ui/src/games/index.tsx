import { useContext } from "solid-js"
import { For } from "solid-js/dom"

import { SteamGame } from "../api"
import { steamPumpContext } from "../context"
import { Game } from "../game"

import "./style.css"

export function Games(): JSX.Element {
  const [state, {}] = useContext(steamPumpContext)

  return (
    <section id="games">
      <h2>Available Games</h2>
      <For each={state.gamesGrouped} fallback={<div>Loading data from {(state.peers() || []).length} peers</div>}>
        {(games: SteamGame[]) => <Game games={games} />}
      </For>
    </section>
  )
}
