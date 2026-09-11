package cli

type cliValueKind uint8

const (
	cliValueNone cliValueKind = iota
	cliValuePath
	cliValueDirectory
	cliValueFormat
	cliValueEncoding
)

type cliOptionSpec struct {
	Spellings   []string
	ValueName   string
	ValueKind   cliValueKind
	Values      []string
	Description string
	Repeatable  bool
}

type cliCommandSpec struct {
	Name          string
	Summary       string
	Description   string
	Usage         string
	Options       []cliOptionSpec
	OperandValues []string
	Notes         []string
	Streams       []string
	Examples      []string
}

var globalOptionSurface = []cliOptionSpec{
	{Spellings: []string{"-q", "--quiet"}, Description: "Suppress informational and success diagnostics."},
	{Spellings: []string{"--silent"}, Description: "Suppress every non-error diagnostic."},
	{Spellings: []string{"--no-color"}, Description: "Disable diagnostic color."},
	{Spellings: []string{"-h", "--help"}, Description: "Show help."},
	{Spellings: []string{"--"}, Description: "End option processing; following paths and operands remain literal."},
}

var orderedCommandSurface = []cliCommandSpec{
	{
		Name:        "encode",
		Summary:     "Encode a subtitle source as Cue JSON.",
		Description: "Encode a SubRip or WebVTT source as Cue JSON while preserving its exact source bytes.",
		Usage:       "cueson [global options] encode [options] INPUT",
		Options: []cliOptionSpec{
			{Spellings: []string{"-o", "--output"}, ValueName: "PATH", ValueKind: cliValuePath, Description: "Write Cue JSON to PATH (default INPUT.cueson.json)."},
			{Spellings: []string{"-f", "--force"}, Description: "Replace an existing regular output file."},
			{Spellings: []string{"--format"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"auto", "srt", "vtt"}, Description: "Select auto, srt, or vtt (default auto)."},
			{Spellings: []string{"--encoding"}, ValueName: "NAME", ValueKind: cliValueEncoding, Values: canonicalEncodingValues(), Description: "Select the source text encoding."},
			{Spellings: []string{"--pretty"}, Description: "Indent Cue JSON output."},
			{Spellings: []string{"--stdout"}, Description: "Write Cue JSON to stdout."},
			{Spellings: []string{"--no-speaker-detection"}, Description: "Disable derived Name: speaker observations."},
		},
		Notes: []string{
			"Format aliases: subrip for srt; webvtt for vtt.",
			"Encoding aliases: utf8; utf-8-bom, utf8-bom, utf-8-sig; utf16le, utf-16-le; utf16be, utf-16-be; windows1252, cp1252; iso8859-1, latin1, latin-1.",
			"WebVTT accepts UTF-8 only.",
			"--output - is equivalent to --stdout. --force requires a filesystem output.",
		},
		Streams: []string{
			"stdout  Cue JSON only when --stdout or --output - is selected; otherwise empty.",
			"stderr  Ordered warnings and errors allowed by the global diagnostic filters.",
		},
		Examples: []string{
			"cueson encode --pretty captions.srt",
			"cueson encode --format vtt --stdout captions.vtt",
		},
	},
	{
		Name:        "restore",
		Summary:     "Restore exact source-envelope bytes.",
		Description: "Restore exact source-envelope bytes without invoking a native format codec.",
		Usage:       "cueson [global options] restore [options] INPUT.cueson.json",
		Options: []cliOptionSpec{
			{Spellings: []string{"-o", "--output"}, ValueName: "PATH", ValueKind: cliValuePath, Description: "Restore one asset to a literal path."},
			{Spellings: []string{"--output-dir"}, ValueName: "DIR", ValueKind: cliValueDirectory, Description: "Restore all assets beneath an existing directory."},
			{Spellings: []string{"-f", "--force"}, Description: "Replace approved existing regular files."},
			{Spellings: []string{"--strict-metadata"}, Description: "Require every captured timestamp to be restored."},
			{Spellings: []string{"--no-metadata"}, Description: "Skip timestamp restoration."},
		},
		Notes: []string{
			"Multi-asset documents require --output-dir. Metadata modes are mutually exclusive.",
		},
		Streams: []string{
			"stdout  Always empty.",
			"stderr  Metadata warnings and errors allowed by the global diagnostic filters.",
		},
		Examples: []string{
			"cueson restore --output restored.srt document.cueson.json",
			"cueson restore --output-dir restored document.cueson.json",
		},
	},
	{
		Name:        "render",
		Summary:     "Render Cue JSON from its structured model.",
		Description: "Render validated Cue JSON from its structured model in its matching native subtitle format.",
		Usage:       "cueson [global options] render [options] INPUT.cueson.json --to FORMAT",
		Options: []cliOptionSpec{
			{Spellings: []string{"--to"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"srt", "vtt"}, Description: "Select srt or vtt."},
			{Spellings: []string{"-o", "--output"}, ValueName: "PATH", ValueKind: cliValuePath, Description: "Write native output to PATH instead of stdout."},
			{Spellings: []string{"-f", "--force"}, Description: "Replace an existing regular output file."},
			{Spellings: []string{"--strict"}, Description: "Reject known non-representable model content."},
		},
		Notes: []string{
			"Format aliases: subrip for srt; webvtt for vtt.",
			"--output - selects stdout. --force requires a filesystem output.",
		},
		Streams: []string{
			"stdout  Canonical native bytes unless a filesystem output is selected.",
			"stderr  Ordered warnings and errors allowed by the global diagnostic filters.",
		},
		Examples: []string{
			"cueson render captions.srt.cueson.json --to srt",
			"cueson render captions.vtt.cueson.json --to vtt --output captions.vtt",
		},
	},
	{
		Name:        "convert",
		Summary:     "Convert SubRip and WebVTT through the common model.",
		Description: "Convert Cue JSON, SubRip, or WebVTT input to the other supported native subtitle format.",
		Usage:       "cueson [global options] convert [options] INPUT --to FORMAT",
		Options: []cliOptionSpec{
			{Spellings: []string{"--to"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"srt", "vtt"}, Description: "Select srt or vtt."},
			{Spellings: []string{"--from"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"auto", "cueson", "srt", "vtt"}, Description: "Select auto, cueson, srt, or vtt (default auto)."},
			{Spellings: []string{"--encoding"}, ValueName: "NAME", ValueKind: cliValueEncoding, Values: canonicalEncodingValues(), Description: "Select the native source text encoding."},
			{Spellings: []string{"-o", "--output"}, ValueName: "PATH", ValueKind: cliValuePath, Description: "Write native output to PATH instead of stdout."},
			{Spellings: []string{"-f", "--force"}, Description: "Replace an existing regular output file."},
			{Spellings: []string{"--strict"}, Description: "Reject conversion when any known loss exists."},
			{Spellings: []string{"--no-speaker-detection"}, Description: "Disable derived speaker observations for native input."},
		},
		Notes: []string{
			"Input aliases: json and cue-json for cueson; subrip for srt; webvtt for vtt.",
			"Encoding aliases: utf8; utf-8-bom, utf8-bom, utf-8-sig; utf16le, utf-16-le; utf16be, utf-16-be; windows1252, cp1252; iso8859-1, latin1, latin-1.",
			"WebVTT accepts UTF-8 only.",
			"--output - selects stdout. --force requires a filesystem output.",
		},
		Streams: []string{
			"stdout  Canonical target bytes unless a filesystem output is selected.",
			"stderr  Ordered source diagnostics, loss warnings, and errors allowed by the global filters.",
		},
		Examples: []string{
			"cueson convert captions.vtt --to srt",
			"cueson convert document.cueson.json --from cueson --to vtt --output captions.vtt",
		},
	},
	{
		Name:        "validate",
		Summary:     "Validate Cue JSON or a native subtitle without output.",
		Description: "Validate Cue JSON, SubRip, or WebVTT completely without creating or printing an output payload.",
		Usage:       "cueson [global options] validate [options] INPUT",
		Options: []cliOptionSpec{
			{Spellings: []string{"--format"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"auto", "cueson", "srt", "vtt"}, Description: "Select auto, cueson, srt, or vtt (default auto)."},
			{Spellings: []string{"--encoding"}, ValueName: "NAME", ValueKind: cliValueEncoding, Values: canonicalEncodingValues(), Description: "Select the native source text encoding."},
		},
		Notes: []string{
			"Format aliases: json and cue-json for cueson; subrip for srt; webvtt for vtt.",
			"Encoding aliases: utf8; utf-8-bom, utf8-bom, utf-8-sig; utf16le, utf-16-le; utf16be, utf-16-be; windows1252, cp1252; iso8859-1, latin1, latin-1.",
			"Cue JSON prohibits --encoding. WebVTT accepts UTF-8 only.",
		},
		Streams: []string{
			"stdout  Always empty.",
			"stderr  One success diagnostic, ordered warnings, and errors allowed by the global filters.",
		},
		Examples: []string{
			"cueson validate document.cueson.json",
			"cueson validate --format vtt captions.vtt",
		},
	},
	{
		Name:        "inspect",
		Summary:     "Inspect safe structural facts about an input.",
		Description: "Inspect validated Cue JSON, SubRip, or WebVTT as a privacy-bounded structural report.",
		Usage:       "cueson [global options] inspect [options] INPUT",
		Options: []cliOptionSpec{
			{Spellings: []string{"--format"}, ValueName: "FORMAT", ValueKind: cliValueFormat, Values: []string{"auto", "cueson", "srt", "vtt"}, Description: "Select auto, cueson, srt, or vtt (default auto)."},
			{Spellings: []string{"--encoding"}, ValueName: "NAME", ValueKind: cliValueEncoding, Values: canonicalEncodingValues(), Description: "Select the native source text encoding."},
			{Spellings: []string{"--json"}, Description: "Write the compact inspection report version 1 as JSON."},
		},
		Notes: []string{
			"Format aliases: json and cue-json for cueson; subrip for srt; webvtt for vtt.",
			"Encoding aliases: utf8; utf-8-bom, utf8-bom, utf-8-sig; utf16le, utf-16-le; utf16be, utf-16-be; windows1252, cp1252; iso8859-1, latin1, latin-1.",
			"Cue JSON prohibits --encoding. WebVTT accepts UTF-8 only.",
		},
		Streams: []string{
			"stdout  One human report, or one compact JSON object with --json.",
			"stderr  Ordered source warnings and errors allowed by the global filters.",
		},
		Examples: []string{
			"cueson inspect document.cueson.json",
			"cueson inspect --json captions.vtt",
		},
	},
	{
		Name:        "schema",
		Summary:     "Print or save the embedded Cue JSON schema.",
		Description: "Print or save the embedded canonical Cue JSON schema and report its version.",
		Usage:       "cueson [global options] schema [options]",
		Options: []cliOptionSpec{
			{Spellings: []string{"--version"}, Description: "Print the embedded schema version."},
			{Spellings: []string{"-o", "--output"}, ValueName: "PATH", ValueKind: cliValuePath, Description: "Write the schema to a literal path instead of stdout."},
			{Spellings: []string{"-f", "--force"}, Description: "Replace an existing regular output file."},
		},
		Notes: []string{
			"--version cannot be combined with --output or --force. --force requires --output.",
		},
		Streams: []string{
			"stdout  Schema or version bytes unless a filesystem output is selected.",
			"stderr  Errors and error-associated short usage.",
		},
		Examples: []string{
			"cueson schema",
			"cueson schema --version",
			"cueson schema --output cueson.schema.json",
		},
	},
	{
		Name:        "version",
		Summary:     "Print the Cueson executable version.",
		Description: "Print the Cueson executable version as one undecorated line.",
		Usage:       "cueson [global options] version",
		Streams: []string{
			"stdout  The executable version followed by one LF.",
			"stderr  Errors only.",
		},
		Examples: []string{"cueson version"},
	},
	{
		Name:          "completion",
		Summary:       "Generate static shell completion.",
		Description:   "Generate a deterministic static completion definition for one explicitly supported shell.",
		Usage:         "cueson [global options] completion bash|zsh|fish|powershell",
		OperandValues: []string{"bash", "zsh", "fish", "powershell"},
		Notes: []string{
			"SHELL is case-sensitive. The command prints a script and never modifies a shell profile.",
		},
		Streams: []string{
			"stdout  One static completion script for the selected shell.",
			"stderr  Errors and error-associated short usage.",
		},
		Examples: []string{
			"cueson completion bash",
			"cueson completion powershell",
		},
	},
}

func canonicalEncodingValues() []string {
	return []string{"utf-8", "utf-8-bom", "utf-16le", "utf-16be", "windows-1252", "iso-8859-1"}
}

func lookupCommandSurface(name string) (cliCommandSpec, bool) {
	for _, command := range orderedCommandSurface {
		if command.Name == name {
			return command, true
		}
	}
	return cliCommandSpec{}, false
}

func commandNames() []string {
	names := make([]string, 0, len(orderedCommandSurface))
	for _, command := range orderedCommandSurface {
		names = append(names, command.Name)
	}
	return names
}
