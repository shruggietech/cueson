export function CueRow({ id, start, end, speaker, text, source }) {
  return <div className="cu-cue-row" role="row"><div className="cu-cue-row__id" role="cell">{id}</div><div className="cu-cue-row__start" role="cell">{start}</div><div className="cu-cue-row__end" role="cell">{end}</div><div className="cu-cue-row__speaker" role="cell">{speaker}</div><div className="cu-cue-row__text" role="cell">{text}</div><div className="cu-cue-row__source" role="cell">{source}</div></div>;
}
