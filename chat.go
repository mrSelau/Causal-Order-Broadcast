// Construido como parte da disciplina de Sistemas Distribuidos
// PUCRS - Escola Politecnica
// Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)

/*
teste
LANCAR N PROCESSOS EM SHELL's DIFERENTES, PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...
go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...
go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001
*/

package main

import "fmt"
import "os"
import "bufio"
//import "strings"

import COBroadcast "COBroadcast"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...")
		fmt.Println("go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...")
		fmt.Println("go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	var registro []string
	addresses := os.Args[1:]
	fmt.Println(addresses[0:])

	cob := COBroadcast.COBroadcast_Module{
		Req: make(chan COBroadcast.COBroadcast_Req_Message),
		Ind: make(chan COBroadcast.COBroadcast_Ind_Message)}

	cob.Init(addresses)

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
			}
			req := COBroadcast.COBroadcast_Req_Message{
				Addresses: addresses[0:],
				Message:   msg}
			cob.Req <- req
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-cob.Ind
			//message := in.Message
			//in.From = message[1]
			registro = append(registro, in.Message)
			//in.Message = message[0]

			// imprime a mensagem recebida na tela
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)
			fmt.Println("--------------------------------")
		}
	}()

	blq := make(chan int)
	<-blq
}