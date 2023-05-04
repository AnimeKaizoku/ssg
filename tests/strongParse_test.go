package tests

import (
	"log"
	"os"
	"testing"

	"github.com/AnimeKaizoku/ssg/ssg/strongParser"
)

type MyConfigStruct struct {
	TheToken        string     `section:"main" key:"the_token"`
	MyToken         string     `key:"the_token"`
	BotId           *int64     `section:"main" key:"bot_id"`
	BotName         *string    `section:"main" key:"bot_name"`
	BotComplex      complex128 `section:"main" key:"bot_complex"`
	TheUsername     string     `section:"telegram" key:"bot_username"`
	BotUsername     string     `section:"telegram"`
	SinglePrefix    rune       `section:"telegram" key:"single_prefix" type:"rune"`
	CmdPrefixed     []rune     `section:"telegram" key:"cmd_prefixes" type:"[]rune"`
	BotOwner        int64      `section:"telegram" key:"bot_owner"`
	OwnerIds        []int64    `section:"telegram" key:"owner_ids"`
	OwnerNumbers    []int32    `section:"telegram" key:"owner_numbers"`
	OwnerNames      []string   `section:"telegram" key:"owner_names"`
	OwnerSupporting []bool     `section:"telegram" key:"owner_supporting"`
	DatabaseUrl     string     `section:"database" key:"url"`
	UseSqlite       bool       `section:"database" key:"use_sqlite" default:"true"`
	APIUrl          string     `key:"api_url"`
	MyUser          string
	MyBool          bool
	OwnerUsernames  []string
}

type MainSectionStruct struct {
	PgDump      string `key:"pg_dump_command"`
	LogChannels string `key:"log_channels"`
}

type ValueSectionStruct struct {
	TheToken        string   `key:"the_token"`
	BotUsername     string   `key:"bot_username"`
	SinglePrefix    rune     `key:"single_prefix" type:"rune"`
	CmdPrefixed     []rune   `key:"cmd_prefixes" type:"[]rune"`
	BotOwner        int64    `key:"bot_owner"`
	OwnerIds        []int64  `key:"owner_ids"`
	OwnerNumbers    []int32  `key:"owner_numbers"`
	OwnerNames      []string `key:"owner_names"`
	OwnerSupporting []bool   `key:"owner_supporting"`
	sectionName     string
}

type MyRuneStruct struct {
	SinglePrefix rune    `section:"telegram" key:"single_prefix" type:"rune"`
	CmdPrefixes  []rune  `section:"telegram" key:"cmd_prefixes" type:"[]rune"`
	shouldIgnore *string `section:"telegram" key:"should_ignore"`
}

const TheStrValue01 = `
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

const TheStrValue02 = `
[main]
pg_dump_command = pg_dump
log_channels = 12454, -124578

[SaitamaRobot]
the_token = 12345:abcd
bot_username = @SaitamaRobot
bot_owner = 123456
owner_ids = 123456, 1234567
owner_names = sayan, aliwoto, sawada
owner_numbers = 1234567, 12345678, 123456789
owner_supporting = true, false, true
single_prefix = !
cmd_prefixes = /, !, #

[KigyoRobot]
the_token = 72345:abcd
bot_username = @kigyorobot
bot_owner = 123456
owner_ids = 8787871234, 1234567
owner_names = sayan, aliwoto, sawada
owner_numbers = 1234567, 12345678, 123456789
owner_supporting = true, false, true
single_prefix = !
cmd_prefixes = /, !, #

[ShellderRobot]
the_token = 82345:abcdefg
bot_username = @ShellderRobot
bot_owner = 123456
owner_ids = 123456, 1234567
owner_names = sayan, aliwoto, sawada
owner_numbers = 1234567, 12345678, 123456789
owner_supporting = true, false, true
single_prefix = !
cmd_prefixes = /, !, #

`

func (v *ValueSectionStruct) SetSectionName(name string) {
	v.sectionName = name
}

func (v *ValueSectionStruct) GetSectionName() string {
	return v.sectionName
}

func TestMainAndArrayParser(t *testing.T) {
	opt := &strongParser.ConfigParserOptions{}
	container, err := strongParser.ParseMainAndArraysStr[MainSectionStruct, ValueSectionStruct](TheStrValue02, opt)
	if err != nil {
		t.Error(err)
		return
	}

	if container == nil {
		t.Error("got nil container")
		return
	}

	if container.Sections[0].BotUsername != "@SaitamaRobot" {
		t.Errorf("expected @SaitamaRobot, got: %s", container.Sections[0].BotUsername)
		return
	}

	if container.Sections[1].OwnerIds[0] != 8787871234 {
		t.Errorf("expected 8787871234, got: %d", container.Sections[1].OwnerIds[0])
		return
	}
}

func TestStrongParser(t *testing.T) {
	myValue := &MyConfigStruct{}
	err := strongParser.ParseStringConfig(myValue, TheStrValue01)
	if err != nil {
		t.Error(err)
		return
	}

	if myValue.TheToken != myValue.MyToken {
		t.Errorf("TheToken should be equal to MyToken, got: %s", myValue.TheToken)
		return
	}

	if myValue.BotUsername != myValue.TheUsername {
		t.Errorf("BotUsername should be equal to MyUsername, got: %s", myValue.BotUsername)
		return
	}

	log.Println(myValue)
}

func TestParseFromEnv(t *testing.T) {
	myValue := &MyConfigStruct{}
	os.Setenv("API_URL", "google.com")
	os.Setenv("my_user", "woto")
	os.Setenv("MY_BOOL", "yes")
	os.Setenv("owner_usernames", "abcd abc123 ab12345 a7634")

	err := strongParser.ParseStringConfigWithOption(myValue, TheStrValue01, &strongParser.ConfigParserOptions{
		ReadEnv: true,
	})
	if err != nil {
		t.Error(err)
		return
	}

	if myValue.APIUrl != "google.com" {
		t.Errorf("APIUrl should be google.com, got: %s", myValue.APIUrl)
		return
	}

	if myValue.MyUser != "woto" {
		t.Errorf("MyUser should be woto, got: %s", myValue.MyUser)
		return
	}

	if !myValue.MyBool {
		t.Error("MyBool should be true")
		return
	}

	if myValue.OwnerUsernames[2] != "ab12345" {
		t.Errorf("OwnerUsernames[0] should be abcd, got: %s", myValue.OwnerUsernames[0])
		return
	}
}

func TestStrongRuneParser(t *testing.T) {
	myValue := &MyRuneStruct{}
	err := strongParser.ParseStringConfig(myValue, TheStrValue01)
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
