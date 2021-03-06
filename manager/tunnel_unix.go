// +build linux darwin

package manager

import (
	"bufio"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/mdsauce/sauced/logger"
)

// Start creates a new tunnel from the metadata and launch arguments
func Start(launchArgs string, wg *sync.WaitGroup, meta Metadata) {
	defer wg.Done()
	args := strings.Split(launchArgs, " ")
	path := args[0]

	if eatLine(path) {
		return
	}
	if vacancy(meta) != true {
		logger.Disklog.Warnf("Too many tunnels open.  Not opening %s \n %v", meta.Pool, launchArgs)
		return
	}

	manufacturedArgs := setDefaults(args)
	meta.Owner = GetOwner(strings.Join(manufacturedArgs, " "))
	logger.Disklog.Debug("Created new set of args with sensible defaults that will be passed to exec.Command: ", manufacturedArgs)
	// tunnel is actually launched here.  new process is spawned
	rand.Seed(time.Now().UnixNano())
	wait := rand.Intn(15)
	time.Sleep(time.Duration(wait) * time.Second)
	scCmd := exec.Command(path, manufacturedArgs[1:]...)
	stdout, _ := scCmd.StdoutPipe()
	err := scCmd.Start()
	if err != nil {
		logger.Disklog.Warnf("Something went wrong while starting the SC binary! %v", err)
		return
	}

	logger.Disklog.Infof("Tunnel started as process %d - %s.  Raw launch arguments: %s\n", scCmd.Process.Pid, manufacturedArgs, launchArgs)
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	// this parsing should be moved to its own funciton.
	// everything should be parsed then supplied to the AddTunnel() func
	var tunLog string
	var asgnID string
	for scanner.Scan() {
		m := scanner.Text()
		// should be a func that is unit tested
		if strings.Contains(m, "Log file:") {
			ll := strings.Split(m, " ")
			tunLog = ll[len(ll)-1]
			logger.Disklog.Debugf("Tunnel log started for tunnel: %s \n %s", launchArgs, m)
		}
		// should be a func that is unit tested
		if strings.Contains(m, "Tunnel ID:") {
			idLine := strings.Split(m, " ")
			asgnID = idLine[len(idLine)-1]
			logger.Disklog.Infof("TUNNEL IS ALIVE as process %d with Assigned ID %s. args: %s", scCmd.Process.Pid, asgnID, manufacturedArgs)
		}
		if strings.Contains(m, "Sauce Connect is up") {
			AddTunnel(launchArgs, path, scCmd.Process.Pid, meta, tunLog, asgnID)
		}
	}
	logger.Disklog.Infof("Sauce Connect client with PID %d shutting down!  If you want more details check our logfile %s  Goodbye!", scCmd.Process.Pid, tunLog)
	RemoveTunnel(scCmd.Process.Pid)
	defer scCmd.Wait()
}

// Stop will halt a running process with SIGINT(CTRL-C)
func Stop(Pid int) {
	tunnel, err := os.FindProcess(Pid)
	if err != nil {
		logger.Disklog.Warnf("Process ID %d does not exist or was not accessible for this user. Error: %v", Pid, err)
	} else {
		err := tunnel.Signal(os.Interrupt)
		if err != nil {
			logger.Disklog.Warnf("Problem killing Process %d %v.  The user may not have permissions to send a SIGINT or SIGKILL to the listed process.", Pid, err)
		}
	}
}

// StopTunnelByID will stop a single tunnel that matches a given ID
func StopTunnelByID(assignedID string) {
	tstate := GetLastKnownState()
	tunnel, err := tstate.FindTunnelByID(assignedID)

	if err != nil {
		logger.Disklog.Warn(err)
	} else {
		Stop(tunnel.PID)
	}

}

// StopTunnelsByPool will stop a tunnel pool matching the given pool name
func StopTunnelsByPool(poolName string) {
	tstate := GetLastKnownState()
	tunnels, err := tstate.FindTunnelsByPool(poolName)

	if err != nil {
		logger.Disklog.Warn(err)
	} else {
		for _, tunnel := range tunnels {
			Stop(tunnel.PID)
		}
	}
}

// StopAll will send a kill or SIGINT signal
// to all tunnels that are running.
func StopAll() {
	last := GetLastKnownState()
	for _, tunnel := range last.Tunnels {
		Stop(tunnel.PID)
		// add a stop via REST API func here.  So rest api gets signal as well.
	}
}
