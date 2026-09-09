package main

import (
	"fmt"
	"time"
)

// CicloAbordagem controla o histórico e os lembretes para cada alvo
type CicloAbordagem struct {
	NomeAlvo        string
	TentativaAtual  int        // 1, 2 ou 3
	UltimoContato   time.Time  
	PerdaMensal     float64    // Quanto ele perde por mês no gargalo
}

// DispararAbordagem simula o envio da mensagem com base no momento do ciclo
func (c *CicloAbordagem) DispararAbordagem() {
	switch c.TentativaAtual {
	case 1:
		c.TentativaAtual = 2
		c.UltimoContato = time.Now()
		fmt.Printf("[ABORDAGEM 1] Para: %s\n", c.NomeAlvo)
		fmt.Println("Gostaria de usar meu sistema?")
		fmt.Printf("Seu ganho seria de aproximadamente 25%% de otimização no seu fluxo atual.\n\n")

	case 2:
		c.TentativaAtual = 3
		c.UltimoContato = time.Now()
		// 30 dias depois: calcula a economia acumulada no período
		economia30Dias := c.PerdaMensal * 0.25
		fmt.Printf("[ABORDAGEM 2 - 30 dias] Para: %s\n", c.NomeAlvo)
		fmt.Printf("Você teria economizado R$ %.2f (25%%) nos últimos 30 dias usando meu sistema.\n", economia30Dias)
		fmt.Println("Ainda deseja otimizar seu processo?\n")

	case 3:
		c.TentativaAtual = 4 // Marca como encerrado
		// 60 dias após os 30 (total de 90 dias): balanço trimestral final
		economia90Dias := (c.PerdaMensal * 3) * 0.25
		fmt.Printf("[ABORDAGEM 3 - FIM DO CICLO / 90 dias] Para: %s\n", c.NomeAlvo)
		fmt.Printf("Por fim... nos últimos 90 dias, você teria economizado R$ %.2f.\n", economia90Dias)
		fmt.Println("Oportunidade encerrada. A infraestrutura está disponível caso decida evoluir.")
		fmt.Println("--------------------------------------------------\n")
	}
}

func main() {
	// Simulando um alvo detectado pelo robô
	alvo := CicloAbordagem{
		NomeAlvo:       "Sistema Comercial / Comércio Parceiro",
		TentativaAtual: 1,
		PerdaMensal:    2000.00, // Exemplo de perda mensal estimada no gargalo
	}

	// 1º Contato imediato
	alvo.DispararAbordagem()

	// (Simulando o passar de 30 dias para o segundo contato)
	alvo.DispararAbordagem()

	// (Simulando o passar de mais 60 dias para o balanço final de 90 dias)
	alvo.DispararAbordagem()
	
	// Tentativa seguinte (já encerrado, o robô ignora)
	if alvo.TentativaAtual > 3 {
		fmt.Println("[ROBÔ] Alvo finalizado. Nenhuma nova mensagem será enviada.")
	}
}
