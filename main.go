package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {

	canal1 := make(chan ViaCep)
	canal2 := make(chan BrasilApi)

	// Exemplo aleatório de CEP para teste
	cep := "01001000"

	go func() {
		req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		}
		var dados ViaCep
		err = json.Unmarshal(res, &dados)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
		}

		//Simulando um delay para testar o time out
		//time.Sleep(2 * time.Second) // Simulando um delay

		canal1 <- dados
	}()

	go func() {

		req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler resposta: %v\n", err)
		}
		var dados BrasilApi
		err = json.Unmarshal(res, &dados)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
		}

		// Simulando um delay para testar o time out
		//time.Sleep(2 * time.Second)

		canal2 <- dados
	}()

	select {
	case endereco := <-canal1:
		fmt.Println("Api ViaCep:")
		fmt.Printf("Dados: %+v\n", endereco)

	case endereco := <-canal2:
		fmt.Println("Api BrasilApi:")
		fmt.Printf("Dados: %+v\n", endereco)

	case <-time.After(1 * time.Second):
		fmt.Println("Erro de Time out")
	}

}
