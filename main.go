package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func proxyViaCep(w http.ResponseWriter, r *http.Request) {
	// Pega o caminho repassado, ignorando o domínio do Render
	targetURL := "https://viacep.com.br" + r.URL.Path

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Disfarça a requisição
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Erro ao contatar ViaCEP: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Repassa os cabeçalhos originais
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", proxyViaCep)

	// O Render exige que a aplicação escute na porta definida pela variável de ambiente PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback para rodar localmente
	}

	log.Println("Proxy rodando na porta " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}