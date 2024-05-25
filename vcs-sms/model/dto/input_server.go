package dto

type InputServer struct {
	Name   string `validate:"required"`
	IPv4   string `validate:"ipv4,required"`
	Status int    `validate:"gte=0,lte=1"`
}
