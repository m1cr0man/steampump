import prettyBytes from "pretty-bytes"
import { Wrapped } from "solid-js/types/state"

import { GameTransfer } from "../api"
import { ProgressBar } from "../progressbar"

export function Transfer({transfer}: {transfer: Wrapped<GameTransfer> }): JSX.Element {
  const colour: string =
    (!transfer || transfer.status === "Failed") ? "darkred" :
      (transfer.status === "Successful") ? "seagreen" :
        "cadetblue"

  const text: string =
  (transfer === undefined) ? "Missing" :
    (transfer.status !== "Running") ? transfer.status :
      `${prettyBytes(transfer.bytesDone)} / ${prettyBytes(transfer.bytesTotal)} (${transfer.files} files)`

  return (
    <figure class="game-transfer">
      <div class="game-transfer-info">
        From {transfer.peerName}
      </div>
      <ProgressBar
        value={transfer.progress}
        text={text}
        colour={colour}
      />
    </figure>
  )
}
