package multicast

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	f "practice1/functions"
	"time"
)

// SendGroupM send message to ip multicast and wait ack
func SendGroupM(chanAck chan f.Ack, connect *f.Conn) error {
	var red *net.UDPAddr
	var connection *net.UDPConn
	var encoder *gob.Encoder
	var buffer bytes.Buffer
	var msm f.Message
	var ok bool
	var err error
	var bufferAck []f.Ack

	// target := ""
	// delay, _ := time.ParseDuration("0s")
	// inf := "Me mataron"
	id := connect.GetId()

	// Update vClock and make a copy
	vector := connect.GetVector()
	vector.Tick(id)
	connect.SetClock(vector)
	copyVector := vector.Copy()

	// Check if it has a target
	// log.Println("[SendGroupM] Check if it has a target &&&&&&&&&&&&&&&&&&&")
	if len(connect.GetKill()) > 0 && len(connect.GetDelays()) > 0 {
		// log.Println("[SendGroupM]pppppppppppppppppppppppppppppppppppppppppppp")
		// target = connect.GetTarget(0)
		// delay = connect.GetDelay(0)
		// inf = "He disparado"
		// connect.SetKill()
		// connect.SetDelay()
		msm = f.Message{
			To:     f.MulticastAddress,
			From:   id,
			Targ:   connect.GetTarget(0),
			Data:   "kill",
			Vector: copyVector,
			Delay:  connect.GetDelay(0),
		}
		// log.Println("[SendGroupM] Creo el  MSM en if", msm)
	} else {
		delay, _ := time.ParseDuration("0s")
		// Created message to send
		msm = f.Message{
			To:     f.MulticastAddress,
			From:   id,
			Targ:   "",
			Data:   "am dead",
			Vector: copyVector,
			Delay:  delay,
		}
		// log.Println("[SendGroupM] Creo el  MSM en else", msm)

	}

	// Creating red connection
	red, err = net.ResolveUDPAddr("udp", f.MulticastAddress)
	f.Error(err, "SendGroupM error ResolveUDPAddr connection \n")

	connection, err = net.DialUDP("udp", nil, red)
	f.Error(err, "SendGroupM error DialUDP connection \n")
	defer connection.Close()

	// Send msm to ip multicast 3 times
	go func() {
		log.Println("[SendGroupM] FOR Send msm to multicast three times ")
		for i := 0; i < 3; i++ {
			log.Println("[SendGroupM] ENVIO Numero ", i, " al ip ")
			encoder = gob.NewEncoder(&buffer)
			err = encoder.Encode(&msm)
			f.Error(err, "SendGroupM encoder error \n")
			_, err = connection.Write(buffer.Bytes())
			// f.Error(err, "Error al enviar el msm")
			time.Sleep(200 * time.Millisecond)
		}
		log.Println("[SendGroupM] ENVIO ", msm)
	}()

	log.Println("[SendGroupM] 77 VOY A RECIBIR ACK")
	// dictAck := make(map[string]f.Ack)
	deadline := time.Now().Add(2 * time.Second)
	// log.Println("[SendGroupM] 79  FOR Time  0000000000000 ", deadline)
	for time.Now().Before(deadline) {
		log.Println("[SendGroupM] 81 ")
		pack, _ := <-chanAck

		log.Println("[SendGroupM] 82 me llego ACK ")
		if connect.GetId() != pack.GetOrigen() {
			// dictAck[pack.GetOrigen()] = pack
			bufferAck, ok = f.AddAcks(bufferAck, pack)
		}

		if len(bufferAck) == len(connect.GetIds())-1 {
			break
		}

	}

	// log.Println("[SendGroupM] CHEQUEOS LOS ACKS ")
	pendCheck, chec := f.CheckAcks(bufferAck, connect)

	// TODO Call Receive
	// log.Println("[SendGroupM] IMPRIMO A VER SI FALTAN ACK ", chec, " Y LOS ACKS ", pendCheck)
	if !chec {
		//Necesito enviar tres veces
		// log.Println("[SendGroupM] me faltan ACK los envio rirectamente")
		go func() {
			for i := 0; i < 3; i++ {
				for _, v := range pendCheck {
					log.Println("[SendGroupM] envio a ", v)
					go SendM(msm, v)
				}
				time.Sleep(200 * time.Millisecond)
			}
		}()

		deadline2 := time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline2) {
			// for i := 0; i < nless; i++ {
			// log.Println("[SendGroupM] FOR RECEIVE ACK for 3 seconds ")
			pack := <-chanAck
			if connect.GetId() != pack.GetOrigen() {
				// dictAck[pack.GetOrigen()] = pack
				bufferAck, ok = f.AddAcks(bufferAck, pack)
			}
		}
	}

	log.Println("[SendGroupM] |||||| Fin send Group |||| ", ok)

	return err

}