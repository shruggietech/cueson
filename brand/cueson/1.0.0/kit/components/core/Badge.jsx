export function Badge({ tone = "neutral", children }) {
  return <span className={`cu-badge cu-badge--${tone}`}>{children}</span>;
}
