// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SYSU532/agenda/entity"
	"github.com/SYSU532/agenda/log"
	"github.com/spf13/cobra"
)

var (
	rpTitle         string
	rpParticipators []string
)

// rpCmd represents the rp command
var rpCmd = &cobra.Command{
	Use:   "rp",
	Short: "Remove participator",
	Long: fmt.Sprintf(`Use this command to remove participators from a meeting
	using a already logged in user.
	Usage: %v rp [-t title] [-p participator1] [-p participator2] ...`, os.Args[0]),
	Run: func(cmd *cobra.Command, args []string) {
		// Write init lOG
		log.WriteLog("Invoke remove participant command to clean special participants in your meetings", 1)
		userinfo, err := entity.GetCurrentUser()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fail to remove participator: %v\n", err)
			log.WriteLog("Error when geeting current user, maybe you are not logged in", 0)
			return
		}
		reader := bufio.NewReader(os.Stdin)
		if rpTitle == "" {
			fmt.Print("Enter the meeting title: ")
			title, _ := reader.ReadString('\n')
			rpTitle = title[:len(title)-1]
		}
		if len(rpParticipators) == 1 && rpParticipators[0] == "" {
			rpParticipators = []string{}
			fmt.Print("Enter the number of participators: ")
			var partNum uint
			fmt.Scan(&partNum)
			for i := uint(0); i < partNum; i++ {
				var part string
				fmt.Printf("Enter participator %d: ", i)
				fmt.Scan(&part)
				rpParticipators = append(rpParticipators, part)
			}
		}
		log.WriteLog(fmt.Sprintf("user %s begin to remove participants from meeting %s, target participants: %v", userinfo. Username, rpTitle, rpParticipators), 1)
		err = entity.CheckBeforeModP(rpTitle, userinfo.Username)
		if err == nil {
			for _, part := range rpParticipators {
				err = entity.RmParticipant(rpTitle, part)
				if err != nil {
					tmp := fmt.Sprintf("Fail to remove participant %v: %v", part, err)
					fmt.Println(tmp)
					log.WriteLog(tmp, 0)
				}
			}
			if err == nil {
				fmt.Println("Successfully removed participant(s)")
				log.WriteLog("Successfully removed participant(s)", 1)
			}
		} else {
			tmp := fmt.Sprintf("Fail to remove participant %v", err)
			fmt.Fprintf(os.Stderr, tmp, "\n")
			log.WriteLog(tmp, 0)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(rpCmd)

	rpCmd.Flags().StringVarP(&rpTitle, "title", "t", "", "The title of the meeting")
	rpCmd.Flags().StringArrayVarP(&rpParticipators, "participators", "p", []string{""}, "All the participators to be removed")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
