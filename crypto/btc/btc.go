package btc

import (
	"net"
	"fmt"
	"bytes"
	"encoding/gob"
	"bufio"
	"io"

	uErr "xox_tokens/errors"
//	dbTypes "xox_tokens/db/types"
	"xox_tokens/types"
	"xox_tokens/constants/messageTypes"
_	"xox_tokens/db"
_	"xox_tokens/db/models/transaction"
_	"xox_tokens/db/models/notHandledTransaction"
_	"xox_tokens/db/models/buyer"
_	"xox_tokens/crypto"
_	"xox_tokens/crypto/eth"
_	"log"
_	"time"
)

// Dial starts connection with BTCService via tcp.
func Dial(port uint, walletAddress string) (err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil { return }

	var message bytes.Buffer
	enc := gob.NewEncoder(&message)
	enc.Encode(types.BTCServiceReq{
		Type: messageTypes.WATCH_ADDRESS,
		Address: walletAddress,
	})

	response := bufio.NewReader(conn)
	_, err = conn.Write(append(message.Bytes(), '\n'))
	if err != nil { return }

	for {
		_, err := response.ReadBytes(byte('\n'))
		switch err {
		case nil:
			//go handleMessage(line)
			fmt.Println("handleMessage") // TODO: Delete this line
		case io.EOF:
			return uErr.Combine(nil, uErr.ErrConnectBTCService)
		default:
			return err
		}
	}
}
