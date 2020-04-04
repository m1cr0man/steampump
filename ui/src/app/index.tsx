import { SteamPumpProvider } from "../context"
import { Games } from "../games"

import "./style.css"

export function App(): JSX.Element {
  return (
    <SteamPumpProvider>
      <Games />
    </SteamPumpProvider>
  )
}
