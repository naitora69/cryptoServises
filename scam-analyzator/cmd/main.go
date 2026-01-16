package main

import (
	"fmt"
	"log"
	orchestrascama "scam-analyzator-service/internal/scamanalyzator/orchestra-scama"
)

func main() {
	// TODO сделать конечную точку для GET spaceId
	spaceID := "stgdao.eth"

	res, err := orchestrascama.FinalRes(spaceID)

	if err != nil {
		log.Println(err)
	}

	resStr := fmt.Sprintf("Найден ли скам в пр-ве %s = ageFlag - %v, countHoldersFLag - %v\n", spaceID, res.TokenAgeScam, res.TokenHoldersCountScam)

	// пока вывод просто в консоль
	// TODO сделать конечную точку для вывода ответа
	fmt.Println(resStr)
}
