package tests

import (
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

	log.Println(result.Stdout)
}
