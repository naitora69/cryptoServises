package main

import (
	"fmt"
	"log"
	"scam-analyzator-service/internal/scamanalyzator"
)

func main() {
	// TODO сделать конечную точку для GET spaceId
	spaceID := "stgdao.eth"

	tokenIds, err := scamanalyzator.GetTokenId(spaceID)
	if err != nil {
		log.Println("Не найден токены для данного пр-ва")
	}
	res := scamanalyzator.FinalRes(tokenIds)

	resStr := fmt.Sprintf("Найден ли скам в пр-ве %s = %v\n", spaceID, res)

	// пока вывод просто в консоль
	// TODO сделать конечную точку для вывода ответа
	fmt.Println(resStr)
}
