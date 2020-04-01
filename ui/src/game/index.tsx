import prettyBytes from "pretty-bytes"
import { createEffect, createResource, useContext } from "solid-js"
import { For, Match, Switch } from "solid-js/dom"

import { FindHeaderImage } from "../api/game"
import { CopyGame } from "../api/mesh"
import { Peer, SteamGame } from "../api/types"
import { steamPumpContext } from "../context"

import "./style.css"

function HeaderImage(props: {game: SteamGame; peers: Peer[]}): JSX.Element {
  const [data, loadData] = createResource<string | undefined>(undefined)
  loadData(FindHeaderImage(props.peers, props.game.appID))

  return (
    <img src={data()} alt={props.game.name} />
  )
}

function getGameStatus(game: SteamGame | undefined, peers: Peer[]): string {
  if (!game) {
    return "Not installed"
  }

  const needsUpdate: boolean = game.stateFlags !== 4
  const updating: boolean = needsUpdate && game.bytesDownloaded > 0
  const outOfDate: boolean = peers.some((peer) => {
    const peerGame: SteamGame | undefined = peer.getGame(game.appID)

    return (peerGame  || false) && +peerGame.buildID > +game.buildID
  })
  const progress: string = `${prettyBytes(game.bytesToDownload)} / ${prettyBytes(game.bytesDownloaded)}`
  let status: string = (outOfDate) ? "Update available" : (updating) ? "Updating" : (needsUpdate) ? "Update required" : "Up To Date"
  if (needsUpdate || outOfDate) { status += `, ${progress}` }

  return status
}

export function Game(props: {game: SteamGame}): JSX.Element {
  const [state, {}] = useContext(steamPumpContext)

  const localPeer: Peer | undefined = (state.peers() || []).filter((peer) => peer.name === "localhost")
    .pop()

  let localGame: SteamGame | undefined
  if (localPeer) { localGame = localPeer.getGame(props.game.appID) }

  return (
    <figure class="game">
      <figcaption><HeaderImage game={props.game} peers={state.peers() || []} /></figcaption>
      <section>
        <h3>Game Info</h3>
        <p><b>Name: </b>{props.game.name}</p>
        <p><b>AppID: </b>{props.game.appID}</p>
        <p><b>Size: </b>{prettyBytes(props.game.sizeOnDisk)}</p>
        <p><b>Status: </b>{getGameStatus(localGame, state.peers() || [])}</p>
      </section>
      <section>
        <h3>Peers</h3>
        <For each={(state.peers() || []).filter((p) => p.name !== "localhost")}>
          {(peer: Peer) => {
            const game: SteamGame | undefined = peer.getGame(props.game.appID)

            if (props.game.appID === 581320) { console.log(game) }

            return (<p>
              <b>{ peer.name }: </b>{getGameStatus(game, state.peers() || [])}
            </p>)
          }}
        </For>
      </section>
      <section>
        <h3 style="padding-bottom: 0.5em;">Actions</h3>
        <For each={(state.peers() || []).filter((p) => p.name !== "localhost")}>
          {(peer: Peer) => {
            const game: SteamGame | undefined = peer.getGame(props.game.appID)

            if (game && game.stateFlags === 4) {
              return (
                <a onClick={() => CopyGame(peer, game)}>Copy from {peer.name}</a>
              )
            }
          }}
        </For>
      </section>
    </figure>
  )
}
