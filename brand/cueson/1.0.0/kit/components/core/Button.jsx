export function Button({ variant = "primary", size = "md", children, ...props }) {
  return <button className={`cu-button cu-button--${variant} cu-button--${size}`} type="button" {...props}>{children}</button>;
}
