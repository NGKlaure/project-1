package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os/exec"
	"strings"

	//"net"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"golang.org/x/crypto/ssh"
)

const (
	remoteUser = "agent1"
	remoteHost = "192.168.56.102"
	port       = " 22"
)

type Machine struct {
	Name string
	//Pid  int
}

type Process struct {
	//Name string
	Pid int
	Cpu float64
}

var pw string

func connect() (*ssh.Client, *ssh.Session) {

	if pw == "" {
		fmt.Println("password")
		fmt.Scan(&pw)
		fmt.Print("\n")
	}

	sshConfig := &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{ssh.Password(pw)},
		//AuthMethod{ssh.password(pw)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// start client connection to ssh server
	connection, err := ssh.Dial("tcp", remoteHost+";"+port, sshConfig)
	if err != nil {
		connection.Close()
		panic(err)
	}
	session, err := connection.NewSession()
	if err != nil {
		session.Close()
		panic(err)
	}
	return connection, session

}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(0)
	}
}

func index(response http.ResponseWriter, request *http.Request) {
	machine := Machine{}

	temp, _ := template.ParseFiles("template/index.html")

	/* connection, session := connect()

	defer connection.Close()
	defer session.Close()*/
	hostname, err := os.Hostname()
	handleErr(err)
	machine.Name = hostname //get host name
	//machine.Pid = "1234"
	//out,_ := session.CombinedOutput()

	temp.Execute(response, machine)
}

func GetInterfaceInfos(w http.ResponseWriter, r *http.Request) {

	// get interfaces MAC/hardware address

	interfStat, err := net.Interfaces()
	handleErr(err)

	html := "<html>Interfaces MAC and Hardware address infos " + "<br>"
	html = html + "<br>"

	for _, interf := range interfStat {
		html = html + "------------------------------------------------------<br>"
		html = html + "Interface Name: " + interf.Name + "<br>"

		if interf.HardwareAddr != "" {
			html = html + "Hardware(MAC) Address: " + interf.HardwareAddr + "<br>"
		}

		for _, flag := range interf.Flags {
			html = html + "Interface behavior or flags: " + flag + "<br>"
		}

		for _, addr := range interf.Addrs {
			html = html + "IPv6 or IPv4 addresses: " + addr.String() + "<br>"

		}

	}

	html = html + "</html>"

	w.Write([]byte(html))

}

func GetCPUData(w http.ResponseWriter, r *http.Request) {

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	handleErr(err)
	percentage, err := cpu.Percent(0, true)
	handleErr(err)

	html := "<html>CPU infos " + "<br>"

	html = html + "<br>"
	html = html + "CPU index number: " + strconv.FormatInt(int64(cpuStat[0].CPU), 10) + "<br>"
	html = html + "VendorID: " + cpuStat[0].VendorID + "<br>"
	html = html + "Family: " + cpuStat[0].Family + "<br>"
	html = html + "Number of cores: " + strconv.FormatInt(int64(cpuStat[0].Cores), 10) + "<br>"
	html = html + "Model Name: " + cpuStat[0].ModelName + "<br>"
	html = html + "Speed: " + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz <br>"

	for idx, cpupercent := range percentage {
		html = html + "Current CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%<br>"
	}

	html = html + "</html>"

	w.Write([]byte(html))

}

func GetDiskData(w http.ResponseWriter, r *http.Request) {

	diskStat, err := disk.Usage("/")
	handleErr(err)

	html := "<html> Disk infos " + "<br>"
	html = html + "<br>"

	html = html + "Disk Path: " + diskStat.Path + "  <br>"
	html = html + "Disk File system type: " + diskStat.Fstype + "  <br>"
	html = html + "Total disk space: " + strconv.FormatUint(diskStat.Total, 10) + " bytes <br>"
	html = html + "Used disk space: " + strconv.FormatUint(diskStat.Used, 10) + " bytes<br>"
	html = html + "Free disk space: " + strconv.FormatUint(diskStat.Free, 10) + " bytes<br>"
	html = html + "Percentage disk space usage: " + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%<br>"

	html = html + "</html>"

	w.Write([]byte(html))

}

func GetHostInfos(w http.ResponseWriter, r *http.Request) {
	runtimeOS := runtime.GOOS

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	handleErr(err)

	//html := "<html>OS : " + runtimeOS + "<br>"
	html := "<html> Host infos " + "<br>"
	html = html + "<br>"
	html = html + " runtime OS : " + runtimeOS + "<br>"

	html = html + "Hostname: " + hostStat.Hostname + "<br>"
	html = html + "Uptime: " + strconv.FormatUint(hostStat.Uptime, 10) + "<br>"

	html = html + "OS: " + hostStat.OS + "<br>"
	html = html + "Platform: " + hostStat.Platform + "<br>"

	// the unique hardware id for this machine
	html = html + "Host ID(uuid): " + hostStat.HostID + "<br>"
	html = html + "</html>"

	w.Write([]byte(html))

}

func GetProcInfos(w http.ResponseWriter, r *http.Request) {

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	handleErr(err)
	//get running proccesses
	miscStat, err := load.Misc()
	handleErr(err)

	//html := "<html>OS : " + runtimeOS + "<br>"
	html := "<html>Processes infos " + "<br>"
	html = html + "<br>"

	html = html + "total Number of processes: " + strconv.FormatUint(hostStat.Procs, 10) + "<br>"
	html = html + "Number of processes running: " + strconv.FormatInt(int64(miscStat.ProcsRunning), 10) + "<br>"
	html = html + "Number of blocked  prossesses: " + strconv.FormatInt(int64(miscStat.ProcsBlocked), 10) + "<br>"

	html = html + "</html>"

	w.Write([]byte(html))

}

func GetMemoryInfos(w http.ResponseWriter, r *http.Request) {
	runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	handleErr(err)

	html := "<html>Memory infos " + "<br>"
	html = html + "<br>"
	html = html + "OS : " + runtimeOS + "<br>"
	html = html + "Total memory: " + strconv.FormatUint(vmStat.Total, 10) + " bytes <br>"
	html = html + "Free memory: " + strconv.FormatUint(vmStat.Free, 10) + " bytes<br>"
	html = html + "Percentage memory used : " + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%<br>"
	html = html + "</html>"

	w.Write([]byte(html))

}

func PrintProcInfos(w http.ResponseWriter, r *http.Request) {

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	handleErr(err)
	//get running proccesses
	miscStat, err := load.Misc()
	handleErr(err)

	//html := "<html>OS : " + runtimeOS + "<br>"
	html := "<html>Processes infos " + "<br>"
	html = html + "<br>"

	html = html + "total Number of processes: " + strconv.FormatUint(hostStat.Procs, 10) + "<br>"
	html = html + "Number of processes running: " + strconv.FormatInt(int64(miscStat.ProcsRunning), 10) + "<br>"
	html = html + "Number of blocked  prossesses: " + strconv.FormatInt(int64(miscStat.ProcsBlocked), 10) + "<br>"

	processList, err := ps.Processes()
	if err != nil {
		log.Println("ps.Processes() Failed, are you using windows?")
		return
	}

	for x := range processList {

		var process ps.Process
		process = processList[x]
		html = html + "------------------------------------------------------------------------------------------------------------<br>"
		html = html + "process  PID: " + strconv.Itoa(process.Pid()) + "        executable Name: " + process.Executable() + "<br>"
		//log.Printf("%d\t%s\n", process.Pid(), process.Executable())

		// do os.* stuff on the pid
	}

	html = html + "</html>"

	w.Write([]byte(html))

}

func PrintProcCPUInfos(w http.ResponseWriter, r *http.Request) {

	//html := "<html>OS : " + runtimeOS + "<br>"
	html := "<html>Processes infos " + "<br>"
	html = html + "<br>"

	//connection, session := connect()
	//out, _ := session.CombinedOutput(cmd)
	//defer connection.Close()
	//defer session.Close()

	cmd := exec.Command("ps", "aux")
	//cmd := executeCmd("ps aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		//log.Println(len(ft), ft)

		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)
		if err != nil {
			log.Fatal(err)
		}
		processes = append(processes, &Process{pid, cpu})
	}
	for _, p := range processes {
		html = html + "------------------------------------------------------------------------------<br>"
		html = html + "process with PID: " + strconv.Itoa(p.Pid) + "  takes : " + strconv.FormatFloat(p.Cpu, 'f', 2, 64) + "% of  the cpu<br>"
		//log.Println("Process ", p.Pid, " takes ", p.Cpu, " % of the CPU")
	}
	html = html + "</html>"

	w.Write([]byte(html))

}

func Killpform(response http.ResponseWriter, request *http.Request) {
	proc := Process{}

	machine := Machine{}
	temp, _ := template.ParseFiles("template/killp.html")

	hostname, err := os.Hostname()
	handleErr(err)
	machine.Name = hostname //get host name
	proc.Pid = 1234
	//machine.Pid = request.FormValue("Pid")

	temp.Execute(response, machine)
}

func Formsubmit(w http.ResponseWriter, r *http.Request) {

	temp, _ := template.ParseFiles("template/killp.html")

	proc := Process{}
	arg1 := "kill"
	//getpid from the html form as a string
	getpid := r.FormValue("pid")

	//convert into and int
	i1, err := strconv.Atoi(getpid)
	handleErr(err)
	proc.Pid = i1

	//convert into a string to use in exec
	machinepid := strconv.Itoa(proc.Pid)

	cmd := exec.Command(arg1, machinepid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Cannot find process")
		os.Exit(1)
	}
	fmt.Printf("Status is: %s", string(out))

	temp.Execute(w, proc)
}

//
func killproc() {
	fmt.Println("enter a procces id that you want to kill")
	var inputPid int
	fmt.Scanln(&inputPid)
	arg1 := "kill"
	arg2 := strconv.Itoa(inputPid)
	cmd := exec.Command(arg1, arg2)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Cannot find process")
		os.Exit(1)
	}
	fmt.Printf("Status is: %s", string(out))
}

func executeCmd(cmd string) []byte {
	connection, session := connect()
	out, _ := session.CombinedOutput(cmd)
	defer connection.Close()
	defer session.Close()
	return out

}

func handler() {
	http.HandleFunc("/", index)

	http.HandleFunc("/getCPUdata", GetCPUData)
	http.HandleFunc("/getDiskdata", GetDiskData)
	http.HandleFunc("/getHostInfos", GetHostInfos)
	http.HandleFunc("/getProcInfos", GetProcInfos)
	http.HandleFunc("/getMemoryInfos", GetMemoryInfos)
	http.HandleFunc("/getInterfaceInfos", GetInterfaceInfos)
	http.HandleFunc("/PrintProcInfos", PrintProcInfos)
	http.HandleFunc("/PrintProc CPUInfos", PrintProcCPUInfos)

	http.HandleFunc("/killpform", Killpform)
	http.HandleFunc("/formsubmit", Formsubmit)

	http.ListenAndServe(":7000", nil)
}

func main() {

	handler()

	//killproc()
	//exec.Command("ssh", "agent1@192.168.56.102", "ls").Run()

}
