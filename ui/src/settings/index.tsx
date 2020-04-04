import { useContext } from "solid-js"

import { steamPumpContext } from "../context"

export function Settings(): JSX.Element {
  const [state, {getPeers}] = useContext(steamPumpContext)

  return (
    <section id="settings">
      <h2>Settings</h2>
      <p>Coming soon tm</p>
    </section>
  s
}
