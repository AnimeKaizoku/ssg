package strongStringGo

// the prefex values for commands.
const (
	COMMAND_PREFIX1 = "!"
	COMMAND_PREFIX2 = "/"
	SUDO_PREFIX1    = ">"
	FLAG_PREFIX     = "--"
)

const (
	JA_FLAG       = "〰〰"
	JA_STR        = "❞" // start character (") for string in japanese.
	JA_EQUALITY   = "＝" // equal character (＝) for string in japanese.
	JA_DDOT       = "：" // equal character (＝) for string in japanese.
	BACK_STR      = "\\\""
	BACK_FLAG     = "\\--"
	BACK_EQUALITY = "\\="
	BACK_DDOT     = "\\:"
)

// the base constant values.
const (
	BaseIndex    = 0 // number 0
	BaseOneIndex = 1 // number 1
)

// additional constants which are not actually used in
// this package, but may be useful in another packages.
const (
	BaseIndexStr    = "0"  // number 0
	BaseOneIndexStr = "1"  // number 1
	DotStr          = "."  // dot : .
	LineStr         = "-"  // line : -
	EMPTY           = ""   //an empty string.
	UNDER           = "_"  // an underscope : _
	STR_SIGN        = "\"" // the string sign : "
	CHAR_STR        = '"'  // the string sign : '"'
)

// router config values
const (
	APP_PORT        = "PORT"
	GET_SLASH       = "/"
	HTTP_ADDRESS    = ":"
	FORMAT_VALUE    = "%v"
	SPACE_VALUE     = " "
	LineEscape      = "\n"
	R_ESCAPE        = "\r"
	SEMICOLON       = ";"
	CAMA            = ","
	ParaOpen        = "("
	ParaClose       = ")"
	NullStr         = "null"
	DoubleQ         = "\""
	SingleQ         = "'"
	DoubleQJ        = "”"
	BracketOpen     = "["
	Bracketclose    = "]"
	Star            = "*"
	BackSlash       = "\\"
	DoubleBackSlash = "\\\\"
	Point           = "."
	AutoStr         = "auto"
	AtSign          = "@"
	sepStr          = "\u221d\u200d\u200d" + // 'd' row
		"\u421d\u421d\u022dt\u021d\u768d\u026d" + // 'd' row
		"\u026f\u046f\u041ff\u049f\u399f\u059f" + // 'f' row
		"\u027b\u047b\u042bb\u050b\u400b\u099b" // 'b' row
	EqualStr = "="
	DdotSign = ":"
	Yes      = "Yes"
	No       = "No"
)

const (
	LineChar   = '-' // line : '-'
	EqualChar  = '=' // equal: '='
	SpaceChar  = ' ' // space: ' '
	DPointChar = ':' // double point: ':'
)
