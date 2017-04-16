package main

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
		Build: []Build{},
	}
}

func (build *BuildFile) load(path string) *BuildFile {
	build.File = path
	build.Directory = filepath.Dir(path)

	data, err := ioutil.ReadFile(build.File)

	if err != nil {
		log.Fatal(err)
	}

	yaml.Unmarshal(data, build)

	return build.determine()
}

func (build *BuildFile) determine() *BuildFile {
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

	build.Verify.Include = buildPaths(build.Root, build.Cwd, build.Verify.Include)
	build.Verify.Exclude = buildPaths(build.Root, build.Cwd, build.Verify.Exclude)

	build.Package.Include = buildPaths(build.Root, build.Cwd, build.Package.Include)
	build.Package.Exclude = buildPaths(build.Root, build.Cwd, build.Package.Exclude)

	build.loadConstants()
	build.executeConstraints()

	return build
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
			}).Debug("Hit cache")

			build.Constant[i].Result = value.(string)

			continue;
		} else {
			log.WithFields(log.Fields{
				"constant": constant.Constant,
				"cmd":      constant.Command,
			}).Debug("Doesn't hit cache")
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

		build.Constant[i].Result = value.(string)

		c.Set(constant.Command, value.(string), cache.DefaultExpiration)
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
			givenVersion, _ := version.NewVersion(args[0].(string))
			constraint, _ := version.NewConstraint(args[1].(string))

			return (bool)(constraint.Check(givenVersion)), nil
		},
	}

	for i, constraint := range build.Constraint {
		expression, err := govaluate.NewEvaluableExpressionWithFunctions(constraint.Condition, functions)
		parameters := make(map[string]interface{})
		parameters["BOB_VERSION"] = APP_VERSION
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
			build.Constraint[i].Hash, _ = hashValue(constraint.Name, result.(string))
			break
		case bool:
			build.Constraint[i].Result = result.(bool)
			build.Constraint[i].Hash, _ = hashValue(constraint.Name, strconv.FormatBool(result.(bool)))
			break
		}
	}

	return nil
}
