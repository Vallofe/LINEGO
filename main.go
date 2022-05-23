package main

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"
	"fmt"
	"time"

	"./Auth"
	"./LINE"
	"golang.org/x/xerrors"
)

var client *Auth.LINEClient

var (
	AppName string = ""
	Token   string = ""
	ctx = context.Background()
)

func Speed(op *LINE.Operation) error {
	message := op.Message
	to := message.To
	if message.ToType == LINE.MIDType_USER { to = message.Get_from()}
	msg := new(LINE.Message)
	msg.To = to
	msg.Text = "SpeedTest..."
	start := time.Now()
	_, err := client.Talk.SendMessage(ctx, 0, msg)
	end := time.Now()
	if err != nil {
		return err
	}
	msg.Text = fmt.Sprintf("%d ms", end.Sub(start).Milliseconds())
	_, err = client.Talk.SendMessage(ctx, 0, msg)
	if err != nil {
		return err
	}
	return nil
}


func OpType_RECEIVE_MESSAGE(op *LINE.Operation) error {
	message := op.Message
	switch message.Text {
	case "speed":
		err := Speed(op)
		if err != nil {
			return err
		}
	}
	return nil
}

func Main_Loop() {
	const (
		count = 100
		sep   = "\x1e"
	)
	var (
		localRev      int64
		globalRev     int64
		individualRev int64
	)

	for {
		operations, err := client.Poll.FetchOps(ctx, localRev, count, globalRev, individualRev)
		if err != nil {
			if strings.Contains(err.Error(), "server sent GOAWAY and closed the connection") {
				continue
			}
			log.Fatalf("failed to call fetchOps: %+v\n", err)
		}

		for _, op := range operations {
			switch op.Type {
			case LINE.OpType_END_OF_OPERATION:
				if op.Param1 != "" {
					individualRevString, err := strconv.Atoi(strings.Split(op.Param1, sep)[0])
					if err != nil {
						log.Fatalf("failed to get individualRev: %+v\n", err)
					}
					individualRev = int64(individualRevString)
				}

				if op.Param2 != "" {
					globalRevString, err := strconv.Atoi(strings.Split(op.Param2, sep)[0])
					if err != nil {
						log.Fatalf("failed to get individualRev: %+v\n", err)
					}
					globalRev = int64(globalRevString)
				}

				continue

			case LINE.OpType_RECEIVE_MESSAGE:
				err := OpType_RECEIVE_MESSAGE(op)
				if err != nil {
					log.Printf("%+v\n", err)
				}
			}

			localRev = int64(math.Max(float64(localRev), float64(op.Revision)))
		}
	}
}

func Generated(accessToken string) (*Auth.LINEClient, error) {
	cfg := Auth.Config {
		Host:                      "https://legy-jp.line.naver.jp",
		TalkService:               "/S4",
		PollService: 			   "/P4",
		UserAgent:                 "LLA/2.14.0 F-01H 6.0.1",
		LINEApp:                   AppName,
		AccessToken:               accessToken,
	}
	client, err := Auth.NewLINEClient(cfg)
	fmt.Println(client)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate LINE client: %w", err)
	}
	return client, nil
}

func main() {
	exec, err := Generated(Token)
	if err != nil { log.Fatalf("failed to generate LINE client: %+v\n", err) }
	client = exec
	Main_Loop()
}