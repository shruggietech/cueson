export function SectionHeading({ eyebrow, title, description }) {
  return <header className="cu-section-heading">{eyebrow ? <div className="cu-eyebrow">{eyebrow}</div> : null}<h2 className="cu-section-heading__title">{title}</h2>{description ? <p className="cu-section-heading__description">{description}</p> : null}</header>;
}
