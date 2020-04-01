import { Peer, SteamGame } from "./types"

export async function GetPeers(): Promise<Peer[]> {
  const res: Response = await fetch("http://localhost:9771/mesh")
  const peers: Peer[] = (await res.json() as any[]).map((peer) => Peer.from_json(peer))
  peers.push(new Peer("localhost"))

  return peers
}

export async function CopyGame(peer: Peer, game: SteamGame): Promise<void> {
  const req: Request = new Request(`http://localhost:9771/copy/${peer.name}/${game.appID}`, {
    method: "POST",
  })
  await fetch(req)
}
