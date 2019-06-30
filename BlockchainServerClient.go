package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// Block representa la unidad de almacenamiento de la cadena
type Block struct {
	Index          int
	Timestamp      string
	MedicalHistory MedicalHistory
	Hash           string
	PrevHash       string
}

//MedicalHistory Estructura de almacenamiento de una historia clínica
type MedicalHistory struct {
	History string
}

// Blockchain Estructura de tipo array de bloques que conforman la cadena
var Blockchain []Block

func calculateHash(block Block) string {
	data := string(block.Index) + block.Timestamp + string(block.MedicalHistory.History) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, History MedicalHistory) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.MedicalHistory = History
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock Block) bool {
	//Se realizan 3 validaciones diferentes
	//Primero se analiza si es que el indice no coincide lo que significa que el bloque que se trata de agregar no pertenece a una cadena de menor o mayor tamaño
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	//Segundo, se analiza que el hash corresponda con el hash anterior
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	//Finalmente, se compara el hash del bloque con el hash obtenido al calcular el hash del lado del servidor
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	//No existe un protocolo de consenso o de votación, simplemente se valida si la cadena es de diferente tamaño para simular dicho proceso
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
func handleConn(conn net.Conn) {
	//El servidor espera el input del usuario para setear el historial médico
	defer conn.Close()
	io.WriteString(conn, "Enter a new Medical History:")

	scanner := bufio.NewScanner(conn)
	//Creamos la estructura de tipo historia clínica y llamamos a la función generar bloque
	go func() {
		for scanner.Scan() {
			history := scanner.Text()
			medicalhistory := MedicalHistory{
				History: history,
			}
			//La generación del bloque requiere de la blockchain para realizar las validaciones correspondientes
			newBlock, err := generateBlock(Blockchain[len(Blockchain)-1], medicalhistory)
			if err != nil {
				log.Println(err)
				continue
			}
			//Se valida el bloque y la nueva cadena
			if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
				newBlockchain := append(Blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			//Se intenta escribir en el canal del servidor encargado de administrar la admisión de bloques
			bcServer <- Blockchain
			io.WriteString(conn, "\nEnter a new Medical History:")
		}
	}()
	//Simulación del protocolo de consenso de la blockchain
	go func() {
		for {
			//Se espera 30 segundos y se emite un broadcast con la nueva cadena válida
			time.Sleep(30 * time.Second)
			output, err := json.Marshal(Blockchain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(output))
		}
	}()
	//Mientras el canal este abierto se imprime el valor de los bloques agregados
	for _ = range bcServer {
		spew.Dump(Blockchain)
	}
}

//bcServer canal para impedir que todos los clientes agreguen bloques al mismo tiempo
var bcServer chan []Block

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bcServer = make(chan []Block)

	//El bloque Genesis siempre esta vacío
	t := time.Now()
	medicalhistory := MedicalHistory{
		History: "",
	}
	genesisBlock := Block{0, t.String(), medicalhistory, "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	server, err := net.Listen("tcp", ":"+os.Getenv("ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}
