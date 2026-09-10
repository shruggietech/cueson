export function Card({ children, ...props }) {
  return <div className="cu-card" {...props}>{children}</div>;
}
