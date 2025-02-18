package cli

import (
	"strings"

	"github.com/mitchellh/cli"
)

type Account struct {
	UI cli.Ui
}

// MarkDown implements cli.MarkDown interface
func (a *Account) MarkDown() string {
	items := []string{
		"# Account",
		"The ```account``` command groups actions to interact with accounts:",
		"- [```account new```](./account_new.md): Create a new account in the Eirene client.",
		"- [```account list```](./account_list.md): List the wallets in the Eirene client.",
		"- [```account import```](./account_import.md): Import an account to the Eirene client.",
	}

	return strings.Join(items, "\n\n")
}

// Help implements the cli.Command interface
func (a *Account) Help() string {
	return `Usage: eirene account <subcommand>

  This command groups actions to interact with accounts.
  
  List the running deployments:

    $ eirene account new
  
  Display the status of a specific deployment:

    $ eirene account import
    
  List the imported accounts in the keystore:
    
    $ eirene account list`
}

// Synopsis implements the cli.Command interface
func (a *Account) Synopsis() string {
	return "Interact with accounts"
}

// Run implements the cli.Command interface
func (a *Account) Run(args []string) int {
	return cli.RunResultHelp
}
