package main

import (
    "bufio"
    "context"
    "encoding/binary"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/goburrow/modbus"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter device IP: ")
    ip, _ := reader.ReadString('\n')
    ip = strings.TrimSpace(ip)

    fmt.Print("Enter port: ")
    port, _ := reader.ReadString('\n')
    port = strings.TrimSpace(port)

    fmt.Print("Enter start address: ")
    startAddressStr, _ := reader.ReadString('\n')
    startAddress, _ := strconv.Atoi(strings.TrimSpace(startAddressStr))

    fmt.Print("Enter register count: ")
    registerCountStr, _ := reader.ReadString('\n')
    registerCount, _ := strconv.Atoi(strings.TrimSpace(registerCountStr))

    //fmt.Print("Enter timeout (ms): ")
    //timeoutStr, _ := reader.ReadString('\n')
    //timeoutMs, _ := strconv.Atoi(strings.TrimSpace(timeoutStr))

    client := modbus.TCPClient(fmt.Sprintf("%s:%s", ip, port))

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                results, err := client.ReadHoldingRegisters(uint16(startAddress), uint16(registerCount*2))
                if err != nil {
                    log.Printf("Error reading registers: %v", err)
                    time.Sleep(1 * time.Second)
                    continue
                }

                // Convert raw bytes to floats and print them
                for i := 0; i < registerCount*2; i += 4 {
                    if i+4 > len(results) {
                        log.Printf("Unexpected end of data at index %d", i)
                        break
                    }
                    raw := binary.BigEndian.Uint32(results[i : i+4])
                    value := math.Float32frombits(raw)
                    fmt.Printf("Register %d: %.4f\n", startAddress+(i/2), value)
                }

                time.Sleep(1 * time.Second)
            }
        }
    }()

    fmt.Println("Type 'exit' to quit.")
    for {
        input, _ := reader.ReadString('\n')
        if strings.TrimSpace(input) == "exit" {
            cancel()
            fmt.Println("Exiting...")
            break
        }
    }
}
