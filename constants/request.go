package constants

type RecipeRequest struct {
	Title       string `json:"title"`
	MakingTime  string `json:"preparation_time"`
	Serves      string `json:"serves"`
	Ingredients string `json:"ingredients"`
	Cost        int    `json:"cost"`
}
