package main

type userRequest struct {
	Name string `json:"name"`
}

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}
