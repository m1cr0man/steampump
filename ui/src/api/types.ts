import prettyBytes from "pretty-bytes"

const APP_STATE_FLAGS: string[] = [
  "Invalid",
  "Uninstalled",
  "Update Required",
  "Fully Installed",
  "Encrypted",
  "Locked",
  "Files Missing",
  "AppRunning",
  "Files Corrupt",
  "Update Running",
  "Update Paused",
  "Update Started",
  "Uninstalling",
  "Backup Running",
  "Reconfiguring",
  "Validating",
  "Adding Files",
  "Preallocating",
  "Downloading",
  "Staging",
  "Committing",
  "Update Stopping",
].reverse()

export class SteamGame {

  public static from_json(peer: Peer, {
    appid,
    name,
    stateflags,
    installdir,
    lastupdated,
    updateresult,
    sizeondisk,
    buildid,
    bytestodownload,
    bytesdownloaded,
    autoupdatebehavior,
  }: any): SteamGame {
    return new SteamGame(
      appid,
      name,
      stateflags,
      installdir,
      lastupdated,
      updateresult,
      sizeondisk,
      buildid,
      bytestodownload,
      bytesdownloaded,
      autoupdatebehavior,
      peer,
    )
  }

  public constructor(
    public appID: number,
    public name: string,
    public stateFlags: number,
    public installDir: string,
    public lastUpdated: number,
    public updateResult: number,
    public sizeOnDisk: number,
    public buildID: string,
    public bytesToDownload: number,
    public bytesDownloaded: number,
    public autoUpdateBehavior: number,
    public peer?: Peer,
  ) {
  }

  public getStates(): string[] {
    let state: string = this.stateFlags.toString(2)
    state = "0".repeat(APP_STATE_FLAGS.length - state.length - 1) + state

    return APP_STATE_FLAGS.filter((_, i) => state.charAt(i) === "1")
  }

  public async loadHeaderImage(): Promise<string | undefined> {
    if (this.peer) { return this.peer.loadHeaderImage(this.appID) }
  }
}

export class GameTransfer {

  public static from_json(
    sourcePeer: Peer,
    { status, appid, bytes_done, bytes_total, files, peer, dest }: any,
  ): GameTransfer {
    return new GameTransfer(sourcePeer, status, appid, bytes_done, bytes_total, files, peer.name, dest)
  }

  public constructor(
    public sourcePeer: Peer,
    public status: string,
    public appID: number,
    public bytesDone: number,
    public bytesTotal: number,
    public files: number,
    public peerName: string,
    public dest: string,
  ) {}

  public get progress(): number {
    return (this.bytesTotal > 0) ? this.bytesDone * 100 / this.bytesTotal : 100
  }
}

export class Peer {

  public static from_json({ name }: Peer): Peer {
    return new Peer(name)
  }
  private games: SteamGame[]
  private transfers: GameTransfer[]

  public constructor(
    public name: string,
  ) {
    this.games = []
    this.transfers = []
  }

  public getGame(appID: SteamGame["appID"]): SteamGame | undefined {
    for (const game of this.games) {
      if (game.appID === appID) {
        return game
      }
    }
  }

  public getGames(): SteamGame[] {
    return this.games
  }

  public getTransfers(): GameTransfer[] {
    return this.transfers
  }

  public async loadGames(): Promise<void> {
    try {
      const res: Response = await this.apiGet("games")

      this.games = (await res.json() as any[]).map((game) =>
        SteamGame.from_json(this, game),
      )
    } catch (err) {
      console.log(`Failed to fetch games for ${this.name}: ${err}`)
    }
  }

  public async loadHeaderImage(appID: SteamGame["appID"]): Promise<string | undefined> {
    try {
      const res: Response = await this.apiGet(`games/${appID}/images/header?encode=base64`)

      if (res.status >= 299) { return }

      return `data:image/jpeg;base64,${await res.text()}`
    } catch (err) {
      console.log(`Failed to fetch header image from ${this.name}: ${err}`)
    }
  }

  public async loadTransfers(): Promise<void> {
    try {
      const res: Response = await this.apiGet("mesh/copy")

      this.transfers = (await res.json() as any[]).map((transfer) =>
        GameTransfer.from_json(this, transfer),
      )
    } catch (err) {
      console.log(`Failed to fetch transfers from ${this.name}: ${err}`)
    }
  }

  private apiGet(url: string): Promise<Response> {
    const headers: Headers = new Headers()
    if (this.name !== "localhost") { headers.append("X-Peer", this.name) }
    const request: Request = new Request(`http://localhost:9771/${url}`, {
      headers,
    })

    return fetch(request)
  }
}

export class MultiPeerGame {
  public ref: SteamGame
  private peers: Peer[]

  public constructor(
    private values: Array<{ game: SteamGame; peer: Peer}>,
  ) {
    this.ref = values[0].game
    this.peers = values.map((m) => m.peer)
  }

  public addPeer(peer: Peer, game: SteamGame): MultiPeerGame {
    return new MultiPeerGame(
      [...this.values, {peer, game}],
    )
  }

  public getGameFrom(peer: Peer): SteamGame | undefined {
    return this.values.filter((p) => p.peer === peer).map((m) => m.game).pop()
  }

  public getPeers(): Peer[] {
    return this.peers.slice()
  }

  public getStatusFrom(targetPeer: Peer): string {
    const game: SteamGame | undefined = this.getGameFrom(targetPeer)
    if (!game) { return "Not installed" }

    const progress: string = `${prettyBytes(game.bytesToDownload)} / ${prettyBytes(game.bytesDownloaded)}`
    const needsUpdate: boolean = game.stateFlags !== 4
    const updating: boolean = needsUpdate && game.bytesDownloaded > 0
    let outOfDate: boolean = false
    for (const peerGame of this.values.map((m) => m.game)) {
      outOfDate = outOfDate || +peerGame.buildID > +game.buildID
    }

    let status: string =
      (outOfDate) ? "Update available" :
      (updating) ? "Updating" :
      (needsUpdate) ? "Update required" :
      "Up To Date"
    if (needsUpdate || outOfDate) { status += `, ${progress}` }

    return status
  }
}
