package main

import (
	"crypto/sha256"
	"fmt"
	"math"
	"time"
)

// RSMPacket representa um bloco fatiado do payload original com integridade criptográfica
type RSMPacket struct {
	SessionID   uint64
	ChunkIndex  int
	TotalChunks int
	Payload     []byte
	Checksum    [32]byte
}

// PathMetrics armazena a telemetria composta e o histórico de tendência de cada meio
type PathMetrics struct {
	Name          string
	Latency       time.Duration
	PacketLoss    float64 // Percentual atual (0.0 a 100.0)
	PreviousLoss  float64 // Histórico do ciclo anterior para calcular tendência
	PreviousLatMs float64 // Histórico de latência anterior
	Active        bool
}

// CalcularScorePreditivo avalia o mérito combinando estado atual + tendência de degradação (Preditivo)
func CalcularScorePreditivo(m PathMetrics) (float64, string) {
	if !m.Active {
		return 0.0, "Inativa"
	}
	
	latMs := float64(m.Latency.Milliseconds())
	if latMs <= 0 {
		latMs = 1.0
	}

	// 1. Análise Preditiva de Tendência (Derivada de Piora)
	trendLat := latMs - m.PreviousLatMs
	trendLoss := m.PacketLoss - m.PreviousLoss

	pFail := 0.0
	statusTendenciaEstavel := "Estável 🟢"

	if trendLat > 5.0 || trendLoss > 0.1 {
		pFail = 25.0 // Penalidade preventiva se houver tendência de piora
		statusTendenciaEstavel = "Degradação Preditiva Detectada ⚠️"
	} else if trendLat < 0 {
		statusTendenciaEstavel = "Otimizando/Melhorando 🚀"
	}

	// 2. Função de Pontuação Composta + Fator Preditivo Preservado
	score := 1000.0 / (latMs + (m.PacketLoss * 50.0) + pFail)
	return math.Max(score, 0.0), statusTendenciaEstavel
}

// ExecutarTransporteInteligente executa o roteamento adaptativo e preditivo
func ExecutarTransporteInteligente(rawPayload []byte, sessionID uint64) string {
	rotas := []PathMetrics{
		{Name: "TERRA_FIBRA", Latency: 15 * time.Millisecond, PacketLoss: 0.1, PreviousLoss: 0.1, PreviousLatMs: 14, Active: true},
		{Name: "AGUA_SUBMARINO", Latency: 42 * time.Millisecond, PacketLoss: 0.4, PreviousLoss: 0.6, PreviousLatMs: 48, Active: true},
		{Name: "ESPACO_SATELITE", Latency: 260 * time.Millisecond, PacketLoss: 2.5, PreviousLoss: 1.2, PreviousLatMs: 230, Active: true},
	}

	fmt.Println("\n🧠 [TRANSPORTE PREDITIVO] Analisando Telemetria e Tendências de Rede:")
	
	var melhorRota string
	maiorScore := -1.0

	for _, r := range rotas {
		score, status := CalcularScorePreditivo(r)
		fmt.Printf("  -> Rota: %-15s | Latência: %-5v | Perda: %5.1f%% | Status: %-32s | Score: %.2f\n", 
			r.Name, r.Latency, r.PacketLoss, status, score)

		if score > maiorScore {
			maiorScore = score
			melhorRota = r.Name
		}
	}

	fmt.Printf("🏆 [DECISÃO PREDITIVA] Rota blindada selecionada: %s\n", melhorRota)

	// Fatiamento UDP mantido integralmente e corrigido
	totalChunks := 3
	chunkSize := len(rawPayload) / totalChunks
	if chunkSize == 0 {
		chunkSize = len(rawPayload)
	}

	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == totalChunks-1 || end > len(rawPayload) {
			end = len(rawPayload)
		}

		chunkData := rawPayload[start:end]
		checksum := sha256.Sum256(chunkData)

		packet := RSMPacket{
			SessionID:   sessionID,
			ChunkIndex:  i + 1,
			TotalChunks: totalChunks,
			Payload:     chunkData,
			Checksum:    checksum,
		}

		fmt.Printf("  📦 [UDP MUX PREDITIVO] Bloco %d/%d (Tamanho: %d bytes) despachado via %s -> Checksum: %x...\n", 
			packet.ChunkIndex, packet.TotalChunks, len(packet.Payload), melhorRota, packet.Checksum[:4])
		
		if end == len(rawPayload) {
			break
		}
	}

	return melhorRota
}
