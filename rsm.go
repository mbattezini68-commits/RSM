package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	PolygonChainID = 137
	CreatorFeePct  = 0.03 // 3% perpétuo para o criador
)

type RSMEngine struct {
	Client          *ethclient.Client
	PrivateKey      *ecdsa.PrivateKey
	CreatorAddress  common.Address
	ContractAddress common.Address
}

func NewRSMEngine() (*RSMEngine, error) {
	rpcURL := os.Getenv("RSM_RPC_URL")
	if rpcURL == "" {
		rpcURL = "https://polygon-bor-rpc.publicnode.com"
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao nó RPC da Polygon: %v", err)
	}

	privKeyHex := os.Getenv("RSM_PRIVATE_KEY")
	if privKeyHex == "" {
		return nil, fmt.Errorf("RSM_PRIVATE_KEY não configurada no ambiente")
	}

	privateKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("chave privada inválida: %v", err)
	}

	creatorAddr := common.HexToAddress(os.Getenv("RSM_CREATOR_ADDRESS"))
	contractAddr := common.HexToAddress(os.Getenv("RSM_CONTRACT_ADDRESS"))

	return &RSMEngine{
		Client:          client,
		PrivateKey:      privateKey,
		CreatorAddress:  creatorAddr,
		ContractAddress: contractAddr,
	}, nil
}

func (engine *RSMEngine) ExecutarCicloAutonomo(ciclo int, ganhoBruto float64) {
	fmt.Printf("\n========================================================\n")
	fmt.Printf("🔄 [RSM DAEMON] Ciclo Autônomo #%d - Otimização de Gargalos\n", ciclo)
	fmt.Printf("========================================================\n")

	if ganhoBruto <= 0 {
		fmt.Println("💤 [POTRO EM REPOUSO] Nenhum gargalo monetizado neste batimento. Gás preservado.")
		return
	}

	parteCriador := ganhoBruto * CreatorFeePct
	liquidoRestante := ganhoBruto - parteCriador

	premioUsuario := liquidoRestante * 0.25
	fundoAprimoramento := liquidoRestante * 0.25
	
	// A regra ajustada: os 3% do criador saem da cota de queima (50% - 3% = 47% efetivos do total bruto)
	tokensParaQueima := (ganhoBruto * 0.50) - parteCriador

	fmt.Printf("💰 [EFICIÊNCIA CAPTURADA!] Ganho Bruto: $%.4f\n", ganhoBruto)
	fmt.Printf("   -> 🛡️ Royalty do Criador (3%%): $%.4f -> Direto para: %s\n", parteCriador, engine.CreatorAddress.Hex())
	fmt.Printf("   -> 🎁 Desconto/Prêmio ao Usuário (25%%): $%.4f\n", premioUsuario)
	fmt.Printf("   -> ⚙️ Aprimoramento do Sistema (25%%): $%.4f\n", fundoAprimoramento)
	fmt.Printf("   -> 🔥 Queima Deflacionária Ajustada (47%%): $%.4f\n", tokensParaQueima)

	engine.TransmitirExecucaoOnChain(ganhoBruto, tokensParaQueima)
}

func (engine *RSMEngine) TransmitirExecucaoOnChain(ganho float64, queima float64) {
	if engine.PrivateKey == nil {
		fmt.Println("❌ [ERRO CRÍTICO] Chave privada não inicializada.")
		return
	}

	ctx := context.Background()
	publicKey := engine.PrivateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := engine.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		fmt.Printf("⚠️ [AVISO RPC] Nonce indisponível: %v\n", err)
		return
	}

	value := big.NewInt(0)
	gasLimit := uint64(65000)
	gasPrice, err := engine.Client.SuggestGasPrice(ctx)
	if err != nil {
		gasPrice = big.NewInt(30000000000)
	}

	data := []byte(fmt.Sprintf("RSM_EXEC_SUCCESS:GAIN=%.2f:BURN=%.2f", ganho, queima))
	tx := types.NewTransaction(nonce, engine.ContractAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(PolygonChainID)), engine.PrivateKey)
	if err != nil {
		fmt.Printf("❌ [ERRO] Falha ao assinar transação: %v\n", err)
		return
	}

	err = engine.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		fmt.Printf("⚠️ [TRANSMISSÃO POLYGON] Falha ao enviar: %v\n", err)
		return
	}

	fmt.Printf("🚀 [SUCESSO NA POLYGON] Transação registrada na Mainnet!\n")
	fmt.Printf("🔗 TxHash: %s\n", signedTx.Hash().Hex())
}

func main() {
	fmt.Println("========================================================")
	fmt.Println("🤖 RSM AGENT DAEMON - OTIMIZADOR DE GARGALOS (24/7)")
	fmt.Println("========================================================")

	engine, err := NewRSMEngine()
	if err != nil {
		fmt.Printf("❌ [FALHA DE INICIALIZAÇÃO]: %v\n", err)
		return
	}

	defer engine.Client.Close()
	fmt.Println("✅ [CONEXÃO ESTABELECIDA] Pronto para capturar valor na Polygon...")

	ciclo := 0
	for {
		ciclo++
		ganhoRealDoServico := 0.0 // Modo passivo aguardando otimizações reais
		engine.ExecutarCicloAutonomo(ciclo, ganhoRealDoServico)
		time.Sleep(30 * time.Second)
	}
}
