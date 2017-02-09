package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"github.com/taskcluster/runlib/win32"
)

func main() {
	Run(
		os.Args[1],
		os.Args[2],
	)
}

func Run(username, password string) {

	fmt.Println("APPDATA test")

	var err error
	var user syscall.Handle
	var name *uint16
	var env uintptr
	var pinfo *win32.ProfileInfo

	name, err = syscall.UTF16PtrFromString(username)
	if err != nil {
		panic(err)
	}

	pinfo = &win32.ProfileInfo{
		Size:     uint32(unsafe.Sizeof(*pinfo)),
		Flags:    win32.PI_NOUI,
		Username: name,
	}

	err = win32.CreateEnvironmentBlock(&env, user, false)
	if err != nil {
		panic(err)
	}

	user, err = win32.LogonUser(
		syscall.StringToUTF16Ptr(username),
		syscall.StringToUTF16Ptr("."),
		syscall.StringToUTF16Ptr(password),
		win32.LOGON32_LOGON_INTERACTIVE,
		win32.LOGON32_PROVIDER_DEFAULT,
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		if user != syscall.Handle(0) && user != syscall.InvalidHandle {
			win32.CloseHandle(user)
		}
	}()

	err = win32.LoadUserProfile(user, pinfo)
	if err != nil {
		panic(err)
	}

	defer func() {
		if pinfo.Profile != syscall.Handle(0) && pinfo.Profile != syscall.InvalidHandle {
			for {
				err := win32.UnloadUserProfile(user, pinfo.Profile)
				if err == nil {
					break
				}
				log.Printf("%v", err)
			}
		}
	}()

	err = win32.CreateEnvironmentBlock(&env, user, false)
	if err != nil {
		panic(err)
	}
	defer win32.DestroyEnvironmentBlock(env)

	var varStartOffset uint
	for {
		envVar := syscall.UTF16ToString((*[1 << 15]uint16)(unsafe.Pointer(env + uintptr(varStartOffset)))[:])
		if envVar == "" {
			break
		}
		fmt.Println(envVar)
		// in UTF16, each rune takes two bytes, as does the trailing uint16(0)
		varStartOffset += uint(2 * (utf8.RuneCountInString(envVar) + 1))
	}
}
