package builder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
	"xray-telegram/entity"

	"github.com/google/uuid"
)

// SetServerIP  Get preferred outbound ip of this machine
func (b *Builder) SetServerIP() *Builder {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("error during the SetServerIP ", err)
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	b.ServerIP = localAddr.IP.String()
	return b
}

// SetSettingsFile returns the settings file
func (b *Builder) SetSettingsFile() *Builder {

	// Open our jsonFile
	jsonFile, err := os.Open("./setting.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("error open setting file", err)
		return nil
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error ReadAll json file", err)
		return nil
	}

	setting := entity.Setting{}
	// we unmarshal our byteArray which contains our
	err = json.Unmarshal(byteValue, &setting)
	if err != nil {
		fmt.Println("error unmarshal setting", err)
		return nil
	}

	b.Setting = setting
	return b

}

func randomNumber(max int) int {
	// Set the seed value for the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate a random Int type number between 0 and 9
	return 0 + r.Intn(max-0)
}

// SetConfigurations sets the xray configuration
func (b *Builder) SetConfigurations() *Builder {

	//Random  make value
	hosts := []string{"mashhad1.irancell.ir", "shiraz1.irancell.ir", "tabriz1.irancell.ir", "speedtest1.irancell.ir", "ahvaz1.irancell.ir", "esfahan1.irancell.ir", "server-9889.prod.hosts.ooklaserver.net", "server-10076.prod.hosts.ooklaserver.net", "server-9795.prod.hosts.ooklaserver.net", "server-4317.prod.hosts.ooklaserver.net"}
	hostSelected1 := hosts[randomNumber(len(hosts))]
	hostSelected2 := hosts[randomNumber(len(hosts))]
	hostSelected3 := hosts[randomNumber(len(hosts))]
	hostSelected4 := hosts[randomNumber(len(hosts))]
	hostSelected := hostSelected1 + "," + hostSelected2 + "," + hostSelected3 + "," + hostSelected4

	ports := []int{8081}
	portSelected := ports[randomNumber(len(ports))]

	methods := []string{"GET", "POST"}
	methodSelected := methods[randomNumber(len(methods))]

	path := []string{
		"/upload",
		"/download",
		"/api/v1/users",
		"/api/v2/users",
		"/blocks/v1/users",
		"/negotiation/users/v1",
		"/books/v1/handler",
		"/shop/buyer/v2",
		"/help-desk/v1/help",
		"/apple/users",
		"/humbles/users",
		"/dirty/class/v1",
		"/small-house/v2/users",
		"/low/users/schools/v1",
		"/high/users/schools/v2",
		"/users/v1/low/schools",
		"/hi-users/v2/high/schools",
		"/api-api/c2/v1/user",
		"/book/users/v1",
		"/rest/v2/shop",
		"/humble/v3/school/thief",
	}
	pathSelected := path[randomNumber(len(path))]

	contextLength := strconv.Itoa(100 + randomNumber(100))

	message := []string{"OK", "Not Found", "Bad Request", "Forbidden", "Internal Server Error", "Service Unavailable"}
	messageSelected := message[randomNumber(len(message))]

	statuses := []string{"200", "202", "404", "400", "403", "500", "503"}
	statusSelected := statuses[randomNumber(len(statuses))]

	if !b.Setting.RandomHeader {
		portSelected = 8081
		methodSelected = "GET"
		pathSelected = "/download"
		contextLength = "109"
		messageSelected = "OK"
		statusSelected = "200"
	}

	//Random  make value

	b.newVmess.Inbounds = make([]entity.Inbound, 1)

	var inbound entity.Inbound
	inbound.Listen = nil
	inbound.Port = portSelected
	inbound.Protocol = "vmess"
	inbound.Settings.Clients = make([]entity.Client, 1)
	inbound.Settings.Clients[0].Email = b.Setting.ChannelName
	inbound.Settings.Clients[0].ID = uuid.New().String()

	inbound.Sniffing.Enabled = true
	inbound.Sniffing.DestOverride = []string{"http", "tls", "quic", "fakedns"}

	inbound.StreamSettings.Network = "tcp"
	inbound.StreamSettings.Security = "none"

	inbound.StreamSettings.Sockopt.AcceptProxyProtocol = false
	inbound.StreamSettings.Sockopt.Mark = 0
	inbound.StreamSettings.Sockopt.Tproxy = "off"
	inbound.StreamSettings.Sockopt.TCPFastOpen = true

	inbound.StreamSettings.TCPSettings.AcceptProxyProtocol = false
	inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host[0] = hostSelected
	inbound.StreamSettings.TCPSettings.Header.Request.Method = methodSelected
	inbound.StreamSettings.TCPSettings.Header.Request.Path = []string{pathSelected}
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.Connection = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.Connection[0] = "keep-alive"

	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentLength = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentLength[0] = contextLength
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentType = make([]string, 1)
	inbound.StreamSettings.TCPSettings.Header.Response.Headers.ContentType[0] = "text/html"

	inbound.StreamSettings.TCPSettings.Header.Response.Reason = messageSelected
	inbound.StreamSettings.TCPSettings.Header.Response.Status = statusSelected
	inbound.StreamSettings.TCPSettings.Header.Response.Version = "1.1"

	inbound.StreamSettings.TCPSettings.Header.Type = "http"

	port := "80"

	inbound.Tag = "inbound-" + port

	code := "{\"add\":\"" + b.ServerIP + "\",\"aid\":\"0\",\"host\":\"" + inbound.StreamSettings.TCPSettings.Header.Request.Headers.Host[0] + "\",\"id\":\"" + inbound.Settings.Clients[0].ID + "\",\"net\":\"tcp\",\"path\":\"" + pathSelected + "\",\"port\":\"" + port + "\",\"ps\":\"" + b.Setting.ChannelName + "\",\"scy\":\"auto\",\"sni\":\"\",\"tls\":\"\",\"type\":\"http\",\"v\":\"2\"}"
	base64code := base64.StdEncoding.EncodeToString([]byte(code))
	StringConfig := "vmess://" + base64code

	b.StringConfigZero = StringConfig
	b.newVmess.Inbounds[0] = inbound

	return b

}

// SetBlock block Iranian and Chinese and porn websites
func (b *Builder) SetBlock() *Builder {

	b.newVmess.Log.Loglevel = "warning"
	b.newVmess.Routing.DomainStrategy = "IPOnDemand"
	b.newVmess.Routing.Rules = make([]entity.Rule, 2)
	b.newVmess.Routing.Rules[0] = entity.Rule{
		Type:        "field",
		IP:          []string{"geoip:cn", "geoip:ir"},
		OutboundTag: "block",
	}

	b.newVmess.Routing.Rules[1] = entity.Rule{
		Type:        "field",
		Domain:      []string{"geosite:category-porn"},
		OutboundTag: "block",
	}

	b.newVmess.Outbounds = make([]entity.Outbound, 2)
	b.newVmess.Outbounds[0] = entity.Outbound{
		Tag:      "direct",
		Protocol: "freedom",
		Settings: struct{}{},
	}
	b.newVmess.Outbounds[1] = entity.Outbound{
		Tag:      "block",
		Protocol: "blackhole",
		Settings: struct{}{},
	}

	return b

}

func (b *Builder) Save() *Builder {

	if b.StringConfigZero == "" {
		fmt.Println(" string config zero is empty")
		return nil
	}

	//save new Reality in file
	err := WriteFile("./config.json", b.newVmess)
	if err != nil {
		log.Fatal("error during the WriteFile ", err)
		return nil
	}

	return b

}
