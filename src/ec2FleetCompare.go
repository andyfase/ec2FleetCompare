package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
	"time"
	"net/http"
	"os"
	"strconv"
	"encoding/json"
	"errors"
	"regexp"
	"io/ioutil"
)

var cacheDir = ".ec2FleetCompare"
var ec2PricesURL string = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json";

type Instance struct {
	Name        string
	Mem         int
	Cpu         int
	Price       float64
	Os					string
	CpuClock		string
	NetworkType int
	NetworkDesc	string
	Description string
	Sku					string
	Region			string
}

type Ec2 struct {
	Demand []Instance
	Spot   []Instance
}

func printError(s string) {
	fmt.Println("***************************** ERROR ********************************************")
	fmt.Printf("ERROR: %s\n", s)
	fmt.Println("********************************************************************************\n\n")
}

func getJson(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func downloadDemandPrices (s *Ec2, region string) error {
	var data map[string]interface{}
	if err := getJson(ec2PricesURL, &data); err != nil {
		return err
	}
	serverTypes, ok := data["products"].(map[string]interface{})
	if ! ok {
		return errors.New("Type assertion failed on products")
	}

	r_mem := regexp.MustCompile(`(\d+)\s+GiB`)

	for server, serverSpecs := range serverTypes {
		serverSpecs, ok := serverSpecs.(map[string]interface{})
		if !ok {
			return errors.New("Type assertion failed on Server Specs")
		}

		// make sure this is actually a EC2 server JSON object
		family, ok := serverSpecs["productFamily"].(string)
		if !ok || family != "Compute Instance" {
			continue
		}

		serverAttributes, ok := serverSpecs["attributes"].(map[string]interface{})
		if ! ok {
			return errors.New("Type assertion failed on attributes")
		}

		// just process Linux and Windows instances (ignore specific RHEL, SUSE etc)
		os, ok := serverAttributes["operatingSystem"].(string)
		if !ok || (os != "Linux" && os != "Windows") {
			continue
		}

		// drop anything thats bring your own license
		license, ok := serverAttributes["licenseModel"].(string)
		if !ok || (license == "Bring your own license") {
			continue
		}

		// drop anything having pre-installed software
		software, ok := serverAttributes["preInstalledSw"].(string)
		if !ok || (software != "" && software != "NA") {
			continue
		}

		mem := r_mem.FindStringSubmatch(serverAttributes["memory"].(string))

		var i Instance
		i.Cpu, _  = strconv.Atoi(serverAttributes["vcpu"].(string))
		i.CpuClock, ok = serverAttributes["clockSpeed"].(string)
		i.Name, ok = serverAttributes["instanceType"].(string)
		i.Region, ok = serverAttributes["location"].(string)
		i.NetworkDesc, ok = serverAttributes["networkPerformance"].(string)
		i.Os = os
		i.Sku = server

		// set memory based on output of regex
		if (len(mem) >= 2) {
			i.Mem, _ = strconv.Atoi(mem[1])
		} else {
			i.Mem = 0 // basically could not match memory
		}

		// set networkType code based on networkDesc
		switch i.NetworkDesc {
		case `10 Gigabit`:
			i.NetworkType = 1
		case `High`:
			i.NetworkType = 2
		case `Moderate`:
			i.NetworkType = 3
		default:
			i.NetworkType = 4
		}

		// add instance to out ec2 struct under the demand array (of instances)
		s.Demand = append(s.Demand, i)
	}
	return nil
}

func getPrices(s *Ec2, forceDownload bool, region string) error {

	// get homedir
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// check for presense of $HOME/.ec2FleetCompare directory, if doesnt exist create it
	if _, err = os.Stat(home + "/" + cacheDir); err != nil {
		if err = os.Mkdir(home + "/" + cacheDir, 0755); err != nil {
			return err
		}
	}

	// get cache files metaData
	metaCache, err := os.Stat(home + "/" + cacheDir + "/ec2.cache")
	if forceDownload || err != nil || metaCache.ModTime().Before(time.Now().AddDate(0, 0, -1)) {
		// file either doesnt exist or is over 1 day old, hence needs replacing.
		fmt.Println("Local price cache either doensnt exist or is to old")
		fmt.Println("Downloading update ...  this is ~50MB so could take some time... ")
		if err = downloadDemandPrices(s, region); err != nil {
			return err
		}

		// write processed response to cache
		b, _ := json.Marshal(s)
		if err = ioutil.WriteFile(home + "/" + cacheDir + "/ec2.cache", b, 0644); err != nil {
			return err
		}
	} else {
		// our cache file is present and not to old!
		b, _ := ioutil.ReadFile(home + "/" + cacheDir + "/ec2.cache")
		if err := json.Unmarshal(b, s); err != nil {
			return err
		}
	}

	return nil
}

func main() {

	app := cli.NewApp()
	app.Name = "EC2 Instance fleet Compare"
	app.Usage = "Use this app to find the cheapest price for a single or set of EC2 instances given your CPU, memory or network requirements. Given a minimum or maximum fleet size and the required resources across the fleet this app will find the cheapest EC2 instances that will fulfil your requirements."
	app.Version = "1.0.0"

	var minCPU, minMem, minNetwork, region string
	var forceDownload bool
	app.Commands = []cli.Command{
		{
			Name:  "node",
			Usage: "Perform comparisons based on a single EC2 instance.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "region, r",
					Value:       "",
					Usage:       "The EC2 region to perform price checks on",
					Destination: &region,
				},
				cli.StringFlag{
					Name:        "cpu, c",
					Value:       "",
					Usage:       "Minimum CPU virtual cores required",
					Destination: &minCPU,
				},
				cli.StringFlag{
					Name:        "mem, m",
					Value:       "",
					Usage:       "Minimum memoy (in GBi) required",
					Destination: &minMem,
				},
				cli.StringFlag{
					Name:        "network, n",
					Value:       "",
					Usage:       "Minimum network speed required (low, medium, high, gbit)",
					Destination: &minNetwork,
				},
				cli.BoolFlag{
					Name:        "force, f",
					Usage:       "Force download of latest version of AWS EC2 pricing file",
					Destination: &forceDownload,
				},
			},
			Action: func(c *cli.Context) error {
				var prices Ec2
				err := getPrices(&prices, forceDownload, region)
				if err != nil {
					printError(err.Error())
					return err
				}
				return nil
			},
		},
	}
	app.Run(os.Args)
}
