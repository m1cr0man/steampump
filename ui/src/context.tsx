import { createContext, createEffect, createResource, createState, reconcile } from "solid-js"
import { Context } from "solid-js/types/signal"
import { Wrapped } from "solid-js/types/state"

import { GetPeers, MultiPeerGame, Peer, SteamGame } from "./api"

function sortStrings(a: string, b: string): number {
  return (a < b) ? -1 : (a > b) ? 1 : 0
}

export interface IStoreState {
    games: SteamGame[]
    peers(): Peer[] | undefined
}

export interface IStoreMutators {
  getGamesGrouped(): SteamGame[][]
}

export type IContext = [Wrapped<IStoreState>, IStoreMutators]

export const steamPumpContext: Context<IContext> = createContext<IContext>([
  { peers: () => [], games: [] },
  {
    getGamesGrouped(): SteamGame[][] {
      return []
    },
  },
])

export function SteamPumpProvider(props: {children: any}): JSX.Element {
    const [peers, loadPeers] = createResource<Peer[]>([])
    const [state, setState] = createState<IStoreState>({ peers, games: [] })
    loadPeers((async () =>
      (await GetPeers()).sort((a, b) =>
        (a.name === "localhost") ? -1 : sortStrings(a.name, b.name),
      )
    )())

    const mutators: IStoreMutators = {
      getGamesGrouped(): SteamGame[][] {
        return [] // TODO
      },
    }

    createEffect(async () => {
      console.log("Loading games")

      await Promise.all((state.peers() || []).map((peer) => (async () => {
        await peer.loadGames()
        setState("games", [...state.games, ...peer.getGames()])
      })()))

      console.log("Done")
    })

    const store: [Wrapped<IStoreState>, IStoreMutators] = [
        state,
        mutators,
    ]

    return (
      <steamPumpContext.Provider value={store}>
        {props.children}
      </steamPumpContext.Provider>
    )
}
