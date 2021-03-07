package cmd

import (
	authApi "canal/api/auth"
	customersApi "canal/api/customers"
	projectsApi "canal/api/projects"
	"canal/util"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strings"
)

// customerAddCmd represents the customer add command
var customerAddCmd = &cobra.Command{
	Use:   "add email:<email> name:<name> phone:<phone>",
	Short: "Adds a new customer",
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := email(args)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		project, err := util.CurrentProject()

		email, _ := email(args)
		name, _ := name(args)
		phone, _ := phone(args)

		if err != nil {
			token, err := util.UserToken()
			if err != nil {
				util.PrintlnError(err)
				return
			}

			projects, err := projectsApi.ProjectList(token)
			if err != nil {
				util.PrintlnError(err)
				return
			}
			var projectNames []string
			for i := range projects {
				projectNames = append(projectNames, projects[i].Id)
			}

			util.PrintlnInfo("please, first select a project you have access to")
			prompt := promptui.Select{
				Items: projectNames,
			}
			_, selectedProject, err := prompt.Run()
			if err != nil {
				util.PrintlnError(err)
				return
			}

			projectToken, err := authApi.LoginProject(token, selectedProject)
			if err != nil {
				util.PrintlnError(err)
				return
			}

			err = util.StoreProjectToken(util.ProjectName(selectedProject), projectToken)
			if err != nil {
				util.PrintlnError(err)
				return
			}

			err = util.UseProject(util.ProjectName(selectedProject))
			if err != nil {
				util.PrintlnError(err)
				return
			}

			project, err = util.CurrentProject()
			if err != nil {
				util.PrintlnError(err)
				return
			}
		}

		util.PrintlnInfo(fmt.Sprintf("waiting %v Canal", color.CyanString("for")))
		fmt.Printf("Adding %v... ", email)

		token, err := util.ProjectToken(util.ProjectName(project))
		if err != nil {
			util.PrintlnError(err)
			return
		}

		err = customersApi.AddCustomer(token, customersApi.Customer{
			Email:    email,
			Name:     name,
			LastName: name,
			Phone:    phone,
		})
		if err != nil {
			util.PrintlnError(err)
			return
		}

		fmt.Printf(" %v!", color.CyanString("done!"))
	},
}

func email(args []string) (string, error) {
	return argValue(args, "email")
}

func name(args []string) (string, error) {
	return argValue(args, "name")
}

func phone(args []string) (string, error) {
	return argValue(args, "phone")
}

func argValue(args []string, argName string) (string, error) {
	for i := range args {
		prefix := argName + ":"
		if strings.HasPrefix(args[i], prefix) {
			return strings.TrimPrefix(args[i], prefix), nil
		}
	}
	return "", errors.New(argName + "not provided")
}

func init() {
	customerCmd.AddCommand(customerAddCmd)
}
