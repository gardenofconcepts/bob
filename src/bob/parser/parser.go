package parser

import (
	"bytes"
	"github.com/Knetic/govaluate"
	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-version"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"github.com/patrickmn/go-cache"
	"time"
	"bob/path"
	"bob/hash"
	"bob/builder"
)

var c = cache.New(15*time.Minute, 60*time.Minute)

func Parser() *BuildFile {
	return &BuildFile{
		Priority: 0,
		Name:     "Unknown",
		Verify: Verify{
			Include: []string{},
			Exclude: []string{},
		},
		Package: Package{
			Include: []string{},
			Exclude: []string{},
		},
		Build: []builder.Build{},
	}
}

func (build *BuildFile) Load(path string) *BuildFile {
	build.File = path
	build.Directory = filepath.Dir(path)

	data, err := ioutil.ReadFile(build.File)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, build)

	build.determine()
	build.loadConstants();
	build.executeConstraints();

	return build
}

func (build *BuildFile) determine() error {
	if len(build.Cwd) == 0 {
		build.Cwd = build.Directory
	} else if !filepath.IsAbs(build.Cwd) {
		build.Cwd = filepath.Join(build.Directory, build.Cwd)
		build.Cwd, _ = filepath.Abs(build.Cwd)
	}

	if len(build.Root) == 0 {
		build.Root = build.Directory
	} else if !filepath.IsAbs(build.Root) {
		build.Root = filepath.Join(build.Directory, build.Root)
		build.Root, _ = filepath.Abs(build.Root)
	}

	build.Verify.Include = path.MakePathsRelative(build.Root, build.Cwd, build.Verify.Include)
	build.Verify.Exclude = path.MakePathsRelative(build.Root, build.Cwd, build.Verify.Exclude)

	build.Package.Include = path.MakePathsRelative(build.Root, build.Cwd, build.Package.Include)
	build.Package.Exclude = path.MakePathsRelative(build.Root, build.Cwd, build.Package.Exclude)

	return nil
}

func (build *BuildFile) loadConstants() error {
	for i, constant := range build.Constant {
		log.WithFields(log.Fields{
			"constant": constant.Constant,
			"cmd":      constant.Command,
			"cwd":      build.Cwd,
		}).Info("Checking constant")

		value, found := c.Get(constant.Command)
		if found {
			log.WithFields(log.Fields{
				"constant": constant.Constant,
				"cmd":      constant.Command,
				"value":    value.(string),
			}).Debug("Hit cache")

			build.Constant[i].Result = value.(string)

			continue
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd := exec.Command("/bin/bash", "-c", constant.Command)
		cmd.Dir = build.Cwd
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		if err != nil {
			log.Fatal("Error while running command", err, stdout.String(), stderr.String())

			return err
		}

		log.Debugf("Result: %q\n", stdout.String())

		value = strings.Trim(stdout.String(), " \n")

		log.WithFields(log.Fields{
			"constant": constant.Constant,
			"cmd":      constant.Command,
			"result":   value.(string),
		}).Debug("Doesn't hit cache")

		build.Constant[i].Result = value.(string)

		c.Set(constant.Command, value.(string), cache.NoExpiration)
	}

	return nil
}

func (build *BuildFile) executeConstraints() error {

	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))

			return (float64)(length), nil
		},
		"version_compare": func(args ...interface{}) (interface{}, error) {
			givenVersion, e := version.NewVersion(args[0].(string))

			if e != nil {
				log.WithFields(log.Fields{
					"version": givenVersion,
				}).Fatal(e)
			}

			constraint, e := version.NewConstraint(args[1].(string))

			if e != nil {
				log.WithFields(log.Fields{
					"version": givenVersion,
				}).Fatal(e)
			}

			return (bool)(constraint.Check(givenVersion)), nil
		},
	}

	for i, constraint := range build.Constraint {
		expression, err := govaluate.NewEvaluableExpressionWithFunctions(constraint.Condition, functions)
		parameters := make(map[string]interface{})
		//parameters["BOB_VERSION"] = '0'
		parameters["OS"] = runtime.GOOS
		parameters["ARCH"] = runtime.GOARCH

		for _, constant := range build.Constant {
			parameters[constant.Constant] = constant.Result
		}

		result, err := expression.Evaluate(parameters)

		if err != nil {
			log.Fatal("Condition failed", result)

			return err
		}

		switch result.(type) {
		case string:
			build.Constraint[i].Result = true
			build.Constraint[i].ResultString = result.(string)
			build.Constraint[i].Hash, _ = hash.Value(constraint.Name, result.(string))
			break
		case bool:
			build.Constraint[i].Result = result.(bool)
			build.Constraint[i].Hash, _ = hash.Value(constraint.Name, strconv.FormatBool(result.(bool)))
			break
		}
	}

	return nil
}
