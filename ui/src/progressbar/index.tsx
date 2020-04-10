import "./style.css"

export function ProgressBar(
  {value, text, colour}: {colour?: string; text?: string; value: number },
): JSX.Element {
  const style: {[index: string]: string} = {"width": `${value}%`, "background-color": colour || "olivegreen"}

  return (
    <div class="progress-bar">
      <div class="progress-bar-remaining" style={style}></div>
      <div class="progress-bar-text">
        {text || value}
      </div>
    </div>
  )
}
