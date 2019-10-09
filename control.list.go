package main

//Direction is a component used for dynamic routing
type Direction struct {
	Name   string `yaml:"name"`
	Target string `yaml:"upstream"`
}

//Directions is a listing of objects that should be used for building reverse proxy targets
type Directions []Direction
