package main

import (
    "github.com/GeertJohan/go.rice"
)

//go:generate rice embed-go

var _ = rice.MustFindBox("static").HTTPBox()
