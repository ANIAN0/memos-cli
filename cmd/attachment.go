package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ANIAN0/memos-cli/internal/client"
	"github.com/ANIAN0/memos-cli/pkg/httpclient"
	"github.com/ANIAN0/memos-cli/pkg/output"
)

var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage attachments",
	Long:  `Upload, list, get, and delete attachments on Memos.`,
}

var attachmentUploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload an attachment",
	Long:  `Upload a file as an attachment to Memos.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		ac := &client.AttachmentClient{C: c}
		att, err := ac.Upload(cmd.Context(), filePath)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)
		return out.PrintObject(att)
	},
}

var attachmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List attachments",
	Long:  `List all attachments on Memos.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		ac := &client.AttachmentClient{C: c}
		atts, err := ac.List(cmd.Context(), attachmentPageSize)
		if err != nil {
			return err
		}

		mode := output.ModeText
		if jsonMode {
			mode = output.ModeJSON
		}
		out := output.New(mode)

		items := make([]any, len(atts))
		for i, a := range atts {
			items[i] = a
		}
		return out.PrintList(items)
	},
}

var attachmentGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Download an attachment",
	Long:  `Download an attachment to a local file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		ac := &client.AttachmentClient{C: c}
		return ac.Get(cmd.Context(), id, attachmentOutput)
	},
}

var attachmentDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an attachment",
	Long:  `Delete an attachment by ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		c := httpclient.New(cfg.InstanceURL,
			httpclient.WithTimeout(getTimeout()),
			httpclient.WithVerbose(verbose),
			httpclient.WithToken("Bearer " + cfg.AccessToken),
			httpclient.WithAuthHeader("Authorization"),
		)

		ac := &client.AttachmentClient{C: c}
		return ac.Delete(cmd.Context(), id)
	},
}

var (
	attachmentPageSize int
	attachmentOutput   string
)

func init() {
	attachmentListCmd.Flags().IntVar(&attachmentPageSize, "page-size", 0, "page size")
	attachmentGetCmd.Flags().StringVar(&attachmentOutput, "output", "", "output file path (default: original filename)")

	attachmentCmd.AddCommand(attachmentUploadCmd)
	attachmentCmd.AddCommand(attachmentListCmd)
	attachmentCmd.AddCommand(attachmentGetCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)

	rootCmd.AddCommand(attachmentCmd)
}