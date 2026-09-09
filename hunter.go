package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	C  = 3e8
	PI = math.Pi
)

type AlvoInfraEspacial struct {
	NomeAlvo      string  `json:"OBJECT_NAME"`
	NoradID       int     `json:"NORAD_CAT_ID"`
	FreqHz        float64 
	DistReal      float64 
	BandaHz       float64 
	JanelaSeg     float64 
	NivelRuido    float64 
	Tentativa     int     
}

type SateliteAPI struct {
	NomeAlvo string `json:"OBJECT_NAME"`
	NoradID  int    `json:"NORAD_CAT_ID"`
}

func BuscarAlvoRealDoMundo(noradID int) (*AlvoInfraEspacial, error) {
	url := fmt.Sprintf("https://celestrak.org/NORAD/elements/gp.php?CATNR=%d&FORMAT=json", noradID)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com a API externa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status code inesperado: %d", resp.StatusCode)
	}

	var resultado []SateliteAPI
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&resultado); err != nil {
		return nil, fmt.Errorf("falha ao decodificar JSON real: %v", err)
	}

	if len(resultado) == 0 {
		return nil, fmt.Errorf("nenhum alvo encontrado para o NORAD ID %d", noradID)
	}

	alvoReal := &AlvoInfraEspacial{
		NomeAlvo:   resultado[0].NomeAlvo,
		NoradID:    resultado[0].NoradID,
		FreqHz:     8450000000.0,
		DistReal:   550000.0,     
		BandaHz:    2000000.0,    
		JanelaSeg:  600.0,        
		NivelRuido: 720.0,        
		Tentativa:  1,
	}

	return alvoReal, nil
}

func (a *AlvoInfraEspacial) calculaRaioMM() float64 {
	return math.Sqrt(C * a.DistReal * 1e6 / (a.FreqHz * 4))
}

func (a *AlvoInfraEspacial) calculaVolumeBolhaComRuido() (vol, raio float64) {
	raioMM := a.calculaRaioMM()
	raioM := raioMM / 1000.0
	volBase := (4.0/3.0) * PI * math.Pow(raioM, 3)
	
	fatorRuido := (1000.0 - a.NivelRuido) / 1000.0
	if a.NivelRuido > 900 { 
		fatorRuido = 0.1 
	}
	
	vol = volBase * fatorRuido
	raio = raioMM
	return vol, raio
}

func (a *AlvoInfraEspacial) podeEstourar(wVol, economiaBits float64) bool {
	if a.BandaHz <= 0 || a.JanelaSeg <= 0 || a.NivelRuido < 0 || a.NivelRuido > 1000 {
		fmt.Printf("[DAEMON SEGURANÇA] Parâmetros de rede inválidos para o alvo real.\n")
		return false
	}
	
	if wVol == 0 || economiaBits == 0 {
		fmt.Printf("[DAEMON SEGURANÇA] W_vol=0 ou economia zerada. Potro morre sem valor, zero gás.\n")
		return false
	}
	
	return true
}

func (a *AlvoInfraEspacial) ProcessarMetricas() (wVol, economiaBits, tempoEconomizado, volBolha float64) {
	volBolha, raioMM := a.calculaVolumeBolhaComRuido()
	volCilindro := PI * math.Pow(raioMM/1000.0, 2) * a.DistReal * ((1000.0 - a.NivelRuido) / 1000.0)
	
	wVol = math.Abs(volCilindro - volBolha)

	bitsSem := a.BandaHz * a.JanelaSeg * 0.5 
	bitsCom := bitsSem * 3.5                 
	economiaBits = math.Abs(bitsSem - bitsCom)
	tempoEconomizado = economiaBits / a.BandaHz

	return wVol, economiaBits, tempoEconomizado, volBolha
}

func (a *AlvoInfraEspacial) ExecutarCadencia() {
	wVol, economiaBits, tempoEco, volBolha := a.ProcessarMetricas()

	if !a.podeEstourar(wVol, economiaBits) {
		fmt.Printf("[HUNTER] Alvo %s abortado pela segurança off-chain. Nenhuma transação enviada.\n\n", a.NomeAlvo)
		return
	}

	switch a.Tentativa {
	case 1:
		a.Tentativa = 2
		fmt.Printf("=== [HUNTER MUNDO REAL - ABORDAGEM 1] Alvo: %s (NORAD: %d) ===\n", a.NomeAlvo, a.NoradID)
		fmt.Printf("Volume Bolha P_vol: %.2e m3 | Raio Fresnel: %.2f m\n", volBolha, a.calculaRaioMM()/1000.0)
		fmt.Println("Gostaria de usar meu sistema de roteamento otimizado em Go?")
		fmt.Printf("Ganho técnico: %.2f s economizados na janela de %.0f s (%.1f%% de ganho de banda).\n\n", 
			tempoEco, a.JanelaSeg, (tempoEco/a.JanelaSeg)*100)

	case 2:
		a.Tentativa = 3
		acumulado30Dias := economiaBits * 30
		fmt.Printf("=== [HUNTER MUNDO REAL - ABORDAGEM 2 / 30 DIAS] Alvo: %s ===\n", a.NomeAlvo)
		fmt.Printf("Nos últimos 30 dias, sua infraestrutura perdeu aproximadamente %.0f bits em sobrecarga evitável.\n", acumulado30Dias)
		fmt.Println("Ainda deseja otimizar o escoamento de dados da sua operação?\n")

	case 3:
		a.Tentativa = 4 
		acumulado90Dias := economiaBits * 90
		mintSugerido := wVol * 1e6 * 0.97 
		
		fmt.Printf("=== [HUNTER MUNDO REAL - ABORDAGEM 3 / FIM DO CICLO 90 DIAS] Alvo: %s ===\n", a.NomeAlvo)
		fmt.Printf("Balanço trimestral: desperdício acumulado de %.0f bits.\n", acumulado90Dias)
		fmt.Printf("PROVA NA PRÁTICA: W_vol > 0. Mint sugerido de %.0f BOLHA se converter agora.\n", mintSugerido)
		fmt.Println("Oportunidade encerrada. Ciclo finalizado. A infraestrutura em Go permanece em repouso.")
		fmt.Println("--------------------------------------------------\n")
	}
}

func main() {
	fmt.Println("=== MOTOR THE COLT: CONECTOR DE MUNDO REAL ATIVADO ===")
	fmt.Println("Buscando dados reais de satélite/telemetria via API externa...\n")

	alvoReal, err := BuscarAlvoRealDoMundo(25544)
	if err != nil {
		fmt.Printf("[ERRO CRÍTICO] Falha ao obter alvo do mundo real: %v\n", err)
		return
	}

	fmt.Printf("[SUCESSO] Alvo real carregado da rede: %s (NORAD ID: %d)\n\n", alvoReal.NomeAlvo, alvoReal.NoradID)

	alvoReal.ExecutarCadencia() 
	alvoReal.ExecutarCadencia() 
	alvoReal.ExecutarCadencia() 
}
