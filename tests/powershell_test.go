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

	PSCode02 = `
	function Get-GoFirstTenProcesses {
		return Get-Process -Name "*go*" | Select-Object -First 10
	}

	$ok = Get-GoFirstTenProcesses
	
	$ok[0].GetType().FullName

	Write-Output "We are done."
	Write-Output "Now!"
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
	//Write-Output ("ok" + $PSVersionTable.PSVersion.Major.ToString() + "`n`n" + $PSVersionTable.PSVersion)
	// result := ssg.RunPowerShell("$PSVersionTable.PSVersion")
	result := ssg.RunPowerShell(PSCode02)
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

func TestPowerShell04(t *testing.T) {
	result := ssg.RunPowerShell("Write-Output (\"ok\" + $PSVersionTable.PSVersion.Major.ToString() + \"`n`n\" + $PSVersionTable.PSVersion)")
	result.PurifyPowerShellOutput()

	fmt.Println(result.Stdout)
}
