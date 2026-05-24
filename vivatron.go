package main

import (
        "encoding/binary"
        "fmt"
        "net"
        "os"
        "strconv"
        "time"
)

func calculateChecksum(data []byte) byte {
        var sum int
        for _, b := range data {
                sum += int(b)
        }
        // In Go, conversion to byte automatically truncates higher bits,
        // but we use & 0xFF to preserve the logic from Python
        return byte(sum & 0xFF)
}

func sendCommand(ip string, payload []byte) {
        // Header and example sequence number (0x01)
        packet := []byte{0xb0, 0xb1, 0xb2, 0xb3, 0x00, 0x01, 0x02, 0x01}

        // Payload length (2 bytes, big-endian)
        lengthBytes := make([]byte, 2)
        binary.BigEndian.PutUint16(lengthBytes, uint16(len(payload)))
        packet = append(packet, lengthBytes...)

        // Actual command
        packet = append(packet, payload...)

        // Checksum
        packet = append(packet, calculateChecksum(packet))

        // TCP connection configuration with 3-second timeout
        address := fmt.Sprintf("%s:5577", ip)
        conn, err := net.DialTimeout("tcp", address, 3*time.Second)
        if err != nil {
                fmt.Printf("Connection error: %v\n", err)
                return
        }
        defer conn.Close()

        // Set timeout for write operation as well
        err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
        if err != nil {
                fmt.Printf("Error setting write timeout: %v\n", err)
                return
        }

        // Send packet
        _, err = conn.Write(packet)
        if err != nil {
                fmt.Printf("Send error: %v\n", err)
        }
}

func main() {
        if len(os.Args) < 3 {
                fmt.Println("Usage: program <ip> <on|off|color> [r g b w]")
                os.Exit(1)
        }

        ip := os.Args[1]
        action := os.Args[2]

        switch action {
        case "on":
                sendCommand(ip, []byte{0x71, 0x23})
        case "off":
                sendCommand(ip, []byte{0x71, 0x24})
        case "color":
                if len(os.Args) == 7 {
                        r, errR := strconv.Atoi(os.Args[3])
                        g, errG := strconv.Atoi(os.Args[4])
                        b, errB := strconv.Atoi(os.Args[5])
                        w, errW := strconv.Atoi(os.Args[6])

                        // Simple argument format error handling
                        if errR != nil || errG != nil || errB != nil || errW != nil {
                                fmt.Println("Error: Color values must be integers.")
                                os.Exit(1)
                        }

                        sendCommand(ip, []byte{0xeb, 0x01, 0x00, 0x04, byte(r), byte(g), byte(b), byte(w)})
                } else {
                        fmt.Println("Error: 'color' action requires 4 values: r, g, b, w")
                        os.Exit(1)
                }
        default:
                fmt.Println("Error: Unknown action. Available actions are: on, off, color")
                os.Exit(1)
        }
}
