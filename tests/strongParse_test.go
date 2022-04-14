package tests

import (
	"log"
	"testing"

	"github.com/AnimeKaizoku/ssg/ssg/strongParser"
)

type MyConfigStruct struct {
	TheToken        string     `section:"main" key:"the_token"`
	BotId           *int64     `section:"main" key:"bot_id"`
	BotName         *string    `section:"main" key:"bot_name"`
	BotComplex      complex128 `section:"main" key:"bot_complex"`
	BotUsername     string     `section:"telegram" key:"bot_username"`
	SinglePrefix    rune       `section:"telegram" key:"single_prefix"`
	CmdPrefixed     []rune     `section:"telegram" key:"cmd_prefixed"`
	BotOwner        int64      `section:"telegram" key:"bot_owner"`
	OwnerIds        []int64    `section:"telegram" key:"owner_ids"`
	OwnerNumbers    []int32    `section:"telegram" key:"owner_numbers"`
	OwnerNames      []string   `section:"telegram" key:"owner_names"`
	OwnerSupporting []bool     `section:"telegram" key:"owner_supporting"`
	DatabaseUrl     string     `section:"database" key:"url"`
	UseSqlite       bool       `section:"database" key:"use_sqlite" default:"true"`
}

type MyRuneStruct struct {
	SinglePrefix rune    `section:"telegram" key:"single_prefix" type:"rune"`
	CmdPrefixes  []rune  `section:"telegram" key:"cmd_prefixes" type:"[]rune"`
	shouldIgnore *string `section:"telegram" key:"should_ignore"`
}

const TheStrValue = `
[main]
the_token = 12345:abcd
bot_id = 202012345
bot_name = kigyo
bot_complex = 1.2+3.4i

[telegram]
bot_username = @kigyorobot
bot_owner = 123456
owner_ids = 123456, 1234567
owner_names = sayan, aliwoto, sawada
owner_numbers = 1234567, 12345678, 123456789
owner_supporting = true, false, true
single_prefix = !
cmd_prefixes = /, !, #

[database]
url = postgres://user:pass@localhost:5432/dbname
# use_sqlite = true

`

func TestStrongParser(t *testing.T) {
	myValue := &MyConfigStruct{}
	err := strongParser.ParseStringConfig(myValue, TheStrValue)
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(myValue)
}

func TestStrongRuneParser(t *testing.T) {
	myValue := &MyRuneStruct{}
	err := strongParser.ParseStringConfig(myValue, TheStrValue)
	if err != nil {
		t.Error(err)
		return
	}

	if myValue.shouldIgnore != nil {
		t.Error("should ignore should be nil")
		return
	}

	log.Println(myValue)
}
