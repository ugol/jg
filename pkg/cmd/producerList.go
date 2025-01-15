// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var producerListCmd = &cobra.Command{
	Use:   "list",
	Short: "describes available producers",
	Long: "describes available producers. Example usage:\n" +
		"jr producer list",
	Run: func(cmd *cobra.Command, args []string) {

		noColor, _ := cmd.Flags().GetBool("nocolor")

		var Green = ""
		var Reset = ""
		if !noColor {
			Green = "\033[32m"
			Reset = "\033[0m"
		}

		fmt.Println()
		fmt.Println("List of JR producers:")
		fmt.Println()

		fmt.Printf("%sConsole *%s (--output = stdout)\n", Green, Reset)
		fmt.Printf("%sKafka%s (--output = kafka)\n", Green, Reset)
		fmt.Printf("%sHTTP%s (--output = http)\n", Green, Reset)
		fmt.Printf("%sRedis%s (--output = redis)\n", Green, Reset)
		fmt.Printf("%sRedis HASH%s (--output = redishash)\n", Green, Reset)
		fmt.Printf("%sRedis JSON%s (--output = redisjson)\n", Green, Reset)
		fmt.Printf("%sMongodb%s (--output = mongo)\n", Green, Reset)
		fmt.Printf("%sElastic%s (--output = elastic)\n", Green, Reset)
		fmt.Printf("%sS3%s (--output = s3)\n", Green, Reset)
		fmt.Printf("%sGCS%s (--output = gcs)\n", Green, Reset)
		fmt.Printf("%sAZBlobStorage%s (--output = azblobstorage)\n", Green, Reset)
		fmt.Printf("%sAZCosmosDB%s (--output = azcosmosdb)\n", Green, Reset)
		fmt.Printf("%sCassandra%s (--output = cassandra)\n", Green, Reset)
		fmt.Printf("%sLUA Script%s (--output = luascript)\n", Green, Reset)
		fmt.Printf("%sWASM Function%s (--output = wasm)\n", Green, Reset)
		fmt.Printf("%sAWS DynamoDB%s (--output = awsdynamodb)\n", Green, Reset)
		fmt.Printf("%sWAMP Topic%s (--output = wamp)\n", Green, Reset)
		fmt.Println()

	},
}

func init() {
	producerCmd.AddCommand(producerListCmd)
	producerListCmd.Flags().BoolP("nocolor", "n", false, "Do not color output")
}
