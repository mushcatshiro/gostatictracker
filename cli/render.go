
package cli

/*
import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mushcatshiro/gostatictracker/render"
	"github.com/spf13/cobra"
)

var (
	rcMonth int
	rcYear  int
	rDir    string
	rDebug  bool
)

var renderCalendarCmd = &cobra.Command{
	Use:   "rc",
	Short: "Render calendar",
	Run:   cliRenderCalendar,
}

var renderGanttCmd = &cobra.Command{
	Use:   "rg",
	Short: "Render gantt",
	Run:   cliRenderGantt,
}

var renderListCmd = &cobra.Command{
	Use:   "rl",
	Short: "Render list",
	Run:   cliRenderList,
}

func renderDirectoryValidation(rDir string) error {
	info, err := os.Stat(rDir)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Directory specified does not exists: %s", rDir)
	}
	if err != nil {
		return fmt.Errorf("Error during directory check %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Expects a valid path instead got %s", rDir)
	}
	return nil
}

func cliRenderCalendar(cmd *cobra.Command, args []string) {
	if err := renderDirectoryValidation(rDir); err != nil {
		log.Fatalf("%v", err)
	}
	render.RenderCalendar(rcMonth, rcYear, rDir, conn)
}

func cliRenderGantt(cmd *cobra.Command, args []string) {
	if err := renderDirectoryValidation(rDir); err != nil {
		log.Fatalf("%v", err)
	}
	render.RenderGantt(rDir, conn, rDebug)
}

func cliRenderList(cmd *cobra.Command, args []string) {
	if err := renderDirectoryValidation(rDir); err != nil {
		log.Fatalf("%v", err)
	}
	render.RenderList(rDir, conn)
}
*/
