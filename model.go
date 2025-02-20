package main

type Guess struct {
	BotId int
	Value int
}

type CoordinationMsg struct {
	WinnerId      int
	Suggestion    int // 1 -> bigger, -1 -> smaller
	ReadyForGuess bool
}

func winnerMsg(id int) CoordinationMsg {
	return CoordinationMsg{id, 0, false}
}
func isWinnerMsg(msg CoordinationMsg) bool {
	return msg.WinnerId > -1
}

func biggerSuggestionMsg() CoordinationMsg {
	return CoordinationMsg{-1, 1, false}
}
func smallerSuggestionMsg() CoordinationMsg {
	return CoordinationMsg{-1, -1, false}
}
func isSuggestionMsg(msg CoordinationMsg) bool {
	return msg.Suggestion != 0
}
func isBiggerSuggestionMsg(msg CoordinationMsg) bool {
	return msg.Suggestion == 1
}
func isSmallerSuggestionMsg(msg CoordinationMsg) bool {
	return msg.Suggestion == -1
}

func readyForGuessMsg() CoordinationMsg {
	return CoordinationMsg{-1, 0, true}
}
func isReadyForGuessMsg(msg CoordinationMsg) bool {
	return msg.ReadyForGuess
}
