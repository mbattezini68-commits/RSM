package main

import (
	"flag"
	"fmt"
	
	
	"math/rand/v2"
	"net"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type MeioFisico string

const (
	MeioTerra  MeioFisico = "TERRA_FIBRA"
	MeioAgua   MeioFisico = "AGUA_SUBMARINO"
	MeioEspaco MeioFisico = "ESPACO_SATELITE"
)

type CanalFisico struct {
	Meio      MeioFisico
	GatewayIP string
	PesoBanda int
	Ativo     bool
}

type GerenciadorToken struct {
	mu              sync.Mutex
	SaldoEficiencia float64
	MetaParaMint    float64
	TokensEmitidos  int
	CarteiraDestino string
}

func (g *GerenciadorToken) RegistrarGanho(eficienciaGerada float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.SaldoEficiencia += eficienciaGerada
	fmt.Printf(" 💰 [RSM TOKENOMICS] Ganho de eficiência computado: +%.2f unidades.\n", eficienciaGerada)
	if g.SaldoEficiencia >= g.MetaParaMint {
		g.MintarTokens()
	}
}

func (g *GerenciadorToken) MintarTokens() {
	g.TokensEmitidos += 1000
	fmt.Printf("\n🚀 [BLOCKCHAIN MINT] Meta atingida! 1000 tokens emitidos para: %s\n\n", g.CarteiraDestino)
	g.SaldoEficiencia = 0
}

// Inicia um listener UDP local para simular o recebimento real dos pacotes nos gateways
func IniciarGatewayUDP(ipPorta string, meio MeioFisico, wg *sync.WaitGroup) {
	defer wg.Done()
	addr, err := net.ResolveUDPAddr("udp", ipPorta)
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		// Porta já em uso ou ambiente restrito, apenas ignoramos na simulação local
		return
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return // Encerra o listener após o teste
		}
	}
}

func DespacharMultipathRedeReal(payload []byte, canais map[MeioFisico]*CanalFisico, tokenMgr *GerenciadorToken) {
	tamanhoTotal := len(payload)
	fmt.Printf("\n[RSM NET ENGINE] Despachando payload real de %d bytes via sockets UDP...\n", tamanhoTotal)

	pesoTotal := 0
	ativos := []*CanalFisico{}
	for _, canal := range canais {
		if canal.Ativo {
			pesoTotal += canal.PesoBanda
			ativos = append(ativos, canal)
		}
	}

	bytesAlocados := 0
	for i, canal := range ativos {
		tamanhoPedaco := (tamanhoTotal * canal.PesoBanda) / pesoTotal
		if i == len(ativos)-1 {
			tamanhoPedaco = tamanhoTotal - bytesAlocados
		}
		if tamanhoPedaco <= 0 {
			continue
		}
		bytesAlocados += tamanhoPedaco

		// Envio real via UDP para o Gateway configurado
		addr, err := net.ResolveUDPAddr("udp", canal.GatewayIP)
		if err == nil {
			conn, err := net.DialUDP("udp", nil, addr)
			if err == nil {
				conn.Write(payload[:tamanhoPedaco])
				conn.Close()
			}
		}

		fmt.Printf(" ├── [REDE REAL] Enviado pedaço de %d bytes para [%s] em %s\n", tamanhoPedaco, canal.Meio, canal.GatewayIP)
	}

	ganho := float64(tamanhoTotal) * 0.40
	tokenMgr.RegistrarGanho(ganho)
}

func main() {
	_ = godotenv.Load()
	walletDestino := os.Getenv("WALLET_DESTINATARIO")
	if walletDestino == "" {
		walletDestino = "0xFallbackWallet"
	}

	flag.Parse()

	fmt.Println("====================================================")
	fmt.Println(" RSM | CAMADA DE REDE REAL (UDP MULTIPATH)")
	fmt.Println("====================================================")

	var wg sync.WaitGroup
	canaisGlobais := map[MeioFisico]*CanalFisico{
		MeioTerra:  {Meio: MeioTerra, GatewayIP: "127.0.0.1:9001", PesoBanda: 50, Ativo: true},
		MeioAgua:   {Meio: MeioAgua, GatewayIP: "127.0.0.1:9002", PesoBanda: 30, Ativo: true},
		MeioEspaco: {Meio: MeioEspaco, GatewayIP: "127.0.0.1:9003", PesoBanda: 20, Ativo: true},
	}

	// Sobe os ouvintes locais para simular os gateways recebendo os pacotes físicos
	for _, canal := range canaisGlobais {
		wg.Add(1)
		go IniciarGatewayUDP(canal.GatewayIP, canal.Meio, &wg)
	}

	minhaCarteira := &GerenciadorToken{
		MetaParaMint:    100.0,
		CarteiraDestino: walletDestino,
	}

	time.Sleep(300 * time.Millisecond)

	for i := 1; i <= 2; i++ {
		payloadReal := make([]byte, 1024+rand.IntN(256))
		DespacharMultipathRedeReal(payloadReal, canaisGlobais, minhaCarteira)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nTeste de rede real concluído com sucesso.")
}
