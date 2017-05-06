package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"os/exec"
	"syscall"

	"github.com/orivej/e"
)

func getRegs(pid int) syscall.PtraceRegs {
	var regs syscall.PtraceRegs
	err := syscall.PtraceGetRegs(pid, &regs)
	e.Exit(err)
	return regs
}

func readString(pid int, addr uint64) string {
	var chunk [64]byte
	var buf []byte
	for {
		n, err := syscall.PtracePeekData(pid, uintptr(addr), chunk[:])
		if err != syscall.EIO {
			e.Print(err)
		}
		end := bytes.IndexByte(chunk[:n], 0)
		if end != -1 {
			buf = append(buf, chunk[:end]...)
			return string(buf)
		}
		buf = append(buf, chunk[:n]...)
		addr += uint64(n)
	}
}

func readStrings(pid int, addr uint64) []string {
	var buf [8]byte
	var res []string
	for {
		n, err := syscall.PtracePeekData(pid, uintptr(addr), buf[:])
		e.Exit(err)
		saddr := binary.LittleEndian.Uint64(buf[:])
		if saddr == 0 {
			return res
		}
		res = append(res, readString(pid, saddr))
		addr += uint64(n)
	}
}

func resume(pid, signal int) {
	err := syscall.PtraceSyscall(pid, signal)
	e.Exit(err)
}

func waitForSyscall() (int, int, bool) {
	var wstatus syscall.WaitStatus
	for {
		pid, err := syscall.Wait4(-1, &wstatus, syscall.WALL, nil)
		e.Exit(err)
		switch {
		case wstatus.Stopped() && wstatus.StopSignal()&0x80 != 0:
			return pid, 0, true
		case wstatus.Stopped() && wstatus.TrapCause() > 0:
			return pid, wstatus.TrapCause(), true
		case wstatus.Exited() || wstatus.Signaled():
			return pid, 0, false
		case wstatus.Stopped():
			resume(pid, int(wstatus.StopSignal()))
		default:
			panic(wstatus)
		}
	}
}

func trace(tracee string, argv []string) (*os.Process, error) {
	cmd := exec.Command(tracee, argv...) //#nosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	return cmd.Process, err
}
