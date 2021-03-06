package cmd

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var ta *jwtauth.JWTAuth
var configuration Config


type Config struct {
	Directory string `yaml:"directory" json:"directory,omitempty"`
	Token string `yaml:"token" json:"token,omitempty"`
	RootURL string `mapstructure:"root_url" yaml:"root" json:"root,omitempty"`
	Port int `yaml:"port" json:"port,omitempty"`
}

func GetConfig() *Config {
	return &configuration
}




//FocalCmdf is the primary cobra command executed to launch the service
var FocalCmd = &cobra.Command{
	Use: "focal",
	Short: "PAM authenticated proxy",
	Long: "A PAM authenticated proxy that supports JWT authentication and header authentication downstream",
	Example: `focal --token <token> --directory /path/to/directory.yml`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetEnvPrefix("focal")
		viper.AutomaticEnv()


		err := viper.Unmarshal(&configuration)


		config := GetConfig()

		if err != nil {
			log.Fatalf("Unable to marshall viper into the designated struct. Details are: %s", err)
		}

		ta = jwtauth.New("HS256", []byte(config.Token), nil)


		directory, err := buildDirectory(config)

		if err != nil {
			log.Fatalf("Unable to read the provided directory file: %s", err)
		}

		r := Routes(directory)

		log.Infof("Listening on port %d", config.Port)
		http.ListenAndServe(":" + strconv.Itoa(config.Port), r)
	},
}


func init(){
	const tokenIdentifier string = "token"
	FocalCmd.Flags().StringP(tokenIdentifier,"t","","The secret token used for security operations with JWT")
	viper.BindPFlag(tokenIdentifier, FocalCmd.Flags().Lookup(tokenIdentifier))

	const directoryIdentifier string = "directory"
	FocalCmd.Flags().StringP(directoryIdentifier,"d","/etc/focal/directory.yml","The directory file containing instructions for what Focal should reverse proxy")
	viper.BindPFlag(directoryIdentifier,  FocalCmd.Flags().Lookup(directoryIdentifier))

	const portIdentifier string = "port"
	FocalCmd.Flags().IntP(portIdentifier,"p",9666,"The port on which focal will run")
	viper.BindPFlag(portIdentifier, FocalCmd.Flags().Lookup(portIdentifier))

	const rootIdentifier string = "root_url"
	FocalCmd.Flags().String(rootIdentifier, "/protected", "the root url on which Focal will reside.")
	viper.BindPFlag(rootIdentifier, FocalCmd.Flags().Lookup(rootIdentifier))
}


func Execute() error {
	return FocalCmd.Execute()
}

func buildDirectory(config *Config) (Directions, error) {
	if _, err := os.Stat(config.Directory); err == nil {
		log.Info("Located a directory file to parse")

		contents, err := ioutil.ReadFile(config.Directory)
		if err != nil {
			log.Error(err)
			return Directions{}, err
		}

		directory := Directions{}

		err = yaml.Unmarshal(contents, &directory)
		if err != nil {
			return Directions{}, err
		}

		return directory, nil
	}

	return Directions{}, fmt.Errorf("The directory file could not be located or accessed")
}

func badAuthResponse(w http.ResponseWriter, r *http.Request) {
	config := GetConfig()

	if r.Header.Get("content-type") == "application/json" {
		log.Error("No token is present and identified as a JSON request")
		log.Error("Request identified as ", r.Header.Get("content-type"))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//For interactive Sessions, re-direct to / for login
	log.Info("No token present, but it appears to be an interactive session. Redirecting to / to login")
	http.Redirect(w, r, config.RootURL+"/", http.StatusTemporaryRedirect)
	return
}




//Direction is a component used for dynamic routing
type Direction struct {
	Name   string `yaml:"name"`
	Target string `yaml:"upstream"`
	Type   string `yaml:"type"`
}

//Directions is a listing of objects that should be used for building reverse proxy targets
type Directions []Direction

