package forward

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net"
	"sync"
)

type Forward struct {
	udpConn    sync.Map
	dstAddr    string
	packetConn net.PacketConn
	aesStream  cipher.Stream
}

func Start(listenAddr, toAddr, aesKey, aesIv string) {
	key, _ := hex.DecodeString(aesKey)
	iv, _ := hex.DecodeString(aesIv)

	block1, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	pc, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pc.Close()

	stream := cipher.NewCTR(block1, iv[:aes.BlockSize])

	forward := &Forward{
		dstAddr:    toAddr,
		packetConn: pc,
		aesStream:  stream,
	}

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := pc.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn, err := forward.connectUdp(clientAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		forward.aesStream.XORKeyStream(buffer[0:n], buffer[0:n])
		conn.Write(buffer[0:n])
	}
}

func (f *Forward) connectUdp(clientAddr net.Addr) (net.Conn, error) {
	v, ok := f.udpConn.Load(clientAddr.String())
	if !ok {
		conn, err := net.Dial("udp", f.dstAddr)
		if err != nil {
			return nil, err
		}
		f.udpConn.Store(clientAddr, conn)
		go f.forwardUdp(clientAddr, conn)
		return conn, nil
	}
	conn := v.(net.Conn)
	return conn, nil
}

func (f *Forward) forwardUdp(clientAddr net.Addr, udpConn net.Conn) {
	defer f.udpConn.Delete(clientAddr.String())
	buffer := make([]byte, 1024)
	for {
		n, err := udpConn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.aesStream.XORKeyStream(buffer[0:n], buffer[0:n])
		_, err = f.packetConn.WriteTo(buffer[0:n], clientAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
