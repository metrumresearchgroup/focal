package main

import (
	"github.com/metrumresearchgroup/focal/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)





func main(){
	viper.SetEnvPrefix("focal")
	viper.AutomaticEnv()

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}





