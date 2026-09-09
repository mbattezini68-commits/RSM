package main

import (
	"fmt"
	"time"
)

// CicloDadosInfra controla o histórico e os lembretes de tráfego pesado para cada alvo
type CicloDadosInfra struct {
	NomeAlvo        string  // Ex: Agencia Espacial / Operadora / Provedor de Cloud
	TentativaAtual  int     // 1, 2 ou 3
	UltimoContato   time.Time  
	VolumeTbMensal  float64 // Volume de dados gerado/trafegado por mês em Terabytes
	CustoPorTb      float64 // Custo atual por Terabyte na infraestrutura deles
}

// DispararAbordagemInfra envia a proposta técnica e financeira focada em tráfego de dados
func (c *CicloDadosInfra) DispararAbordagemInfra() {
	switch c.TentativaAtual {
	case 1:
		c.TentativaAtual = 2
		c.UltimoContato = time.Now()
		fmt.Printf("[HUNTER-INFRA] Destinatário: %s\n", c.NomeAlvo)
		fmt.Printf("Volume mapeado: %.1f TB/mês de tráfego de dados.\n", c.VolumeTbMensal)
		fmt.Println("Gostaria de usar meu sistema de roteamento otimizado em Go?")
		fmt.Printf("Seu ganho estimado seria de 30%% de redução na latência e no custo de banda.\n\n")

	case 2:
		c.TentativaAtual = 3
		c.UltimoContato = time.Now()
		// 30 dias depois: calcula o desperdício financeiro em tráfego de dados
		custoMensalAtual := c.VolumeTbMensal * c.CustoPorTb
		desperdicio30Dias := custoMensalAtual * 0.30
		fmt.Printf("[HUNTER-INFRA - 30 dias] Destinatário: %s\n", c.NomeAlvo)
		fmt.Printf("Nos últimos 30 dias, você perdeu aproximadamente R$ %.2f em sobrecarga de tráfego que poderia ter sido otimizada.\n", desperdicio30Dias)
		fmt.Println("Ainda deseja otimizar o escoamento de dados da sua infraestrutura?\n")

	case 3:
		c.TentativaAtual = 4 // Marca como encerrado
		// 60 dias após os 30 (total de 90 dias): balanço trimestral final de tráfego
		custoMensalAtual := c.VolumeTbMensal * c.CustoPorTb
		desperdicio90Dias := (custoMensalAtual * 3) * 0.30
		fmt.Printf("[HUNTER-INFRA - FIM DO CICLO / 90 dias] Destinatário: %s\n", c.NomeAlvo)
		fmt.Printf("Balanço trimestral: nos últimos 90 dias, o custo excedente de banda acumulou R$ %.2f.\n", desperdicio90Dias)
		fmt.Println("Oportunidade encerrada. A infraestrutura descentralizada em Go permanece em repouso.")
		fmt.Println("--------------------------------------------------\n")
	}
}

func main() {
	// Exemplo de um alvo de grande volume de dados detectado pelo robô
	alvoInfra := CicloDadosInfra{
		NomeAlvo:       "Centro de Telemetria / Infraestrutura de Dados de Alta Escala",
		TentativaAtual: 1,
		VolumeTbMensal: 150.0, // Exemplo: 150 Terabytes trafegados por mês
		CustoPorTb:     800.0, // Exemplo de custo por Terabyte
	}

	// 1º Contato imediato
	alvoInfra.DispararAbordagemInfra()

	// (Simulando a passagem de 30 dias)
	alvoInfra.DispararAbordagemInfra()

	// (Simulando o balanço final de 90 dias)
	alvoInfra.DispararAbordagemInfra()
}
