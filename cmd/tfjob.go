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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kubeflow/crd-validation/pkg/config"
	"github.com/kubeflow/crd-validation/pkg/crd"
)

func init() {
	RootCmd.AddCommand(tfjobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tfjobCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tfjobCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// tfjobCmd represents the tfjob command
var tfjobCmd = &cobra.Command{
	Use:   "tfjob",
	Short: "Generate TFJob CRD definition",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		generateTFJob()
	},
}

func generateTFJob() {
	original := config.NewCustomResourceDefinition("tfjob")
	var outputDir string
	if viper.Get("global") != nil {
		outputDir = viper.Get("global").(map[string]interface{})["output"].(string)
	}
	generator := crd.NewTFJobGenerator(outputDir)
	final := generator.Generate(original)
	generator.Export(final)
}
