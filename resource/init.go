package resource

import (
	"github.com/google/uuid"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/models"
	"log"
	"net"
	"strings"
)

// tracker
var (
	NodeSession     *models.Session
	FilePieceMapper *models.PieceMapper
	NodeInfoMapper  *models.NodeInfoMapper
)

// node
var (
	// local info
	PieceLocalMapper *models.PieceLocalMapper
	TaskChan         chan *models.SeedInfo
	NodeId           string
	Addr             string
)

func InitBaseNodeInfo() {
	// nodeId is uuid
	guid := uuid.New()
	NodeId = guid.String()

	// get private ip
	ipList := GetIntranetIp()
	if len(ipList) > 0 {
		Addr = ipList[0] + ":" + extractPort(config.GetOpts().Web.ListenAddress)
	}
	log.Println("NodeId", NodeId, "Addr", Addr)
}

func extractPort(addr string) string {
	tempList := strings.Split(addr, ":")
	return tempList[1]
}

func GetIntranetIp() []string {
	result := make([]string, 0)
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return result
	}

	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}

		}
	}
	return result
}
