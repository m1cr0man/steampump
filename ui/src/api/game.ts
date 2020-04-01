import { Peer, SteamGame } from "./types"

export async function FindHeaderImage(peers: Peer[], appID: SteamGame["appID"]): Promise<string | undefined> {
  const targetPeer: Peer | undefined = peers.filter((peer) => !!peer.getGame(appID))
    .pop()

  if (!targetPeer) { return }

  return targetPeer.loadHeaderImage(appID)
}
