package pack_size

type OrderRequestDto struct {
	Quantity uint `json:"quantity"`
}

type PackCountResponseDto struct {
	Size  uint `json:"size"`
	Count uint `json:"count"`
}
