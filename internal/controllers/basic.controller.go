package controllers

import (
	"github.com/Sinanaas/gotth-financial-tracker/internal/managers"
)

type BasicController struct {
	BM managers.BasicManager
}

func NewBasicController(bm managers.BasicManager) BasicController {
	return BasicController{bm}
}
