package model

type Input struct {
	String      string   `json:"string" binding:"required"`
	ArrayString []string `json:"arrayString" binding:"required"`
}

type Pikmin struct {
	ID    string `json:"id" bson:"_id"`
	Color string `json:"color" bson:"color" binding:"required"`
	Head  string `json:"head" bson:"head" binding:"required"`
}

type InputByColor struct {
	Color string `json:"color" bson:"color" binding:"required"`
}

type InputByPikminsID struct {
	IDS []string `json:"ids" binding:"required"`
}
