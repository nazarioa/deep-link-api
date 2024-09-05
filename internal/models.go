package internal

type (
	Link struct {
		Destination string `json:"destination" validate:"required"`
		CreatedAt   string `json:"created_at" validate:"required"`

		Fingerprint  string `json:"fingerprint"`
		ID           int    `json:"id"`
		MemberIdHash string `json:"member_id_hash"`
	}
)
