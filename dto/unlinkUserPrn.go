package dto

type PostUnlinkUserPrnDto struct {
	Email    string `json:"email" validate:"required"`
	BranchId int    `json:"branchId" validate:"required"`
	Prn      string `json:"prn" validate:"required"`
}
