package flag

import (
	"flag"
)

type Flag struct {
	Fresh       bool
	Seeder      bool
	SeederModel string
}

var FlagVars = &Flag{}

func Parse(args []string) error {
	flagSet := flag.NewFlagSet("paskihub-be", flag.ContinueOnError)
	fresh := flagSet.Bool("fresh", false, "Dropping all database tables before running new migration")
	seeder := flagSet.Bool("seeder", false, "Inserting all seeders to database")
	seederModel := flagSet.String("model", "", "Specify seeder model to run\nMust be used with flag seeder\nAvailable models can be looked up at pgsql_conn.go file func Seeder")

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	FlagVars = &Flag{
		Fresh:       *fresh,
		Seeder:      *seeder,
		SeederModel: *seederModel,
	}
	return nil
}
