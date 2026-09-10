export function Textarea({ id, label, required = false, error, ...props }) {
  return (
    <div className="cu-field">
      <label className="cu-field__label" htmlFor={id}>{label}{required ? <span className="cu-field__required" aria-hidden="true"> *</span> : null}</label>
      <textarea className="cu-field__control" id={id} required={required} {...props} />
      {error ? <p className="cu-field__error" role="alert">{error}</p> : null}
    </div>
  );
}
