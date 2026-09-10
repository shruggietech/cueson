export function Select({ id, label, required = false, error, children, ...props }) {
  return (
    <div className="cu-field">
      <label className="cu-field__label" htmlFor={id}>{label}{required ? <span className="cu-field__required" aria-hidden="true"> *</span> : null}</label>
      <select className="cu-field__control" id={id} required={required} {...props}>{children}</select>
      {error ? <p className="cu-field__error" role="alert">{error}</p> : null}
    </div>
  );
}
