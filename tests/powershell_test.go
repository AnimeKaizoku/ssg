package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/AnimeKaizoku/ssg/ssg"
)

const (
	PSCode01 = `
	$theWrite = Write-Output -InputObject "test"
	
	Write-Output $theWrite.GetType()
	
	`
)

func TestPowerShell01(t *testing.T) {
	finishedChan := make(chan bool)
	result := ssg.RunPowerShellAsyncWithChan(PSCode01, finishedChan)

	<-finishedChan

	result.PurifyPowerShellOutput()

	fmt.Println(result.Stdout)
}

func TestPowerShell02(t *testing.T) {
	result := ssg.RunPowerShell("$PSVersionTable.PSVersion")
	result.PurifyPowerShellOutput()

	fmt.Println(result.Stdout)
}

func TestPowerShell03(t *testing.T) {
	finishedChan := make(chan bool)
	result := ssg.RunPowerShellAsyncWithChan("$PSVersionTable.PSVersion", finishedChan)

	<-finishedChan

	result.PurifyPowerShellOutput()

	log.Println(result.Stdout)
}
