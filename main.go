package main

import "user/app/cmd"

func main() {
	cmd.Main.AddCommand(&cmd.GormDto)
	cmd.Main.Execute()
}
