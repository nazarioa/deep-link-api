package internal

type (
	Link struct {
		Destination  string `json:"destination" validate:"required"`
		Fingerprint  string `json:"fingerprint" validate:"required"`
		MemberIdHash string `json:"member_id_hash" validate:"required"`

		CreatedAt string `json:"created_at"`
		ID        int    `json:"id"`
	}
)
