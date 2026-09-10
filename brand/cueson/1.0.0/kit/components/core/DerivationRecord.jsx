export function DerivationRecord({ method, input, output, confidence, provenance }) {
  return <div className="cu-derivation-record" role="row"><div className="cu-derivation-record__method" role="cell">{method}</div><div className="cu-derivation-record__input" role="cell">{input}</div><div className="cu-derivation-record__output" role="cell">{output}</div><div className="cu-derivation-record__confidence" role="cell">{confidence}</div><div className="cu-derivation-record__provenance" role="cell">{provenance}</div></div>;
}
