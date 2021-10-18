// StrongStringGo Project
// Copyright (C) 2021 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package strongStringGo

// the prefex values for commands.
const (
	COMMAND_PREFIX1 = "!"
	COMMAND_PREFIX2 = "/"
	SUDO_PREFIX1    = ">"
	FLAG_PREFIX     = "--"
)

const (
	JA_Flag       = "〰\u200d；〰"
	JA_Str        = "❞\u200d；" // start character (") for string in japanese.
	JA_Equality   = "＝\u200d；" // equal character (＝) in japanese.
	JA_Ddot       = "：\u200d；" // ddot character (:) in japanese.
	JA_Cama       = "、\u200d；" // cama character (,) in japanese.
	JA_RealStr    = "\uff4e"   // the real str
	JA_BrOpen     = "「\u200d；" // the real str
	JA_BrClose    = "」\u200d；" // the real str
	BACK_Str      = `\"`
	BACK_Flag     = `\--`
	BACK_Equality = `\=`
	BACK_Ddot     = `\:`
	BACK_Cama     = `\,`
	BACK_BrOpen   = `\[`
	BACK_BrClose  = `\]`
)

// the base constant values.
const (
	BaseIndex      = 0  // number 0
	BaseOneIndex   = 1  // number 1
	BaseTwoIndex   = 2  // number 2
	BaseThreeIndex = 3  // number 2
	Base4Bit       = 4  // number 8
	Base8Bit       = 8  // number 8
	Base16Bit      = 16 // number 16
	Base32Bit      = 32 // number 32
	Base64Bit      = 64 // number 64
	BaseTimeOut    = 40 // 40 seconds
	BaseTen        = 10 // 10 seconds
)

// additional constants which are not actually used in
// this package, but may be useful in another packages.
const (
	BaseIndexStr    = "0" // number 0
	BaseOneIndexStr = "1" // number 1
	DotStr          = "." // dot : .
	LineStr         = "-" // line : -
	EMPTY           = ""  //an empty string.
	UNDER           = "_" // an underscope : _
	STR_SIGN        = `"` // the string sign : "
	CHAR_STR        = '"' // the string sign : '"'
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
	EqualStr        = "="
	DdotSign        = ":"
	Yes             = "Yes"
	No              = "No"
	OrRegexp        = "|" // the or string sign: "|"
)

const (
	LineChar         = '-' // line : '-'
	EqualChar        = '=' // equal: '='
	SpaceChar        = ' ' // space: ' '
	DPointChar       = ':' // double point: ':'
	BracketOpenChar  = '[' // bracket open: '['
	BracketcloseChar = ']' // bracket close: ']'
	CamaChar         = ',' // cama: ','
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
