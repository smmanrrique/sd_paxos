package comunicacionCausal

import (
	"fmt"
	"net"
	f "practice1/functions"
)

func ReceiveGroup(connect *f.Conn) error {
	var err error
	var listener net.Listener
	// var wg sync.WaitGroup
	var n = len(connect.GetIds())
	id := connect.GetId()
	// bufferMsm := make(chan Message, n)
	bufferMsm := make(chan f.Message)
	defer close(bufferMsm)

	listener, err = net.Listen("tcp", connect.GetPort())
	f.Error(err, "Listen Error")
	defer listener.Close()
	fmt.Printf("EL VALOR N:       %d \n", n)
	for i := 0; i < n; i++ {
		// wg.Add(1)
		fmt.Printf("FOR:       %d \n", i)
		// go Receive(bufferMsm, listener, &wg)
		go Receive(bufferMsm, listener, id)

		msm, _ := <-bufferMsm
		// select {
		// case msm, _ := <-bufferMsm:
		fmt.Println("targets en receive group:  ", connect.GetTarget(1))
		if id != msm.To {
			if len(connect.GetKill()) > 0 {
				fmt.Println("Contenido del mensaje recibido:", msm)
				go SendGroup(connect)
			}
			// case <-time.After(5 * time.Second):
			fmt.Println("TIME OUT receive CIERRO CONNECTION EN RECEIVE ")
			// 	close(bufferMsm)
			// 	break
		}

	}

	// wg.Wait()

	fmt.Println("--------------------------")
	return err

}
