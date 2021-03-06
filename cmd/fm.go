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
/*
* CMD -- Finding Meeting
 */

package cmd

import (
	"bufio"
	"fmt"
	"github.com/SYSU532/agenda/entity"
	"github.com/SYSU532/agenda/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	fmTitle                string
	fmStartTime, fmEndTime string
)

// fmCmd represents the fm command
var fmCmd = &cobra.Command{
	Use:   "fm",
	Short: "Find meeting within speical conditions",
	Long: fmt.Sprintf(`Find special meetings within input range that has association with the login user.

Usage: %v fm [-t title] [-s startTime] [-e endTime]`, os.Args[0]),

	Run: func(cmd *cobra.Command, args []string) {
		// Write init lOG
		log.WriteLog("Invoke find meeting command to search conditional meetings you create or join in", 1)
		var (
			result []entity.Meeting
			err    error
		)
		format := "2006-01-02 15:04"
		reader := bufio.NewReader(os.Stdin)
		_, err0 := entity.GetCurrentUser()
		if err0 != nil {
			fmt.Fprintf(os.Stderr, "Fail to find meetings: %v\n", err0)
			log.WriteLog("Error when geeting current user, maybe you are not logged in", 0)
			return
		}
		// If user do enter title, just search with title
		if fmTitle != "" {
			result, err = entity.FindMeetingByTitle(fmTitle)
		} else {
			// If user do not enter title, search with time interval
			if fmStartTime == "" {
				fmt.Printf("Enter the start time of meetings interval: ")
				fmStartTime, _ = reader.ReadString('\n')
				fmStartTime = fmStartTime[:len(fmStartTime)-1]
			}
			if fmEndTime == "" {
				fmt.Printf("Enter the end time of meetings interval: ")
				fmEndTime, _ = reader.ReadString('\n')
				fmEndTime = fmEndTime[:len(fmEndTime)-1]
			}
			// Parsing string to time.Time
			startTime, _ := time.Parse(format, fmStartTime)
			endTime, _ := time.Parse(format, fmEndTime)
			result, err = entity.FindMeetingsByTimeInterval(startTime, endTime)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if len(result) == 0 {
			tmp := "No any meeting you create or join satisfies the searching conditions!"
			fmt.Println(tmp)
			log.WriteLog(tmp, 0)
			os.Exit(1)
		}
		fmt.Printf("|%-20v|%-20v|%-20v|%-20v|%-20v|\n", "Title", "StartTime", "EndTime", "Creator", "Participant")
		fmt.Println("|--------------------|--------------------|--------------------|--------------------|--------------------|")
		for _, ele := range result {
			fmt.Printf("|%-20v|%-20v|%-20v|%-20v|", ele.Title, ele.StartTime, ele.EndTime, ele.Creator)
			var strRes string
			for i, pa := range ele.Participant {
				strRes += pa
				if i != len(ele.Participant)-1 {
					strRes += ", "
				}
			}
			fmt.Printf("%-20v|\n", strRes)
		}
		log.WriteLog("Finding meetings done! Print target meetings..", 1)
	},
}

func init() {
	rootCmd.AddCommand(fmCmd)

	fmCmd.Flags().StringVarP(&fmTitle, "title", "t", "", "The title of meeting you want to find")
	fmCmd.Flags().StringVarP(&fmStartTime, "startTime", "s", "", "The start time of meetings you want to find")
	fmCmd.Flags().StringVarP(&fmEndTime, "endTime", "e", "", "The end time of meetings you want to find")
}
