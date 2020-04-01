import { useContext } from "solid-js"
import { For } from "solid-js/dom"

import { SteamGame } from "../api"
import { steamPumpContext } from "../context"
import { Game } from "../game"

import "./style.css"

export function Games(): JSX.Element {
  const [state, {addGames}] = useContext(steamPumpContext)

  return (
    <section id="games">
      <h2>Available Games</h2>
      <For each={state.games} fallback={<div>Loading data from {(state.peers() || []).length} peers</div>}>
        {(game: SteamGame) => <Game game={game} />}
      </For>
    </section>
  )
}
