package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	log "github.com/mgutz/logxi/v1"
	"../proto"
	"encoding/json"
	"github.com/skorobogatov/input"
)

var lconn, rconn *net.TCPConn
var parent *net.TCPConn
var l = false
var r = false
var first = true
var choose ,ch int
var encl, encr, encp  *json.Encoder

//подключение к род вершине
func connectParent(c string){
	for {
		if addr, err := net.ResolveTCPAddr("tcp",c); err != nil {
			log.Error("addr error")
		}else if conn, err := net.DialTCP("tcp",nil,addr); err != nil{
			log.Error("conn error")
		}else{
			parent = conn
			encp = json.NewEncoder(parent)
			interact(conn)
		}
	}
}
//создание запроса
func interact(conn *net.TCPConn){
	defer conn.Close()
	encoder, decoder := json.NewEncoder(conn), json.NewDecoder(conn)
	for {
		if first && ch != 0 {
			first = false
			send_request(encoder,"my addr",serverAddr1)
		}
		// Чтение команды из стандартного потока ввода
		fmt.Print("Enter your command \n")
		command := input.Gets()
		// Отправка запроса.
		switch command {
		case "quit":
			send_request(encoder, "quit", nil)
			return
		case "add":
			var n proto.MapPeer
			fmt.Print("Key := ")
			fmt.Scan(&n.Key)
			fmt.Println()
			fmt.Print("Value := ")
			fmt.Scan(&n.Value)
			fmt.Println()
			voc[n.Key] = n.Value
			n.Side = side
			if len(voc) == 0 {
				fmt.Println("Voc is empty :(")
			}else {
				fmt.Println("your voc :)")
				for key, value := range (voc) {
					fmt.Println("KEY ", key, " VALUE ", value)
				}
			}
			send_request(encoder, "add",n)
			if l {
				n.Side = 3
				send_request(encl,"add",n)
			}
			if r {
				n.Side = 3
				send_request(encr,"add",n)
			}
		case "delete":
			var n proto.MapPeer
			fmt.Print("Key := ")
			fmt.Scan(&n.Key)
			fmt.Println()
			delete(voc, n.Key)
			n.Side = side
			if len(voc) == 0 {
				fmt.Println("Voc is empty :(")
			}else {
				fmt.Println("your voc :)")
				for key, value := range (voc) {
					fmt.Println("KEY ", key, " VALUE ", value)
				}
			}
			send_request(encoder, "delete",n)
			if l {
				n.Side = 3
				send_request(encl,"delete",n)
			}
			if r {
				n.Side = 3
				send_request(encr,"delete",n)
			}
		case "check":
			if len(voc) == 0 {
				fmt.Println("Voc is empty :(")
			}else {
				for key, value := range (voc) {
					fmt.Println("KEY ", key, " VALUE ", value)
				}
			}
			continue
		case "find":
			var word string
			fmt.Print("input key please: ")
			fmt.Scan(&word)
			for key,value := range(voc){
				if key == word {
					fmt.Println("ok , here's the key :",value)
					break
				}
			}
			continue
		default:
			fmt.Printf("error:unknown command\n")
			continue
		}

		// Получение ответа.
		var resp proto.Response
		if err := decoder.Decode(&resp); err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		// Вывод ответа в стандартный поток вывода.
		switch resp.Status {
		case "ok":
			fmt.Printf("ok\n")
		case "failed":
			if resp.Data == nil {
				fmt.Printf("error: data field is absent in response\n")
			} else {
				var errorMsg string
				if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
					fmt.Printf("error: malformed data field in response\n")
				} else {
					fmt.Printf("failed: %s\n", errorMsg)
				}
			}
		default:
			log.Error("server error")
			fmt.Printf("error: server reports unknown status %q\n", resp.Status)
		}
	}
}
//обработка запроса
func handleRequest(req *proto.Request, encoder *json.Encoder) bool {
	switch req.Command {
	case "quit":
		respond(encoder,"ok", nil)
		return true
	case "my addr":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var addr string
			if err := json.Unmarshal(*req.Data, &addr); err != nil {
				errorMsg = "malformed data field"
			} else {
				if ad, err := net.ResolveTCPAddr("tcp", addr); err != nil {
					log.Error("addr error")
				} else if conn, err := net.DialTCP("tcp", nil, ad); err != nil {
					log.Error("dial error")
				} else {
					if (!l) {
						l = !l
						lconn = conn
						encl = json.NewEncoder(lconn)
						send_request(encl,"voc",voc)
						break
					} else {
						r = !r
						rconn = conn
						encr = json.NewEncoder(rconn)
						send_request(encr,"voc",voc)
						break
					}
				}
			}
		}
		if errorMsg == "" {
			respond(encoder,"ok", nil)
		} else {
			respond(encoder,"failed", errorMsg)
		}
	case "voc":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			if err := json.Unmarshal(*req.Data, &voc); err != nil {
				errorMsg = "malformed data field"
			} else {
				if len(voc) != 0 {
					fmt.Print("\nHey, here's your voc: ")
					for key, value := range voc {
						fmt.Println("KEY:", key, " VALUE:", value, ";\n ")
					}
				}
			}
			if errorMsg == "" {
				respond(encoder, "ok", nil)
			} else {
				fmt.Println("error")
				respond(encoder, "failed", errorMsg)
			}
		}
	case "delete":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var n proto.MapPeer
			if err := json.Unmarshal(*req.Data, &n); err != nil {
				errorMsg = "malformed data field"
			} else {
				delete(voc, n.Key)
				var curSide = n.Side
				if len(voc) == 0 {
					fmt.Println("Updated Voc is empty :(")
				}else {
					fmt.Println("Upd voc :)")
					for key, value := range (voc) {
						fmt.Println("KEY ", key, " VALUE ", value)
					}
				}
				if ch != 0 && curSide != 3 {
					n.Side = side
					send_request(encp, "delete",n)
				}
				if l && curSide != 1 {
					n.Side = 3
					send_request(encl, "delete",n)
				}
				if r && curSide != 0 {
					n.Side = 3
					send_request(encr, "delete", n)
				}
				if errorMsg == "" {
					respond(encoder, "ok", nil)
				} else {
					fmt.Println("error")
					respond(encoder, "failed", errorMsg)
				}
			}
		}

	case "add":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var n proto.MapPeer
			if err := json.Unmarshal(*req.Data, &n); err != nil {
				errorMsg = "malformed data field"
			} else {
				voc[n.Key] = n.Value
				var curSide = n.Side
				fmt.Println("Updated voc:")
				for key,value := range voc{
					fmt.Println("KEY:",key," VALUE:",value)
				}
				if ch != 0 && curSide != 3 {
					n.Side = side
					send_request(encp, "add",n)
				}
				if l && curSide != 1 {
					n.Side = 3
					send_request(encl, "add",n)
				}
				if r && curSide != 0 {
					n.Side = 3
					send_request(encr, "add", n)
				}
				if errorMsg == "" {
					respond(encoder, "ok", nil)
				} else {
					fmt.Println("error")
					respond(encoder, "failed", errorMsg)
				}
			}
		}
	default:
		//
	}
	return false
}
//отправка запроса
func send_request(encoder *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&proto.Request{command, &raw})
}
//обслуживание пира
func serve(conn *net.TCPConn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {

		} else {
			if handleRequest(&req,encoder) {
				log.Info("shutting down connection")
				break
			}
		}
	}
}
//отправка ответа
func respond(encoder *json.Encoder,status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&proto.Response{status, &raw})
}

var ipArr = [...]string{"127.0.0.1:6100","127.0.0.1:6010","127.0.0.1:6011","127.0.0.1:6012","127.0.0.1:6013","127.0.0.1:6014","127.0.0.1:6015","127.0.0.1:6016","127.0.0.1:6017","127.0.0.1:6018"}

var serverAddr1 string
var side = 1
var ll = false
var rr = false
var pp = false
var voc = make(map[string]string)
//реализация начального состояния пира
func main() {
	var (
		serverAddrStr string
		parentAddrStr string
		helpFlag      bool
	)
	fmt.Scan(&choose)
	ch = choose
	flag.StringVar(&serverAddrStr, "addr", ipArr[choose], "set server IP address and port")
	flag.BoolVar(&helpFlag, "help", false, "print options list")
	if flag.Parse(); helpFlag {
		fmt.Fprint(os.Stderr, "server [options]\n\nAvailable options:\n")
		flag.PrintDefaults()
	} else if serverAddr, err := net.ResolveTCPAddr("tcp", serverAddrStr); err != nil {
		log.Error("addr error")
	} else {
		serverAddr1 = serverAddrStr
		var i = 0
		if choose != 0 {
			for i = 0; i < 4; i++ {
				if (2*i) + 1 == choose {
					choose = i
					side = 1
					break
				}else{
					if (2*i) + 2 == choose{
						side = 0
						choose = i
						break
					}
				}
			}
			fmt.Println("Parent",choose)
			parentAddrStr = ipArr[choose]
			go connectParent(parentAddrStr)
		}else{
			side = 3
		}
		if listener, err := net.ListenTCP("tcp", serverAddr); err != nil {
			log.Error("listener error")
		} else {
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					fmt.Println("w8 for connection")
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())
					// Запуск go-программы для обслуживания клиентов.
					if !pp && ch != 0{
						pp = !pp
						go serve(conn)
					}else if !ll {
						ll = true
						go serve(conn)
					} else if !rr {
						rr = true
						go serve(conn)
					}
				}
			}
		}
	}
}
