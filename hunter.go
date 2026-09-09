package main

import (
	"fmt"
	"math"
)

const (
	C  = 3e8
	PI = math.Pi
)

// AlvoInfraEspacial representa a infraestrutura pesada de dados (satélites, telemetria, redes de alta escala)
type AlvoInfraEspacial struct {
	NomeAlvo      string
	EdgeID        string  // Identificador da borda (ex: DSN_GOLDSTONE)
	NoradID       int     // ID do satélite / infra
	FreqHz        float64 // Frequência em Hz (ex: X-band 8.45 GHz)
	DistReal      float64 // Distância em metros (mínimo 384.400km ou LEO)
	BandaHz       float64 // Largura de banda em Hz
	JanelaSeg     float64 // Janela de tempo em segundos
	NivelRuido    float64 // Ruído de 0 a 1000
	Tentativa     int     // 1, 2 ou 3 (Cadência comercial)
}

// 1. MOTOR DE CÁLCULO DE RAIO E VOLUME (Física de Propagação)
func (a *AlvoInfraEspacial) calculaRaioMM() float64 {
	return math.Sqrt(C * a.DistReal * 1e6 / (a.FreqHz * 4))
}

func (a *AlvoInfraEspacial) calculaVolumeBolhaComRuido() (vol, raio float64) {
	raioMM := a.calculaRaioMM()
	raioM := raioMM / 1000.0
	volBase := (4.0/3.0) * PI * math.Pow(raioM, 3)
	
	fatorRuido := (1000.0 - a.NivelRuido) / 1000.0
	if a.NivelRuido > 900 { 
		fatorRuido = 0.1 // Muito ruído = bolha menor mas mais densa
	}
	
	vol = volBase * fatorRuido
	raio = raioMM
	return vol, raio
}

// 2. TRAVA DE SEGURANÇA DA BORDA (Garante zero desperdício de gás off-chain)
func (a *AlvoInfraEspacial) podeEstourar(wVol, economiaBits float64) bool {
	if a.BandaHz <= 0 || a.JanelaSeg <= 0 || a.NivelRuido < 0 || a.NivelRuido > 1000 {
		fmt.Printf("[DAEMON SEGURANÇA] Parâmetros de rede inválidos no Edge %s.\n", a.EdgeID)
		return false
	}
	
	if wVol == 0 || economiaBits == 0 {
		fmt.Printf("[DAEMON SEGURANÇA] W_vol=0 ou economia zerada. Potro morre sem valor, zero gás.\n")
		return false
	}
	
	return true
}

// 3. PROCESSAMENTO MATEMÁTICO COMPLETO
func (a *AlvoInfraEspacial) ProcessarMetricas() (wVol, economiaBits, tempoEconomizado, volBolha float64) {
	volBolha, raioMM := a.calculaVolumeBolhaComRuido()
	volCilindro := PI * math.Pow(raioMM/1000.0, 2) * a.DistReal * ((1000.0 - a.NivelRuido) / 1000.0)
	
	wVol = math.Abs(volCilindro - volBolha)

	bitsSem := a.BandaHz * a.JanelaSeg * 0.5 // Eficiência base com ruído
	bitsCom := bitsSem * 3.5                 // Ganho volumétrico da bolha
	economiaBits = math.Abs(bitsSem - bitsCom)
	tempoEconomizado = economiaBits / a.BandaHz

	return wVol, economiaBits, tempoEconomizado, volBolha
}

// 4. ROBÔ CAÇADOR COM CADÊNCIA COMERCIAL DE 3 TOQUES
func (a *AlvoInfraEspacial) ExecutarCadencia() {
	wVol, economiaBits, tempoEco, volBolha := a.ProcessarMetricas()

	if !a.podeEstourar(wVol, economiaBits) {
		fmt.Printf("[HUNTER] Alvo %s abortado pela segurança off-chain. Nenhuma transação enviada.\n\n", a.NomeAlvo)
		return
	}

	switch a.Tentativa {
	case 1:
		a.Tentativa = 2
		fmt.Printf("=== [HUNTER - ABORDAGEM 1] Alvo: %s (NORAD: %d) ===\n", a.NomeAlvo, a.NoradID)
		fmt.Printf("Volume Bolha P_vol: %.2e m3 | Raio Fresnel: %.2f m\n", volBolha, a.calculaRaioMM()/1000.0)
		fmt.Println("Gostaria de usar meu sistema de roteamento otimizado em Go?")
		fmt.Printf("Ganho técnico: %.2f s economizados na janela de %.0f s (%.1f%% de ganho de banda).\n\n", 
			tempoEco, a.JanelaSeg, (tempoEco/a.JanelaSeg)*100)

	case 2:
		a.Tentativa = 3
		acumulado30Dias := economiaBits * 30
		fmt.Printf("=== [HUNTER - ABORDAGEM 2 / 30 DIAS] Alvo: %s ===\n", a.NomeAlvo)
		fmt.Printf("Nos últimos 30 dias, sua infraestrutura perdeu aproximadamente %.0f bits em sobrecarga evitável.\n", acumulado30Dias)
		fmt.Println("Ainda deseja otimizar o escoamento de dados da sua operação?\n")

	case 3:
		a.Tentativa = 4 // Fim do ciclo
		acumulado90Dias := economiaBits * 90
		mintSugerido := wVol * 1e6 * 0.97 // 3% royalty deduzido / 47% burn aplicado
		
		fmt.Printf("=== [HUNTER - ABORDAGEM 3 / FIM DO CICLO 90 DIAS] Alvo: %s ===\n", a.NomeAlvo)
		fmt.Printf("Balanço trimestral: desperdício acumulado de %.0f bits.\n", acumulado90Dias)
		fmt.Printf("PROVA NA PRÁTICA: W_vol > 0. Mint sugerido de %.0f BOLHA se converter agora.\n", mintSugerido)
		fmt.Println("Oportunidade encerrada. Ciclo finalizado. A infraestrutura em Go permanece em repouso.")
		fmt.Println("--------------------------------------------------\n")
	}
}

func main() {
	fmt.Println("=== MOTOR THE COLT: INFRAESTRUTURA ESPACIAL E DADOS ===")
	fmt.Println("Modo de auditoria e prospecção ativado com segurança off-chain.\n")

	alvo := AlvoInfraEspacial{
		NomeAlvo:   "Centro de Telemetria Espacial",
		EdgeID:     "DSN_GOLDSTONE",
		NoradID:    25544,
		FreqHz:     8450000000.0, // X-band 8.45 GHz
		DistReal:   550000.0,     // 550 km LEO
		BandaHz:    1000000.0,    // 1 MHz
		JanelaSeg:  600.0,        // 10 minutos
		NivelRuido: 750.0,        // Ruído alto
		Tentativa:  1,
	}

	alvo.ExecutarCadencia() // Toque 1
	alvo.ExecutarCadencia() // Toque 2 (30 dias)
	alvo.ExecutarCadencia() // Toque 3 (90 dias)
}
