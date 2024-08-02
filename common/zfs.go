package common

type ZpoolCreateRequest struct {
	Name string `json:"name"`
}

type ZPoolResponse struct {
	Name string `json:"name"`
}

type ZpoolListResponse struct {
	Pools []ZPoolResponse `json:"pools"`
}
