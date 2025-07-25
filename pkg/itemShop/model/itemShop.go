package model

type (
	Item struct { //จังหวะนี้จะใช้ json เพราะเราต้องสร้าง fill ให้ฝั่ง client
		ID          uint64 `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Picture     string `json:"picture"`
		Price       uint   `json:"price"`
	}

	ItemFilter struct {
		Name        string `query:"name" validate:"omitempty,max=64"` //method get ใช้ boddy ไม่ได้ แต่ใช้ query
		Description string `query:"description" validate:"omitempty,max=128"`
		Paginate
	}

	Paginate struct {
		Page int64 `query:"page" validate:"omitempty,min=1"`
		Size int64 `query:"size" validate:"omitempty,min=1,max=20"`
	}

	ItemResult struct {
		Items    []*Item        `json:"items"`
		Paginate PaginateResult `json:"paginate"`
	}

	PaginateResult struct {
		Page      int64 `json:"page"`
		TotalPage int64 `json:"totalPage"`
	}
)
