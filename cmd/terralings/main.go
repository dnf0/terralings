package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dnf0/terralings/exercises"
	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/doctor"
	"github.com/dnf0/terralings/internal/lsp"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/search"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/tour"
	"github.com/dnf0/terralings/internal/tui"
	"github.com/dnf0/terralings/internal/ui"
	"github.com/dnf0/terralings/internal/watcher"
	"github.com/spf13/cobra"
)

// Version is the current release version of terralings.
const Version = "v0.2.0"

var (
	binOverride        string
	stateOverride      string
	hintIndex          int
	initForce          bool
	resetDir           string
	watchJSON          bool
	watchInteractive   bool
	tourStep           int
	tourNonInteractive bool
	tourJSON           bool
	doctorJSON         bool
)

// NewRootCmd constructs and returns the root Cobra command and its subcommands.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "terralings",
		Short:        "Terralings - Interactive CLI learning environment for Terraform & OpenTofu",
		Long:         "Terralings guides you through hands-on exercises to master Terraform & OpenTofu HCL and workflows.",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatBanner())
			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatFirstRunWelcome())
			_ = cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVar(&binOverride, "bin", "", "Custom path to tofu or terraform binary")
	rootCmd.PersistentFlags().StringVar(&stateOverride, "state", "", "Custom path to state file (default .terralings/state.json)")

	// watch command
	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Start continuous interactive watch mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				bin = ""
			}
			if watchInteractive {
				return tui.RunTUI(cmd.Context(), bin, "exercises", stateOverride, os.Stdin, os.Stdout)
			}
			if bin == "" {
				bin, err = detector.DetectBinary(binOverride)
				if err != nil {
					return err
				}
			}
			store, err := state.NewStore(stateOverride)
			if err != nil {
				return err
			}
			if watchJSON {
				return watcher.RunWatchJSON(context.Background(), runner.NewRunner(bin), manifest.GetManifest().AllExercises(), store, "exercises", os.Stdout)
			}
			return watcher.RunWatchWithStore(context.Background(), runner.NewRunner(bin), manifest.GetManifest().AllExercises(), store, "exercises", os.Stdout)
		},
	}
	watchCmd.Flags().BoolVar(&watchJSON, "json", false, "Emit structured NDJSON stream of evaluation events")
	watchCmd.Flags().BoolVarP(&watchInteractive, "interactive", "i", false, "Start interactive full-screen TUI dashboard")

	completeExerciseNames := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var completions []string
		for _, ex := range manifest.GetManifest().AllExercises() {
			if strings.HasPrefix(ex.Name, toComplete) {
				completions = append(completions, fmt.Sprintf("%s\t%s", ex.Name, ex.Title))
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	// run command
	runCmd := &cobra.Command{
		Use:               "run [exercise_name]",
		Short:             "Run verification on a single exercise",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeExerciseNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				return err
			}

			ex := manifest.GetExerciseByName(args[0])
			if ex == nil {
				return fmt.Errorf("exercise '%s' not found", args[0])
			}

			r := runner.NewRunner(bin)
			res := r.Run(*ex)
			fmt.Fprint(cmd.OutOrStdout(), ui.FormatResult(res))

			store, err := state.NewStore(stateOverride)
			if err == nil {
				_ = store.RecordAttempt(ex.Name, ex.ChapterName, res.Passed)
			}

			if !res.Passed {
				return fmt.Errorf("exercise %s did not pass", ex.Name)
			}
			return nil
		},
	}

	// hint command
	hintCmd := &cobra.Command{
		Use:               "hint [exercise_name]",
		Short:             "Show progressive hints for an exercise",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeExerciseNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			ex := manifest.GetExerciseByName(args[0])
			if ex == nil {
				return fmt.Errorf("exercise '%s' not found", args[0])
			}

			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatHint(ex, hintIndex))

			store, err := state.NewStore(stateOverride)
			if err == nil {
				_ = store.RecordHint(ex.Name, ex.ChapterName, hintIndex+1)
			}
			return nil
		},
	}
	hintCmd.Flags().IntVarP(&hintIndex, "index", "i", 0, "Zero-based index of the hint to display")

	// stats command
	statsCmd := &cobra.Command{
		Use:   "stats",
		Short: "Display progress and learning analytics",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := state.NewStore(stateOverride)
			if err != nil {
				return err
			}
			m := manifest.GetManifest()
			summary := store.GetAnalytics(m)
			fmt.Fprint(cmd.OutOrStdout(), ui.FormatAnalytics(summary))
			return nil
		},
	}

	// list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all curriculum chapters and exercises",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatBanner())
			m := manifest.GetManifest()
			fmt.Fprint(cmd.OutOrStdout(), ui.FormatChapterList(m, nil))
		},
	}

	// verify command
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify and evaluate all exercises across the curriculum",
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				return err
			}

			m := manifest.GetManifest()
			all := m.AllExercises()
			r := runner.NewRunner(bin)

			statuses := make(map[string]models.ExerciseStatus)
			completedCount := 0

			for _, ex := range all {
				res := r.Run(ex)
				if res.Passed {
					statuses[ex.Name] = models.StatusCompleted
					completedCount++
				} else if res.HasNotDoneMarker {
					statuses[ex.Name] = models.StatusInProgress
				} else {
					statuses[ex.Name] = models.StatusFailed
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatBanner())
			fmt.Fprint(cmd.OutOrStdout(), ui.FormatChapterList(m, statuses))
			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatProgress(completedCount, len(all)))

			if completedCount == len(all) {
				fmt.Fprintln(cmd.OutOrStdout(), "\n🎉 Congratulations! You have completed all Terralings exercises! 🎉")
			}

			return nil
		},
	}

	completeSearchQueries := func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		completions, directive := completeExerciseNames(cmd, args, toComplete)
		if len(args) != 0 {
			return completions, directive
		}
		m := manifest.GetManifest()
		for _, ch := range m.Chapters {
			if strings.HasPrefix(ch.Name, toComplete) {
				completions = append(completions, fmt.Sprintf("%s\tChapter: %s", ch.Name, ch.Title))
			}
		}
		return completions, directive
	}

	// search command
	searchCmd := &cobra.Command{
		Use:               "search [query]",
		Short:             "Search curriculum exercises by concept, keyword, or chapter",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeSearchQueries,
		Run: func(cmd *cobra.Command, args []string) {
			m := manifest.GetManifest()
			results := search.SearchExercises(m, args[0])
			fmt.Fprint(cmd.OutOrStdout(), ui.FormatSearchResults(args[0], results))
		},
	}

	// completion command
	completionCmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Aliases:   []string{"completions"},
		Short:     "Generate shell completion scripts",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]
			switch shell {
			case "bash":
				return rootCmd.GenBashCompletionV2(cmd.OutOrStdout(), true)
			case "zsh":
				return rootCmd.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return rootCmd.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell '%s', choose from bash, zsh, fish, powershell", shell)
			}
		},
	}

	// version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print Terralings version and detected binary information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "terralings %s\n", Version)
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Detected binary: none found")
				return
			}

			ver, err := detector.GetBinaryVersion(bin)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Detected binary: %s\n", bin)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Detected binary: %s (%s)\n", bin, ver)
		},
	}

	// init command
	initCmd := &cobra.Command{
		Use:   "init [target_dir]",
		Short: "Initialize the complete Terralings exercise curriculum into a directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDir := "exercises"
			if len(args) > 0 {
				targetDir = args[0]
			}

			if err := exercises.ExtractAll(targetDir, initForce); err != nil {
				return err
			}

			m := manifest.GetManifest()
			all := m.AllExercises()
			fmt.Fprintf(cmd.OutOrStdout(), "✨ Successfully initialized %d exercises into '%s'!\n\n", len(all), targetDir)
			fmt.Fprintln(cmd.OutOrStdout(), ui.FormatFirstRunWelcome())
			fmt.Fprintln(cmd.OutOrStdout(), "To get started, run:")
			fmt.Fprintln(cmd.OutOrStdout(), "  terralings watch")
			return nil
		},
	}
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing files if target directory is not empty")

	// reset command
	resetCmd := &cobra.Command{
		Use:               "reset [exercise_name]",
		Short:             "Reset an exercise back to its initial starting code",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeExerciseNames,
		RunE: func(cmd *cobra.Command, args []string) error {
			exerciseName := args[0]
			ex := manifest.GetExerciseByName(exerciseName)
			if ex == nil {
				return fmt.Errorf("exercise '%s' not found", exerciseName)
			}

			baseDir := resetDir
			if baseDir == "" {
				baseDir = "exercises"
			}

			if err := exercises.ResetExercise(exerciseName, baseDir); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "🔄 Reset exercise '%s' (%s) back to original template.\n", ex.Name, ex.Path)
			return nil
		},
	}
	resetCmd.Flags().StringVarP(&resetDir, "dir", "d", "exercises", "Base exercises directory")

	// lsp command
	lspCmd := &cobra.Command{
		Use:   "lsp",
		Short: "Start Language Server Protocol (LSP) daemon over stdio",
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				bin = ""
			}
			store, _ := state.NewStore(stateOverride)
			m := manifest.GetManifest()
			r := runner.NewRunner(bin)
			srv := lsp.NewServer(r, m, store)
			return srv.RunWithContext(cmd.Context(), os.Stdin, os.Stdout)
		},
	}

	// tui command
	tuiCmd := &cobra.Command{
		Use:   "tui",
		Short: "Start interactive full-screen terminal dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			bin, err := detector.DetectBinary(binOverride)
			if err != nil {
				bin = ""
			}
			return tui.RunTUI(cmd.Context(), bin, "exercises", stateOverride, os.Stdin, os.Stdout)
		},
	}

	// tour command
	tourCmd := &cobra.Command{
		Use:   "tour",
		Short: "Start the interactive guided onboarding tour",
		RunE: func(cmd *cobra.Command, args []string) error {
			t := tour.NewTour(cmd.OutOrStdout(), cmd.InOrStdin())
			t.NonInteractive = tourNonInteractive
			t.JSONMode = tourJSON
			return t.Run(cmd.Context(), tourStep)
		},
	}
	tourCmd.Flags().IntVar(&tourStep, "step", 0, "Specific tour step to render (1-5)")
	tourCmd.Flags().BoolVar(&tourNonInteractive, "non-interactive", false, "Render tour steps without interactive prompt")
	tourCmd.Flags().BoolVar(&tourJSON, "json", false, "Emit tour content as structured JSON")

	// doctor command
	doctorCmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostics to verify environment and workspace readiness",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, _ := os.Getwd()
			report := doctor.RunDiagnostics(cwd, binOverride, stateOverride)
			if doctorJSON {
				data, err := json.MarshalIndent(report, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), doctor.FormatReport(report))
			return nil
		},
	}
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Emit diagnostic report as JSON")

	rootCmd.AddCommand(watchCmd, runCmd, hintCmd, statsCmd, listCmd, verifyCmd, versionCmd, initCmd, resetCmd, searchCmd, completionCmd, lspCmd, tuiCmd, tourCmd, doctorCmd)
	return rootCmd
}

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
