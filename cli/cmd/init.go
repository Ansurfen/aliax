// Copyright 2025 The Aliax Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"aliax/internal/aos"
	"aliax/internal/aslices"
	bashast "aliax/internal/ast/bash"
	psast "aliax/internal/ast/powershell"
	"aliax/internal/cfg"
	"aliax/internal/shell"
	"aliax/internal/style"
	bashtoken "aliax/internal/token/bash"
	token "aliax/internal/token/powershell"
	"errors"
	"path/filepath"
	"runtime"
	"time"

	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"aliax/internal/log"
	"github.com/ansurfen/globalenv"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// initCmdParameter stores parameters for the "init" command.
type initCmdParameter struct {
	global   bool
	force    bool
	verbose  bool
	template string
	save     bool
}

var (
	initParameter initCmdParameter
	initCmd       = &cobra.Command{
		Use:   "init",
		Short: "Initialize the aliax workspace and generate execution scripts",
		Long: `The "init" command scans the aliax configuration file and generates necessary execution scripts.
It creates platform-specific scripts in the "run-scripts" directory for alias commands and extensions.
If the --global (-g) flag is set, it applies configurations globally.`,
		Example: "  aliax init\n  aliax init --global",
		Run: func(cmd *cobra.Command, args []string) {
			var start time.Time
			if initParameter.verbose {
				start = time.Now()
				log.SetLevel(log.DebugLevel)
			}
			if len(initParameter.template) > 0 {
				config = filepath.Join(aos.TemplatePath, initParameter.template+".yaml")
			}
			var file cfg.Aliax
			err := aos.ReadYAML(config, &file)
			if err != nil {
				log.WithError(err).Fatalf("fail to parse file")
			}
			if initParameter.save {
				path := filepath.Join(aos.TemplatePath, filepath.Base(config))
				output, err := aos.Create(path)
				if err != nil {
					log.WithError(err).Fatal("fail to create file")
				}
				err = yaml.NewEncoder(output).Encode(file)
				if err != nil {
					log.WithError(err).Fatal("backuping template")
				}
			}
			err = aos.MkdirAll("run-scripts/bash", 0755)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					log.WithError(err).Warn("making run-scripts directory")
				} else {
					log.WithError(err).Fatal("making run-scripts directory")
				}
			}

			extBins := map[string]string{}

			for name, ext := range file.Extend {
				if len(ext.Bin) == 0 {
					ext.Bin, err = shell.LookPath(name)
					if err != nil {
						log.WithError(err).WithField("suggestion", "please make sure the executable exists,\neither by adding the bin field to the YAML to manually indicate the path,\nor by moving the extend to command").Fatal("looking path")
					}
				} else {
					if aos.IsWindows {
						ext.Bin = envResolver.apply(ext.Bin, func(matched string) string {
							return fmt.Sprintf("$env:%s", matched)
						})
					} else {
						ext.Bin = envResolver.apply(ext.Bin, func(matched string) string {
							return fmt.Sprintf("$%s", matched)
						})
					}
				}
				extBins[name] = ext.Bin
				file.Extend[name].Bin = ext.Bin
			}

			builder := &runScriptsBuilder{}

			if len(file.RunPath) == 0 {
				file.RunPath = "run-scripts"
			}

			if len(file.Executable) != 0 {
				executable = file.Executable
			}

			err = builder.generateScriptExtension(file.RunPath, file.Extend)
			if err != nil {
				log.WithError(err).Fatal("generating extension script")
			}

			err = builder.generateCommand(file.RunPath, file.Command)
			if err != nil {
				log.WithError(err).Fatal("generating command script")
			}

			if initParameter.global {
				if err = setGlobal(); err != nil {
					log.WithError(err).Fatal("setting for the global")
				} else {
					data, err := getAliaxPath()
					if err != nil && !errors.Is(err, errAliaxPathNotFound) {
						log.WithError(err).Fatal("setting environment")
					}
					for k, v := range extBins {
						data[k] = v
					}
					output, err := setAliaxPath(data)
					if err != nil {
						log.WithField("output", string(output)).WithError(err).Fatal("setting environment")
					}
					log.Info("setting for the global")
				}
			}
			if initParameter.verbose {
				duration := time.Since(start)
				log.Debugf("executing succeded after %s", style.Bold(duration.String()))
			}
		},
	}
)

type runScriptsBuilder struct{}

func (s *runScriptsBuilder) generateScriptExtension(dir string, cmds map[string]*cfg.Command) error {
	for name, cmd := range cmds {
		_, err := s.generatePowershellExtension(dir, name, cmd)
		if err != nil {
			log.WithError(err).Fatal("generating powershell script")
		}

		sh, err := s.generateBashExtension(dir, name, cmd)
		if err != nil {
			log.WithError(err).Fatal("generating bash script")
		}

		target, err := filepath.Abs(sh.filename())
		if err != nil {
			log.WithError(err).Fatal("invalid path")
		}
		link := filepath.Join(filepath.Dir(target), "bash", strings.TrimSuffix(filepath.Base(target), ".sh"))
		err = s.createSymbolLink(target, link)
		if err != nil {
			log.WithError(err).WithField("suggestion", suggestionSymbolLinkError).Fatalf("creating symbol link for %s", target)
		}
	}
	return nil
}

func (s *runScriptsBuilder) generatePowershellExtension(dir, name string, cmd *cfg.Command) (*psScriptBuilder, error) {
	psBuilder, err := newPsScriptBuilder(dir, name, cmd)
	if err != nil {
		log.WithError(err).Error("creating powershell builder")
		return nil, err
	}
	defer psBuilder.close()
	psBuilder.node.Append(psast.AssignStatement(
		psast.RefRaw(executable),
		psast.String(psBuilder.cmd.Bin),
	))

	psBuilder.node.Append(psBuilder.generateExtension(name, name, 0, cmd)...)

	psBuilder.node.Append(&psast.CallStmt{
		Op:   token.BITAND,
		Func: psast.RefRaw(executable),
		Recv: []psast.Expr{psast.RefRaw("args")},
	})
	psast.Print(psBuilder.node, psBuilder.file)
	return psBuilder, nil
}

func (s *runScriptsBuilder) generateBashExtension(dir, name string, cmd *cfg.Command) (*bashScriptBuilder, error) {
	bashBuiler, err := newBashScriptBuilder(dir, name, cmd)
	if err != nil {
		log.WithError(err).Error("creating bash builder")
		return nil, err
	}
	defer bashBuiler.close()
	bashBuiler.node.Append(bashast.AssignStatement(
		bashast.Identifier(executable),
		bashast.String(cmd.Bin),
	), bashast.RawStmt(`args=("$@")`))

	bashBuiler.node.Append(bashBuiler.generateExtension(name, name, 0, cmd)...)

	bashBuiler.node.Append(&bashast.CallStmt{
		Func: &bashast.RefExpr{X: bashast.Identifier(executable)},
		Recv: []bashast.Expr{bashast.Raw(`"${args[@]}"`)},
	})

	bashast.Print(bashBuiler.node, bashBuiler.file)
	return bashBuiler, nil
}

func (s *runScriptsBuilder) createSymbolLink(target, link string) error {
	var err error
	if ok, _ := aos.Exist(link); ok && initParameter.force {
		log.WithField("link", link).Info("deleting old symbol link")
		err = aos.Remove(link)
		if err != nil {
			log.WithError(err).Error("fail to delete file")
			return err
		}
	}

	if aos.IsWindows {
		err = shell.Run("cmd", "/C", "mklink", link, target)
	} else {
		err = shell.Run("ln", "-s", target, link)
	}

	if err != nil {
		return err
	}
	return nil
}

func (s *runScriptsBuilder) generateCommand(dir string, cmds map[string]*cfg.Command) error {
	for name, cmd := range cmds {
		err := cmd.Preload(name)
		if err != nil {
			return err
		}

		_, err = s.generatePowershellCommand(dir, name, cmd)
		if err != nil {
			log.WithError(err).Fatal("generating powershell script")
		}

		sh, err := s.generateBashCommand(dir, name, cmd)
		if err != nil {
			log.WithError(err).Fatal("generating bash script")
		}

		target, err := filepath.Abs(sh.filename())
		if err != nil {
			log.WithError(err).Fatal("invalid path")
		}
		link := filepath.Join(filepath.Dir(target), "bash", strings.TrimSuffix(filepath.Base(target), ".sh"))
		err = s.createSymbolLink(target, link)
		if err != nil {
			log.WithError(err).WithField("suggestion", suggestionSymbolLinkError).Fatalf("creating symbol link for %s", target)
		}

	}
	return nil
}

func (s *runScriptsBuilder) generatePowershellCommand(dir, name string, cmd *cfg.Command) (*psScriptBuilder, error) {
	psBuilder, err := newPsScriptBuilder(dir, name, cmd)
	if err != nil {
		log.WithError(err).Error("creating powershell builder")
		return nil, err
	}
	defer psBuilder.close()

	psBuilder.node.Append(psBuilder.generateCommand(name, name, 0, cmd)...)

	psast.Print(psBuilder.node, psBuilder.file)
	return psBuilder, nil
}

func (s *runScriptsBuilder) generateBashCommand(dir, name string, cmd *cfg.Command) (*bashScriptBuilder, error) {
	bashBuiler, err := newBashScriptBuilder(dir, name, cmd)
	if err != nil {
		log.WithError(err).Error("creating bash builder")
		return nil, err
	}
	defer bashBuiler.close()
	bashBuiler.node.Append(bashast.RawStmt(`args=("$@")`))

	bashBuiler.node.Append(bashBuiler.generateCommand(name, name, 0, cmd)...)

	// bashBuiler.node.Append(bashast.CallStatement("cat", fmt.Sprintf("<<EOF\n%s\nEOF", cmd.HelpCmd(name))))

	bashast.Print(bashBuiler.node, bashBuiler.file)
	return bashBuiler, nil
}

type flagType uint8

const (
	flagTypeString flagType = iota
	flagTypeBool
)

type bashScriptBuilder struct {
	file  *os.File
	cmd   *cfg.Command
	node  *bashast.File
	ident string
}

func newBashScriptBuilder(dir, name string, cmd *cfg.Command) (*bashScriptBuilder, error) {
	filename := filepath.Join(dir, name+".sh")

	fp, err := aos.Create(filename)
	if err != nil {
		if errors.Is(err, os.ErrExist) && !initParameter.force {
			log.WithError(err).
				WithField("file", filename).
				WithField("suggestion", suggestionCleanWorkspace).
				Fatal("file already exist")
		} else {
			return nil, err
		}
	}

	builder := &bashScriptBuilder{
		file:  fp,
		node:  &bashast.File{},
		cmd:   cmd,
		ident: name,
	}

	builder.node.Append(
		bashast.Docs("!/bin/bash"),
		bashast.Docs(copyright),
		bashast.RawStmt("set -e"))

	return builder, nil
}

func (b *bashScriptBuilder) filename() string {
	return b.file.Name()
}

func (b *bashScriptBuilder) close() error {
	return b.file.Close()
}

func (b *bashScriptBuilder) generateExtension(ident, cmdName string, level int, cmd *cfg.Command) []bashast.Stmt {
	subCommand := b.collectSubCommand4Extend(ident, level, cmd)

	if level > 0 {
		ifStmt := bashast.IfStatement()
		ifStmt.Cond = &bashast.BinaryExpr{X: &bashast.RefExpr{X: &bashast.IndexExpr{X: &bashast.Ident{Name: "args"}, Key: &bashast.BasicExpr{Kind: bashtoken.NUMBER, Value: "0"}}}, Op: bashtoken.EQ, Y: &bashast.BasicExpr{Kind: bashtoken.STRING, Value: cmdName}}
		ifStmt.Body.List = append(ifStmt.Body.List, bashast.RawStmt(`args=("${args[@]:1}")`))
		ifStmt.Body.List = append(ifStmt.Body.List, b.buildBlockSmt(subCommand, ident, cmd)...)
		return []bashast.Stmt{ifStmt}
	}

	return b.buildBlockSmt(subCommand, ident, cmd)
}

func (b *bashScriptBuilder) generateCommand(ident, cmdName string, level int, cmd *cfg.Command) []bashast.Stmt {
	if level == 0 {
		cmd.SetName(cmdName)
	} else {
		cmd.SetName(strings.Join(strings.Split(ident, "_"), " "))
	}
	subCommand := b.collectSubCommand4Command(ident, level, cmd)

	if level > 0 {
		ifStmt := bashast.IfStatement()
		ifStmt.Cond = &bashast.BinaryExpr{X: &bashast.RefExpr{X: &bashast.IndexExpr{X: &bashast.Ident{Name: "args"}, Key: &bashast.BasicExpr{Kind: bashtoken.NUMBER, Value: "0"}}}, Op: bashtoken.EQ, Y: &bashast.BasicExpr{Kind: bashtoken.STRING, Value: cmdName}}
		ifStmt.Body.List = append(ifStmt.Body.List, bashast.RawStmt(`args=("${args[@]:1}")`))
		ifStmt.Body.List = append(ifStmt.Body.List, b.buildBlockStmt4Command(subCommand, ident, cmd)...)
		return []bashast.Stmt{ifStmt}
	}

	return b.buildBlockStmt4Command(subCommand, ident, cmd)
}

func (b *bashScriptBuilder) buildBlockStmt4Command(subCommand []bashast.Stmt, ident string, cmd *cfg.Command) (bs []bashast.Stmt) {
	bs = b.buildBlockSmt(subCommand, ident, cmd)
	if !cmd.DisableHelp {
		bs = append(bs, bashast.CallStatement("cat",
			fmt.Sprintf("<<EOF\n%s\nEOF", cmd.HelpCmd(cmd.Name()))),
			bashast.CallStatement("exit"))
	}
	return bs
}

func (b *bashScriptBuilder) collectSubCommand4Extend(ident string, level int, cmd *cfg.Command) []bashast.Stmt {
	subCommand := []bashast.Stmt{}
	for name, subcmd := range cmd.Command {
		subcmd.SetName(name)
		subCommand = append(subCommand, b.generateExtension(fmt.Sprintf("%s_%s", ident, name), name, level+1, subcmd)...)
	}
	return subCommand
}

func (b *bashScriptBuilder) collectSubCommand4Command(ident string, level int, cmd *cfg.Command) []bashast.Stmt {
	subCommand := []bashast.Stmt{}
	for name, subcmd := range cmd.Command {
		subCommand = append(subCommand, b.generateCommand(fmt.Sprintf("%s_%s", ident, name), name, level+1, subcmd)...)
	}
	return subCommand
}

func (b *bashScriptBuilder) buildArgsStmt(ident string, subCommand, bs []bashast.Stmt) []bashast.Stmt {
	for i, sc := range subCommand {
		if i == 0 {
			if i != len(subCommand)-1 {
				bs = append(bs, bashast.AssignStatement(bashast.Identifier(fmt.Sprintf("temp_args_%s", ident)), bashast.Raw(`("${args[@]}")`)))
			}
			bs = append(bs, sc)
		} else {
			bs = append(bs, bashast.AssignStatement(bashast.Identifier("args"), bashast.Raw(fmt.Sprintf(`("${temp_args_%s[@]}")`, ident))))
			bs = append(bs, sc)
		}
	}
	return bs
}

func (b *bashScriptBuilder) buildFlagDict(ident string, bs []bashast.Stmt, cmd *cfg.Command) (typeDict map[string]flagType, stmts []bashast.Stmt) {
	typeDict = make(map[string]flagType)
	for _, flag := range cmd.Flags {
		flagIdent := fmt.Sprintf("%s_%s", ident, flag.Name)
		switch flag.Type {
		case "string":
			typeDict[flagIdent] = flagTypeString
			bs = append(bs, bashast.AssignStatement(bashast.Identifier(flagIdent), bashast.String("")))
		case "bool":
			typeDict[flagIdent] = flagTypeBool
			bs = append(bs, bashast.AssignStatement(bashast.Identifier(flagIdent), bashast.FALSE))
		}
	}
	return typeDict, bs
}

func (b *bashScriptBuilder) buildBlockSmt(subCommand []bashast.Stmt, ident string, cmd *cfg.Command) (bs []bashast.Stmt) {
	bs = b.buildArgsStmt(ident, subCommand, bs)

	typeDict, bs := b.buildFlagDict(ident, bs, cmd)
	bs = append(bs, bashast.RawStmt("non_matched_args=()"))

	if len(cmd.Flags) > 0 {
		forStmt := bashast.ForStatement(
			bashast.BinaryExpression(bashast.Identifier("i"), bashtoken.ASSIGN, bashast.Number(0)),
			bashast.BinaryExpression(bashast.Identifier("i"), bashtoken.LT, bashast.Raw("${#args[@]}")),
			bashast.IncDecExpression(bashast.Identifier("i"), true),
		)
		bs = append(bs, forStmt)
		forStmt.Body.Append(b.collectFlagStmt(ident, cmd))

		bs = b.buildMatchStmt(ident, cmd, bs, typeDict)
	} else {
		for _, matchCase := range cmd.Match {
			matchCase.Run = indexResolver.apply(matchCase.Run, func(matched string) string {
				i, _ := strconv.Atoi(matched)
				i--
				return fmt.Sprintf(`"$($args[%d])"`, i)
			})
			matchCase.Run = envResolver.apply(matchCase.Run, func(matched string) string {
				return fmt.Sprintf("$env:%s", matched)
			})
			switch pattern := matchCase.Pattern.(type) {
			case string:
				if pattern == "_" || len(pattern) == 0 {
					bs = append(bs, bashast.CallStatement(matchCase.Run), bashast.CallStatement("exit"))
				}
			}
		}
	}
	return
}

func (b *bashScriptBuilder) collectFlagStmt(ident string, cmd *cfg.Command) bashast.Stmt {
	switchStmt := &bashast.SwitchStmt{
		Cond: bashast.String("${args[i]}"),
		Default: &bashast.CaseStmt{
			Body: &bashast.BlockStmt{
				List: []bashast.Stmt{
					bashast.RawStmt("non_matched_args+=(\"${args[i]}\")"),
				},
			},
		},
	}

	for _, flag := range cmd.Flags {
		flagIdent := fmt.Sprintf("%s_%s", ident, flag.Name)
		alias := []string{}
		if len(flag.Alias) == 0 {
			alias = append(alias, flag.Name)
		} else {
			alias = append(alias, flag.Alias...)
		}
		for i, a := range alias {
			alias[i] = regexp.QuoteMeta(a)
		}
		rule := strings.Join(alias, "|")
		caseStmt := bashast.CaseStatement(bashast.Identifier(rule))
		switchStmt.Cases = append(switchStmt.Cases, caseStmt)
		switch flag.Type {
		case "string":
			caseStmt.Body.Append(
				bashast.AssignStatement(bashast.Identifier(flagIdent), bashast.String("${args[i+1]}")),
				bashast.RawStmt("((i++))"))
		case "bool":
			caseStmt.Body.Append(bashast.AssignStatement(bashast.Identifier(flagIdent), bashast.TRUE))
			// caseStmt.Body.List = append(caseStmt.Body.List, &psast.ExprStmt{
			// 	X: &psast.IncDecExpr{
			// 		X:  &psast.RefExpr{X: &psast.Ident{Name: "i"}},
			// 		Op: token.Inc,
			// 	},
			// })
		}
	}
	return switchStmt
}

func (b *bashScriptBuilder) buildMatchStmt(ident string, cmd *cfg.Command, bs []bashast.Stmt, typeDict map[string]flagType) []bashast.Stmt {
	type sortedMatchCase struct {
		weight int
		names  []string
		body   string
	}

	match := []sortedMatchCase{}
	var defaultMatchCase *sortedMatchCase
	for _, matchCase := range cmd.Match {
		if len(matchCase.Platform) > 0 && matchCase.Platform != "bash" {
			continue
		}
		names := []string{}
		matchCase.Run = indexResolver.apply(matchCase.Run, func(matched string) string {
			i, _ := strconv.Atoi(matched)
			i--
			return fmt.Sprintf(`"$($args[%d])"`, i)
		})
		matchCase.Run = namedResolver.apply(matchCase.Run, func(matched string) string {
			// shellcheck: Double quote to prevent globbing and word splitting.
			return fmt.Sprintf("$%s_%s", ident, matched)
		})
		matchCase.Run = envResolver.apply(matchCase.Run, func(matched string) string {
			return fmt.Sprintf("$%s", matched)
		})
		switch pattern := matchCase.Pattern.(type) {
		case string:
			if pattern == "_" || len(pattern) == 0 {
				defaultMatchCase = &sortedMatchCase{
					body: matchCase.Run,
				}
				continue
			}
			names = append(names, fmt.Sprintf("%s_%s", ident, pattern))
		case []any:
			for _, v := range pattern {
				if v, ok := v.(string); ok {
					names = append(names, fmt.Sprintf("%s_%s", ident, v))
				}
			}
		}
		if len(names) != 0 {
			match = append(match, sortedMatchCase{weight: len(names), names: names, body: matchCase.Run})
		}
	}

	if len(match) == 0 {
		return bs
	}

	sort.Slice(match, func(i, j int) bool {
		return match[i].weight > match[j].weight
	})

	matchStmt := bashast.IfStatement()
	bs = append(bs, matchStmt)
	for i, c := range match {
		var cases bashast.Expr
		for _, name := range c.names {
			var cond bashast.Expr
			switch typeDict[name] {
			case flagTypeString:
				cond = bashast.Raw(fmt.Sprintf(`-n "$%s"`, name))
			case flagTypeBool:
				cond = bashast.BinaryExpression(
					bashast.RefRaw(name),
					bashtoken.EQ,
					bashast.TRUE,
				)
			}
			if cases == nil {
				cases = cond
			} else {
				cases = bashast.BinaryExpression(cases, bashtoken.AND, cond)
			}
		}

		matchStmt.Cond = cases
		if len(c.body) > 0 {
			lines := strings.Split(c.body, "\n")
			for i, line := range lines {
				if i == len(lines)-1 && len(line) == 0 {
					continue
				}
				matchStmt.Body.Append(bashast.CallStatement(line))
			}
		}
		matchStmt.Body.Append(bashast.CallStatement("exit"))
		if i != len(match)-1 {
			ifstmt := bashast.IfStatement()
			matchStmt.Else = ifstmt
			matchStmt = ifstmt
		}
	}

	if defaultMatchCase != nil {
		matchStmt.Else = bashast.BlockStatement(bashast.RawStmt(defaultMatchCase.body))
	}
	return bs
}

type psScriptBuilder struct {
	file  *os.File
	cmd   *cfg.Command
	node  *psast.File
	ident string
}

var (
	suggestionCleanWorkspace = fmt.Sprintf("\nplease run %s to clean up workspace, or add %s flag to execute forcibly",
		style.Keyword("aliax clean"),
		style.Keyword("-f"))
	suggestionSymbolLinkError = fmt.Sprintf("\nif the error is that the file already exists,\nrun %s and try again",
		style.Keyword("aliax clean"))
)

func newPsScriptBuilder(dir, name string, cmd *cfg.Command) (*psScriptBuilder, error) {
	filename := filepath.Join(dir, name+".ps1")

	fp, err := aos.Create(filename)
	if err != nil {
		if errors.Is(err, os.ErrExist) && !initParameter.force {
			log.WithError(err).
				WithField("file", filename).
				WithField("suggestion", suggestionCleanWorkspace).
				Fatal("file already exist")
		} else {
			return nil, err
		}
	}

	builder := &psScriptBuilder{
		file:  fp,
		node:  &psast.File{},
		cmd:   cmd,
		ident: name,
	}

	builder.node.Append(psast.Docs(copyright))

	return builder, nil
}

func (b *psScriptBuilder) close() error {
	return b.file.Close()
}

func nodeString(node psast.Node) string {
	buf := strings.Builder{}
	psast.Print(node, &buf)
	return buf.String()
}

func stmtString(stmts []psast.Stmt) []string {
	res := []string{}
	for _, s := range stmts {
		buf := strings.Builder{}
		psast.Print(s, &buf)
		res = append(res, buf.String())
	}
	return res
}

func (b *psScriptBuilder) generateExtension(ident, cmdName string, level int, cmd *cfg.Command) []psast.Stmt {
	log.Tracef("enter generateExtension.%s", cmdName)
	subCommand := b.collectSubCommand4Extend(ident, level, cmd)

	if level > 0 {
		ifStmt := psast.IfStatement()
		ifStmt.Cond = psast.BinaryExpression(psast.RefExpression(psast.IndexExpression(psast.Identifier("args"), psast.Number(0))), token.EQ, psast.String(cmdName))
		ifStmt.Body.List = append(ifStmt.Body.List, &psast.AssignStmt{Lhs: psast.RefRaw("args"), Rhs: &psast.RefExpr{X: &psast.IndexExpr{X: psast.Identifier("args"), Key: &psast.BinaryExpr{X: psast.Number(1), Op: token.DOUBLE_DOT, Y: &psast.RefExpr{X: &psast.SelectorExpr{X: psast.Identifier("args"), Sel: psast.Identifier("Length")}}}}}})
		ifStmt.Body.List = append(ifStmt.Body.List, b.buildBlockSmt(subCommand, ident, cmd)...)

		log.WithField("output", nodeString(ifStmt)).Tracef("exit generateExtension.%s", cmdName)
		return []psast.Stmt{ifStmt}
	}

	log.WithField("output", stmtString(subCommand)).Tracef("exit generateExtension.%s", cmdName)
	return b.buildBlockSmt(subCommand, ident, cmd)
}

func (b *psScriptBuilder) generateCommand(ident, cmdName string, level int, cmd *cfg.Command) []psast.Stmt {
	if level == 0 {
		cmd.SetName(cmdName)
	} else {
		cmd.SetName(strings.Join(strings.Split(ident, "_"), " "))
	}
	subCommand := b.collectSubCommand4Command(ident, level, cmd)

	if level > 0 {
		ifStmt := psast.IfStatement()
		ifStmt.Cond = psast.BinaryExpression(psast.RefExpression(psast.IndexExpression(psast.Identifier("args"), psast.Number(0))), token.EQ, psast.String(cmdName))
		ifStmt.Body.List = append(ifStmt.Body.List, &psast.AssignStmt{Lhs: psast.RefRaw("args"), Rhs: &psast.RefExpr{X: &psast.IndexExpr{X: psast.Identifier("args"), Key: &psast.BinaryExpr{X: psast.Number(1), Op: token.DOUBLE_DOT, Y: &psast.RefExpr{X: &psast.SelectorExpr{X: psast.Identifier("args"), Sel: psast.Identifier("Length")}}}}}})
		ifStmt.Body.List = append(ifStmt.Body.List, b.buildBlockStmt4Command(subCommand, ident, cmd)...)
		return []psast.Stmt{ifStmt}
	}

	return b.buildBlockStmt4Command(subCommand, ident, cmd)
}

func (b *psScriptBuilder) buildBlockStmt4Command(subCommand []psast.Stmt, ident string, cmd *cfg.Command) []psast.Stmt {
	bs := b.buildBlockSmt(subCommand, ident, cmd)
	if !cmd.DisableHelp {
		bs = append(bs,
			psast.CallStatement(token.None, "Write-Host", psast.String(cmd.HelpCmd(cmd.Name()))),
			psast.CallStatement(token.None, "exit"))
	}
	return bs
}

func (b *psScriptBuilder) collectSubCommand4Extend(ident string, level int, cmd *cfg.Command) []psast.Stmt {
	subCommand := []psast.Stmt{}
	for name, subcmd := range cmd.Command {
		log.Tracef("collectSubCommand4Extend.%s", name)
		subcmd.SetName(name)
		subCommand = append(subCommand, b.generateExtension(fmt.Sprintf("%s_%s", ident, name), name, level+1, subcmd)...)
	}
	return subCommand
}

func (b *psScriptBuilder) collectSubCommand4Command(ident string, level int, cmd *cfg.Command) []psast.Stmt {
	subCommand := []psast.Stmt{}
	for name, subcmd := range cmd.Command {
		subCommand = append(subCommand, b.generateCommand(fmt.Sprintf("%s_%s", ident, name), name, level+1, subcmd)...)
	}
	return subCommand
}

func (b *psScriptBuilder) buildBlockSmt(subCommand []psast.Stmt, ident string, cmd *cfg.Command) (bs []psast.Stmt) {
	log.Tracef("enter buildBlockSmt.%s", ident)
	bs = b.buildArgsStmt(ident, subCommand, bs)

	log.WithField("output", stmtString(bs)).Tracef("")
	typeDict, bs := b.buildFlagDict(ident, bs, cmd)

	bs = append(bs, psast.AssignStatement(psast.RefRaw("non_matched_args"), psast.Raw("@()")))
	log.WithField("output", stmtString(bs)).Tracef("")

	if len(cmd.Flags) > 0 {
		forStmt := psast.ForStatement(
			psast.BinaryExpression(psast.RefRaw("i"), token.ASSIGN, psast.Number(0)),
			psast.BinaryExpression(psast.RefRaw("i"), token.LT, psast.RefExpression(psast.SelectorExpression(psast.Identifier("args"), psast.Identifier("Length")))),
			psast.IncDecExpression("i", true),
		)
		bs = append(bs, forStmt)

		forStmt.Body.Append(b.collectFlagStmt(ident, cmd))

		bs = b.buildMatchStmt(ident, cmd, bs, typeDict)
	} else {
		for _, matchCase := range cmd.Match {
			matchCase.Run = indexResolver.apply(matchCase.Run, func(matched string) string {
				i, err := strconv.Atoi(matched)
				if err != nil {
					log.WithError(err).Fatal("evaluating value")
				}
				i--
				return fmt.Sprintf(`"$($args[%d])"`, i)
			})
			matchCase.Run = envResolver.apply(matchCase.Run, func(matched string) string {
				return fmt.Sprintf("$env:%s", matched)
			})
			switch pattern := matchCase.Pattern.(type) {
			case string:
				if pattern == "_" || len(pattern) == 0 {
					bs = append(bs,
						psast.CallStatement(token.None, matchCase.Run),
						psast.CallStatement(token.None, "exit"))
				}
			}
		}
	}
	log.Tracef("exit buildBlockSmt.%s", ident)
	return
}

func (b *psScriptBuilder) buildArgsStmt(ident string, subCommand, bs []psast.Stmt) []psast.Stmt {
	for i, sc := range subCommand {
		if i == 0 {
			if i != len(subCommand)-1 {
				bs = append(bs, psast.AssignStatement(psast.RefRaw(fmt.Sprintf("temp_args_%s", ident)), psast.RefRaw("args")))
			}
			bs = append(bs, sc)
		} else {
			bs = append(bs, psast.AssignStatement(psast.RefRaw("args"), psast.RefRaw(fmt.Sprintf("temp_args_%s", ident))))
			bs = append(bs, sc)
		}
	}
	return bs
}

func (b *psScriptBuilder) buildFlagDict(ident string, bs []psast.Stmt, cmd *cfg.Command) (typeDict map[string]flagType, stmts []psast.Stmt) {
	stmts = bs
	typeDict = make(map[string]flagType)
	for _, flag := range cmd.Flags {
		flagIdent := fmt.Sprintf("%s_%s", ident, flag.Name)
		switch flag.Type {
		case "string":
			typeDict[flagIdent] = flagTypeString
			stmts = append(stmts, psast.AssignStatement(psast.RefRaw(flagIdent), psast.NULL))
		case "bool":
			typeDict[flagIdent] = flagTypeBool
			stmts = append(stmts, psast.AssignStatement(psast.RefRaw(flagIdent), psast.FALSE))
		}
	}
	return
}

func (b *psScriptBuilder) collectFlagStmt(ident string, cmd *cfg.Command) psast.Stmt {
	switchStmt := &psast.SwitchStmt{
		Mode: psast.MatchModeRegex,
		Cond: psast.RefExpression(psast.IndexExpression(psast.Identifier("args"), psast.RefRaw("i"))),
		Default: &psast.CaseStmt{
			Body: psast.BlockStatement(&psast.ExprStmt{
				X: psast.BinaryExpression(
					psast.RefRaw("non_matched_args"),
					token.ADD_ASSIGN,
					psast.Raw("$args[$i]"),
				),
			}),
		},
	}

	for _, flag := range cmd.Flags {
		flagIdent := fmt.Sprintf("%s_%s", ident, flag.Name)
		alias := []string{}
		if len(flag.Alias) == 0 {
			alias = append(alias, flag.Name)
		} else {
			alias = append(alias, flag.Alias...)
		}
		aslices.MapInPlace(alias, regexp.QuoteMeta)
		rule := strings.Join(alias, "|")
		caseStmt := psast.CaseStatement(psast.String(rule))
		switchStmt.Cases = append(switchStmt.Cases, caseStmt)
		switch flag.Type {
		case "string":
			caseStmt.Body.Append(
				psast.AssignStatement(
					psast.RefRaw(flagIdent),
					psast.IndexExpression(
						psast.RefRaw("args"), psast.BinaryExpression(psast.RefRaw("i"), token.ADD, psast.Number(1)))))
			caseStmt.Body.Append(&psast.ExprStmt{X: psast.IncDecExpression("i", true)})
		case "bool":
			caseStmt.Body.Append(psast.AssignStatement(psast.RefRaw(flagIdent), psast.TRUE))
			// caseStmt.Body.List = append(caseStmt.Body.List, &psast.ExprStmt{
			// 	X: &psast.IncDecExpr{
			// 		X:  &psast.RefExpr{X: &psast.Ident{Name: "i"}},
			// 		Op: token.Inc,
			// 	},
			// })
		}
	}
	return switchStmt
}

func (b *psScriptBuilder) buildMatchStmt(ident string, cmd *cfg.Command, bs []psast.Stmt, typeDict map[string]flagType) []psast.Stmt {
	match := []sortedMatchCase{}
	var defaultMatchCase *sortedMatchCase
	for _, matchCase := range cmd.Match {
		if len(matchCase.Platform) > 0 && matchCase.Platform != "powershell" {
			continue
		}
		names := []string{}
		matchCase.Run = indexResolver.apply(matchCase.Run, func(matched string) string {
			i, _ := strconv.Atoi(matched)
			i--
			return fmt.Sprintf(`"$($args[%d])"`, i)
		})
		matchCase.Run = namedResolver.apply(matchCase.Run, func(matched string) string {
			return fmt.Sprintf("$%s_%s", ident, matched)
		})
		matchCase.Run = envResolver.apply(matchCase.Run, func(matched string) string {
			return fmt.Sprintf("$env:%s", matched)
		})
		switch pattern := matchCase.Pattern.(type) {
		case string:
			if pattern == "_" || len(pattern) == 0 {
				defaultMatchCase = &sortedMatchCase{
					body: matchCase.Run,
				}
				continue
			}
			names = append(names, fmt.Sprintf("%s_%s", ident, pattern))
		case []any:
			for _, v := range pattern {
				if v, ok := v.(string); ok {
					names = append(names, fmt.Sprintf("%s_%s", ident, v))
				}
			}
		}
		if len(names) != 0 {
			match = append(match, sortedMatchCase{weight: len(names), names: names, body: matchCase.Run})
		}
	}
	if len(match) == 0 {
		return bs
	}

	sort.Slice(match, func(i, j int) bool {
		return match[i].weight > match[j].weight
	})

	matchStmt := psast.IfStatement()
	bs = append(bs, matchStmt)
	for i, c := range match {
		var cases psast.Expr
		for _, name := range c.names {
			var cond psast.Expr
			switch typeDict[name] {
			case flagTypeString:
				cond = psast.BinaryExpression(
					psast.NULL,
					token.NE,
					psast.RefRaw(name),
				)
			case flagTypeBool:
				cond = psast.BinaryExpression(
					psast.RefRaw(name),
					token.NE,
					psast.FALSE,
				)
			}

			if cases == nil {
				cases = cond
			} else {
				cases = psast.BinaryExpression(cases, token.AND, cond)
			}
		}
		matchStmt.Cond = cases
		if len(c.body) > 0 {
			lines := strings.Split(c.body, "\n")
			for i, line := range lines {
				if i == len(lines)-1 && len(line) == 0 {
					continue
				}
				matchStmt.Body.Append(psast.CallStatement(token.None, line))
			}
		}
		matchStmt.Body.Append(psast.CallStatement(token.None, "exit"))
		if i != len(match)-1 {
			ifstmt := psast.IfStatement()
			matchStmt.Else = ifstmt
			matchStmt = ifstmt
		}
	}

	if defaultMatchCase != nil {
		matchStmt.Else = &psast.BlockStmt{
			List: []psast.Stmt{
				&psast.ExprStmt{X: psast.Identifier(defaultMatchCase.body)},
			},
		}
	}

	return bs
}

var executable = "executable"

const copyright = " Code generated by [aliax](github.com/ansurfen/aliax). DO NOT EDIT"

func init() {
	aliaxCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&initParameter.global, "global", "g", false, "Apply the initialization globally, affecting the entire system instead of the current project")
	initCmd.PersistentFlags().BoolVarP(&initParameter.force, "force", "f", false, "Force the initialization, bypassing confirmation prompts")
	initCmd.PersistentFlags().BoolVarP(&initParameter.verbose, "verbose", "v", false, "Enable verbose output")
	initCmd.PersistentFlags().StringVarP(&initParameter.template, "template", "t", "", "Specify a template to use for initialization")
	initCmd.PersistentFlags().BoolVarP(&initParameter.save, "save", "s", false, "Backup the current executed YAML to the template directory")
}

var (
	namedResolver = newResolver(`\{\{\s*\.(\w+)\s*\}\}`)
	indexResolver = newResolver(`\{\{\s*\$(\d+)\s*\}\}`)
	envResolver   = newResolver(`\{\{\s*\$\w+\.(\w+)\s*\}\}`)
)

type resolver struct {
	pattern *regexp.Regexp
}

func newResolver(str string) *resolver {
	return &resolver{
		pattern: regexp.MustCompile(str),
	}
}

func (r *resolver) apply(raw string, hd func(matched string) string) string {
	return r.pattern.ReplaceAllStringFunc(raw, func(s string) string {
		matched := r.pattern.FindAllStringSubmatch(s, -1)
		if len(matched) > 0 && len(matched[0]) > 1 {
			return hd(matched[0][1])
		}
		return s
	})
}

type sortedMatchCase struct {
	weight int
	names  []string
	body   string
}

func setGlobal() error {
	var (
		value string
		err   error
	)
	if runtime.GOOS == "windows" {
		value, err = globalenv.Get("Path")
	} else {
		value, err = globalenv.Get("PATH")
	}
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	path = filepath.Join(path, "run-scripts")
	if !strings.Contains(value, path) {
		var output []byte
		if runtime.GOOS == "windows" {
			output, err = globalenv.Set("Path", fmt.Sprintf("%s;%s", path, value))
		} else {
			output, err = globalenv.Set("PATH", fmt.Sprintf("$PATH:%s", path))
		}
		if err != nil {
			log.WithField("output", string(output)).Error("setting Path")
			return err
		}
	}
	log.WithField("path", path).Warn("the environment variable is already set")
	return nil
}
