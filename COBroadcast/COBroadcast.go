package COBroadcast

import "fmt"
import "strconv"

import "strings"
import BestEffortBroadcast "BestEffortBroadcast"

type COBroadcast_Req_Message struct {
	Addresses []string
	Message   string
}

type COBroadcast_Ind_Message struct {
	From    string
	Message string
}

type COBroadcast_Module struct {
	Ind      chan COBroadcast_Ind_Message
	Req      chan COBroadcast_Req_Message

	beb BestEffortBroadcast.BestEffortBroadcast_Module

	countAddresses int
	addresses []string

	countPending int
	pending []string
	rank []string

	lsn int
	v []string
}

func (module COBroadcast_Module) Init(addresses [] string) {

	fmt.Println("Init COBEB!")
	
	module.beb = BestEffortBroadcast.BestEffortBroadcast_Module{
		Req: make(chan BestEffortBroadcast.BestEffortBroadcast_Req_Message),
		Ind: make(chan BestEffortBroadcast.BestEffortBroadcast_Ind_Message)}

	module.countAddresses = len(addresses)
	module.addresses = addresses

	module.countPending = 0

	
	module.v = make([]string, module.countAddresses)
	for i := 0; i < module.countAddresses; i++ {//usar make
		module.v[i] = "0"
	}

	module.beb.Init(addresses[0])
	module.Start()
	fmt.Println("-------------------------------")
}

func (module COBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.COBroadcast(y)
				module.lsn = module.lsn + 1 //lsn = lsn + 1
			case y := <-module.beb.Ind:
				
				splitMessage := strings.Split(y.Message, ":")
				module.pending = append(module.pending, splitMessage[0])
				module.rank = append(module.rank,splitMessage[1])
				module.countPending += 1
				
				value := module.Deliver(BEB2PLink2COB(y))

				if (value > -1){
					module.pending = append(module.pending[:value], module.pending[value+1:]...)
					module.rank = append(module.rank[:value], module.rank[value+1:]...)
					module.countPending -= 1
				}
				
			}
		}
	}()
}

func (module COBroadcast_Module) COBroadcast(message COBroadcast_Req_Message) {

	w := module.v // w = v

	a := message.Addresses[0]
	aux :=len(a) - 1
	ranks := a[aux:]

	rank, err := strconv.ParseInt(ranks, 0, 64)
	_ = err

	w[rank] = strconv.Itoa(module.lsn) //w[rank]=lsn

	msg := w[0]
	for i := 1; i < module.countAddresses; i++ {
		msg += " " + w[i]
	}
	port := strings.Split(message.Addresses[0], ":")
	msg += ":" + port[1] + ":" + message.Message

	req:= BestEffortBroadcast.BestEffortBroadcast_Req_Message{
		Addresses: message.Addresses[0:],
		Message:   msg}
	module.beb.Req <- req
}

func (module COBroadcast_Module) Deliver(message COBroadcast_Ind_Message) (int) {

	//fmt.Println("Received1 '" + message.Message + "' from " + message.From)
	//fmt.Println(message.Message)

	for i := 0; i < module.countPending; i++ {
		w := strings.Split(module.pending[i], " ")
		isBigger := false
		for j := 0; j < len(w); j++ {
			wnum, err1 := strconv.ParseInt(w[j], 0, 64)
			_ = err1

			vnum, err2 := strconv.ParseInt(module.v[j], 0, 64)
			_ = err2

			if(wnum > vnum){
				isBigger = true
				break
			}
		}
		if(isBigger == false){
			a := module.rank[i]
			aux := len(a) - 1
			ranks := a[aux:]
			rank, err := strconv.ParseInt(ranks, 0, 64)
			_ = err
			
			vnum, err1 := strconv.ParseInt(w[rank], 0, 64)
			_ = err1
			
			vnum += 1

			module.v[rank] = strconv.Itoa(int(vnum))


			//message.Message = module.pending[i]

			v:= ""+module.v[0]
			for k := 1; k < len(module.v); k++ {
				v += " " + module.v[k]
			}
			message.Message = v
			
			message.From = module.rank[i]

			
			module.Ind <- message
			return i
		}
	}

	message.Message = "nothing"
	message.From = "nobody"
	module.Ind <- message
	return -1
	//fmt.Println("# End CoB Received")
}

/*func COB2BEB2PLink(message COBroadcast_Req_Message) BestEffortBroadcast.BestEffortBroadcast_Req_Message {

	return BestEffortBroadcast.BestEffortBroadcast_Req_Message{
		Addresses: addresses[0:],
		Message:   msg}
}*/

func BEB2PLink2COB(message BestEffortBroadcast.BestEffortBroadcast_Ind_Message) COBroadcast_Ind_Message {

	return COBroadcast_Ind_Message{
		From:    message.From,
		Message: message.Message}
}

