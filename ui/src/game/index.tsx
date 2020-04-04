import prettyBytes from "pretty-bytes"
import { createResource, useContext } from "solid-js"
import { For } from "solid-js/dom"

import { FindHeaderImage } from "../api/game"
import { CopyGame } from "../api/mesh"
import { Peer, SteamGame } from "../api/types"
import { steamPumpContext } from "../context"

import "./style.css"

interface IGameStatus {
  colclass: string
  human: string
  installed: boolean
  needsUpdate: boolean
  outOfDate: boolean
  progress: string
  updating: boolean
}

function HeaderImage({games}: {games: SteamGame[]}): JSX.Element {
  const [data, loadData] = createResource<string | undefined>(undefined)
  loadData(FindHeaderImage(games))

  return (
    <img src={data()} alt={games[0].name} />
  )
}

function getGameStatus(game: SteamGame | undefined, games: SteamGame[]): IGameStatus {
  if (!game) {
    return {
      colclass: "status-grey",
      human: "Not installed",
      installed: false,
      needsUpdate: false,
      outOfDate: false,
      progress: "0 B / 0 B",
      updating: false,
    }
  }

  const needsUpdate: boolean = game.stateFlags !== 4
  const updating: boolean = needsUpdate && game.bytesDownloaded > 0
  const outOfDate: boolean = games.some((g) => +g.buildID > +game.buildID)
  const progress: string = `${prettyBytes(game.bytesToDownload)} / ${prettyBytes(game.bytesDownloaded)}`
  let human: string = (outOfDate) ? "Update available" : (updating) ? "Updating" : (needsUpdate) ? "Update required" : "Up To Date"
  if (updating) { human += `, ${progress}` }

  const colclass: string =
    (outOfDate) ? "status-orange" :
      (needsUpdate) ? "status-red" :
        "status-green"

  return {
    colclass,
    human,
    installed: true,
    needsUpdate,
    outOfDate,
    progress,
    updating,
  }
}

function PeerPill({peer, games}: {games: SteamGame[]; peer: Peer }): JSX.Element {
  const game: SteamGame | undefined = peer.getGame(games[0].appID)
  const status: IGameStatus = getGameStatus(game, games)

  return (
    <figure
      class={`peerpill ${status.colclass} ${game && !status.needsUpdate && "clickable"}`}
      onClick={() => game && !status.needsUpdate && CopyGame(peer, game)}
    >
      <figcaption>{peer.name}</figcaption>
      <p>{status.human}</p>
    </figure>
  )
}

export function Game({games}: {games: SteamGame[]}): JSX.Element {
  const [_, {getPeers}] = useContext(steamPumpContext)

  const ref: SteamGame = games[0]
  const localGame: SteamGame | undefined = games.filter((g) => g.peer && g.peer.name === "localhost").pop()
  const status: IGameStatus = getGameStatus(localGame, games)

  return (
    <figure class="game">
      <figcaption><HeaderImage games={games} /></figcaption>
      <section class="game-info">
        <h3>Game Info</h3>
        <dl>
          <dt>Name:</dt>
          <dd>{ref.name}</dd><br />
          <dt>AppID:</dt>
          <dd>{ref.appID}</dd><br />
          <dt>Size:</dt>
          <dd>{prettyBytes(ref.sizeOnDisk)}</dd><br />
          <dt>Status:</dt>
          <dd>{status.human}</dd>
        </dl>
      </section>
      <section class="game-peers">
        <h3>Peers</h3>
        <For each={getPeers().filter((p) => p.name !== "localhost")}>
          {(peer: Peer) => (
            <PeerPill peer={peer} games={games}/>
          )}
        </For>
      </section>
    </figure>
  )
}
