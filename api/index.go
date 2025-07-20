package handler

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
)

type requestBody struct {
	Strength        int
	Endurance       int
	DiceNumber      int
	TouchDifficulty int
	ArmorSave       int
	InvuSave        int
	RunNumber       int
}

type params struct {
	DiceNumber      int
	TouchDifficulty int
	HurtDifficulty  int
	ArmorSave       int
	InvuSave        int
	RunNumber       int
}

func Handler(w http.ResponseWriter, r *http.Request) {
	params, err := getParamsFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results := runAll(params)

	fmt.Fprintf(w, "<h1>Here are the results of all the runs : %d</h1>", results)
}

func getParamsFromRequest(r *http.Request) (params, error) {
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return params{}, fmt.Errorf("invalid JSON body: %w", err)
	}

	hurtDifficulty := getDifficulty(body.Strength, body.Endurance)

	return params{
		TouchDifficulty: body.TouchDifficulty,
		HurtDifficulty:  hurtDifficulty,
		ArmorSave:       body.ArmorSave,
		InvuSave:        body.InvuSave,
		RunNumber:       body.RunNumber,
		DiceNumber:      body.DiceNumber,
	}, nil
}

func rollDice(sides int) int {
	return rand.IntN(sides) + 1
}

func rollDices(numberOfDices int, sides int) []int {
	var dicesResults []int

	for i := 0; i < numberOfDices; i++ {
		dicesResults = append(dicesResults, rollDice(sides)+1)
	}

	return dicesResults
}

func getDifficulty(strength int, endurance int) int {
	if strength >= endurance*2 {
		return 2
	} else if strength > endurance {
		return 3
	} else if strength*2 <= endurance {
		return 6
	} else if strength < endurance {
		return 5
	}

	return 4
}

func getNumberOfTouches(params params, numberOfDices int) int {
	dicesResults := rollDices(numberOfDices, 6)
	numberOfTouches := 0

	for _, result := range dicesResults {
		if result >= params.TouchDifficulty {
			numberOfTouches++
		}
	}

	return numberOfTouches
}

func getNumberOfHurts(params params, numberOfDices int) int {
	dicesResults := rollDices(numberOfDices, 6)
	numberOfHurts := 0

	for _, result := range dicesResults {
		if result >= params.HurtDifficulty {
			numberOfHurts++
		}
	}

	return numberOfHurts
}

func getNumberOfHurtsAfterSave(params params, numberOfHurts int) int {
	numberOfHurtsAfterArmorSave := 0
	if params.ArmorSave >= 1 {
		for i := 0; i < numberOfHurts; i++ {
			if rollDice(6) < params.ArmorSave {
				numberOfHurtsAfterArmorSave++
			}
		}
	} else {
		numberOfHurtsAfterArmorSave = numberOfHurts
	}

	numberOfHurtsAfterInvuSave := 0
	if params.InvuSave >= 1 {
		for i := 0; i < numberOfHurtsAfterArmorSave; i++ {
			if rollDice(6) < params.InvuSave {
				numberOfHurtsAfterInvuSave++
			}
		}
	} else {
		numberOfHurtsAfterInvuSave = numberOfHurtsAfterArmorSave
	}

	return numberOfHurtsAfterInvuSave
}

func runOnce(params params) int {
	numberOfTouches := getNumberOfTouches(params, params.DiceNumber)
	numberOfHurts := getNumberOfHurts(params, numberOfTouches)
	finalHurts := getNumberOfHurtsAfterSave(params, numberOfHurts)

	return finalHurts
}

func runAll(params params) []int {
	var allResults []int

	for i := 0; i < params.RunNumber; i++ {
		allResults = append(allResults, runOnce(params))
	}

	return allResults
}
