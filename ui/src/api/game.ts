import { Peer, SteamGame } from "./types"

export async function FindHeaderImage(games: SteamGame[]): Promise<string | undefined> {
  for (const game of games) {
    if (game.peer) {
      const url: string | undefined = await game.loadHeaderImage()
      if (url) { return url }
    }
  }
}
