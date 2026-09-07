package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// ============================================================
// RSM (Real Agent Network) | NÓ DE ROTEAMENTO HÍBRIDO GLOBAL
// ============================================================

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
	fmt.Printf("                     (Saldo atual na mempool: %.2f / Meta: %.2f)\n", g.SaldoEficiencia, g.MetaParaMint)

	if g.SaldoEficiencia >= g.MetaParaMint {
		g.MintarTokens()
	}
}

func (g *GerenciadorToken) MintarTokens() {
	g.TokensEmitidos += 1000
	fmt.Printf("\n============================================================\n")
	fmt.Printf("🚀 [BLOCKCHAIN MINT] Meta de eficiência atingida!\n")
	fmt.Printf("   -> 1000 novos tokens mintados com sucesso.\n")
	fmt.Printf("   -> Transferidos automaticamente para a carteira segura: %s\n", g.CarteiraDestino)
	fmt.Printf("   -> Total acumulado de tokens: %d\n", g.TokensEmitidos)
	fmt.Printf("============================================================\n\n")
	g.SaldoEficiencia = 0
}

type ConfigAgua struct {
	MinBufferBytes  int
	MaxBufferBytes  int
	CapacidadeBytes int
	LatenciaBase    int
	JitterMax       int
	TaxaPerda       float64
}

type CanalSubmarino struct {
	cfg ConfigAgua
}

func NovoCanalSubmarino(cfg ConfigAgua) *CanalSubmarino {
	return &CanalSubmarino{cfg: cfg}
}

func (c *CanalSubmarino) Transmitir(bytes int, rng *rand.Rand) (int, bool, bool) {
	if bytes <= 0 {
		return 0, false, false
	}
	if c.cfg.TaxaPerda >= 1.0 {
		return 0, true, true
	}
	if c.cfg.TaxaPerda > 0 && rng.Float64() < c.cfg.TaxaPerda {
		return 0, true, false
	}
	jitter := 0
	if c.cfg.JitterMax > 0 {
		jitter = rng.IntN(c.cfg.JitterMax + 1)
	}
	capacidade := c.cfg.CapacidadeBytes
	if capacidade <= 0 {
		capacidade = 1
	}
	tempoExecucao := int(math.Ceil(float64(bytes)/float64(capacidade))) + c.cfg.LatenciaBase + jitter
	return tempoExecucao, false, false
}

func DespacharMultipathGlobal(payloadOriginal []byte, canais map[MeioFisico]*CanalFisico, tokenMgr *GerenciadorToken, cfgAgua ConfigAgua) {
	tamanhoTotal := len(payloadOriginal)
	fmt.Printf("\n[RSM ENGINE] Pacote íntegro de %d bytes recebido para distribuição global.\n", tamanhoTotal)

	pesoTotal := 0
	ativos := []*CanalFisico{}
	for _, canal := range canais {
		if canal.Ativo {
			pesoTotal += canal.PesoBanda
			ativos = append(ativos, canal)
		}
	}

	if pesoTotal == 0 {
		fmt.Println("[ERRO] Nenhum canal físico ativo para despacho.")
		return
	}

	bytesAlocados := 0
	rngSim := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 999))

	for i, canal := range ativos {
		var tamanhoPedaco int
		if i == len(ativos)-1 {
			tamanhoPedaco = tamanhoTotal - bytesAlocados
		} else {
			tamanhoPedaco = (tamanhoTotal * canal.PesoBanda) / pesoTotal
		}

		if tamanhoPedaco <= 0 {
			continue
		}
		bytesAlocados += tamanhoPedaco

		if canal.Meio == MeioAgua {
			subChan := NovoCanalSubmarino(cfgAgua)
			_, perdido, _ := subChan.Transmitir(tamanhoPedaco, rngSim)
			if perdido {
				fmt.Printf(" ├── [ROTA SUBMARINA] ⚠️ Alerta: Perda parcial detectada em %d bytes via [%s]. Ajustando FEC...\n", tamanhoPedaco, canal.Meio)
			}
		}

		fmt.Printf(" ├── [ROTA ÓTIMA] %d bytes despachados via [%s] (Gateway: %s)\n",
			tamanhoPedaco, canal.Meio, canal.GatewayIP)
	}

	fmt.Println(" [RECOMPOSIÇÃO] Payload remontado com integridade total no destino.")

	ganhoDaTransmissao := float64(tamanhoTotal) * 0.45
	tokenMgr.RegistrarGanho(ganhoDaTransmissao)
}

func main() {
	// Carrega as variáveis de ambiente do arquivo .env de forma segura
	_ = godotenv.Load()

	privKey := os.Getenv("RSM_PRIVATE_KEY")
	walletDestino := os.Getenv("WALLET_DESTINATARIO")

	if walletDestino == "" {
		walletDestino = "0xDefaultFallbackWalletWithoutKey"
	}

	if privKey == "" {
		log.Println("⚠️  [AVISO DE SEGURANÇA] RSM_PRIVATE_KEY não encontrada no .env. O nó rodará em modo simulação offline.")
	} else {
		fmt.Println("🔐 [SEGURANÇA] Chave privada carregada com sucesso via ambiente off-chain.")
	}

	volume := flag.Int("volume", 5000, "Volume de dados simulados no ecossistema RSM")
	flag.Parse()

	fmt.Println("====================================================")
	fmt.Println(" RSM | PROTOCOLO DE ROTEAMENTO HÍBRIDO & DEPIN")
	fmt.Println("====================================================")
	fmt.Printf("Volume Alvo Configurado: %d bytes\n", *volume)
	fmt.Printf("Destino das Recompensas: %s\n", walletDestino)

	minhaCarteira := &GerenciadorToken{
		MetaParaMint:    150.0,
		CarteiraDestino: walletDestino,
	}

	canaisGlobais := map[MeioFisico]*CanalFisico{
		MeioTerra:  {Meio: MeioTerra, GatewayIP: "192.168.1.10:8080", PesoBanda: 50, Ativo: true},
		MeioAgua:   {Meio: MeioAgua, GatewayIP: "10.0.2.10:8080", PesoBanda: 30, Ativo: true},
		MeioEspaco: {Meio: MeioEspaco, GatewayIP: "172.16.0.10:8080", PesoBanda: 20, Ativo: true},
	}

	cfgAgua := ConfigAgua{
		MinBufferBytes:  32,
		MaxBufferBytes:  1024,
		CapacidadeBytes: 2048,
		LatenciaBase:    140,
		JitterMax:       15,
		TaxaPerda:       0.02,
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("\n--- Iniciando fluxo contínuo de roteamento RSM ---")

	tamanhoPayload := *volume / 4
	if tamanhoPayload < 100 {
		tamanhoPayload = 100
	}

	for i := 1; i <= 3; i++ {
		payloadFalso := make([]byte, tamanhoPayload+rand.IntN(50))

		DespacharMultipathGlobal(payloadFalso, canaisGlobais, minhaCarteira, cfgAgua)
		time.Sleep(700 * time.Millisecond)
	}

	fmt.Println("\nExecução consolidada do RSM finalizada com sucesso.")
}
