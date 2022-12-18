package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"muzzammil.xyz/jsonc"
)

// runs the given command while attaching stdin, stdout and stderr
func execCmd(command []string, capture bool) (string, error) {
	var stdout []byte
	var err error

	cwd, _ := os.Getwd()
	cmd := exec.Command(command[0], command[1:]...)
	log.Info().Str("workdir", cwd).Str("command", cmd.String()).Send()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if capture {
		stdout, err = cmd.Output()
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	return strings.TrimSpace(string(stdout)), err
}

// return the md5 hash for a string
func md5sum(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash
}

type resolve func(string) string

type matchStore struct {
	Old string
	New string
}

// replace all matches of the regexp in the string with the resolved value
func replace(s string, r *regexp.Regexp, f resolve) string {
	if r.MatchString(s) {
		matches := r.FindAllStringSubmatchIndex(s, -1)
		matchesStore := []matchStore{}
		for _, match := range matches {
			matchesStore = append(matchesStore, matchStore{s[match[0]:match[1]], f(s[match[2]:match[3]])})
		}
		for _, store := range matchesStore {
			s = strings.Replace(s, store.Old, store.New, 1)
		}
	}

	return s
}

// resolve ${localEnv:VARIABLE_NAME}
func resolveLocalEnv(s string) string {
	regexpLocalEnv := regexp.MustCompile(`\${localEnv:([[:word:]]+)}`)

	return replace(s, regexpLocalEnv, os.Getenv)
}

// resolve ${containerEnv:VARIABLE_NAME}
func resolveContainerEnv(e Engine, s string) string {
	regexpLocalEnv := regexp.MustCompile(`\${containerEnv:([[:word:]]+)}`)

	return replace(s, regexpLocalEnv, e.ResolveEnv)
}

// PRERUN UTILS

func (d *DevContainer) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

	switch rootVerbose {
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
}

func (d *DevContainer) ParseConfig() {
	// return JSONC as JSON
	_, j, err := jsonc.ReadFromFile(filepath.Join(rootConfigDir, "devcontainer.json"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read devcontainer settings")
	}

	// pass data to viper
	d.Config = viper.New()
	d.Config.SetConfigType("json")
	if err := d.Config.ReadConfig(bytes.NewBuffer(j)); err != nil {
		log.Fatal().Err(err).Msg("cannot read json")
	}

	log.Debug().Str("devcontainer", fmt.Sprintf("%+v", d.Config)).Send()
}

func (d *DevContainer) SetAliases() {
	// set aliases
	d.Config.RegisterAlias("dockerfile", "build.dockerfile")
	d.Config.RegisterAlias("context", "build.context")
}

func (d *DevContainer) NormalizeTypes() {
	// convert multi-types settings to one type
	if value, ok := d.Config.Get("build.cacheFrom").(string); ok {
		d.Config.Set("build.cacheFrom", []string{value})
	}
	if value, ok := d.Config.Get("dockerComposeFile").(string); ok {
		d.Config.Set("dockerComposeFile", []string{value})
	}
}

func (d *DevContainer) SetDefaults() {
	// set defaults values
	d.ConfigDir = rootConfigDir
	d.WorkingDirectoryPath, _ = os.Getwd()
	d.WorkingDirectoryName = filepath.Base(d.WorkingDirectoryPath)
	d.Config.SetDefault("build.context", ".")
	d.Config.SetDefault("name", d.WorkingDirectoryName)
	d.Config.SetDefault("overrideCommand", true)
	d.Config.SetDefault("updateRemoteUsedUID", true)
	d.Config.SetDefault("workspaceFolder", "/workspace")
	d.Config.SetDefault("workspaceMount", "type=bind,source="+d.WorkingDirectoryPath+",target="+d.Config.GetString("workspaceFolder")+",consistency=cached")
}

func (d *DevContainer) CheckConfig() {
	// check required and conflicting settings
	if !d.Config.IsSet("image") && !d.Config.IsSet("build.dockerfile") && !d.Config.IsSet("dockerComposeFile") {
		log.Fatal().Msg("one of these settings is required: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if (d.Config.IsSet("image") && d.Config.IsSet("build.dockerfile")) ||
		(d.Config.IsSet("image") && d.Config.IsSet("dockerComposeFile")) ||
		(d.Config.IsSet("build.dockerfile") && d.Config.IsSet("dockerComposeFile")) {
		log.Fatal().Msg("one of these settings conflicts with another one: 'image', 'build.dockerfile', 'dockerComposeFile'")
	}
	if d.Config.IsSet("dockerComposeFile") && !d.Config.IsSet("service") {
		log.Fatal().Msg("'service' setting is required when using 'dockerComposeFile'")
	}
}

func (d *DevContainer) SetEngine() {
	// determine container engine
	switch {
	case d.Config.IsSet("image") || d.Config.IsSet("build.dockerfile"):
		d.Engine = &Docker{}
	case d.Config.IsSet("dockerComposeFile"):
		d.Engine = &DockerCompose{}
	default:
		log.Fatal().Msg("cannot determine devcontainer engine")
	}

	// initialize engine
	if err := d.Engine.Init(d); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
	log.Debug().Str("engine", fmt.Sprintf("%+v", d.Engine)).Send()
}

func (d *DevContainer) ResolveVars() {
	// resolve:
	// - ${localEnv:VARIABLE_NAME}
	// - ${localWorkspaceFolder}
	// - ${containerWorkspaceFolder}
	// - ${localWorkspaceFolderBasename}
	// - ${containerWorkspaceFolderBasename}
	// cf. https://containers.dev/implementors/json_reference/#variables-in-devcontainerjson
	keys := []string{
		"build.args", "build.cacheFrom", "build.context", "build.dockerfile",
		"build.target", "containerEnv", "containerUser", "dockerComposeFile",
		"forwardPorts", "image", "mounts", "name", "remoteEnv", "remoteUser",
		"runArgs", "runServices", "service", "workspaceFolder", "workspaceMount",
	}
	for _, key := range keys {
		switch d.Config.Get(key).(type) {
		case string:
			// resolve strings
			d.Config.Set(key, resolveLocalEnv(d.Config.GetString(key)))
			d.Config.Set(key, strings.ReplaceAll(d.Config.GetString(key), "${localWorkspaceFolder}", d.WorkingDirectoryPath))
			d.Config.Set(key, strings.ReplaceAll(d.Config.GetString(key), "${containerWorkspaceFolder}", d.WorkingDirectoryPath))
			d.Config.Set(key, strings.ReplaceAll(d.Config.GetString(key), "${localWorkspaceFolderBasename}", d.WorkingDirectoryName))
			d.Config.Set(key, strings.ReplaceAll(d.Config.GetString(key), "${containerWorkspaceFolderBasename}", d.WorkingDirectoryName))
		case []interface{}:
			// resolve slices of strings
			d.Config.Set(key, lo.Map(d.Config.GetStringSlice(key), func(v string, _ int) string { return resolveLocalEnv(v) }))
			d.Config.Set(key, lo.Map(d.Config.GetStringSlice(key), func(v string, _ int) string {
				return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
			}))
			d.Config.Set(key, lo.Map(d.Config.GetStringSlice(key), func(v string, _ int) string {
				return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.WorkingDirectoryPath)
			}))
			d.Config.Set(key, lo.Map(d.Config.GetStringSlice(key), func(v string, _ int) string {
				return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
			}))
			d.Config.Set(key, lo.Map(d.Config.GetStringSlice(key), func(v string, _ int) string {
				return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", d.WorkingDirectoryName)
			}))
		case map[string]interface{}:
			// resolve maps of strings
			d.Config.Set(key, lo.MapValues(d.Config.GetStringMapString(key), func(v string, _ string) string { return resolveLocalEnv(v) }))
			d.Config.Set(key, lo.MapValues(d.Config.GetStringMapString(key), func(v string, _ string) string {
				return strings.ReplaceAll(v, "${localWorkspaceFolder}", d.WorkingDirectoryPath)
			}))
			d.Config.Set(key, lo.MapValues(d.Config.GetStringMapString(key), func(v string, _ string) string {
				return strings.ReplaceAll(v, "${containerWorkspaceFolder}", d.WorkingDirectoryPath)
			}))
			d.Config.Set(key, lo.MapValues(d.Config.GetStringMapString(key), func(v string, _ string) string {
				return strings.ReplaceAll(v, "${localWorkspaceFolderBasename}", d.WorkingDirectoryName)
			}))
			d.Config.Set(key, lo.MapValues(d.Config.GetStringMapString(key), func(v string, _ string) string {
				return strings.ReplaceAll(v, "${containerWorkspaceFolderBasename}", d.WorkingDirectoryName)
			}))
		}
	}
}
