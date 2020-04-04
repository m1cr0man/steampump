import prettyBytes from "pretty-bytes"
import { createEffect, createResource, useContext } from "solid-js"
import { For, Match, Switch } from "solid-js/dom"

import { FindHeaderImage } from "../api/game"
import { CopyGame } from "../api/mesh"
import { Peer, SteamGame } from "../api/types"
import { steamPumpContext } from "../context"

import "./style.css"

function HeaderImage({games}: {games: SteamGame[]}): JSX.Element {
  const [data, loadData] = createResource<string | undefined>(undefined)
  loadData(FindHeaderImage(games))

  return (
    <img src={data()} alt={games[0].name} />
  )
}

function getGameStatus(game: SteamGame | undefined, games: SteamGame[]): string {
  if (!game) {
    return "Not installed"
  }

  const needsUpdate: boolean = game.stateFlags !== 4
  const updating: boolean = needsUpdate && game.bytesDownloaded > 0
  const outOfDate: boolean = games.some((g) => +g.buildID > +game.buildID)
  const progress: string = `${prettyBytes(game.bytesToDownload)} / ${prettyBytes(game.bytesDownloaded)}`
  let status: string = (outOfDate) ? "Update available" : (updating) ? "Updating" : (needsUpdate) ? "Update required" : "Up To Date"
  if (needsUpdate || outOfDate) { status += `, ${progress}` }

  return status
}

export function Game({games}: {games: SteamGame[]}): JSX.Element {
  const [_, {getPeers}] = useContext(steamPumpContext)

  const ref: SteamGame = games[0]
  const localGame: SteamGame | undefined = games.filter((g) => g.peer && g.peer.name === "localhost").pop()

  return (
    <figure class="game">
      <figcaption><HeaderImage games={games} /></figcaption>
      <section>
        <h3>Game Info</h3>
        <p><b>Name: </b>{ref.name}</p>
        <p><b>AppID: </b>{ref.appID}</p>
        <p><b>Size: </b>{prettyBytes(ref.sizeOnDisk)}</p>
        <p><b>Status: </b>{getGameStatus(localGame, games)}</p>
      </section>
      <section>
        <h3>Peers</h3>
        <For each={getPeers().filter((p) => p.name !== "localhost")}>
          {(peer: Peer) => {
            const game: SteamGame | undefined = peer.getGame(ref.appID)

            if (ref.appID === 581320) { console.log(game) }

            return (<p>
              <b>{ peer.name }: </b>{getGameStatus(game, games)}
            </p>)
          }}
        </For>
      </section>
      <section>
        <h3 style="padding-bottom: 0.5em;">Actions</h3>
        <For each={getPeers().filter((p) => p.name !== "localhost")}>
          {(peer: Peer) => {
            const game: SteamGame | undefined = peer.getGame(ref.appID)

            if (game && game.stateFlags === 4) {
              return (
                <button onClick={() => CopyGame(peer, game)}>Copy from {peer.name}</button>
              )
            }
          }}
        </For>
      </section>
    </figure>
  )
}
