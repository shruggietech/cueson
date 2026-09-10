export function SourceEnvelope({ format, filename, sha256, bytes, preservation }) {
  return <div className="cu-source-envelope" role="row"><div className="cu-source-envelope__format" role="cell">{format}</div><div className="cu-source-envelope__filename" role="cell">{filename}</div><div className="cu-source-envelope__sha256" role="cell">{sha256}</div><div className="cu-source-envelope__bytes" role="cell">{bytes}</div><div className="cu-source-envelope__preservation" role="cell">{preservation}</div></div>;
}
