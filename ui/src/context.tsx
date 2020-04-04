import { createContext, createEffect, createResource, createSignal, createState, onCleanup } from "solid-js"
import { Context } from "solid-js/types/signal"
import { Wrapped } from "solid-js/types/state"

import { GetPeers, MultiPeerGame, Peer, SteamGame } from "./api"

const HOME_PAGE: string = "games"

function sortStrings(a: string, b: string): number {
  return (a < b) ? -1 : (a > b) ? 1 : 0
}

export interface IStoreState {
    games: SteamGame[]
    gamesGrouped: SteamGame[][]
    page(): string
    peers(): Peer[] | undefined
}

export interface IStoreMutators {
  getPeers(): Peer[]
  isPage(page: string): boolean
}

export type IContext = [Wrapped<IStoreState>, IStoreMutators]

export const steamPumpContext: Context<IContext> = createContext<IContext>([
  { peers: () => [], page: () => "", games: [], gamesGrouped: [] },
  {
    getPeers(): Peer[] {
      return []
    },
    isPage(page: string): boolean {
      return false
    },
  },
])

export function SteamPumpProvider(props: {children: any}): JSX.Element {
    const [peers, loadPeers] = createResource<Peer[]>([])
    const [page, setPage] = createSignal(HOME_PAGE)
    const [state, setState] = createState<IStoreState>({ peers, page, games: [], gamesGrouped: [] })
    loadPeers((async () =>
      (await GetPeers()).sort((a, b) =>
        (a.name === "localhost") ? -1 : sortStrings(a.name, b.name),
      )
    )())

    const mutators: IStoreMutators = {
      getPeers(): Peer[] {
        return state.peers() || []
      },
      isPage(qpage: string): boolean {
        console.log(`${page()} ${qpage}`)

        return page() === qpage
      },
    }

    createEffect(async () => {
      console.log("Loading games")

      await Promise.all((mutators.getPeers()).map((peer) => (async () => {
        await peer.loadGames()
        if (peer.getGames().length === 0) { return }
        setState("games", [...state.games, ...peer.getGames()].sort((a, b) => sortStrings(a.name, b.name)))
      })()))

      console.log("Done")
    })

    createEffect(() => {
      console.log("Grouping games")

      // What's really nice is maps preserve insertion order,
      // So the grouped games will be sorted since the ungrouped games already are
      setState("gamesGrouped", Array.from(
        state.games.reduce(
          (mapper, game) => {
            if (mapper.has(game.appID)) {
              mapper.get(game.appID)!.push(game)
            } else {
              mapper.set(game.appID, [game])
            }

            return mapper
          },
          new Map<SteamGame["appID"], SteamGame[]>()).values(),
      ))
    })

    const onPageChange: (e: Event) => void = (_) => setPage(window.location.hash.slice(1))
    // Add page change listener
    window.addEventListener("hashchange", onPageChange)
    onCleanup(() => window.removeEventListener("hashchange", onPageChange))

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
