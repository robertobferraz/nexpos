package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robertobff/food-service/adapter"
	"github.com/robertobff/food-service/application"
	"go.uber.org/fx"
	"os"
	"path/filepath"
	"runtime"
)

func main() {
	if os.Getenv("ENV") != "production" {
		LoadConfig()
	}

	fx.New(
		adapter.Module,
		application.Module,
	).Run()
	// Abre o arquivo allCountries.txt
	//file, err := os.Open("allCountries.txt")
	//if err != nil {
	//	log.Fatalf("Erro ao abrir arquivo: %v", err)
	//}
	//defer file.Close()
	//
	//// Configura o leitor CSV para TSV
	//reader := csv.NewReader(file)
	//reader.Comma = '\t'         // Define tabulação como delimitador
	//reader.LazyQuotes = true    // Permite aspas não escapadas
	//reader.FieldsPerRecord = -1 // Permite número variável de campos por linha
	//
	//// Lê todas as linhas
	//records, err := reader.ReadAll()
	//if err != nil {
	//	log.Fatalf("Erro ao ler arquivo: %v", err)
	//}
	//
	//// Processa as linhas
	//for i, record := range records {
	//	// Verifica se a linha tem pelo menos 11 colunas (para evitar erros de índice)
	//	if len(record) < 11 {
	//		log.Printf("Linha %d inválida: número de colunas insuficiente (%d)", i+1, len(record))
	//		continue
	//	}
	//
	//	// Filtra por Brasil (coluna 8 é o código do país)
	//	if record[8] == "BR" {
	//		fmt.Printf("Cidade: %s, Estado: %s, Lat: %s, Long: %s\n",
	//			record[1],  // Nome do local
	//			record[10], // Divisão administrativa (estado)
	//			record[4],  // Latitude
	//			record[5])  // Longitude
	//	}
	//}
	//
	//fmt.Printf("Processamento concluído. Total de linhas: %d\n", len(records))
}

func LoadConfig() {
	_, b, _, _ := runtime.Caller(0)

	basepath := filepath.Dir(b)

	err := godotenv.Load(fmt.Sprintf("%v/.env", basepath))
	if err != nil {
		panic(err)
	}
}
