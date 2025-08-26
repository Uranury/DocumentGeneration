package models

type RequestBody struct {
	Code   string `json:"code"`
	Format string `json:"format"`
	Data   any    `json:"data"`
}
