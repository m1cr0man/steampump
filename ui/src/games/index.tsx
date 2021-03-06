import { useContext } from "solid-js"
import { For } from "solid-js/dom"

import { SteamGame } from "../api"
import { steamPumpContext } from "../context"
import { Game } from "../game"

import "./style.css"

export function Games(): JSX.Element {
  const [state, {getPeers}] = useContext(steamPumpContext)

  return (
    <section id="games">
      <h2>Available Games</h2>
      <aside>{state.games.length} games on {getPeers().length} peers</aside>
      <For each={state.gamesGrouped} fallback={<div>Loading data from {getPeers().length} peers</div>}>
        {(games: SteamGame[]) => <Game games={games} />}
      </For>
    </section>
  )
}
