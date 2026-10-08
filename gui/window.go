package gui

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/helwan-linux/hpi/installer"
	"github.com/helwan-linux/hpi/i18n"
	packageinfo "github.com/helwan-linux/hpi/package"
)

type Window struct {
	window       *gtk.Window
	titleLabel   *gtk.Label
	statusLabel  *gtk.Label
	packageLabel *gtk.Label
	openBtn      *gtk.Button
	installBtn   *gtk.Button
	cancelBtn    *gtk.Button

	packagePath string
	pkg         *packageinfo.Package
	transaction *installer.Transaction
	resolution  *installer.Resolution

	operationCancel context.CancelFunc
	busy            bool

	tr *i18n.Translator
}

func NewWindow() (*Window, error) {
	window, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, fmt.Errorf("create window: %w", err)
	}

	tr := i18n.New()

	window.SetTitle(tr.Get("app_title"))
	window.SetDefaultSize(620, 420)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetResizable(false)

	result := &Window{
		window:      window,
		transaction: installer.NewTransaction(),
		tr:          tr,
	}

	if err := result.build(); err != nil {
		window.Destroy()
		return nil, err
	}

	window.Connect("destroy", func() {
		if result.operationCancel != nil {
			result.operationCancel()
			result.operationCancel = nil
		}

		gtk.MainQuit()
	})

	return result, nil
}

func (w *Window) build() error {
	mainBox, err := gtk.BoxNew(
		gtk.ORIENTATION_VERTICAL,
		0,
	)
	if err != nil {
		return fmt.Errorf("create main box: %w", err)
	}

	mainBox.SetBorderWidth(24)

	headerBox, err := gtk.BoxNew(
		gtk.ORIENTATION_VERTICAL,
		8,
	)
	if err != nil {
		return fmt.Errorf("create header box: %w", err)
	}

	w.titleLabel, err = gtk.LabelNew(
		w.tr.Get("app_title"),
	)
	if err != nil {
		return fmt.Errorf("create title label: %w", err)
	}

	w.titleLabel.SetHAlign(gtk.ALIGN_START)

	headerBox.PackStart(
		w.titleLabel,
		false,
		false,
		0,
	)

	w.packageLabel, err = gtk.LabelNew(
		w.tr.Get("package_open_prompt"),
	)
	if err != nil {
		return fmt.Errorf("create package label: %w", err)
	}

	w.packageLabel.SetHAlign(gtk.ALIGN_START)
	w.packageLabel.SetLineWrap(true)

	headerBox.PackStart(
		w.packageLabel,
		false,
		false,
		0,
	)

	mainBox.PackStart(
		headerBox,
		false,
		false,
		0,
	)

	separator, err := gtk.SeparatorNew(
		gtk.ORIENTATION_HORIZONTAL,
	)
	if err != nil {
		return fmt.Errorf("create separator: %w", err)
	}

	mainBox.PackStart(
		separator,
		false,
		false,
		20,
	)

	w.statusLabel, err = gtk.LabelNew(
		w.tr.Get("ready"),
	)
	if err != nil {
		return fmt.Errorf("create status label: %w", err)
	}

	w.statusLabel.SetHAlign(gtk.ALIGN_START)
	w.statusLabel.SetLineWrap(true)
	w.statusLabel.SetSelectable(true)

	mainBox.PackStart(
		w.statusLabel,
		true,
		true,
		0,
	)

	buttonBox, err := gtk.BoxNew(
		gtk.ORIENTATION_HORIZONTAL,
		8,
	)
	if err != nil {
		return fmt.Errorf("create button box: %w", err)
	}

	buttonBox.SetHAlign(gtk.ALIGN_END)

	w.openBtn, err = gtk.ButtonNewWithLabel(
		w.tr.Get("open_package"),
	)
	if err != nil {
		return fmt.Errorf("create open button: %w", err)
	}

	w.openBtn.Connect("clicked", func() {
		if err := w.choosePackage(); err != nil {
			w.ShowError(
				w.tr.Get("package_error"),
				err.Error(),
			)
		}
	})

	buttonBox.PackStart(
		w.openBtn,
		false,
		false,
		0,
	)

	w.cancelBtn, err = gtk.ButtonNewWithLabel(
		w.tr.Get("cancel"),
	)
	if err != nil {
		return fmt.Errorf("create cancel button: %w", err)
	}

	w.cancelBtn.SetSensitive(false)

	w.cancelBtn.Connect("clicked", func() {
		if w.operationCancel != nil {
			w.operationCancel()
			w.operationCancel = nil
		}

		w.window.Destroy()
	})

	buttonBox.PackStart(
		w.cancelBtn,
		false,
		false,
		0,
	)

	w.installBtn, err = gtk.ButtonNewWithLabel(
		w.tr.Get("install"),
	)
	if err != nil {
		return fmt.Errorf("create install button: %w", err)
	}

	w.installBtn.SetSensitive(false)

	w.installBtn.Connect("clicked", func() {
		w.install()
	})

	buttonBox.PackStart(
		w.installBtn,
		false,
		false,
		0,
	)

	mainBox.PackEnd(
		buttonBox,
		false,
		false,
		0,
	)

	w.window.Add(mainBox)

	return nil
}

func (w *Window) ShowAll() {
	w.window.ShowAll()
}

func (w *Window) choosePackage() error {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		w.tr.Get("open_package"),
		w.window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		w.tr.Get("cancel"),
		gtk.RESPONSE_CANCEL,
		w.tr.Get("open_package"),
		gtk.RESPONSE_ACCEPT,
	)
	if err != nil {
		return fmt.Errorf("create file chooser: %w", err)
	}

	defer dialog.Destroy()

	filter, err := gtk.FileFilterNew()
	if err != nil {
		return fmt.Errorf("create package filter: %w", err)
	}

	filter.SetName("Arch/Helwan Packages (*.pkg.tar.zst)")
	filter.AddPattern("*.pkg.tar.zst")

	dialog.AddFilter(filter)

	response := dialog.Run()

	if response != gtk.RESPONSE_ACCEPT {
		return nil
	}

	path := dialog.GetFilename()

	if path == "" {
		return fmt.Errorf(
			"%s",
			w.tr.Get("package_path_empty"),
		)
	}

	return w.OpenPackage(path)
}

func (w *Window) OpenPackage(path string) error {
	if path == "" {
		return fmt.Errorf(
			"%s",
			w.tr.Get("package_path_empty"),
		)
	}

	pkg, err := packageinfo.Open(path)
	if err != nil {
		return err
	}

	w.packagePath = path
	w.pkg = pkg
	w.resolution = nil

	displayName := pkg.Metadata.FullName()

	if displayName == "" {
		displayName = filepath.Base(path)
	}

	w.SetPackageName(displayName)

	w.SetInstallEnabled(false)
	w.SetCancelEnabled(true)

	w.SetStatus(
		fmt.Sprintf(
			w.tr.Get("package_loaded")+"\n"+
				w.tr.Get("checking_dependencies"),
			displayName,
		),
	)

	w.checkDependencies()

	return nil
}

func (w *Window) checkDependencies() {
	if w.pkg == nil {
		w.ShowError(
			w.tr.Get("package_error"),
			w.tr.Get("no_package_loaded"),
		)
		return
	}

	if w.busy {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.operationCancel = cancel
	w.busy = true

	progress, err := NewProgress(
		w.window,
		w.tr.Get("checking_dependencies_title"),
	)
	if err != nil {
		cancel()
		w.operationCancel = nil
		w.busy = false

		w.ShowError(
			w.tr.Get("dependency_check_failed_title"),
			err.Error(),
		)
		return
	}

	progress.SetMessage(
		w.tr.Get("checking_package_dependencies"),
	)
	progress.SetIndeterminate(true)
	progress.Show()

	pkg := w.pkg

	go func() {
		resolution, err := w.transaction.Check(
			ctx,
			pkg,
		)

		glib.IdleAdd(func() {
			progress.Close()

			w.operationCancel = nil
			w.busy = false

			if err != nil {
				w.SetInstallEnabled(false)
				w.SetStatus(
					w.tr.Get("dependency_check_failed"),
				)

				w.ShowError(
					w.tr.Get("dependency_check_failed_title"),
					err.Error(),
				)
				return
			}

			w.resolution = resolution

			if resolution.HasMissing() {
				w.SetInstallEnabled(false)
				w.SetStatus(
					w.formatResolution(resolution),
				)

				w.ShowError(
					w.tr.Get("missing_dependencies_title"),
					w.formatMissingDependencies(resolution),
				)
				return
			}

			w.SetInstallEnabled(true)
			w.SetStatus(
				w.formatResolution(resolution),
			)
		})
	}()
}

func (w *Window) install() {
	if w.pkg == nil {
		w.ShowError(
			w.tr.Get("installation_error_title"),
			w.tr.Get("no_package_loaded"),
		)
		return
	}

	if w.resolution == nil {
		w.ShowError(
			w.tr.Get("installation_error_title"),
			w.tr.Get("dependency_check_not_completed"),
		)
		return
	}

	if w.resolution.HasMissing() {
		w.ShowError(
			w.tr.Get("installation_error_title"),
			w.formatMissingDependencies(w.resolution),
		)
		return
	}

	if w.busy {
		return
	}

	confirmed := AskConfirmation(
		w.window,
		w.tr.Get("install_package_title"),
		fmt.Sprintf(
			w.tr.Get("install_confirmation"),
			w.pkg.Metadata.FullName(),
		),
	)

	if !confirmed {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	w.operationCancel = cancel
	w.busy = true

	w.SetInstallEnabled(false)
	w.SetCancelEnabled(true)
	w.SetStatus(
		w.tr.Get("installing_package"),
	)

	progress, err := NewProgress(
		w.window,
		w.tr.Get("installing_package_title"),
	)
	if err != nil {
		cancel()
		w.operationCancel = nil
		w.busy = false

		w.SetInstallEnabled(true)

		w.ShowError(
			w.tr.Get("installation_error_title"),
			err.Error(),
		)
		return
	}

	progress.SetMessage(
		w.tr.Get("installing_package"),
	)
	progress.SetIndeterminate(true)
	progress.Show()

	packagePath := w.packagePath

	go func() {
		_, err := w.transaction.Install(
			ctx,
			packagePath,
		)

		glib.IdleAdd(func() {
			progress.Close()

			w.operationCancel = nil
			w.busy = false

			if err != nil {
				w.SetInstallEnabled(true)
				w.SetStatus(
					w.tr.Get("installation_failed"),
				)

				w.ShowError(
					w.tr.Get("installation_failed_title"),
					err.Error(),
				)
				return
			}

			w.SetInstallEnabled(false)

			successMessage := fmt.Sprintf(
				w.tr.Get("installation_completed"),
				w.pkg.Metadata.FullName(),
			)

			w.SetStatus(successMessage)

			ShowInfo(
				w.window,
				w.tr.Get("installation_complete_title"),
				fmt.Sprintf(
					w.tr.Get("installation_success"),
					w.pkg.Metadata.FullName(),
				),
			)
		})
	}()
}

func (w *Window) formatResolution(
	resolution *installer.Resolution,
) string {
	if resolution == nil {
		return w.tr.Get("dependency_check_failed")
	}

	if len(resolution.Dependencies) == 0 {
		return w.tr.Get("no_dependencies")
	}

	var builder strings.Builder

	builder.WriteString(
		w.tr.Get("dependency_check_completed"),
	)
	builder.WriteString("\n")

	if len(resolution.Installed) > 0 {
		builder.WriteString(
			fmt.Sprintf(
				w.tr.Get("already_installed"),
				len(resolution.Installed),
			),
		)
		builder.WriteString("\n")
	}

	if len(resolution.Available) > 0 {
		builder.WriteString(
			fmt.Sprintf(
				w.tr.Get("available_repositories"),
				len(resolution.Available),
			),
		)
		builder.WriteString("\n")
	}

	if len(resolution.Missing) > 0 {
		builder.WriteString(
			fmt.Sprintf(
				w.tr.Get("missing"),
				len(resolution.Missing),
			),
		)
		builder.WriteString("\n")
	}

	if resolution.CanInstallOffline() {
		builder.WriteString(
			w.tr.Get("ready_to_install"),
		)
	} else if resolution.RequiresInternet() {
		builder.WriteString(
			w.tr.Get("repository_dependencies"),
		)
	}

	return builder.String()
}

func (w *Window) formatMissingDependencies(
	resolution *installer.Resolution,
) string {
	if resolution == nil || len(resolution.Missing) == 0 {
		return w.tr.Get("no_missing_dependencies")
	}

	var builder strings.Builder
	var dependencies strings.Builder

	for _, dependency := range resolution.Missing {
		dependencies.WriteString(
			fmt.Sprintf(
				w.tr.Get("missing_dependency_line"),
				dependency.Dependency.DisplayName(),
			),
		)
		dependencies.WriteString("\n")
	}

	builder.WriteString(
		fmt.Sprintf(
			w.tr.Get("missing_dependencies"),
			dependencies.String(),
		),
	)

	return builder.String()
}

func (w *Window) SetPackageName(name string) {
	if name == "" {
		w.packageLabel.SetText(
			w.tr.Get("package_open_prompt"),
		)
		return
	}

	w.packageLabel.SetText(name)
}

func (w *Window) SetStatus(status string) {
	if status == "" {
		return
	}

	w.statusLabel.SetText(status)
}

func (w *Window) SetInstallEnabled(enabled bool) {
	w.installBtn.SetSensitive(enabled)
}

func (w *Window) SetCancelEnabled(enabled bool) {
	w.cancelBtn.SetSensitive(enabled)
}

func (w *Window) ShowError(title string, message string) {
	ShowError(w.window, title, message)
}
