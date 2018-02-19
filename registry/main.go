package main

func main() {
	registry := NewRegistry()
	registry.StartRegistryServiceAndListen()
}
