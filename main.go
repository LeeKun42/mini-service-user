package main

import "iris-app/app/cmd"

func main() {
	cmd.Main.AddCommand(&cmd.GormDto)
	cmd.Main.Execute()
}
