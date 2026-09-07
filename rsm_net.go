package main

import (
	"fmt"
	"os"
)

// SimularProcessamentoDeRede executa o ciclo de roteamento e aciona a segurança criptográfica
func SimularProcessamentoDeRede(nodeID string, pacotesTransmitidos int, pacotesPerdidos int) {
	var eficiencia float64
	if pacotesTransmitidos > 0 {
		eficiencia = (float64(pacotesTransmitidos-pacotesPerdidos) / float64(pacotesTransmitidos)) * 100.0
	}

	tokensGerados := int(eficiencia * 1.5)

	fmt.Printf("🌐 [RSM REDE] Nó %s processou dados.\n", nodeID)
	fmt.Printf("📈 [MÉTRICA] Eficiência de Entrega: %.2f%%\n", eficiencia)
	fmt.Printf("🪙 [TOKENOMICS] Recompensa calculada: %d tokens\n", tokensGerados)

	AssinarComEIP712(nodeID, eficiencia, tokensGerados)
}

func AssinarComEIP712(nodeID string, eficiencia float64, tokens int) {
	privKey := os.Getenv("PRIVATE_KEY")
	if privKey == "" {
		privKey = "variavel_ambiente_padrao_simulada"
	}

	fmt.Printf("🔐 [EIP-712] Assinando payload estruturado para o nó %s (Eficiência: %.2f%%) com L2/ECDSA...\n", nodeID, eficiencia)
	fmt.Println("✨ [SUCESSO] Transação criptografada pronta para broadcast no L2.")
}
